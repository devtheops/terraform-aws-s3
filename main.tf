data "template_file" "description" {
  template = "This module creates an s3 bucket."
}

resource "aws_s3_bucket" "main" {
  bucket = var.bucket_name
  acl    = var.acl
  policy = var.policy

  versioning {
    enabled = var.enable_versioning
  }

  force_destroy = true

  cors_rule {
    allowed_headers = var.cors_allowed_headers
    allowed_methods = var.cors_allowed_methods
    allowed_origins = var.cors_allowed_origins
    expose_headers  = var.cors_expose_headers
    max_age_seconds = var.cors_max_age_seconds
  }

  tags = merge(local.default_tags, var.tags)

  # Optional default encryption
  dynamic "server_side_encryption_configuration" {
    for_each = local.enable_encryption
    content {
      rule {
        apply_server_side_encryption_by_default {
          kms_master_key_id = var.encryption_kms_key_arn
          sse_algorithm     = local.sse_algorithm
        }
      }
    }
  }

  # Optional lifecycle
  dynamic "lifecycle_rule" {
    for_each = local.enable_lifecycle

    content {
      id      = var.bucket_name
      enabled = true

      expiration {
        days = var.lifecycle_expiration
      }

      noncurrent_version_expiration {
        days = var.lifecycle_noncurrent_expiration
      }
    }
  }

  dynamic "website" {
    for_each = local.enable_website_hosting
    content {
      index_document = var.website_hosting_index_document
      error_document = var.website_hosting_error_document
    }
  }
}

resource "aws_s3_bucket_notification" "bucket_notification" {
  count = length(concat(var.topic_notifications, var.queue_notifications, var.lambda_notifications)) > 0 ? 1 : 0

  bucket = aws_s3_bucket.main.id

  dynamic "topic" {
    for_each = var.topic_notifications
    content {
      topic_arn     = topic.value.topic_arn
      events        = topic.value.events
      filter_prefix = lookup(topic.value, "filter_prefix", null)
      filter_suffix = lookup(topic.value, "filter_suffix", null)
    }
  }

  dynamic "queue" {
    for_each = var.queue_notifications
    content {
      queue_arn     = queue.value.queue_arn
      events        = queue.value.events
      filter_prefix = lookup(queue.value, "filter_prefix", null)
      filter_suffix = lookup(queue.value, "filter_suffix", null)
    }
  }

  dynamic "lambda_function" {
    for_each = var.lambda_notifications
    content {
      lambda_function_arn = lambda_function.value.lambda_function_arn
      events              = lambda_function.value.events
      filter_prefix       = lookup(lambda_function.value, "filter_prefix", null)
      filter_suffix       = lookup(lambda_function.value, "filter_suffix", null)
    }
  }

  depends_on = [aws_lambda_permission.allow_bucket]
}

resource "aws_lambda_permission" "allow_bucket" {
  count         = length(var.lambda_notifications)
  action        = "lambda:InvokeFunction"
  function_name = lookup(var.lambda_notifications[count.index], "lambda_function_arn")
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.main.arn
}

resource "aws_s3_bucket_policy" "main" {
  count  = var.bucket_policy == null ? 0 : 1
  bucket = var.bucket_name
  policy = var.bucket_policy
}
