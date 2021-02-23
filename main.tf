terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}


provider "aws" {
  region  = "ap-southeast-2"
}

resource "aws_ecr_repository" "zipcoingester" {
  name = "zipco-ecr-repo/ingester"
  image_tag_mutability = "MUTABLE"
}
output "ecr_repository_ingester_endpoint" {
    value = aws_ecr_repository.zipcoingester.repository_url
}

resource "aws_ecs_cluster" "ecs_cluster" {
    name  = "zipco-cluster"
}

resource "aws_ecs_task_definition" "elasticserver" {
  family                   = "elasticserver" # Naming our first task
  container_definitions    = <<DEFINITION
  [
    {
      "name": "elasticserver",
      "image": "docker.elastic.co/elasticsearch/elasticsearch:7.11.1",
      "essential": true,
      "portMappings": [
        {
          "containerPort": 9200,
          "hostPort": 9200
        }
      ],
      "memory": 512,
      "environment": [
        {
          "name": "discovery.type",
          "value": "single-node"
        }
      ]
    }
  ]
  DEFINITION
  memory                   = 512
  cpu                      = 256
}

resource "aws_ecs_service" "elasticservice" {
  name            = "elastic-service"
  cluster         = aws_ecs_cluster.ecs_cluster.id
  task_definition = aws_ecs_task_definition.elasticserver.arn
  desired_count   = 1
}


resource "aws_autoscaling_group" "ecs_cluster_instances" {
    # IN progress
}