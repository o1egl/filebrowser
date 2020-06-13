//go:generate go-enum --sql --marshal --lower --names --file $GOFILE
package api

import (
	"context"
	"io"
	"time"

	"github.com/spf13/afero"

	"github.com/filebrowser/filebrowser/v3/backend/image"
)

type TokenService interface {
	Generate(filename, clientIP string, expires time.Duration) (string, error)
	Parse(token, clientIP string) (userID string, err error)
}

type imageHandlers struct {
	root                afero.Fs
	tokenService        TokenService
	imgResizer          ImageResizer
	fileCache           FileCache
	enablePreviewResize bool
	enableThumbnails    bool
	dataStore           userGetterStore
}

type ImageResizer interface {
	FormatFromExtension(ext string) (image.Format, error)
	Resize(ctx context.Context, in io.Reader, width, height int, out io.Writer, options ...image.Option) error
}

type FileCache interface {
	Store(ctx context.Context, key string, value []byte) error
	Load(ctx context.Context, key string) ([]byte, bool, error)
	Delete(ctx context.Context, key string) error
}

/*
ENUM(
thumb
big
)
*/
type PreviewSize int

/*func (h *imageHandlers) Preview(w http.ResponseWriter, r *http.Request) {
	previewSize, err := ParsePreviewSize(r.URL.Query().Get("size"))
	if err != nil {
		rest.SendErrorJSONChi(w, r, http.StatusBadRequest, err, "incorrect preview size", rest.ErrBadRequest)
		return
	}

	clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		rest.SendErrorJSONChi(w, r, http.StatusForbidden, err, "failed to get client ip address", rest.ErrCodeNoPermissions)
		return
	}
	token := r.URL.Query().Get("token")
	userID, err := h.tokenService.Parse(token, clientIP)
	if err != nil {
		rest.SendErrorJSONChi(w, r, http.StatusForbidden, err, "invalid token", rest.ErrCodeNoPermissions)
		return
	}
	user, err := h.dataStore.FindUserByID(r.Context(), userID)
	if err != nil {
		rest.SendErrorJSONChi(w, r, http.StatusInternalServerError, err, "can't load user", rest.ErrCodeInternal)
		return
	}
	if !user.Permissions.CanDownload() {
		rest.SendErrorJSONChi(w, r, http.StatusForbidden, errors.New("access denied"), "access denied", rest.ErrCodeNoPermissions)
		return
	}
	userFs := afero.NewBasePathFs(h.root, user.Scope)
	filePath := filePathFromRequest(r)

	fileInfo, err := file.Stat(userFs, filePath)
	if err != nil {
		rest.SendErrorJSONChi(w, r, http.StatusInternalServerError, err, "can't open file", rest.ErrCodeInternal)
		return
	}

	switch fileInfo.Type {
	case file.TypeImage:
		h.imagePreview(w, r, userFs, filePath, previewSize)
	default:
		rest.SendErrorJSONChi(w, r, http.StatusBadRequest, err, "can't open file", rest.ErrBadRequest)
		return
	}
}

func (h *imageHandlers) imagePreview(w http.ResponseWriter, r *http.Request, userFs afero.Fs, filePath string, previewSize PreviewSize) {
	inline := r.URL.Query().Get("inline") == "true"
	setContentDisposition(w, inline, filepath.Base(filePath))

	cacheKey := previewCacheKey(filePath, previewSize)
	cachedFile, ok, err := h.fileCache.Load(r.Context(), cacheKey)
	if err != nil {
		rest.SendErrorJSONChi(w, r, http.StatusInternalServerError, err, "failed to load cached preview", rest.ErrCodeInternal)
		return
	}
	if ok {
		_, _ = w.Write(cachedFile)
		return
	}

	fd, err := userFs.Open(filePath)
	if err != nil {
		rest.SendErrorJSONChi(w, r, http.StatusInternalServerError, err, "can't open file", rest.ErrCodeInternal)
		return
	}
	defer fd.Close()

	var (
		width   int
		height  int
		options []img.Option
	)

	format, err := h.imgResizer.FormatFromExtension(filepath.Ext(filePath))
	if err != nil {
		rest.SendErrorJSONChi(w, r, http.StatusInternalServerError, err, "unsupported image format", rest.ErrCodeInternal)
		return
	}

	switch {
	case previewSize == PreviewSizeBig && h.enablePreviewResize && format != img.FormatGif:
		width = 1080
		height = 1080
		options = append(options, img.WithMode(img.ResizeModeFit), img.WithQuality(img.QualityMedium))
	case previewSize == PreviewSizeThumb && h.enableThumbnails:
		width = 128
		height = 128
		options = append(options, img.WithMode(img.ResizeModeFill), img.WithQuality(img.QualityLow), img.WithFormat(img.FormatJpeg))
	default:
		fileInfo, err := fd.Stat()
		if err != nil {
			return
		}
		http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), fd)
		return
	}

	buf := &bytes.Buffer{}
	if err := h.imgResizer.Resize(context.Background(), fd, width, height, buf, options...); err != nil {
		return
	}

	go func() {
		if err := h.fileCache.Store(context.Background(), cacheKey, buf.Bytes()); err != nil {
			fmt.Printf("failed to cache resized image: %v", err)
		}
	}()

	_, _ = w.Write(buf.Bytes())

}

func previewCacheKey(fPath string, previewSize PreviewSize) string {
	return fPath + previewSize.String()
}*/
