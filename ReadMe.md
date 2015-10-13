# sprocess cloud storage
Built upon output from [sprocess](https://github.com/hyperboloide/sprocess)

Stream process output to google cloud storage

## installation
```bash
go get github.com/smacken/sprocess-cloudstorage
```

### Outputs, Inputs

**outputs** can save a stream (to a file or an s3 bucket) and **inputs** read this stream back.

These may also allow for deletetion.

#### Google Cloud Storage

Save to a Google Cloud Storage instance

```go
type CloudStore struct {
	Name string
	// Name of the bucket to store files into
	Bucket string
	// The cloud project name
	Project string
	ctx.Context
	service *storage.Service
	insert  *storage.ObjectsInsertCall
	// the link to the stored file
	MediaLink string
	data      *sprocess.Data
	// file permissions
	Acl string
	// optional prefix/suffix to append to each file being saved
	Prefix string
	Suffix string
}
```

e.g. 
```go
cloudStore := &cloudStorage.CloudStore{
	Name:    "cloud",
	Bucket:  "Bucket",
	Project: "Project",
	Acl:     "publicRead",
}
```