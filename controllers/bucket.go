package controllers

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/s3cmd-go/config"
	"github.com/urfave/cli/v2"
)

// Make bucket -- s3cmd-go mb s3://BUCKET
func MakeBucket(config *config.Config, c *cli.Context) error {
	args := c.Args().Slice()

	svc := SessionNew(config)

	u, err := FileURINew(args[0])
	if err != nil || u.Scheme != "s3" {
		return fmt.Errorf("ls requires buckets to be prefixed with s3://")
	}
	if u.Path != "" {
		return fmt.Errorf("Parameter problem: Expecting S3 URI with just the bucket name set instead of '%s'", args[0])
	}

	params := &s3.CreateBucketInput{
		Bucket: aws.String(u.Bucket),
	}
	if _, err := svc.CreateBucket(params); err != nil {
		return err
	}

	fmt.Printf("Bucket 's3://%s/' created\n", u.Bucket)
	return nil
}

// Remove bucket -- s3cmd-go mb s3://BUCKET
func RemoveBucket(config *config.Config, c *cli.Context) error {
	args := c.Args().Slice()

	svc := SessionNew(config)
	u, err := FileURINew(args[0])
	if err != nil || u.Scheme != "s3" {
		return fmt.Errorf("ls requires buckets to be prefixed with s3://")
	}
	if u.Path != "/" {
		return fmt.Errorf("Parameter problem: Expecting S3 URI with just the bucket name set instead of '%s'", args[0])
	}

	params := &s3.DeleteBucketInput{
		Bucket: aws.String(u.Bucket), // Required
	}
	if _, err := svc.DeleteBucket(params); err != nil {
		return err
	}

	fmt.Printf("Bucket 's3://%s/' removed\n", u.Bucket)
	return nil
}

// List objects or buckets -- s3cmd-go ls [s3://BUCKET[/PREFIX]]
func ListBucket(config *config.Config, c *cli.Context) error {
	args := c.Args().Slice()

	svc := SessionNew(config)

	if len(args) == 0 {
		var params *s3.ListBucketsInput
		resp, err := svc.ListBuckets(params)
		if err != nil {
			return err
		}

		for _, bucket := range resp.Buckets {
			fmt.Printf("%s  s3://%s\n", bucket.CreationDate.Format(DATE_FMT), *bucket.Name)
		}
		return nil
	}

	return listBucket(config, svc, args)
}

func listBucket(config *config.Config, svc *s3.S3, args []string) error {
	for _, arg := range args {
		u, err := FileURINew(arg)
		if err != nil || u.Scheme != "s3" {
			return fmt.Errorf("ls requires buckets to be prefixed with s3://")
		}

		_, err = SessionForBucket(config, u.Bucket)
		if err != nil {
			return err
		}

		todo := []string{arg}

		for len(todo) != 0 {
			var item string
			item, todo = todo[0], todo[1:]

			remotePager(config, svc, item, !config.Recursive, func(page *s3.ListObjectsV2Output) {
				for _, item := range page.CommonPrefixes {
					uri := fmt.Sprintf("s3://%s/%s", u.Bucket, *item.Prefix)

					if config.Recursive {
						todo = append(todo, uri)
					} else {
						fmt.Printf("%16s %9s   %s\n", "", "DIR", uri)
					}
				}
				if page.Contents != nil {
					for _, item := range page.Contents {
						fmt.Printf("%16s %9d   s3://%s/%s\n", item.LastModified.Format(DATE_FMT), *item.Size, u.Bucket, *item.Key)
					}
				}
			})
		}
	}

	return nil
}

// List all object in all buckets -- s3cmd-go la
func ListAll(config *config.Config, c *cli.Context) error {
	args := c.Args().Slice()
	if len(args) != 0 {
		return fmt.Errorf("la shouldn't have arguments")
	}

	svc := SessionNew(config)

	var params *s3.ListBucketsInput
	resp, err := svc.ListBuckets(params)
	if err != nil {
		return err
	}

	for _, bucket := range resp.Buckets {
		uri := fmt.Sprintf("s3://%s", *bucket.Name)

		// Shared with "ls"
		listBucket(config, svc, []string{uri})
	}

	return nil
}

// Disk usage by buckets -- [s3://BUCKET[/PREFIX]]
func GetUsage(config *config.Config, c *cli.Context) error {
	args := c.Args().Slice()

	svc := SessionNew(config)

	// If we're not passed any args, we're going to do all S3 buckets
	if len(args) == 0 {
		var params *s3.ListBucketsInput
		resp, err := svc.ListBuckets(params)
		if err != nil {
			return err
		}

		for _, bucket := range resp.Buckets {
			args = append(args, fmt.Sprintf("s3://%s", *bucket.Name))
		}
	}

	// Get the usage for the buckets
	for _, arg := range args {
		// Only do usage on S3 buckets
		u, err := FileURINew(arg)
		if err != nil || u.Scheme != "s3" {
			continue
		}

		var (
			bucketSize, bucketObjs int64
		)

		remotePager(config, svc, arg, false, func(page *s3.ListObjectsV2Output) {
			for _, obj := range page.Contents {
				bucketSize += *obj.Size
				bucketObjs += 1
			}
		})

		fmt.Printf("%d %d objects %s\n", bucketSize, bucketObjs, arg)
	}

	return nil
}
