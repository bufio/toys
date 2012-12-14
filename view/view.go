package view

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

type Writer interface {
	Write(p []byte) (n int, err error)
}

type ViewSet struct {
	diver *template.Template
	page  map[string]*template.Template
}

func NewViewSet() *ViewSet {
	v := &ViewSet{}
	v.page = make(map[string]*template.Template)
	return v
}

type View struct {
	Root    string
	Set     map[string]*ViewSet
	current string
}

func (v *View) Parse(set string) error {
	_, ok := v.Set[set]
	if ok {
		v.current = set
		return nil
	}

	setFolder := filepath.Join(v.Root, set)
	tmpl := template.Must(template.ParseGlob(filepath.Join(setFolder, "shared", "*.tmpl")))
	vs := NewViewSet()
	vs.diver = tmpl

	//parse page
	root, err := os.Open(setFolder)
	if err != nil {
		return err
	}

	files, err := root.Readdir(-1)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() {
			p, err := tmpl.Clone()
			if err != nil {
				continue
			}

			_, err = p.ParseFiles(filepath.Join(setFolder, file.Name()))
			if err == nil {
				vs.page[file.Name()] = p
			}
		}
	}

	v.Set[set] = vs
	v.current = set
	return nil
}

func (v *View) Load(w Writer, pageName string, data interface{}) {
	p, ok := v.Set[v.current].page[pageName]
	if ok {
		p.Execute(w, data)
		return
	}
	fmt.Fprintf(w, "data: %#v", data)
}

func NewView(root string) *View {
	v := &View{}
	v.Root = root
	v.Set = make(map[string]*ViewSet)
	return v
}
