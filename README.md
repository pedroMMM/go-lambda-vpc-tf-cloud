# go-lambda-vpc-tf-cloud
Go Lambda that counts VPCs, keeps the state in S3. All instrumented with Terraform Cloud.

## Requirements

1. **Go** >= 1.14
1. **Terraform** >= 0.12.23
1. A AWS Account
1. A Public Route 53 Hosted Zone

## How to use

### Terraform remote
In the `./terraform/main.tf` modify the `terraform backend` block: 
```hcl
terraform {
  backend "remote" {
    organization = "pedrommm"

    workspaces {
      name = "go-lambda-vpc-tf-cloud"
    }
  }
}
```
Either by removing it to use the local state or by using your Terrafor Cloud account and modifing the `organization`.

### Configuration
Modify `./terraform/config.auto.tfvars` with your personal configuration.

### To Terraform init:
```bash
./tf-init
```
### To Terraform plan:
```bash
./tf-plan
```
### To Terraform apply:
```bash
./tf-apply
```
