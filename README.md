# KYLE_Challenge
Author:

Kyle Speer

January 2023


Contact Info:

610-996-7373

kyle.d.speer@gmail.com

[LinkedIn](https://www.linkedin.com/in/kyle-d-speer/)



## Infrastructure
![Deploy](https://github.com/kspeer825/KYLE_Challenge/blob/main/.github/workflows/deploy_infra.yml/badge.svg)

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

The static webpage is now available at `<CLOUDFRONT_DOMAIN>.cloudfront.net`.

![Alt text](/infra/demo/200_success.png?raw=true "200 Success")

The `index.html` lives in a bucket in s3, but is only accessible via HTTPS connection from the Cloudfront CDN.

![Alt text](/infra/demo/403_forbidden.png?raw=true "403 Forbidden")

This exercise can be taken further:
- Purchase a custom domain.
- Set up a DNS in Route53 connecting a custom domain to the s3.
- Secure the CDN connections with a public SSL Certficate generated in ACM.

### Securing w/ Self Signed Certificiate:

To generate a self-signed Certificat w/ SSL run included script:

```
$ cd infra/openssl
$ bin/setup.sh d31xfsxbx5d9z3.cloudfront.net
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

NOTE: TESTS ARE NOT INCLUDED!  I just ran out of time.

#### Tear Down The Static Site

```
$ terraform destroy -auto-apply -var="aws_region=us-east-2" -var="aws_access_key=${AWS_ACCESS_KEY_ID}" -var="aws_secret_access_key=${AWS_SECRET_ACCESS_KEY}"
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
