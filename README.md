# Basic Setup of a Go application in AWS APP Runner

## Create a GitHub connection

Go to AppRunner [Github Connections] and create a new connection. You will need the `ARN` of this resource for setting up
the CloudFormation Stack

> Make sure your connection has access only to specific repositories

## Create the CloudFormation Stack

This Stack creates a DynamoDB Table, an AppRunner instance role that can access the Table, and an AppRunner Service
that deploys the go backend. 

> We have disabled the automatic deployments, so we can control the deployments to the AppRunner Service from CircleCI

```bash
> # Make sure you have configured aws cli through `aws configure` with proper credentials
> export CONNECTION_ARN="arn:aws:apprunner:eu-west-1:xxxx:connection/xxx/xxx" # <- Copy this from connection created above
> export REPOSITORY_URL="https://github.com/shyamz-22/apprunner-go-runtime-app"
> export APPRUNNER_SERVICE_NAME="apprunner-go-demo"
> aws deploy aws cloudformation deploy
             --stack-name apprunner-demo-stack \
             --template-body app-runner-cfn.yml \
             --capabilities CAPABILITY_IAM \
             --parameters ParameterKey=ServiceARN,ParameterValue=$CONNECTION_ARN \
                          ParameterKey=GitHubRepository,ParameterValue=$REPOSITORY_URL \
                          ParameterKey=AppRunnerServiceName,ParameterValue=$APPRUNNER_SERVICE_NAME
```

## Setup CircleCI

Read [Continuous Integration with CircleCI] for setting up the continuous integration for this project

### Create CircleCI Deployer Role
This role has permission only to start deployment and read only permissions for listing services and operations.

```bash
> # Make sure you have configure aws cli through `aws configure` with proper credentials
> export APPRUNNER_SERVICE_ARN="arn:aws:apprunner:eu-west-1:XXX:service/xxx/xxx" # <- Copy this from service created above
> aws deploy aws cloudformation deploy
             --stack-name apprunner-demo-ci-stack \
             --template-body app-runner-deployer-cfn.yml \
             --capabilities CAPABILITY_IAM \
             --parameters ParameterKey=AppRunnerServiceArn,ParameterValue=$APPRUNNER_SERVICE_ARN
> aws iam create-access-key --user-name CIDeployer
{
    "AccessKey": {
        "UserName": "username",
        "AccessKeyId": "YOUR_ACCESS_KEY_ID",
        "Status": "Active",
        "SecretAccessKey": "YOUR_SECRET_ACCESS_KEY",
        "CreateDate": "YYYY-MM-DDTHH:MM:SSZ"
    }
}
```
Now set up following environment variables for the project in CircleCI

- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS
- AWS_APPRUNNER_SERVICE_ARN
- AWS_REGION

🥳💃💃🥳 Awesome you have made it till here. You are all set now. With every push to this git repository
CI will run the test and when the test is successful the app is deployed to AWS App Runner.

Next Steps:
- Custom Domains
- Observability

[Github Connections]: https://eu-west-1.console.aws.amazon.com/apprunner/home?region=eu-west-1#/connections
[Continuous Integration with CircleCI]: https://circleci.com/blog/setting-up-continuous-integration-with-github/