package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/s3cmd-go/config"
	"github.com/urfave/cli/v2"
)

// Put file to bucket (really 'upload') -- s3cmd-go put FILE [FILE....] s3://BUCKET/PREFIX
func PutObjects(config *config.Config, c *cli.Context) error {
	return CmdCopy(config, c)
}

// Get file from bucket (really 'download') -- s3cmd-go get s3://BUCKET/OBJECT LOCAL_FILE
func GetObjects(config *config.Config, c *cli.Context) error {
	return CmdCopy(config, c)
}

// One command to do it all, since get/put/cp should be able to copy from anywhere to anywhere
//  using standard "cp" command semantics
//
func CmdCopy(config *config.Config, c *cli.Context) error {
	args := c.Args().Slice()

	if len(args) < 2 {
		return fmt.Errorf("Not enought arguments to the copy command")
	}

	dst, args := args[len(args)-1], args[:len(args)-1]

	dst_u, err := FileURINew(dst)
	if err != nil {
		return fmt.Errorf("Invalid destination argument")
	}
	if dst_u.Scheme == "" {
		dst_u.Scheme = "file"
	}
	if dst_u.Path == "" {
		dst_u.Path = "/"
	}

	for _, path := range args {
		u, err := FileURINew(path)
		if err != nil {
			return fmt.Errorf("Invalid destination argument")
		}
		if u.Scheme == "" {
			u.Scheme = "file"
		}
		if err := copyCore(config, u, dst_u); err != nil {
			return err
		}
	}

	return nil
}

// Ok, this probably could just be in CopyCmd()
func copyCore(config *config.Config, src, dst *FileURI) error {
	// svc := SessionNew(config)

	if src.Scheme != "file" && src.Scheme != "s3" {
		return fmt.Errorf("cp only supports local and s3 URLs")
	}
	if dst.Scheme != "file" && dst.Scheme != "s3" {
		return fmt.Errorf("cp only supports local and s3 URLs")
	}

	if config.Recursive {
		if src.Scheme == "s3" {
			// Get the remote file list and start copying
			svc, err := SessionForBucket(config, src.Bucket)
			if err != nil {
				return err
			}

			// For recusive we should assume that the src path ends in '/' since it's a directory
			nsrc := src
			if !strings.HasSuffix(src.Path, "/") {
				nsrc = src.SetPath(src.Path + "/")
			}

			basePath := nsrc.Path

			remotePager(config, svc, nsrc.String(), false, func(page *s3.ListObjectsV2Output) {
				for _, obj := range page.Contents {
					src_path := *obj.Key
					fmt.Printf("src_path=%s  basePath=%s\n", src_path, basePath)
					src_path = src_path[len(basePath):]

					fmt.Printf("new src_path = %s\n", src_path)

					// uri := fmt.Sprintf("/%s", src.Host, *obj.Key)
					dst_path := dst.String()
					if strings.HasSuffix(dst.String(), "/") {
						dst_path += src_path
					} else {
						dst_path += "/" + src_path
					}

					dst_uri, _ := FileURINew(dst_path)
					dst_uri.Scheme = dst.Scheme
					src_uri, _ := FileURINew("s3://" + src.Bucket + "/" + *obj.Key)

					copyFile(config, src_uri, dst_uri, true)
				}
			})
		} else {
			// Get the local file list and start copying
			err := filepath.Walk(src.Path, func(path string, info os.FileInfo, _ error) error {
				if info == nil || info.IsDir() {
					return nil
				}

				dst_path := dst.String()
				if strings.HasSuffix(dst.String(), "/") {
					dst_path += path
				} else {
					dst_path += "/" + path
				}
				dst_uri, _ := FileURINew(dst_path)
				dst_uri.Scheme = dst.Scheme
				src_uri, _ := FileURINew("file://" + path)

				return copyFile(config, src_uri, dst_uri, true)
			})
			if err != nil {
				return err
			}
		}
	} else {
		return copyFile(config, src, dst, false)
	}
	return nil
}

// TODO: Handle --recusrive
func DeleteObjects(config *config.Config, c *cli.Context) error {
	args := c.Args().Slice()

	svc := SessionNew(config)

	buckets := make(map[string][]*s3.ObjectIdentifier, 0)

	for _, path := range args {
		u, err := FileURINew(args[0])

		if err != nil || u.Scheme != "s3" {
			return fmt.Errorf("rm requires buckets to be prefixed with s3://")
		}

		if (u.Path == "" || strings.HasSuffix(u.Path, "/")) && !config.Recursive {
			return fmt.Errorf("Parameter problem: Expecting S3 URI with a filename or --recursive: %s", path)
		}

		objects := buckets[u.Bucket]
		if objects == nil {
			objects = make([]*s3.ObjectIdentifier, 0)
		}
		buckets[u.Bucket] = append(objects, &s3.ObjectIdentifier{Key: u.Key()})
	}

	// FIXME: Limited to 1000 objects, that's that shouldn't be an issue, but ...
	for bucket, objects := range buckets {
		bsvc, err := SessionForBucket(config, bucket)
		if err != nil {
			return err
		}

		if config.Recursive {
			for _, obj := range objects {
				uri := fmt.Sprintf("s3://%s/%s", bucket, *obj.Key)

				remotePager(config, svc, uri, false, func(page *s3.ListObjectsV2Output) {
					olist := make([]*s3.ObjectIdentifier, 0)
					for _, item := range page.Contents {
						olist = append(olist, &s3.ObjectIdentifier{Key: item.Key})

						fmt.Printf("delete: s3://%s/%s\n", bucket, *item.Key)
					}

					if !config.DryRun {
						params := &s3.DeleteObjectsInput{
							Bucket: aws.String(bucket), // Required
							Delete: &s3.Delete{
								Objects: olist,
							},
						}

						_, err := bsvc.DeleteObjects(params)
						if err != nil {
							fmt.Println("Error removing")
						}
					}
				})
			}
		} else if !config.DryRun {
			params := &s3.DeleteObjectsInput{
				Bucket: aws.String(bucket), // Required
				Delete: &s3.Delete{ // Required
					Objects: objects,
				},
			}

			_, err := bsvc.DeleteObjects(params)
			if err != nil {
				return err
			}
		}
		for _, objs := range objects {
			fmt.Printf("delete: s3://%s/%s\n", bucket, *objs.Key)
		}
	}

	return nil
}

func ModifyMetadata(config *config.Config, c *cli.Context) error {
	for _, arg := range c.Args().Slice() {
		u, err := FileURINew(arg)
		if err != nil {
			return fmt.Errorf("Invalid destination argument")
		}
		if u.Scheme != "s3" {
			return fmt.Errorf("only works on S3 objects")
		}
		if err := copyFile(config, u, u, false); err != nil {
			return err
		}
	}
	return nil
}
