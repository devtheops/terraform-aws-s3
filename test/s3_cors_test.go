package test

import (
	"testing"

	utils "gitlab.internal.knowbe4.com/sre/tf-test-utils"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAdditionalAllowedMethods(t *testing.T) {
	// t.Parallel()

	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	awsSession := utils.GetAwsSession(t, awsRegion)

	expectedCorsAllowedMethods := []string{
		"HEAD",
		"GET",
		"POST",
		"PUT",
		"DELETE",
	}

	s3TerraformArgs := NewS3TerraformArgs(t)
	s3TerraformArgs.CorsAllowedMethods = expectedCorsAllowedMethods

	terraformOptions := &terraform.Options{
		TerraformDir: GetS3TerraformDir(),
		EnvVars:      utils.GetDefaultEnvVars(awsRegion),
		Vars:         utils.Map(s3TerraformArgs),
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	bucketName := terraform.Output(t, terraformOptions, "bucket_name")

	AssertBucketCreated(t, terraformOptions, awsSession, bucketName, awsRegion)

	corsInput := s3.GetBucketCorsInput{
		Bucket: &bucketName,
	}

	svc := s3.New(awsSession)
	result, err := svc.GetBucketCors(&corsInput)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	for _, rule := range result.CORSRules {
		utils.AssertSlicesEqual(t, expectedCorsAllowedMethods, rule.AllowedMethods)
	}
}

func TestCorsAllowedHeaders(t *testing.T) {
	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	awsSession := utils.GetAwsSession(t, awsRegion)

	expectedCorsAllowedHeaders := []string{"ETag"}

	tfArgs := NewS3TerraformArgs(t)
	tfArgs.CorsAllowedHeaders = expectedCorsAllowedHeaders

	terraformOptions := &terraform.Options{
		TerraformDir: GetS3TerraformDir(),
		EnvVars:      utils.GetDefaultEnvVars(awsRegion),
		Vars:         utils.Map(tfArgs),
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	AssertBucketCreated(t, terraformOptions, awsSession, tfArgs.BucketName, awsRegion)

	corsInput := s3.GetBucketCorsInput{
		Bucket: &tfArgs.BucketName,
	}

	svc := s3.New(awsSession)

	result, err := svc.GetBucketCors(&corsInput)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	for _, rule := range result.CORSRules {
		utils.AssertSlicesEqual(t, expectedCorsAllowedHeaders, rule.AllowedHeaders)
	}
}

func TestCorsExposeHeaders(t *testing.T) {
	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	awsSession := utils.GetAwsSession(t, awsRegion)

	expectedCorsExposeHeaders := []string{"ETag"}

	tfArgs := NewS3TerraformArgs(t)
	tfArgs.CorsExposeHeaders = expectedCorsExposeHeaders

	terraformOptions := &terraform.Options{
		TerraformDir: GetS3TerraformDir(),
		EnvVars:      utils.GetDefaultEnvVars(awsRegion),
		Vars:         utils.Map(tfArgs),
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	AssertBucketCreated(t, terraformOptions, awsSession, tfArgs.BucketName, awsRegion)

	corsInput := s3.GetBucketCorsInput{
		Bucket: &tfArgs.BucketName,
	}

	svc := s3.New(awsSession)

	result, err := svc.GetBucketCors(&corsInput)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	for _, rule := range result.CORSRules {
		utils.AssertSlicesEqual(t, expectedCorsExposeHeaders, rule.ExposeHeaders)
	}
}
