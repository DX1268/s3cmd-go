# s3cmd-go -- Go version of s3cmd

Command line utility frontend to the [AWS Go SDK](http://docs.aws.amazon.com/sdk-for-go/api/)
for S3.  Inspired by [s3cmd](https://github.com/s3tools/s3cmd) and attempts to be a
drop-in replacement. 

## Features

* Compatible with [s3cmd](https://github.com/s3tools/s3cmd)'s config file
* Supports a subset of s3cmd's commands and parameters
  - including `put`, `get`, `del`, `ls`, `sync`, `cp`
  - commands are much smarter (get, put, cp - can move to and from S3)
* When syncing directories, instead of uploading one file at a time, it 
  uploads many files in parallel resulting in more bandwidth.
* Uses multipart uploads for large files and uploads each part in parallel. This is
  accomplished using the s3manager that comes with the SDK
* More efficent at using CPU and resources on your local machine

## Install

`go get github.com/DX1268/s3cmd-go`

## Configuration

s3cmd-go is compatible with s3cmd's config file, so if you already have that
configured, you're all set. Otherwise you can put this in `~/.s3cfg`:

```ini
[default]
access_key = foo
secret_key = bar
```

You can also point it to another config file with e.g. `$ s3cmd-go --config /path/to/s3cmd.conf ...`.

## Documentation

In general the commands follow `rsync` as a guide for command options or the unix command line 
commands.

### cp

Copy files to and from S3

Example:

```
s3cmd-go cp /path/to/file s3://bucket/key/on/s3
s3cmd-go cp s3://bucket/key/on/s3 /path/to/file
s3cmd-go cp s3://bucket/key/on/s3 s3://another-bucket/some/thing
```

### get

Download a file from S3 -- really an alias for `cp`

### put

Upload a file to S3 -- really an alias for `cp`

### del

Deletes an object or a directory on S3.

Example:

```
s3cmd-go del [--recursive] s3://bucket/key/on/s3/
```

### rm

Alias for `del`

```
s3cmd-go rm [--recursive] s3://bucket/key/on/s3/
```

### sync

Sync a local directory to S3

```
s3cmd-go sync [--delete-removed] /path/to/folder/ s3://bucket/key/on/s3/
```

### mv

Move an object which is already on S3.

Example:

```
s3cmd-go mv s3://sourcebucket/source/key s3://destbucket/dest/key
```

### General Notes about s3cmd commpatability

DONE - 

* s3cmd-go mb s3://BUCKET
* s3cmd-go rb s3://BUCKET
* s3cmd-go ls [s3://BUCKET[/PREFIX]]
* s3cmd-go la
* s3cmd-go put FILE [FILE...] s3://BUCKET[/PREFIX]
* s3cmd-go get s3://BUCKET/OBJECT LOCAL_FILE
* s3cmd-go del s3://BUCKET/OBJECT
* s3cmd-go rm s3://BUCKET/OBJECT
* s3cmd-go du [s3://BUCKET[/PREFIX]]
* s3cmd-go cp s3://BUCKET1/OBJECT1 s3://BUCKET2[/OBJECT2]
* s3cmd-go modify s3://BUCKET1/OBJECT
* s3cmd-go sync LOCAL_DIR s3://BUCKET[/PREFIX] or s3://BUCKET[/PREFIX] LOCAL_DIR
* s3cmd-go info s3://BUCKET[/OBJECT]

TODO - for full compatibility (with s3cmd)

* s3cmd-go restore s3://BUCKET/OBJECT
* s3cmd-go mv s3://BUCKET1/OBJECT1 s3://BUCKET2[/OBJECT2]

* s3cmd-go setacl s3://BUCKET[/OBJECT]
* s3cmd-go setpolicy FILE s3://BUCKET
* s3cmd-go delpolicy s3://BUCKET
* s3cmd-go setcors FILE s3://BUCKET
* s3cmd-go delcors s3://BUCKET
* s3cmd-go payer s3://BUCKET
* s3cmd-go multipart s3://BUCKET [Id]
* s3cmd-go abortmp s3://BUCKET/OBJECT Id
* s3cmd-go listmp s3://BUCKET/OBJECT Id
* s3cmd-go accesslog s3://BUCKET
* s3cmd-go sign STRING-TO-SIGN
* s3cmd-go signurl s3://BUCKET/OBJECT <expiry_epoch|+expiry_offset>
* s3cmd-go fixbucket s3://BUCKET[/PREFIX]
* s3cmd-go ws-create s3://BUCKET
* s3cmd-go ws-delete s3://BUCKET
* s3cmd-go ws-info s3://BUCKET
* s3cmd-go expire s3://BUCKET
* s3cmd-go setlifecycle FILE s3://BUCKET
* s3cmd-go dellifecycle s3://BUCKET
* s3cmd-go cflist
* s3cmd-go cfinfo [cf://DIST_ID]
* s3cmd-go cfcreate s3://BUCKET
* s3cmd-go cfdelete cf://DIST_ID
* s3cmd-go cfmodify cf://DIST_ID
* s3cmd-go cfinvalinfo cf://DIST_ID[/INVAL_ID]
