### Description

No description provided. Please add a data description to the module with a
description.

```tf
data "template_file" "description" {
    template = "This is my super awesome module and it does this stuff..."
}
```

### Variables

| **Name** | **Type** | **Description** | **Default** |
| -------- | ---- | ----------- | ------- |
| acl | `unknown` | The [canned ACL](https://docs.aws.amazon.com/AmazonS3/latest/dev/acl-overview.html#canned-acl) to apply. | `"private"` |
| app | `unknown` | The app name. | *None* |
| bucket_name | `unknown` | The name of the bucket | *None* |
| bucket_region | `bool` | If specified, the AWS region this bucket should reside in. Otherwise, the region used by the callee. | `false` |
| department | `unknown` | The Tag Department. | `"engineering"` |
| enable_replication | `bool` | Turns bucket replication on or off. | `false` |
| enable_lifecycle | `bool` | No description provided. Please add a description to this variable. | `false` |
| enable_versioning | `bool` | Turns versiong of the bucket on or off. | `true` |
| environment | `unknown` | The Tag Environment. | `"development"` |
| lifecycle_expiration | `integer` | How many days to expire the curernt version of an object | `1095` |
| lifecycle_noncurrent_expiration | `integer` | How many days to delete the previous version of an object | `1` |
| policy | `unknown` | Adds a policy. | `""` |
| product | `unknown` | The Tag product. | *None* |
| cors_allowed_headers | `list` | No description provided. Please add a description to this variable. | `["*"]` |
| cors_allowed_methods | `list` | No description provided. Please add a description to this variable. | `["GET"]` |
| cors_allowed_origins | `list` | No description provided. Please add a description to this variable. | `["*"]` |
| cors_expose_headers | `list` | No description provided. Please add a description to this variable. | `["ETag"]` |
| cors_max_age_seconds | `integer` | No description provided. Please add a description to this variable. | `3000` |

### Outputs

| **Name** | **Description** |
| -------- | --------------- |
| bucket_arn | No description provided. Please add a description to this variable. |
| bucket_name | No description provided. Please add a description to this variable. |
| bucket_region | No description provided. Please add a description to this variable. |
| bucket_domain_name | No description provided. Please add a description to this variable. |

### Resources used

* aws_iam_policy
* aws_iam_role
* aws_iam_role_policy_attachment
* aws_s3_bucket