package test

import (
	"testing"

	utils "gitlab.internal.knowbe4.com/sre/tf-test-utils"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestS3VersionedBucket(t *testing.T) {
	//t.Parallel()

	awsRegion := aws.GetRandomStableRegion(t, nil, nil)
	awsSession := utils.GetAwsSession(t, awsRegion)

	enableVersioning := true

	tfArgs := NewS3TerraformArgs(t)
	tfArgs.EnableVersioning = &enableVersioning

	terraformOptions := &terraform.Options{
		TerraformDir: GetS3TerraformDir(),
		EnvVars:      utils.GetDefaultEnvVars(awsRegion),
		Vars:         utils.Map(tfArgs),
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	AssertBucketCreated(t, terraformOptions, awsSession, tfArgs.BucketName, awsRegion)

	aws.AssertS3BucketVersioningExists(t, awsRegion, tfArgs.BucketName)
}
