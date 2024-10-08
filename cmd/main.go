package main

import (
	"context"
	"fmt"
	"log"

	internalAtlas "github.com/PicPay/dbre/automations/atlas/autoscaling/internal/atlas"
	secretmanager "github.com/PicPay/dbre/automations/atlas/autoscaling/pkg/secretManager"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	SecretName string `envconfig:"SECRET_NAME" required:"true"`
}

var config Config

func init() {
	var err error
	config, err = NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
}

func NewConfig() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	return cfg, err
}

func handleRequest(ctx context.Context, event *internalAtlas.MyEvent) (*string, error) {
	if err := validateEvent(event); err != nil {
		return logAndReturnError("validating event", err)
	}

	publicKey, privateKey, err := secretmanager.GetAPIKey(config.SecretName)

	if err != nil {
		return logAndReturnError("getting secret", err)
	}

	sdk, err := internalAtlas.NewClient(publicKey, privateKey)
	if err != nil {
		return logAndReturnError("instantiating new client", err)
	}

	project, err := internalAtlas.GetProjectByName(ctx, sdk, event.Project)
	if err != nil {
		return logAndReturnError("getting project", err)
	}

	err = internalAtlas.AutoScaling(ctx, *project, event, sdk)
	if err != nil {
		return logAndReturnError("auto scaling", err)
	}

	message := fmt.Sprintf("Cluster %s in project %s scaled.", event.Cluster, event.Project)

	log.Println(message)
	return &message, nil
}

func validateEvent(event *internalAtlas.MyEvent) error {
	if event.Project == "" || event.Cluster == "" {
		return fmt.Errorf("project and cluster must not be empty")
	}

	if (event.ElectableSpecs == nil || (event.ElectableSpecs.InstanceSize == "" && event.ElectableSpecs.DiskIOPS == 0 && event.ElectableSpecs.EBSVolumeType == "")) &&
		(event.ReadOnlySpecs == nil || (event.ReadOnlySpecs.InstanceSize == "" && event.ReadOnlySpecs.DiskIOPS == 0 && event.ReadOnlySpecs.EBSVolumeType == "")) &&
		(event.AnalyticsSpecs == nil || (event.AnalyticsSpecs.InstanceSize == "" && event.AnalyticsSpecs.DiskIOPS == 0 && event.AnalyticsSpecs.EBSVolumeType == "")) {
		return fmt.Errorf("at least one scaling parameter (InstanceSize, DiskIOPS, EBSVolumeType) must be provided in any specs")
	}

	return nil
}

func logAndReturnError(action string, err error) (*string, error) {
	log.Printf("Error %s: %v", action, err)
	return nil, fmt.Errorf("error %s: %w", action, err)
}

func main() {
	lambda.Start(handleRequest)
}
