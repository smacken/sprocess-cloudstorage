package sprocess_cloudstorage

import (
	"bytes"
	"errors"
	"github.com/hyperboloide/sprocess"
	ctx "golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
	"google.golang.org/appengine/urlfetch"
	"io"
	"log"
	"net/http"
	"path"
	"strings"
)

var (
	InvalidBucket           = errors.New("Bucket name is undefined")
	InvalidGoogleContext    = errors.New("Context is invalid or not provided")
	InvalidCloudStore       = errors.New("Invalid Cloud store configaration")
	UnavailableCloudService = errors.New("Google Cloud service is unavailable or configured incorrectly")
)

type CloudStore struct {
	Name    string
	Bucket  string // Name of the bucket to store files into
	Project string // The cloud project name
	ctx.Context
	service   *storage.Service
	insert    *storage.ObjectsInsertCall
	MediaLink string // the link to the stored file
	data      *sprocess.Data
	Acl       string // file permissions
	Prefix    string // optional prefix/suffix to append to each file being saved
	Suffix    string
}

func (c *CloudStore) GetName() string {
	return c.Name
}

func (cloud *CloudStore) Start() error {
	if cloud.Bucket == "" {
		return InvalidBucket
	}
	if cloud.Context == nil {
		return InvalidGoogleContext
	}
	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(cloud.Context, storage.DevstorageFullControlScope),
			Base:   &urlfetch.Transport{Context: cloud.Context},
		},
	}

	service, err := storage.New(client)
	if err != nil {
		log.Printf("Unable to get storage service %v", err)
		return err
	}
	cloud.service = service
	if _, err := service.Buckets.Get(cloud.Bucket).Do(); err == nil {
		log.Printf("Got storage bucket %v %v", cloud.Bucket, err)
	} else {
		if _, err := service.Buckets.Insert(cloud.Project, &storage.Bucket{Name: cloud.Bucket}).Do(); err == nil {
			log.Printf("Created bucket: %v", cloud.Bucket)
		} else {
			return err
		}
	}
	return nil
}

func (c *CloudStore) NewWriter(id string, d *sprocess.Data) (io.WriteCloser, error) {
	f := c.getFileName(id, d)
	if c.service == nil || c.Bucket == "" {
		log.Print("no service")
		return nil, InvalidCloudStore
	}
	c.insert = c.service.Objects.Insert(c.Bucket, &storage.Object{Name: f})
	c.data = d
	return c, nil
}

func (c *CloudStore) getFileName(id string, d *sprocess.Data) string {
	f, err := d.Get("filename")
	if err != nil || f.(string) == "" {
		f = ".jpg"
	}
	filename := id
	if c.Prefix != "" {
		filename = c.Prefix + filename
	}
	if c.Suffix != "" {
		filename = filename + c.Suffix
	}
	return filename + path.Ext(strings.ToLower(f.(string)))
}

func (c *CloudStore) Write(p []byte) (n int, err error) {
	if c.insert == nil {
		return 0, UnavailableCloudService
	}
	obj, err := c.insert.Media(bytes.NewReader(p)).PredefinedAcl(c.Acl).Do()
	if err != nil {
		return 0, err
	}
	c.MediaLink = obj.MediaLink
	c.data.Set("medialink", c.MediaLink)
	return int(obj.Size), err
}

func (c *CloudStore) Close() error {
	return nil
}
