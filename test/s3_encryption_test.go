package test

import (
	"testing"

	utils "gitlab.internal.knowbe4.com/sre/tf-test-utils"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestS3BucketEncryptionDisabled(t *testing.T) {
	//t.Parallel()

	awsRegion := aws.GetRandomStableRegion(t, nil, nil)
	awsSession := utils.GetAwsSession(t, awsRegion)

	disableDefaultEncryption := true

	tfArgs := NewS3TerraformArgs(t)
	tfArgs.DisableDefaultEncryption = &disableDefaultEncryption

	terraformOptions := &terraform.Options{
		TerraformDir: GetS3TerraformDir(),
		EnvVars:      utils.GetDefaultEnvVars(awsRegion),
		Vars:         utils.Map(tfArgs),
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	AssertBucketCreated(t, terraformOptions, awsSession, tfArgs.BucketName, awsRegion)

	svc := s3.New(awsSession)

	req := &s3.GetBucketEncryptionInput{
		Bucket: &tfArgs.BucketName,
	}
	_, err := svc.GetBucketEncryption(req)
	assert.Error(t, err)
	if awsError, ok := err.(awserr.RequestFailure); ok {
		assert.Equal(t, "ServerSideEncryptionConfigurationNotFoundError", awsError.Code())
	} else {
		assert.Fail(t, "Unexpected error [%s] received", err.Error())
	}
}
