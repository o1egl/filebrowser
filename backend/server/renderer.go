package server

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"sync"
	"text/template"

	"github.com/gin-gonic/gin/render"

	"github.com/filebrowser/filebrowser/v3/assets"
)

type tplEngine struct {
	tpl  *template.Template
	once sync.Once
}

func (e *tplEngine) Instance(name string, data interface{}) render.Render {
	e.once.Do(func() {
		e.tpl = e.loadTemplate()
	})
	tpl := e.tpl
	return &tplRenderer{
		Template: tpl,
		Name:     name,
		Data:     data,
	}
}

func (e *tplEngine) loadTemplate() *template.Template {
	var files []string
	err := fs.WalkDir(assets.FS(), ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && (strings.HasSuffix(d.Name(), ".js") || strings.HasSuffix(d.Name(), ".html")) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return template.Must(loadTplFiles(template.New("").Delims("[{[", "]}]"), files...))
}

// loadTplFiles is the helper that loads templates from the assets. If the argument
// template is nil, it is created from the first file.
func loadTplFiles(t *template.Template, filenames ...string) (*template.Template, error) {
	if len(filenames) == 0 {
		// Not really a problem, but be consistent.
		return nil, fmt.Errorf("template: no files named in call to ParseFiles")
	}
	for _, filename := range filenames {
		b, err := fs.ReadFile(assets.FS(), filename)
		if err != nil {
			return nil, err
		}
		s := string(b)
		name := filename
		// First template becomes return value if not already defined,
		// and we use that one for subsequent New calls to associate
		// all the templates together. Also, if this file has the same name
		// as t, this file becomes the contents of t, so
		//  t, err := New(name).Funcs(xxx).ParseFiles(name)
		// works. Otherwise we create a new template associated with t.
		var tmpl *template.Template
		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
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
	if val := header["Content-Mode"]; len(val) == 0 {
		header["Content-Mode"] = []string{"text/html; charset=utf-8"}
	}
}
