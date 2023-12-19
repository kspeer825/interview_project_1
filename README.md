# Interview Project 1

Author: Kyle Speer

Date: January 2023

Contact Info: 610-996-7373 | kyle.d.speer@gmail.com | [LinkedIn](https://www.linkedin.com/in/kyle-d-speer/)

Description: A project submission for an interview. The prompt is in two parts.

  1. Infrastructure: Create and deploy a running instance of a `Hello World` web server using a configuration management tool of your choice.
    a. Secure this application and host such that only appropriate ports are publicly exposed and any http requests are redirected to https. This should be automated using a configuration management tool of your choice and you should feel free to use a self-signed certificate for the web server.
	b. Develop and apply automated tests to validate the correctness of the server configuration.

  2. Coding: Solve the following problem in Python https://www.hackerrank.com/challenges/validating-credit-card-number/problem


## Infrastructure
![Deploy](https://github.com/kspeer825/KYLE_Challenge/actions/workflows/deploy_infra.yml/badge.svg)

### DEMO
The infrastructure can be deployed via GH Actions [here](https://github.com/kspeer825/KYLE_Challenge/actions/workflows/deploy_infra.yml).

### Prerequisites
The following tools are needed to stand up this infrastructure:

```
$ openssl version
LibreSSL 2.8.3

$ terraform -version
Terraform v1.3.7
on darwin_arm64

$ aws --version
aws-cli/2.9.15 Python/3.9.11 Darwin/21.6.0 exe/x86_64 prompt/of
```

Additionally, you must have an AWS account and non-root user with proper permissions.


### Standing Up The Static Site

Set AWS Credentials.

```
$ export AWS_ACCESS_KEY_ID=<YOUR-KEY-ID>
$ export AWS_SECRET_ACCESS_KEY=<YOUR-SECRET-KEY>
```

Apply the terraform plan creating a bucket in S3 and a CloudFront distribution (CDN).

```
$ cd infra && terraform init
$ terraform apply -auto-apply -var="aws_region=us-east-2" -var="aws_access_key=${AWS_ACCESS_KEY_ID}" -var="aws_secret_access_key=${AWS_SECRET_ACCESS_KEY}"
```

The static webpage is now available at `<DISTIBUTION_ID>.cloudfront.net`.

![Alt text](/infra/demo/200_success.png?raw=true "200 Success")

The `index.html` lives in a bucket in s3, but is only accessible via HTTPS connection from the Cloudfront CDN.

![Alt text](/infra/demo/403_forbidden.png?raw=true "403 Forbidden")

This exercise can be taken further:
- Purchase a custom domain.
- Set up a DNS in Route53 connecting a custom domain to the s3 bucket.
- Secure the CDN connections with a public SSL Certficate generated in ACM.

### Securing w/ Self Signed Certificiate:

To generate a self-signed Certificate w/ SSL run the included script:

```
$ cd infra/openssl
$ bin/setup.sh <CLOUDFRONT_DOMAIN>
```

Then the certificate `server.pem` and key `serverKey.pem` can be imported to ACM.

That certificate can then be linked to the Cloudfront distribution in order to enable SSL.

Example Output:

```
$ bin/setup.sh d31xfsxbx5d9z3.cloudfront.net
Creating Certificate Authority (CA)...
Generating a 2048 bit RSA private key
.....+++
...........+++
writing new private key to 'rootKey.pem'
-----
Creating private key...
Generating RSA private key, 2048 bit long modulus
..............................................................+++
..................................................+++
e is 65537 (0x10001)
Configuring CSR (Certificate Signing Request)...
Generating CSR...
Configuring SSL (Secure Socket Layer)...
Signature ok
subject=/C=US/ST=Pennsylvania/L=Philadelphia/O=speer/OU=kylespeer/CN=d31xfsxbx5d9z3.cloudfront.net
Getting CA Private Key
```

### Tearing Down The Static Site

```
$ terraform destroy -auto-apply -var="aws_region=us-east-2" -var="aws_access_key=${AWS_ACCESS_KEY_ID}" -var="aws_secret_access_key=${AWS_SECRET_ACCESS_KEY}"
```

### Testing

Tests written in Go with Terratest live under the [test](https://github.com/kspeer825/KYLE_Challenge/blob/main/test) directory.

You must first install Go and ensure the executable is present in `PATH` in order to run the tests.

To setup the test directory:
```
$ cd test
$ go mod init github.com/kspeer825/KYLE_Challenge
$ go mod tidy
```

And to execute the tests:

```
$  go test -v -timeout 30m

=== RUN   TestTerraformCloudfrontS3StaticSite

**********DRYRUN EXECUTION: Skipping terraform init and apply.**********
TestTerraformCloudfrontS3StaticSite 2023-01-20T13:58:46-05:00 retry.go:91: terraform [output -no-color -json s3_bucket_name]
TestTerraformCloudfrontS3StaticSite 2023-01-20T13:58:46-05:00 logger.go:66: Running command terraform with args [output -no-color -json s3_bucket_name]
TestTerraformCloudfrontS3StaticSite 2023-01-20T13:58:46-05:00 logger.go:66: "kyle.speer.infra.challenge"
TestTerraformCloudfrontS3StaticSite 2023-01-20T13:58:46-05:00 s3.go:126: Read contents from s3://kyle.speer.infra.challenge/index.html
TestTerraformCloudfrontS3StaticSite 2023-01-20T13:58:46-05:00 http_helper.go:59: Making an HTTP GET call to URL http://kyle.speer.infra.challenge.s3-website.us-east-2.amazonaws.com
TestTerraformCloudfrontS3StaticSite 2023-01-20T13:58:46-05:00 retry.go:91: terraform [output -no-color -json cloudfront_distribution_domain_name]
TestTerraformCloudfrontS3StaticSite 2023-01-20T13:58:46-05:00 logger.go:66: Running command terraform with args [output -no-color -json cloudfront_distribution_domain_name]
TestTerraformCloudfrontS3StaticSite 2023-01-20T13:58:47-05:00 logger.go:66: "d32st815qezhpt.cloudfront.net"
TestTerraformCloudfrontS3StaticSite 2023-01-20T13:58:47-05:00 http_helper.go:59: Making an HTTP GET call to URL https://d32st815qezhpt.cloudfront.net
--- PASS: TestTerraformCloudfrontS3StaticSite (1.09s)
PASS
ok      github.com/kspeer825/KYLE_Challenge     1.877s
```
Note the above output is executed with `dryRun = true` relying on existing infrastruture.

When `dryRun = false` the plan is actually applied and destroyed. This can take a while to execute as Cloudfront can take up to 20 minutes to fully distribute the CDN.

In practice though execution typically takes 6 - 7 minutes. See abbreviated output below:

```
$  go test -v -timeout 30m

=== RUN   TestTerraformCloudfrontS3StaticSite

**********LIVE EXECUTION**********
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:14:41-05:00 retry.go:91: terraform [init -upgrade=false]
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:14:41-05:00 logger.go:66: Running command terraform with args [init -upgrade=false]
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:14:41-05:00 logger.go:66: Initializing modules...
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:14:41-05:00 logger.go:66:
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:14:41-05:00 logger.go:66: Initializing the backend...
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:14:41-05:00 logger.go:66:
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:14:41-05:00 logger.go:66: Initializing provider plugins...
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:14:41-05:00 logger.go:66: - Reusing previous version of hashicorp/aws from the dependency lock file
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:14:42-05:00 logger.go:66: - Using previously-installed hashicorp/aws v4.50.0
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:14:42-05:00 logger.go:66:
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:14:42-05:00 logger.go:66: Terraform has been successfully initialized!
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:14:42-05:00 logger.go:66:
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:14:42-05:00 logger.go:66: You may now begin working with Terraform. Try running "terraform plan" to see
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:14:42-05:00 logger.go:66: any changes that are required for your infrastructure. All Terraform commands
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:14:42-05:00 logger.go:66: should now work.
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:14:42-05:00 logger.go:66:
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:14:42-05:00 logger.go:66: If you ever set or change modules or backend configuration for Terraform,
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:14:42-05:00 logger.go:66: rerun this command to reinitialize your working directory. If you forget, other
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:14:42-05:00 logger.go:66: commands will detect it and remind you to do so if necessary.
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:14:42-05:00 retry.go:91: terraform [apply -input=false -auto-approve -var aws_access_key=[REDACTED] -var aws_secret_access_key=[REDACTED] -var aws_region=us-east-2 -lock=false]
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:14:42-05:00 logger.go:66: Running command terraform with args [apply -input=false -auto-approve -var aws_access_key=[REDACTED] -var aws_secret_access_key=[REDACTED] -var aws_region=us-east-2 -lock=false]
[ . . . ]
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:17:44-05:00 logger.go:66: Apply complete! Resources: 8 added, 0 changed, 0 destroyed.
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:17:44-05:00 logger.go:66:
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:17:44-05:00 logger.go:66: Outputs:
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:17:44-05:00 logger.go:66:
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:17:44-05:00 logger.go:66: cloudfront_distribution_domain_name = "djiaiskpbnti.cloudfront.net"
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:17:44-05:00 logger.go:66: s3_bucket_name = "kyle.speer.infra.challenge"
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:17:44-05:00 retry.go:91: terraform [output -no-color -json s3_bucket_name]
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:17:44-05:00 logger.go:66: Running command terraform with args [output -no-color -json s3_bucket_name]
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:17:45-05:00 logger.go:66: "kyle.speer.infra.challenge"
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:17:46-05:00 s3.go:126: Read contents from s3://kyle.speer.infra.challenge/index.html
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:17:46-05:00 http_helper.go:59: Making an HTTP GET call to URL http://kyle.speer.infra.challenge.s3-website.us-east-2.amazonaws.com
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:17:46-05:00 retry.go:91: terraform [output -no-color -json cloudfront_distribution_domain_name]
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:17:46-05:00 logger.go:66: Running command terraform with args [output -no-color -json cloudfront_distribution_domain_name]
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:17:47-05:00 logger.go:66: "djiaiskpbnti.cloudfront.net"
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:17:47-05:00 http_helper.go:59: Making an HTTP GET call to URL https://djiaiskpbnti.cloudfront.net
[ . . . ]
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:21:23-05:00 logger.go:66: Destroy complete! Resources: 8 destroyed.
TestTerraformCloudfrontS3StaticSite 2023-01-20T14:21:23-05:00 logger.go:66:
--- PASS: TestTerraformCloudfrontS3StaticSite (402.08s)
PASS
ok      github.com/kspeer825/KYLE_Challenge     402.851s
```

## Coding

Problem from [hackerrank](https://www.hackerrank.com/challenges/validating-credit-card-number/problem).

Solution: [coding/validator.py](https://github.com/kspeer825/KYLE_Challenge/blob/main/coding/validator.py)

Prerequisites:
```
$ python --version
Python 3.10.0
```

Execute validator script and enter input via STDIN.

```
$ python src/validator.py
6
4123456789123456
5123-4567-8912-3456
61234-567-8912-3456
4123356789123456
5133-3367-8912-3456
5123 - 3567 - 8912 - 3456

Valid
Valid
Invalid
Valid
Invalid
Invalid
```
