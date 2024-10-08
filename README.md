# Autoscaling MongoDB Atlas

This repository contains code to automatically scale MongoDB Atlas clusters, following specific business rules for different specifications (`electableSpecs`, `readOnlySpecs`, `analyticsSpecs`).

## Description

"The main function of this project is to handle the auto-scaling of MongoDB clusters based on the parameters received. It uses AWS Lambda to execute these operations automatically, with EventBridge rules triggering the Lambda based on specified schedules."

### Project Structure

- **cmd/**: Contains the main application code.
- **internal/**: Contains internal logic and helper functions, including business rules for synchronizing parameters.
- **pkg/**: Contains utilities, such as the interface with `secretManager`.
- **testdata/**: Contains test payloads used to simulate auto-scaling events.

## Use Cases

Here are potential use cases for this project:

1. **Peak Time Scaling**

**Scenario**: I have an M10 cluster and would like to scale it up to an M30 during peak hours and revert back to an M10 outside of peak times.

*Payload Example*:

```json
{
  "project": "project0",
  "cluster": "Cluster0",
  "electableSpecs": {
    "instanceSize": "M30",
    "AutoScale": {
      "MinInstanceSize": "M10",
      "MaxInstanceSize": "M30"
    }
  }
}
```

2. **Read-Only Scaling**

**Scenario**: I have a cluster with 3 electable nodes and would like to add 2 read-only nodes.

*Payload Example*:

```json
{
  "project": "project0",
  "cluster": "Cluster0",
  "readOnlySpecs": {
    "NodeCount": 2
  }
}
```

3. **Analytics Scaling**

**Scenario**: I have a cluster with 3 electable nodes and would like to add 2 analytics nodes.

*Payload Example*:

```json
{
  "project": "project0",
  "cluster": "Cluster0",
  "analyticsSpecs": {
    "NodeCount": 2,
    "instanceSize": "M30"
  }
}
```

4. **Increase Disk IOPS**

**Scenario**: I have a cluster with 3 electable nodes and would like to increase the disk IOPS to 3250.

*Payload Example*:

```json
{
  "project": "project0",
  "cluster": "Cluster0",
  "electableSpecs": {
      "diskiops": 3250,
      "ebsVolumeType": "PROVISIONED",
      "instanceSize": "M30",
      "AutoScale": {
          "MinInstanceSize": "M30"
      }
  }
}
```

> **Note**: You need to specify the minimum, as provisioned IOPS is only possible on an M30 instance.

## Supported Parameters

The following parameters can be passed in the payload to control the cluster scaling specifications:

```json
{
  "project": "project0",
  "cluster": "Cluster0",
  "electableSpecs": {
    "diskiops": 3250,
    "ebsVolumeType": "PROVISIONED",
    "instanceSize": "M30",
    "NodeCount": 5,
    "AutoScale": {
      "MinInstanceSize": "M30",
      "MaxInstanceSize": "M60"
    }
  },
  "readOnlySpecs": {
    "NodeCount": 5
  },
  "analyticsSpecs": {
    "NodeCount": 5,
    "instanceSize": "M30",
    "AutoScale": {
      "MinInstanceSize": "M30",
      "MaxInstanceSize": "M60"
    }
  }
}
```

### Parameter Details

|**Parameter**|**Required**|**Description**|
|---------------------|----|----------------------------------------|
|**project**|**Yes**|The name of the MongoDB Atlas project.|
|**cluster**|**Yes**|The name of the cluster to be scaled.|
|**electableSpecs**|**No**|Specifications for the electable instance:|
|-> `diskiops`|**No**|Disk IOPS.|
|-> `ebsVolumeType`|**No**|The EBS volume type (e.g., "PROVISIONED" or "STANDARD").|
|-> `instanceSize`|**No**|Instance size (e.g., "M30").|
|-> `NodeCount`|**No**|Number of electable nodes.|
|-> `AutoScale`|**No**|Auto-scaling configuration|
|--> `MinInstanceSize`|**No**|Minimum instance size for auto-scaling.|
|--> `MaxInstanceSize`|**No**|Maximum instance size for auto-scaling.|
|**readOnlySpecs**|**No**|Specifications for read-only instances:|
|-> `NodeCount`|**No**|Number of read-only nodes.|
|**analyticsSpecs**|**No**|Specifications for analytics instances:|
|-> `NodeCount`|**No**|Number of analytics nodes.|
|-> `instanceSize`|**No**|Instance size for analytics nodes.|
|-> `AutoScale`|**No**|Auto-scaling configuration:|
|--> `MinInstanceSize`|**No**|Minimum instance size for auto-scaling.|
|--> `MaxInstanceSize`|**No**|Maximum instance size for auto-scaling.|

> Note: At least one attribute must be specified.

## Synchronization Between Electable, Read-Only, and Analytics Specifications

To ensure consistent performance and prevent errors when scaling MongoDB Atlas clusters, certain attributes across readOnlySpecs and analyticsSpecs must be synchronized with electableSpecs. The synchronization rules are as follows:

**Disk Type (ebsVolumeType)**: The ebsVolumeType in analyticsSpecs and readOnlySpecs must always match the value in electableSpecs. A synchronization function ensures this consistency.

**Disk IOPS (diskIOPS)**: The diskIOPS in analyticsSpecs and readOnlySpecs must always match the value in electableSpecs. A synchronization function ensures this consistency.

**Instance Size**: The instanceSize in readOnlySpecs must always match the value in electableSpecs. A synchronization function ensures this consistency. In the case of analyticsSpecs, the instanceSize can be different from electableSpecs, but it must be explicitly specified.

> **Note**: This project has been specifically designed for clusters with a single region and replica sets. It works by iterating over each region_configs entry within the cluster and synchronizing all specified electable, analytics, and readOnly specifications in the payload. If you have a different requirement, feel free to open an issue or submit a pull request with the necessary changes to accommodate your use case.

### Step 1: Create an API Key on MongoDB Atlas

1. Log in to your **MongoDB Atlas** account.
2. Navigate to **Project Settings** and create an API key with **Programmatic Access**.
   - Ensure the key has sufficient permissions to modify clusters (e.g., `Cluster Manager`).
3. Take note of the **Public Key** and **Private Key**, as you will need them in the next step.

### Step 2: Store the API Key in AWS Secrets Manager

1. Log in to the **AWS Console** and navigate to **AWS Secrets Manager**.
2. Create a new secret with the following key-value pairs:
   - `MONGODB_ATLAS_PUBLIC_KEY`: *Your MongoDB Atlas Public Key*
   - `MONGODB_ATLAS_PRIVATE_KEY`: *Your MongoDB Atlas Private Key*
3. Store the secret and take note of the **ARN** (Amazon Resource Name), which will be needed for the Lambda configuration.

### Step 3: Build the Code for AWS Lambda with Event Rule

Before deploying the project, you need to compile the Go code for the AWS Lambda environment (Linux).

1. Open a terminal in the root directory of your project.
2. Compile the Go binary targeting Linux and package it into a ZIP file using the following commands:

```bash
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
zip function.zip bootstrap
````

Deploy to AWS Lambda: Use the AWS CLI to deploy the function:

```bash
aws lambda update-function-configuration \
--function-name <your-function-name> \
--environment "Variables={SECRET_NAME=<ARN-da-Secret-Manager>}"
```

3. Give the Lambda function permission to access the secret by attaching the following policy to the Lambda execution role:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "secretsmanager:GetSecretValue",
      "Resource": "<ARN-da-Secret-Manager>"
    }
  ]
}
```

4. Create Event Rue with your schedule expression:

```bash
aws events put-rule \
    --name "AutoScalingRule" \
    --schedule-expression "cron(0 12 * * ? *)" \
    --state "ENABLED" \
    --description "Triggers Lambda for auto-scaling MongoDB clusters during peak hours"

aws lambda add-permission \
    --function-name <your-lambda-function-name> \
    --statement-id "AllowEventBridgeInvoke" \
    --action "lambda:InvokeFunction" \
    --principal events.amazonaws.com \
    --source-arn arn:aws:events:<region>:<account-id>:rule/AutoScalingRule

aws events put-targets \
    --rule "AutoScalingRule" \
    --targets "Id"="1","Arn"="arn:aws:lambda:<region>:<account-id>:function:<your-lambda-function-name>", \
              "Input"='{
                  "project": "myProject",
                  "cluster": "myCluster",
                  "electableSpecs": {
                      "instanceSize": "M30",
                      "AutoScale": {
                          "MinInstanceSize": "M10",
                          "MaxInstanceSize": "M30"
                      }
                  },
                  "readOnlySpecs": {
                      "NodeCount": 2
                  },
                  "analyticsSpecs": {
                      "NodeCount": 1,
                      "instanceSize": "M20"
                  }
              }'
```

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.
