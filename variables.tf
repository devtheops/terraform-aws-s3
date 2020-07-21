variable "acl" {
  description = "The [canned ACL](https://docs.aws.amazon.com/AmazonS3/latest/dev/acl-overview.html#canned-acl) to apply."
  type        = string
  default     = "private"
}

variable "app" {
  description = "The app name."
  type        = string
}

variable "service" {
  description = "The service of the application."
  type        = string
}

variable "env" {
  description = "The Tag Environment."
  type        = string
}

variable "bucket_name" {
  description = "The name of the bucket"
  type        = string
  default     = ""
}

variable "enable_lifecycle" {
  type    = bool
  default = false
}

variable "enable_versioning" {
  description = "Turns versiong of the bucket on or off."
  type        = bool
  default     = true
}

variable "lifecycle_expiration" {
  description = "How many days to expire the curernt version of an object"
  type        = number
  default     = 1095
}

variable "lifecycle_noncurrent_expiration" {
  description = "How many days to delete the previous version of an object"
  type        = number
  default     = 1
}

variable "policy" {
  description = "Adds a policy."
  type        = string
  default     = ""
}

variable "cors_allowed_headers" {
  type    = list(string)
  default = ["*"]
}

variable "cors_allowed_methods" {
  type    = list(string)
  default = ["GET"]
}

variable "cors_allowed_origins" {
  type    = list(string)
  default = ["*"]
}

variable "cors_expose_headers" {
  type    = list(string)
  default = ["ETag"]
}

variable "cors_max_age_seconds" {
  type    = number
  default = 3000
}

variable "disable_default_encryption" {
  type    = bool
  default = false
}

variable "encryption_kms_key_arn" {
  type    = string
  default = ""
}

variable "encryption_kms_key_arn_replication" {
  type    = string
  default = ""
}

variable "tags" {
  description = "Additional tags to be added to the bucket"
  type        = map(string)
  default     = {}
}


variable "topic_notifications" {
  description = "List of sns topics to notify. https://www.terraform.io/docs/providers/aws/r/s3_bucket_notification.html#topic"

  type = list(object({
    topic_arn     = string
    events        = list(string) # https://docs.aws.amazon.com/AmazonS3/latest/dev/NotificationHowTo.html#notification-how-to-event-types-and-destinations
    filter_prefix = string
    filter_suffix = string
  }))

  default = []
}

variable "queue_notifications" {
  description = "List of sqs queue to notify. https://www.terraform.io/docs/providers/aws/r/s3_bucket_notification.html#queue"

  type = list(object({
    queue_arn     = string
    events        = list(string) # https://docs.aws.amazon.com/AmazonS3/latest/dev/NotificationHowTo.html#notification-how-to-event-types-and-destinations
    filter_prefix = string
    filter_suffix = string
  }))

  default = []
}

variable "lambda_notifications" {
  description = "List of lambda functions to notify. https://www.terraform.io/docs/providers/aws/r/s3_bucket_notification.html#lambda_function"

  type = list(object({
    lambda_function_arn = string
    events              = list(string) # https://docs.aws.amazon.com/AmazonS3/latest/dev/NotificationHowTo.html#notification-how-to-event-types-and-destinations
    filter_prefix       = string
    filter_suffix       = string
  }))

  default = []
}

variable "enable_website_hosting" {
  description = "Enables static website hosting for bucket"
  type        = bool
  default     = false
}

variable "website_hosting_index_document" {
  description = "Index page for static site hosting"
  default     = "index.html"
}

variable "website_hosting_error_document" {
  description = "Error page for static site hosting"
  default     = "index.html"
}

variable "website_hosting_redirect_all_requests_to" {
  description = "A hostname to redirect all website requests for this bucket to."
  default = null
}


variable "bucket_policy" {
  type        = string
  default     = null
  description = "the bucket policy you want to add to this bucket"
}
