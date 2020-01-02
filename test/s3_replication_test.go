package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	utils "gitlab.internal.knowbe4.com/sre/tf-test-utils"

	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/gruntwork-io/terratest/modules/terraform"

	"github.com/gruntwork-io/terratest/modules/aws"
)

func TestS3BucketReplication(t *testing.T) {
	//t.Parallel()

	awsRegion := aws.GetRandomStableRegion(t, nil, nil)
	awsSession := utils.GetAwsSession(t, awsRegion)

	replicationEnabled := true

	tfArgs := NewS3TerraformArgs(t)
	tfArgs.EnableReplication = &replicationEnabled

	terraformOptions := &terraform.Options{
		TerraformDir: GetS3TerraformDir(),
		EnvVars:      utils.GetDefaultEnvVars(awsRegion),
		Vars:         utils.Map(tfArgs),
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	AssertBucketCreated(t, terraformOptions, awsSession, tfArgs.BucketName, awsRegion)

	svc := s3.New(awsSession)

	req := &s3.GetBucketReplicationInput{
		Bucket: &tfArgs.BucketName,
	}

	result, err := svc.GetBucketReplication(req)
	assert.NoError(t, err)

	for _, rule := range result.ReplicationConfiguration.Rules {
		assert.Equal(t, "Enabled", *rule.Status)
	}
}
