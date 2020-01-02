package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func GetS3TerraformDir() string {
	return "../"
}

func AssertBucketCreated(t *testing.T, terraformOptions *terraform.Options, awsSession *session.Session, expectedBucketName string, expectedBucketRegion string) {
	t.Helper()

	aws.AssertS3BucketExists(t, expectedBucketRegion, expectedBucketName)
	assertDefaultTagsExist(t, awsSession, &expectedBucketName)
	terraformOutput, err := terraform.OutputAllE(t, terraformOptions)
	assert.NoError(t, err)
	assert.Contains(t, terraformOutput, "bucket_arn")
	assert.Contains(t, terraformOutput, "bucket_name")
	assert.Equal(t, expectedBucketName, terraformOutput["bucket_name"])
	assert.Contains(t, terraformOutput, "bucket_region")
	assert.Equal(t, expectedBucketRegion, terraformOutput["bucket_region"])
	assert.Contains(t, terraformOutput, "bucket_domain_name")
}

func assertDefaultTagsExist(t *testing.T, awsSession *session.Session, bucketName *string) {
	t.Helper()

	tags := GetS3Tags(t, awsSession, bucketName)
	bucketTags := map[string]string{}
	for _, tag := range tags {
		bucketTags[*tag.Key] = *tag.Value
	}
	defaultTagKeys := []string{"app", "service", "env"}
	var missingTags []string
	for _, defaultTagKey := range defaultTagKeys {
		if _, exists := bucketTags[defaultTagKey]; !exists {
			missingTags = append(missingTags, defaultTagKey)
		}
	}
	if len(missingTags) > 0 {
		assert.Fail(t, fmt.Sprintf("The [%s] default tags were missing", strings.Join(missingTags, ", ")))
	}
}
