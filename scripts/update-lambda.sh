#!/bin/bash

# The creation of the lambdas is handled in a CDK repo and Lambda's are a pain in managing this way.
# As a result, the first deployment  must only deploy the s3 code which CDK can then use when creating the lambda.
# So first deploy will:
# 1. Deploy the code to S3
# Then I will create the lambda in the CDK repo. Then any subsequent deploys will:
# 1. Deploy the code to s3
# 2. Update the lambda to point to the new code now that it exists

aws lambda get-function --function-name $1 > /dev/null 2>&1

if [ 0 -eq $? ]; then
	aws lambda update-function-code  \
    		--function-name $1  \
    		--s3-bucket kicker-deployments \
		--s3-key $2
else
	echo "Lambda '$1' does not exist - this means that it needs to be created with CDK first"
fi
