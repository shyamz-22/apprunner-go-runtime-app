#!/bin/bash
set -e

service_arn="$1"
deployment_command="aws apprunner start-deployment --service-arn $service_arn --query OperationId --output text"
list_operations_command="aws apprunner list-operations --service-arn $service_arn"
timeout_seconds=900  # Adjust the timeout duration as needed (e.g., 900 seconds = 15 minutes)
aws_profile="$2"  # Optional AWS profile parameter

if [[ -z "$service_arn" ]]; then
    echo "Usage: $0 <service_arn> [aws_profile]"
    exit 1
fi

if [[ -n "$aws_profile" ]]; then
    deployment_command="aws apprunner start-deployment --service-arn $service_arn --query OperationId --output text --profile $aws_profile"
    list_operations_command="aws apprunner list-operations --service-arn $service_arn --profile $aws_profile"
fi

operation_id=$($deployment_command)

if [[ -z "$operation_id" ]]; then
    echo "Failed to start deployment."
    exit 1
fi

echo "Deployment started with operation ID: $operation_id"

function display_progress {
    local width=30
    local elapsed=$1
    local percent=$((elapsed * 100 / timeout_seconds))
    local num_chars=$((percent * width / 100))
    local bar=$(printf "[%-${width}s]" "$(printf '#%.0s' $(seq 1 "$num_chars"))")
    printf "\r%s %d%%" "${bar}" "${percent}"
}

start_time=$(date +%s)

while true; do
    current_time=$(date +%s)
    elapsed_time=$((current_time - start_time))

    if [[ $elapsed_time -ge $timeout_seconds ]]; then
        echo -e "\nTimeout reached. Operation did not complete within $timeout_seconds seconds."
        exit 2
    fi

    operation_status=$($list_operations_command --query "OperationSummaryList[?Id == '$operation_id'].Status" --output text)
    if [[ $operation_status == "SUCCEEDED" ]]; then
        display_progress "$timeout_seconds"
        echo -e "\nDeployment operation $operation_id has succeeded"
        exit 0
    elif [[ $operation_status == "FAILED" ]]; then
        display_progress "$timeout_seconds"
        echo -e "\nDeployment operation $operation_id has failed"
        exit 2
    fi

    display_progress "$elapsed_time"
    sleep 5  # You can adjust the polling interval as needed
done
