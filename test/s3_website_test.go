package test

import (
	"fmt"
	"strings"
	"testing"

	utils "gitlab.internal.knowbe4.com/sre/tf-test-utils"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"

	"github.com/gruntwork-io/terratest/modules/aws"
)

func TestCreateWebsiteBucket(t *testing.T) {
	//t.Parallel()

	awsRegion := aws.GetRandomStableRegion(t, nil, nil)
	awsSession := utils.GetAwsSession(t, awsRegion)

	enableWebsiteHosting := true
	websiteHostingErrorDocument := fmt.Sprintf("%s/index.html", strings.ToLower(random.UniqueId()))
	websiteHostingIndexDocument := fmt.Sprintf("%s.html", strings.ToLower(random.UniqueId()))

	tfArgs := NewS3TerraformArgs(t)
	tfArgs.EnableWebsiteHosting = &enableWebsiteHosting
	tfArgs.WebsiteHostingErrorDocument = websiteHostingErrorDocument
	tfArgs.WebsiteHostingIndexDocument = websiteHostingIndexDocument

	terraformOptions := &terraform.Options{
		TerraformDir: GetS3TerraformDir(),
		EnvVars:      utils.GetDefaultEnvVars(awsRegion),
		Vars:         utils.Map(tfArgs),
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	AssertBucketCreated(t, terraformOptions, awsSession, tfArgs.BucketName, awsRegion)

	svc := s3.New(awsSession)

	req := s3.GetBucketWebsiteInput{Bucket: &tfArgs.BucketName}

	result, err := svc.GetBucketWebsite(&req)
	assert.NoError(t, err)

	assert.NotNil(t, result)
	assert.Equal(t, websiteHostingIndexDocument, *result.IndexDocument.Suffix)
	assert.Equal(t, websiteHostingErrorDocument, *result.ErrorDocument.Key)
}

// TODO: Write test to ensure enable_static_website is disabled by default
