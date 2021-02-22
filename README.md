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
