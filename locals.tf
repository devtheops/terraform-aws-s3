locals {
  enable_encryption = var.disable_default_encryption ? [] : [1]

  sse_algorithm = length(var.encryption_kms_key_arn) > 0 ? "aws:kms" : "AES256"

  enable_lifecycle = var.enable_lifecycle ? [1] : []

  enable_website_hosting = var.enable_website_hosting ? [1] : []

  default_tags = {
    app       = var.app
    service   = var.service
    env       = var.env
    terraform = "true"
  }
}
