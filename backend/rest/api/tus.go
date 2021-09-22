package api

import (
	"fmt"
	"io/ioutil"
	stdlog "log"

	"github.com/filebrowser/filebrowser/v3/log"
	"github.com/tus/tusd/pkg/filestore"
	tusd "github.com/tus/tusd/pkg/handler"
)

func newUploadHandler() *tusd.Handler {
	fStore := filestore.New("./var/uploads")
	composer := tusd.NewStoreComposer()
	fStore.UseIn(composer)

	handler, err := tusd.NewHandler(tusd.Config{
		StoreComposer:           composer,
		MaxSize:                 0,
		BasePath:                "/uploads/",
		NotifyCompleteUploads:   false,
		NotifyTerminatedUploads: false,
		NotifyUploadProgress:    false,
		NotifyCreatedUploads:    false,
		Logger:                  stdlog.New(ioutil.Discard, "", 0),
		RespectForwardedHeaders: true,
		PreUploadCreateCallback: func(hook tusd.HookEvent) error {
			fmt.Println("upload started", hook.Upload.ID, hook.Upload.MetaData["filename"])
			return nil
		},
		PreFinishResponseCallback: func(hook tusd.HookEvent) error {
			fmt.Println("upload finished", hook.Upload.ID, hook.Upload.MetaData["filename"], hook.Upload.PartialUploads)
			return nil
		},
	})
	if err != nil {
		log.Fatalf("Failed to create tus handler: %v", err)
	}

	return handler
}
