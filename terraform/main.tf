# In this scenario I am using the Terraform Cloud just a remote 
# state and lock store.
# 
# More complex cenarios can use the GitHub integrations to trigger
# a plan creation and still require human operators to aprrove them.
terraform {
  backend "remote" {
    organization = "pedrommm"

    workspaces {
      name = "go-lambda-vpc-tf-cloud"
    }
  }
}

provider "aws" {
  region = var.region

  # For simplicity I am using a local User with static Access Key.
  # 
  # In a more complex scenario the developer should request temporary
  # STS creadentials via Vault (or a similar service).
  # 
  # Terraform Cloud can also hold sensative variables that are write 
  # only for any user.
  profile = var.aws_local_profile
}

variable "region" {
  type = string
}

variable "aws_local_profile" {
  type = string
}

variable "vpc_count" {
  type = number
}

variable "owner" {
  type = string
}

data "aws_caller_identity" "current" {}

locals {
  app        = "go-lambda-vpc-tf-cloud"
  account_id = data.aws_caller_identity.current.account_id
}
