package view

import (
	"fmt"
	"github.com/openvn/toys/lang"
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
	root     string
	set      map[string]*ViewSet
	current  string
	funcsMap template.FuncMap
	Resource string
}

func NewView(root string) *View {
	v := &View{}
	v.root = root
	v.set = make(map[string]*ViewSet)
	v.funcsMap = template.FuncMap{}
	v.funcsMap["resource"] = func(uri string) string {
		return v.Resource + uri
	}
	return v
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

func (v *View) SetDefault(set string) error {
	_, ok := v.set[set]
	if !ok {
		err := v.Parse(set)
		if err != nil {
			return err
		}
	}
	v.current = set
	return nil
}

func (v *View) Parse(set string) error {
	setFolder := filepath.Join(v.root, set)

	tmpl := template.Must(template.New("layout.tmpl").Funcs(v.funcsMap).
		ParseGlob(filepath.Join(setFolder, "shared", "*.tmpl")))
	vs := NewViewSet()
	vs.diver = tmpl
	//parse page
	setroot, err := os.Open(setFolder)
	if err != nil {
		return err
	}

	files, err := setroot.Readdir(-1)
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

	v.set[set] = vs
	v.current = set
	return nil
}

func (v *View) Load(w Writer, pageName string, data interface{}) {
	p, ok := v.set[v.current].page[pageName]
	if ok {
		p.ExecuteTemplate(w, "layout.tmpl", data)
		return
	}
	fmt.Fprintf(w, "%#v", data)
}

func (v *View) SetLang(l *lang.Lang) {
	v.funcsMap["lang"] = func(file, key string) string {
		return l.Load(file, key)
	}
	v.funcsMap["langset"] = func(set, file, key string) string {
		return l.LoadSet(set, file, key)
	}
}
