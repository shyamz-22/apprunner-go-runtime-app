# Basic Setup of a Go application in AWS APP Runner

 - [Create a GitHub connection](#create-a-github-connection)
 - [Deploy App Runner Service](#deploy-app-runner-service)
 - [Setup CircleCI](#setup-circleci)
   + [Create CircleCI Deployer User](#create-circleci-deployer-user)
 - [Testing the deployed App](#testing-the-deployed-app)
 - [Cleanup](#cleanup-)
 - [Next Steps:](#next-steps)

## Create a GitHub connection

Go to AppRunner [Connected Accounts] and create a new connection. You will need the `ARN` of this resource for setting up
the CloudFormation Stack

> [!IMPORTANT]  
> Make sure your connection has access only to specific repositories

## Deploy App Runner Service

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
  --parameter-overrides GitHubConnectionArn=$CONNECTION_ARN GitHubRepository=$REPOSITORY_URL AppRunnerServiceName=$APPRUNNER_SERVICE_NAME
```

## Setup CircleCI

Read [Continuous Integration with CircleCI] for setting up the continuous integration for this project

> [!Warning]
> We use [eddiewebb/queue@2.2.1] orb for queuing workflows.
> To enable this you have to allow third party uncertified orbs in `Organization Settings > Security`

### Create CircleCI Deployer Role
This [Role](app-runner-ci-cfn.yml) has an attached policy with permission only to start deployment and read only permissions for listing services and operations.

> [!NOTE]
> You can deploy this stack while the App Runner Service is being created. You only need the Service ARN

```bash
> # Make sure you have configure aws cli through `aws configure` with proper credentials
> export APPRUNNER_SERVICE_ARN="arn:aws:apprunner:eu-west-1:XXX:service/xxx/xxx" # <- Copy this from service created above
> export CIRCLE_CI_OPENID_ORGID="xxx-xxx-xxx"  # you can find this in your CircleCI organization settings
>  TMP_TEMPLATE_FILE=$(mktemp) && \
   sed "s/ORG_ID/$CIRCLE_CI_OPENID_ORGID/g" app-runner-ci-cfn.yml > "$TMP_TEMPLATE_FILE" 
> aws cloudformation deploy \
 --template-file $TMP_TEMPLATE_FILE \
 --stack-name apprunner-demo-ci-stack \
 --disable-rollback \
 --capabilities CAPABILITY_NAMED_IAM \
 --parameter-overrides AppRunnerServiceArn=$APPRUNNER_SERVICE_ARN OrgID=$CIRCLE_CI_OPENID_ORGID && rm $TMP_TEMPLATE_FILE 


```
Now create a context named `aws` in organization settings and set up following environment variables for the project in CircleCI

- AWS_APPRUNNER_ROLE_ARN (output of the ci stack above)
- AWS_APPRUNNER_SERVICE_ARN (App runner service ARN)
- AWS_REGION (eu-west-1)

🥳💃💃🥳 Awesome you have made it till here. You are all set now. With every push to this git repository
CI will run the test and when the test is successful the app is deployed to AWS App Runner.

## Testing the deployed App

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

## Cleanup

```bash
> aws cloudformation delete-stack --stack-name apprunner-demo-stack
> aws cloudformation delete-stack --stack-name apprunner-demo-ci-stack
```


## Next Steps:
- Custom Domains
- Observability


[Connected Accounts]: https://eu-west-1.console.aws.amazon.com/apprunner/home?region=eu-west-1#/connections
[Continuous Integration with CircleCI]: https://circleci.com/blog/setting-up-continuous-integration-with-github/
[eddiewebb/queue@2.2.1]:https://circleci.com/developer/orbs/orb/eddiewebb/queue?version=2.2.1#usage-queue_workflow