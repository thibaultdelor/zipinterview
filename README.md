# ZipCo Interview Project

## Instructions

> Tech Lead - Search & Big Data - Technical Challenge
> 
> We want to implement a golang ingestion server on ECS that will receive http requests representing
> click events from our users. The event is just composed of the user_id. Those events will be ingested
> by our Elasticsearch database (also deployed on ECS).
> Every 3m we have a python lambda function to query the ES index to count what happened during
> that elapsed time:
> - the total number of events and
> - the number of users
> 
> Create a python script that is sending fake events to your API.
> 
> BONUS: create an AWS cloudwatch log filter that will parse the logs and extract the count metrics in
> order to create a monitoring dashboard.
> 
> You will use terraform to build the ECS/cloudwatch infrastructure on AWS and the serverless function
> for the lambda function.


## State

 - Locally we have a docker compose that runs Elastic search an the go API
 - The `ingester` folder contains the go API in go that takes event and store them in Elastic
 - The `pythonlambda` contains an AWS lambda handler that returns the number of events in the last 3 minutes and the number of users
 - An incomplete terraform definition that create a non working ECS cluster
  
### Why stopping here?

I have already invested a lot of my time and can't commit more.

I feel like a lot of my struggles are things that wouldn't be an issue at work :
 - I dont want to use Fargate since there's no free tier
 - I am trying to setup a cluster with just t2.micro to stay in the free tier
 - I would use AWS ElasticService instead of trying to run my own
 - ECS takes time to setup and it's a one off thing (setting up the VPC, 
   the security config, the AMI user_data config, the IAM roles, ...). In a work environment,
   I would expect to have at least a full day to set it up using terraform.

## Local testing

```
# Start the Go API backed by Elastic
docker compose up
# Send an event
curl -X PUT -H "Content-Type: application/json" -d '{"user_id":1}' "http://localhost:8080/clickevents"
# Collect Stats
cd pythonlambda &&  pip install -r requirements.txt && python app.py
```

## Deploying 

First lets push our images
```
$ terraform apply -target="aws_ecr_repository.zipco_ecr_repo"
$ aws ecr get-login-password | docker login --username AWS --password-stdin <ACCOUNT_ID>.dkr.ecr.ap-southeast-2.amazonaws.com
```

Then deploy our whole infra
```
$ terraform apply
```
