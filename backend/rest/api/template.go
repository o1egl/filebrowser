package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/markbates/pkger"
)

func mustParseTemplates() *template.Template {
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

// parsePkgerTplFiles is the helper for the method and function. If the argument
// template is nil, it is created from the first file.
func parsePkgerTplFiles(t *template.Template, filenames ...string) (*template.Template, error) {
	if len(filenames) == 0 {
		// Not really a problem, but be consistent.
		return nil, fmt.Errorf("template: no files named in call to ParseFiles")
	}
	for _, filename := range filenames {
		b, err := readPkgerFile(filename)
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

func readPkgerFile(filename string) ([]byte, error) {
	file, err := pkger.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return b, nil
}
