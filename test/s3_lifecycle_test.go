package test

import (
	"testing"

	utils "gitlab.internal.knowbe4.com/sre/tf-test-utils"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestS3LifecycleExpiration(t *testing.T) {
	//t.Parallel()

	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	lifecycleEnabled := true
	expectedLifecycleExpiration := int64(random.Random(1, 1095))

	tfArgs := NewS3TerraformArgs(t)
	tfArgs.EnableLifecycle = &lifecycleEnabled
	tfArgs.LifecycleExpiration = expectedLifecycleExpiration

	awsSession := utils.GetAwsSession(t, awsRegion)

	terraformOptions := &terraform.Options{
		TerraformDir: GetS3TerraformDir(),
		EnvVars:      utils.GetDefaultEnvVars(awsRegion),
		Vars:         utils.Map(tfArgs),
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	AssertBucketCreated(t, terraformOptions, awsSession, tfArgs.BucketName, awsRegion)

	svc := s3.New(awsSession)

	lifecycleConfigInput := &s3.GetBucketLifecycleConfigurationInput{
		Bucket: &tfArgs.BucketName,
	}

	result, err := svc.GetBucketLifecycleConfiguration(lifecycleConfigInput)
	assert.NoError(t, err)

	for _, rule := range result.Rules {
		assert.Equal(t, expectedLifecycleExpiration, *rule.Expiration.Days)
	}
}

func TestS3LifecycleNonCurrentExpiration(t *testing.T) {
	//t.Parallel()

	awsRegion := aws.GetRandomStableRegion(t, nil, nil)
	awsSession := utils.GetAwsSession(t, awsRegion)

	lifecycleEnabled := true
	expectedNonCurrentExpiration := int64(random.Random(3, 30))

	tfArgs := NewS3TerraformArgs(t)
	tfArgs.EnableLifecycle = &lifecycleEnabled
	tfArgs.LifecycleNonCurrentExpiration = expectedNonCurrentExpiration

	terraformOptions := &terraform.Options{
		TerraformDir: GetS3TerraformDir(),
		EnvVars:      utils.GetDefaultEnvVars(awsRegion),
		Vars:         utils.Map(tfArgs),
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	AssertBucketCreated(t, terraformOptions, awsSession, tfArgs.BucketName, awsRegion)

	svc := s3.New(awsSession)

	req := &s3.GetBucketLifecycleConfigurationInput{
		Bucket: &tfArgs.BucketName,
	}

	result, err := svc.GetBucketLifecycleConfiguration(req)
	assert.NoError(t, err)

	for _, rule := range result.Rules {
		assert.Equal(t, expectedNonCurrentExpiration, *rule.NoncurrentVersionExpiration.NoncurrentDays)
	}
}

// TODO: Test what happens if you don't pass "enable_lifecycle"
