package api

/*const (
	downloadTokenExpire = 10 * time.Minute
)

// protectedHandlers provides router for all requests available for regular users
type protectedHandlers struct {
	root         afero.Fs
	dataStore    protectedStore
	tokenService TokenService
	imgResizer   ImageResizer
	fileCache    FileCache
}

type protectedStore interface {
	FindUserByID(ctx context.Context, userID string) (*store.User, error)
}

type TokenService interface {
	Generate(filename, clientIP string, expires time.Duration) (string, error)
	Parse(token, clientIP string) (userID string, err error)
}

func (h *protectedHandlers) List(w http.ResponseWriter, r *http.Request) {
	user := MustGetUser(r)
	userFs := afero.NewBasePathFs(h.root, user.Scope)
	filePath := filePathFromRequest(r)
	info, err := file.Stat(userFs, filePath)
	if err != nil {
		rest.SendErrorJSONChi(w, r, http.StatusInternalServerError, err, "can't read file", rest.ErrCodeInternal)
		return
	}

	render.JSON(w, r, info)
}

func (h *protectedHandlers) DownloadToken(w http.ResponseWriter, r *http.Request) {
	user := MustGetUser(r)
	if !user.Permissions.CanDownload() {
		rest.SendErrorJSONChi(w, r, http.StatusForbidden, errors.New("access denied"), "access denied", rest.ErrCodeNoPermissions)
		return
	}
	clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		rest.SendErrorJSONChi(w, r, http.StatusForbidden, err, "failed to get client ip address", rest.ErrCodeNoPermissions)
		return
	}
	token, err := h.tokenService.Generate(user.ID, clientIP, downloadTokenExpire)
	if err != nil {
		rest.SendErrorJSONChi(w, r, http.StatusInternalServerError, err, "token generation error", rest.ErrCodeInternal)
		return
	}
	render.PlainText(w, r, token)
}

func (h *protectedHandlers) Download(w http.ResponseWriter, r *http.Request) {
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

	f, err := userFs.Open(filePath)
	if err != nil {
		rest.SendErrorJSONChi(w, r, http.StatusInternalServerError, err, "can't open file", rest.ErrCodeInternal)
		return
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		rest.SendErrorJSONChi(w, r, http.StatusInternalServerError, err, "can't read file info", rest.ErrCodeInternal)
		return
	}

	http.ServeContent(w, r, info.Name(), info.ModTime(), f)
}*/
