package main

import (
	"github.com/hyperboloide/sprocess"
	"github.com/labstack/echo"
	cloudStorage "github.com/smacken/sprocess-cloudstorage"
	"google.golang.org/appengine"
	"log"
	"net/http"
)

var process *sprocess.HTTP

func initProcess(r *http.Request) {
	ctx := appengine.NewContext(r)
	cloudStoreLarge := &cloudStorage.CloudStore{
		Name:    "cloud",
		Bucket:  "Bucket",
		Project: "Project",
		Acl:     "publicRead",
	}
	cloudStoreLarge.Context = ctx
	if err := cloudStoreLarge.Start(); err != nil {
		log.Fatal(err)
	}

	imageLarge := &sprocess.Image{
		Operation: sprocess.ImageThumbnail,
		Height:    600,
		Width:     600,
		Output:    "jpg",
		Name:      "image",
	}
	if err := imageLarge.Start(); err != nil {
		log.Fatal(err)
	}

	largeImage := &sprocess.Tee{
		Encoders: []sprocess.Encoder{imageLarge},
		Output:   cloudStoreLarge,
		Name:     "tee",
	}

	if err := largeImage.Start(); err != nil {
		log.Fatal(err)
	}

	// process -> tee -> imageLarge -> cloudstorage
	//		   		  -> thumbnail  -> cloudstorage

	process = &sprocess.HTTP{
		Encoders: []sprocess.Encoder{imageLarge},
		Output:   cloudStoreLarge,
		//Input:    cloud,
		//Delete: uploads,
	}
}

func Upload(c *echo.Context) error {
	initProcess(c.Request())
	if data, err := process.Encode(c.Response(), c.Request(), sprocess.GenId()); err != nil {
		log.Fatal(err)
		return err
	} else {
		str, ok := data["medialink"].(string)
		if !ok {
			/* act on str */
		}
		return c.String(http.StatusOK, str)
	}
}
