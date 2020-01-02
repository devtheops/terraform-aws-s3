package test

import (
	"fmt"
	"strings"
	"testing"

	utils "gitlab.internal.knowbe4.com/sre/tf-test-utils"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestVarTagsOverrideDefaultTags(t *testing.T) {
	// t.Parallel()

	appOverride := strings.ToLower(random.UniqueId())
	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	awsSession := utils.GetAwsSession(t, awsRegion)

	s3TerraformArgs := NewS3TerraformArgs(t)
	s3TerraformArgs.Tags = map[string]string{
		"app": appOverride,
	}

	terraformOptions := &terraform.Options{
		TerraformDir: GetS3TerraformDir(),
		EnvVars:      utils.GetDefaultEnvVars(awsRegion),
		Vars:         utils.Map(s3TerraformArgs),
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	bucketName := terraform.Output(t, terraformOptions, "bucket_name")

	AssertBucketCreated(t, terraformOptions, awsSession, bucketName, awsRegion)

	tags := GetS3Tags(t, awsSession, &bucketName)

	for _, tag := range tags {
		if *tag.Key == "app" {
			assert.Equal(t, appOverride, *tag.Value)
		}
	}
}

func TestAdditionalTags(t *testing.T) {
	// t.Parallel()

	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	awsSession := utils.GetAwsSession(t, awsRegion)

	tags := map[string]string{
		"terratest": "true",
		"foo":       "bar",
	}

	s3TerraformArgs := NewS3TerraformArgs(t)

	s3TerraformArgs.Tags = tags
	terraformOptions := &terraform.Options{
		TerraformDir: GetS3TerraformDir(),
		EnvVars:      utils.GetDefaultEnvVars(awsRegion),
		Vars:         utils.Map(s3TerraformArgs),
	}
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	AssertBucketCreated(t, terraformOptions, awsSession, s3TerraformArgs.BucketName, awsRegion)

	bucketTags := map[string]string{}
	for _, tag := range GetS3Tags(t, awsSession, &s3TerraformArgs.BucketName) {
		bucketTags[*tag.Key] = *tag.Value
	}

	var missingAdditionalTags []string
	for key, val := range tags {
		tagValue, exists := bucketTags[key]
		if !exists {
			missingAdditionalTags = append(missingAdditionalTags, key)
		} else if exists && tagValue != val {
			assert.Fail(t, fmt.Sprintf("Tag [%s] has incorrect value [%s], expected [%s]", key, tagValue, val))
		}
	}
	if len(missingAdditionalTags) > 0 {
		assert.Fail(t, fmt.Sprintf("The bucket is missing additional tags [%s]",
			strings.Join(missingAdditionalTags, ", ")))
	}
}

func TestDefaultTags(t *testing.T) {
	// t.Parallel()

	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	awsSession := utils.GetAwsSession(t, awsRegion)

	s3TerraformArgs := NewS3TerraformArgs(t)

	terraformOptions := &terraform.Options{
		TerraformDir: GetS3TerraformDir(),
		EnvVars:      utils.GetDefaultEnvVars(awsRegion),
		Vars:         utils.Map(s3TerraformArgs),
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	AssertBucketCreated(t, terraformOptions, awsSession, s3TerraformArgs.BucketName, awsRegion)
}
