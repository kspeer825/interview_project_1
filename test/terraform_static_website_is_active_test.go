package test_static_site

import (
	"fmt"
	"os"
	"testing"
	"regexp"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/aws"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/stretchr/testify/assert"
)

var awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
var awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
var awsRegion = "us-east-2"
var dryRun = false

func TestTerraformCloudfrontS3StaticSite(t *testing.T) {
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
	fileContents, fileReadError := os.ReadFile("../infra/www/index.html")
	if fileReadError != nil {
		panic(fileReadError)
	}
	expectedWebsiteContents := string(fileContents)
	expectedDomainFormat, _ := regexp.Compile(".+\\.cloudfront\\.net$")


	// Verify s3 bucket exists
	actualBucketName := terraform.Output(t, terraformOptions, "s3_bucket_name")
	assert.Equal(t, expectedBucketName, actualBucketName)

	// Verify web contents in s3 bucket match local version
	bucketContents := aws.GetS3ObjectContents(t, awsRegion, expectedBucketName, "index.html")
	assert.Equal(t, expectedWebsiteContents, bucketContents)

	// Verify s3 bucket website endpoint is inaccessible
	_, err := http_helper.HttpGet(t, expectedBucketWebsiteEndpoint, nil)
	assert.Contains(t, err, expectedErrorMsg)

	// Verify Cloudfront distribution exists
	actualDomain := fmt.Sprintf("https://%s", terraform.Output(t, terraformOptions, "cloudfront_distribution_domain_name"))
	assert.Regexp(t, expectedDomainFormat, actualDomain)

	// Verify Cloudfront endpoint is accessible
	code, contents := http_helper.HttpGet(t, actualDomain, nil)
	assert.Equal(t, 200, code)

	// Verify the hosted website html matches the local version
	assert.Equal(t, expectedWebsiteContents, contents)
}