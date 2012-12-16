package view

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
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
	Root     string
	Set      map[string]*ViewSet
	current  string
	funcsMap template.FuncMap
}

func (v *View) AddFunc(name string, f interface{}) error {
	if r := reflect.TypeOf(f); r.Kind() == reflect.Func {
		if r.NumOut() > 2 {
			return fmt.Errorf("view: %s", "function must have no more than 2 output parameter")
		}
		if r.NumOut() == 2 {
			o := r.Out(1)
			_, ok := o.MethodByName("Error")
			if !ok {
				return fmt.Errorf("view: %s", "function must have the last output parameter implements error")
			}
		}
		v.funcsMap[name] = f
		return nil
	}
	return fmt.Errorf("view: %s", "AddFunc require a valid function")
}

func (v *View) Parse(set string) error {
	_, ok := v.Set[set]
	if ok {
		v.current = set
		return nil
	}

	setFolder := filepath.Join(v.Root, set)

	tmpl := template.Must(template.New("layout.tmpl").Funcs(v.funcsMap).ParseGlob(filepath.Join(setFolder, "shared", "*.tmpl")))
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
			//read file
			b, err := ioutil.ReadFile(filepath.Join(setFolder, file.Name()))
			if err != nil {
				continue
			}
			_, err = p.Parse(string(b))
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
		p.ExecuteTemplate(w, "layout.tmpl", data)
		return
	}
	fmt.Fprintf(w, "%#v", data)
}

func NewView(root string) *View {
	v := &View{}
	v.Root = root
	v.Set = make(map[string]*ViewSet)
	v.funcsMap = template.FuncMap{}

	return v
}
