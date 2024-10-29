package main

import (
	"context"
	"fmt"
	"github.com/filebrowser/filebrowser/api/gen"
	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	"github.com/ogen-go/ogen/ogenerrors"
	"net/http"
	"os"
	"time"
)

func main() {
	server, err := gen.NewServer(&Handler{}, gen.WithErrorHandler(ErrorHandler))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	httpServer := http.Server{
		ReadHeaderTimeout: time.Second,
		Addr:              ":8080",
		Handler:           server,
	}

	httpServer.ListenAndServe()
}

func ErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	code := ogenerrors.ErrorCode(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	e := jx.GetEncoder()
	e.ObjStart()
	e.FieldStart("code")
	e.Int(code)
	e.FieldStart("message")
	e.StrEscape(err.Error())
	e.ObjEnd()

	_, _ = w.Write(e.Bytes())
}

var _ gen.Handler = (*Handler)(nil)

type Handler struct {
}

func (h *Handler) V1FilesListGet(ctx context.Context, params gen.V1FilesListGetParams) (*gen.FileGroup, error) {
	return &gen.FileGroup{
		Name: gen.NewOptString("foo"),
		Files: []gen.File{
			{
				Name:     gen.NewOptString("file1"),
				Size:     gen.NewOptInt(100),
				Modified: gen.NewOptString("2021-01-01T00:00:00Z"),
				IsDir:    gen.NewOptBool(false),
				Permissions: gen.NewOptFilePermissions(gen.FilePermissions{
					Read:   gen.NewOptBool(true),
					Write:  gen.NewOptBool(true),
					Delete: gen.NewOptBool(false),
					Move:   gen.NewOptBool(false),
					Share:  gen.NewOptBool(true),
				}),
			},
		},
	}, errors.New("some error")
}

func (h *Handler) NewError(ctx context.Context, err error) *gen.ErrorStatusCode {
	return &gen.ErrorStatusCode{
		StatusCode: 500,
		Response: gen.Error{
			Code:    gen.NewOptInt(500),
			Message: gen.NewOptString(err.Error()),
		},
	}
}
