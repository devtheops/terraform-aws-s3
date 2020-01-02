// Provider for DR account

provider "aws" {
  alias  = "kb4-disasterrecovery"
  region = "us-west-1"

  assume_role {
    role_arn = "arn:aws:iam::406119160266:role/OrganizationAccountAccessRole"
  }
}

