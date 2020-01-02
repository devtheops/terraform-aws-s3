package test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
)

func GetS3Tags(t *testing.T, awsSession *session.Session, bucketName *string) []*s3.Tag {
	svc := s3.New(awsSession)
	bucketInput := &s3.GetBucketTaggingInput{
		Bucket: bucketName,
	}
	result, err := svc.GetBucketTagging(bucketInput)
	assert.NoError(t, err)
	return result.TagSet
}
