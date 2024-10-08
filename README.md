# Autoscaling MongoDB Atlas

This repository contains code to automatically scale MongoDB Atlas clusters, following specific business rules for different specifications (`electableSpecs`, `readOnlySpecs`, `analyticsSpecs`).

## Description

The main function of this project is to handle the auto-scaling of MongoDB clusters based on the parameters received. It uses AWS Lambda to execute these operations automatically.

### Project Structure

- **cmd/**: Contains the main application code.
- **internal/**: Contains internal logic and helper functions, including business rules for synchronizing parameters.
- **pkg/**: Contains utilities, such as the interface with `secretManager`.
- **testdata/**: Contains test payloads used to simulate auto-scaling events.

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

|**Parameter**|**Description**|
|---------------------|--------------------------------------------|
|**project**|The name of the MongoDB Atlas project.|
|**cluster**|The name of the cluster to be scaled.|
|**electableSpecs**|Specifications for the electable instance:|
|-> `diskiops`|Disk IOPS.|
|-> `ebsVolumeType`|The EBS volume type (e.g., "PROVISIONED" or "STANDARD").|
|-> `instanceSize`|Instance size (e.g., "M30").|
|-> `NodeCount`|Number of electable nodes.|
|-> `AutoScale`|Auto-scaling configuration:|
|--> `MinInstanceSize`|Minimum instance size for auto-scaling.|
|--> `MaxInstanceSize`|Maximum instance size for auto-scaling.|
|**readOnlySpecs**|Specifications for read-only instances:|
|-> `NodeCount`|Number of read-only nodes.|
|**analyticsSpecs**|Specifications for analytics instances:|
|-> `NodeCount`|Number of analytics nodes.|
|-> `instanceSize`|Instance size for analytics nodes.|
|-> `AutoScale`|Auto-scaling configuration:|
|--> `MinInstanceSize`|Minimum instance size for auto-scaling.|
|--> `MaxInstanceSize`|Maximum instance size for auto-scaling.|

### Explanation of the Tables

**Parameter**: Lists the parameter names.

**Description**: Provides details about each parameter and its function.

## Business Rules

There are specific rules governing how parameters are synchronized across different specifications:

Disk Type (ebsVolumeType): The disk type (ebsVolumeType) of analyticsSpecs must always match the value in electableSpecs. To ensure this, a synchronization function automatically adjusts this value.

Instance Size Synchronization: Some specs, such as analyticsSpecs, can have different instanceSize values. Synchronization between specifications occurs only when needed.
Auto-Scaling: If specified, the minimum and maximum instance sizes for auto-scaling must align with the provided specs.

This project is designed to run on AWS Lambda and perform auto-scaling operations on MongoDB Atlas clusters. Follow the steps below to configure and deploy the project.

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

### Step 3: Build the Code for AWS Lambda

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


## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.
