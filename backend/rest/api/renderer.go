package api

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"text/template"

	"github.com/markbates/pkger"

	"github.com/gin-gonic/gin/render"
)

type tplEngine struct {
	Reload bool

	tpl  *template.Template
	once sync.Once
}

func (e *tplEngine) Instance(name string, data interface{}) render.Render {
	e.once.Do(func() {
		e.tpl = e.loadTemplate()
	})
	tpl := e.tpl
	if e.Reload {
		tpl = e.loadTemplate()
	}
	return &tplRenderer{
		Template: tpl,
		Name:     name,
		Data:     data,
	}
}

func (e *tplEngine) loadTemplate() *template.Template {
	var files []string
	const distPath = "/frontend/dist"
	err := pkger.Walk(distPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".js") || strings.HasSuffix(info.Name(), ".html")) {
			fileName := path[strings.Index(path, distPath):]
			files = append(files, fileName)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return template.Must(parsePkgerTplFiles(template.New("").Delims("[{[", "]}]"), files...))
}

type tplRenderer struct {
	Template *template.Template
	Name     string
	Data     interface{}
}

// Render (HTML) executes template and writes its result with custom ContentType for response.
func (r *tplRenderer) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	if r.Name == "" {
		return r.Template.Execute(w, r.Data)
	}
	return r.Template.ExecuteTemplate(w, r.Name, r.Data)
}

// WriteContentType (HTML) writes HTML ContentType.
func (r *tplRenderer) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{"text/html; charset=utf-8"}
	}
}
