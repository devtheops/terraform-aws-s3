package test

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	utils "gitlab.internal.knowbe4.com/sre/tf-test-utils"
)

// TFArgsModel
type S3TerraformArgs struct {
	Acl                            string            `structs:"acl,omitempty"`
	App                            string            `structs:"app,omitempty"`
	BucketName                     string            `structs:"bucket_name,omitempty"`
	CorsAllowedHeaders             []string          `structs:"cors_allowed_headers,omitempty"`
	CorsAllowedMethods             []string          `structs:"cors_allowed_methods,omitempty"`
	CorsAllowedOrigins             []string          `structs:"cors_allowed_origins,omitempty"`
	CorsExposeHeaders              []string          `structs:"cors_expose_headers,omitempty"`
	CorsMaxAgeSeconds              int64             `structs:"cors_max_age_seconds,omitempty"`
	DisableDefaultEncryption       *bool             `structs:"disable_default_encryption,omitempty"`
	EnableLifecycle                *bool             `structs:"enable_lifecycle,omitempty"`
	EnableReplication              *bool             `structs:"enable_replication,omitempty"`
	EnableVersioning               *bool             `structs:"enable_versioning,omitempty"`
	EnableWebsiteHosting           *bool             `structs:"enable_website_hosting,omitempty"`
	EncryptionKmsKeyArn            string            `structs:"encryption_kms_key_arn,omitempty"`
	EncryptionKmsKeyArnReplication string            `structs:"encryption_kms_key_arn_replication,omitempty"`
	Environment                    string            `structs:"environment,omitempty"`
	LifecycleExpiration            int64             `structs:"lifecycle_expiration,omitempty"`
	LifecycleNonCurrentExpiration  int64             `structs:"lifecycle_noncurrent_expiration,omitempty"`
	Policy                         string            `structs:"policy,omitempty"`
	Service                        string            `structs:"service,omitempty"`
	Tags                           map[string]string `structs:"tags,omitempty"`
	WebsiteHostingErrorDocument    string            `structs:"website_hosting_error_document,omitempty"`
	WebsiteHostingIndexDocument    string            `structs:"website_hosting_index_document,omitempty"`
}

type S3TerraformArg func(arg S3TerraformArgs)

func AppOption(app string) S3TerraformArg {
	return func(args S3TerraformArgs) {
		args.App = app
	}
}

func ServiceOption(service string) S3TerraformArg {
	return func(args S3TerraformArgs) {
		args.Service = service
	}
}

func EnvironmentOption(environment string) S3TerraformArg {
	return func(args S3TerraformArgs) {
		args.Environment = environment
	}
}

func BucketNameOption(bucketName string) S3TerraformArg {
	return func(args S3TerraformArgs) {
		args.BucketName = bucketName
	}
}

func NewS3TerraformArgs(t *testing.T, options ...S3TerraformArg) S3TerraformArgs {
	s3Args := S3TerraformArgs{}
	for _, option := range options {
		option(s3Args)
	}
	if s3Args.App == "" {
		s3Args.App = strings.ToLower(random.UniqueId())
	}
	if s3Args.Service == "" {
		s3Args.Service = strings.ToLower(random.UniqueId())
	}
	if s3Args.Environment == "" {
		s3Args.Environment = "GoTest"
	}
	if s3Args.BucketName == "" {
		s3Args.BucketName = utils.GetTestIdentifier(t)
	}
	return s3Args
}
