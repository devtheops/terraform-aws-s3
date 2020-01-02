output "bucket_arn" {
  value       = aws_s3_bucket.main.arn
  description = "The ARN of the bucket. Will be of format arn:aws:s3:::bucketname."
}

output "bucket_name" {
  value       = aws_s3_bucket.main.id
  description = "The name of the bucket."
}

output "bucket_region" {
  value       = aws_s3_bucket.main.region
  description = "The region of the bucket."
}

output "bucket_domain_name" {
  value       = aws_s3_bucket.main.bucket_domain_name
  description = "The domain name of the bucket."
}
