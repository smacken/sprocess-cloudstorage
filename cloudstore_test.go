package sprocess_cloudstorage_test

import (
	cloudStorage "github.com/smacken/sprocess-cloudstorage"
	"testing"
)

func TestCloudStoreInstantiation(t *testing.T) {
	t.Log("Should instantiate a cloudstore output")
	c := &cloudStorage.CloudStore{
		Name:    "cloud",
		Bucket:  "Bucket",
		Project: "Project",
		Acl:     "publicRead",
	}

	if c.Bucket != "Bucket" {
		t.Errorf("Expected context to have fields set: %v", c)
	}
}

func BenchmarkCloudStoreCtor(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = &cloudStorage.CloudStore{
			Name:    "cloud",
			Bucket:  "Bucket",
			Project: "Project",
			Acl:     "publicRead",
		}
	}
}
