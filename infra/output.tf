############################################################
# AWS S3 Bucket
############################################################

output "s3_bucket_name" {
  value = module.website_s3_bucket.s3_bucket_id
}

output "cloudfront_distribution_domain_name" {
  value = module.website_cf_cdn.cloudfront_distribution_domain_name
}
