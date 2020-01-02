// Replication configuration access policy
data "template_file" "replication_access_policy" {
  count    = var.enable_replication ? 1 : 0
  template = file("${path.module}/files/replication_access_policy.json")

  vars = {
    bucket_arn             = "arn:aws:s3:::${var.bucket_name}"
    replication_bucket_arn = "arn:aws:s3:::${var.bucket_name}-rep"
  }
}

resource "aws_iam_policy" "replication_access_policy" {
  count  = var.enable_replication ? 1 : 0
  name   = "policy-s3-rep-${var.bucket_name}"
  policy = data.template_file.replication_access_policy[0].rendered
}

// Replication configuration assume policy

data "template_file" "replication_assume_policy" {
  count    = var.enable_replication ? 1 : 0
  template = file("${path.module}/files/replication_assume_policy.json")
}

// Replication configuration role

resource "aws_iam_role" "replication_role" {
  count              = var.enable_replication ? 1 : 0
  name               = "s3-rep-${var.bucket_name}"
  assume_role_policy = data.template_file.replication_assume_policy[0].rendered
}

// Attach replication configuration access policy to replication role

resource "aws_iam_role_policy_attachment" "attach" {
  count      = var.enable_replication ? 1 : 0
  role       = aws_iam_role.replication_role[0].name
  policy_arn = aws_iam_policy.replication_access_policy[0].arn
}

// Bucket policy for replicated bucket
data "aws_caller_identity" "current" {
}

data "template_file" "replication_bucket_policy" {
  count    = var.enable_replication ? 1 : 0
  template = file("${path.module}/files/replication_bucket_policy.json")

  vars = {
    replication_bucket_arn = "arn:aws:s3:::${var.bucket_name}-rep"
    account_id             = data.aws_caller_identity.current.account_id
  }
}

// Create replicated S3 bucket and attach policy
resource "aws_s3_bucket" "replication_bucket" {
  count    = var.enable_replication ? 1 : 0
  provider = aws.kb4-disasterrecovery
  bucket   = "${var.bucket_name}-rep"
  acl      = "private"
  region   = "us-west-1"
  policy   = data.template_file.replication_bucket_policy[0].rendered

  versioning {
    enabled = true
  }

  cors_rule {
    allowed_headers = var.cors_allowed_headers
    allowed_methods = var.cors_allowed_methods
    allowed_origins = var.cors_allowed_origins
    expose_headers  = var.cors_expose_headers
    max_age_seconds = var.cors_max_age_seconds
  }

  tags = {
    project    = var.app
    role       = "s3"
    env        = var.environment
    terraform  = "true"
    app        = var.app
  }

  # Optional default encryption
  dynamic "server_side_encryption_configuration" {
    for_each = local.enable_encryption
    content {
      rule {
        apply_server_side_encryption_by_default {
          kms_master_key_id = var.encryption_kms_key_arn_replication
          sse_algorithm     = local.sse_algorithm
        }
      }
    }
  }
}
