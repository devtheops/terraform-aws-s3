package test

import (
	"fmt"
	"strings"
	"testing"

	utils "gitlab.internal.knowbe4.com/sre/tf-test-utils"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestCreateBucket(t *testing.T) {
	// t.Parallel()

	expectedBucketRegion := aws.GetRandomStableRegion(t, nil, nil)

	tfArgs := NewS3TerraformArgs(t)

	awsSession := utils.GetAwsSession(t, expectedBucketRegion)

	terraformOptions := &terraform.Options{
		TerraformDir: GetS3TerraformDir(),
		EnvVars:      utils.GetDefaultEnvVars(expectedBucketRegion),
		Vars:         utils.Map(tfArgs),
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	AssertBucketCreated(t, terraformOptions, awsSession, tfArgs.BucketName, expectedBucketRegion)
}

func TestCreateBucketWithoutBucketName(t *testing.T) {
	// t.Parallel()

	expectedApp := fmt.Sprintf("terratest-created-resource-%s", strings.ToLower(random.UniqueId()))
	expectedService := strings.ToLower(random.UniqueId())
	expectedEnvironment := "Go Test"
	expectedBucketRegion := aws.GetRandomStableRegion(t, nil, nil)

	terraformOptions := &terraform.Options{
		TerraformDir: GetS3TerraformDir(),
		EnvVars:      utils.GetDefaultEnvVars(expectedBucketRegion),
		Vars: map[string]interface{}{
			"app":         expectedApp,
			"service":     expectedService,
			"environment": expectedEnvironment,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)
}
