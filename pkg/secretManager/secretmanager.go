package secretmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func GetAPIKey(secretName string) (string, string, error) {
	region := "us-east-1"

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Fatal(err.Error())
	}

	var secretMap map[string]string

	if err := json.Unmarshal([]byte(*result.SecretString), &secretMap); err != nil {
		log.Fatal(err.Error())
	}

	publicKey, okPublic := secretMap["MONGODB_ATLAS_PUBLIC_KEY"]
	if !okPublic {
		return "", "", fmt.Errorf("public key not found in secret")
	}

	privateKey, okPrivate := secretMap["MONGODB_ATLAS_PRIVATE_KEY"]
	if !okPrivate {
		return "", "", fmt.Errorf("private key not found in secret")
	}

	return publicKey, privateKey, nil
}
