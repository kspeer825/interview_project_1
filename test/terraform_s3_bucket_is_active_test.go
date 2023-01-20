package test_s3

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/aws"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/stretchr/testify/assert"
)

var awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
var awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
var awsRegion = "us-east-2"
var dryRun = true

func TestTerraformS3Bucket(t *testing.T) {
	// Test Setup
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../infra",
		Vars: map[string]interface{}{
			"aws_region": awsRegion,
			"aws_access_key": awsAccessKey,
			"aws_secret_access_key": awsSecretAccessKey,
		},
	})
	if ( dryRun ) {
		fmt.Printf("\n**********DRYRUN EXECUTION: Skipping terraform init and apply.**********\n" )
	} else {
		// Teardown
		defer terraform.Destroy(t, terraformOptions)
		fmt.Printf("\n**********LIVE EXECUTION**********\n" )
		// Execute
		terraform.InitAndApply(t, terraformOptions)
	}


	// Set Expectations
	expectedBucketName := "kyle.speer.infra.challenge"
	expectedBucketWebsiteEndpoint := fmt.Sprintf(
		"http://%s.s3-website.%s.amazonaws.com",
		expectedBucketName,
		awsRegion,
	)
	expectedErrorMsg := "403 Forbidden"
	expectedWebsiteMsg := "Hello World!"


	// Verify s3 bucket exists
	actualBucketName := terraform.Output(t, terraformOptions, "s3_bucket_name")
	assert.Equal(t, expectedBucketName, actualBucketName)

	// Verify s3 bucket contains index.html w/ proper contents
	contents := aws.GetS3ObjectContents(t, awsRegion, expectedBucketName, "index.html")
	fmt.Print(contents)
	assert.Contains(t, contents, expectedWebsiteMsg)

	// Verify s3 bucket website endpoint is inaccessible via http
	_, err := http_helper.HttpGet(t, expectedBucketWebsiteEndpoint, nil)
	fmt.Print(err)
	assert.Contains(t, err, expectedErrorMsg)
}