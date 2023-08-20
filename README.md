# Basic Setup of a Go application in AWS APP Runner

## Create a GitHub connection

Go to AppRunner [Github Connections] and create a new connection. You will need the `ARN` of this resource for setting up
the CloudFormation Stack

> [!IMPORTANT]  
> Make sure your connection has access only to specific repositories

## Create the CloudFormation Stack

This [Stack](app-runner-cfn.yml) creates a DynamoDB Table, an AppRunner instance role that can access the Table, and an AppRunner Service
that deploys the go backend. 

> [!NOTE]  
> We have disabled the automatic deployments, so we can control the deployments to the AppRunner Service from CircleCI

```bash
> # Make sure you have configured aws cli through `aws configure` with proper credentials
> export CONNECTION_ARN="arn:aws:apprunner:eu-west-1:xxxx:connection/xxx/xxx" # <- Copy this from connection created above
> export REPOSITORY_URL="https://github.com/shyamz-22/apprunner-go-runtime-app"
> export APPRUNNER_SERVICE_NAME="apprunner-go-demo"
> aws cloudformation deploy \
      --stack-name apprunner-demo-stack \
      --template-file app-runner-cfn.yml \
      --capabilities CAPABILITY_IAM \
      --disable-rollback \
      --parameter-overrides GitHubConnectionArn=$CONNECTION_ARN \
                   GitHubRepository=$REPOSITORY_URL \
                   AppRunnerServiceName=$APPRUNNER_SERVICE_NAME
```

## Setup CircleCI

Read [Continuous Integration with CircleCI] for setting up the continuous integration for this project

### Create CircleCI Deployer User
This [user](app-runner-deployer-cfn.yml) has an attached policy with permission only to start deployment and read only permissions for listing services and operations.

> [!NOTE]
> You can deploy this stack while the App Runner Service is being created

```bash
> # Make sure you have configure aws cli through `aws configure` with proper credentials
> export APPRUNNER_SERVICE_ARN="arn:aws:apprunner:eu-west-1:XXX:service/xxx/xxx" # <- Copy this from service created above
> aws cloudformation deploy \
       --stack-name apprunner-demo-ci-stack \
       --template-file app-runner-deployer-cfn.yml \
       --capabilities CAPABILITY_NAMED_IAM \
       --disable-rollback \
       --parameter-overrides AppRunnerServiceArn=$APPRUNNER_SERVICE_ARN
> aws iam create-access-key --user-name CIDeployer
#{
#    "AccessKey": {
#        "UserName": "username",
#        "AccessKeyId": "YOUR_ACCESS_KEY_ID",
#        "Status": "Active",
#        "SecretAccessKey": "YOUR_SECRET_ACCESS_KEY",
#        "CreateDate": "YYYY-MM-DDTHH:MM:SSZ"
#    }
#}
```
> [!WARNING]  
> Before deleting the stack make sure to deactivate and delete all associated security credentials

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


## Testing the Go app

This app generates short codes for a given url, and when referenced via the short code redirects to the page.

```bash
> export APPRUNNER_SERVICE_ARN="arn:aws:apprunner:eu-west-1:XXX:service/xxx/xxx" # <- Copy this from service created above
> export APP_URL=https://$(aws apprunner describe-service --service-arn $APPRUNNER_SERVICE_ARN --output text --query Service.ServiceUrl)/app/
> curl -X POST \
       -d 'https://dev.to/shyamala_u/the-case-of-disappearing-metrics-in-kubernetes-1kdh' \
      $APP_URL
# {"ShortCode":"a19e9737"}
> echo $APP_URL
> curl -i https://xxx.yyy.awsapprunner.com/app/a19e9737
HTTP/1.1 302 Found
content-length: 0
date: Sun, 20 Aug 2023 21:18:01 GMT
location: https://dev.to/shyamala_u/the-case-of-disappearing-metrics-in-kubernetes-1kdh # <- Redirected to the page
x-envoy-upstream-service-time: 4
server: envoy
```


[Github Connections]: https://eu-west-1.console.aws.amazon.com/apprunner/home?region=eu-west-1#/connections
[Continuous Integration with CircleCI]: https://circleci.com/blog/setting-up-continuous-integration-with-github/