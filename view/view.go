// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package view implements a basic template system on top of html/template package.
The package require a special structed folder like this example:

	path/to/template/
	├── christmas
	│   ├── page1.tmpl
	│   ├── page2.tmpl
	│   └── shared
	│       ├── menu.tmpl
	│       └── layout.tmpl
	└── default
	    ├── page1.tmpl
	    ├── page2.tmpl
	    └── shared
	        ├── menu.tmpl
	        └── layout.tmpl

Assuming your template folder located at path/to/template. There is some rules with this folder:

	In the template folder should contain some sub-folder called view-set.
	Each view-set shloud contain some xyz.tmpl (where xyz is the file name) and a folder name "shared".
	In the "shared" folder must contain a file "layout.tmpl" (all the template file must end with .tmpl).

There is some rules for "xyz.tmpl" files:

	All the contain of the file must in the {{define "page"}} ... {{end}} block.
	You can call the shared content (insert the content of the tmpl file in shared folder) by {{template "menu.tmpl"}} etc.

The "layout.tmpl" in shared folder is the main layout. The content of "xyz.tmpl" files should be
embedded in this file. You must put {{template "page" .}} some where in this file.

For more detail, see the tutorial.
*/
package view

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/kidstuff/toys/locale"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
)

// ViewSet represents the view-set sub folder.
type ViewSet struct {
	diver *template.Template
	page  map[string]*template.Template
}

// NewViewSet return a new ViewSet
func NewViewSet() *ViewSet {
	v := &ViewSet{}
	v.page = make(map[string]*template.Template)
	return v
}

// View manages the whole template system.
type View struct {
	root     string
	set      map[string]*ViewSet
	current  string
	funcsMap template.FuncMap
	resource string
}

// NewView returns a new View with the given location of the template folder.
func NewView(root string) *View {
	v := &View{}
	v.root = root
	v.set = make(map[string]*ViewSet)
	v.funcsMap = template.FuncMap{}
	v.funcsMap["resource"] = func(uri string) string {
		return v.resource + uri
	}
	v.funcsMap["equal"] = func(a, b interface{}) bool {
		return a == b
	}
	v.funcsMap["plus"] = func(a, b int) int {
		return a + b
	}
	v.funcsMap["indent"] = func(s string, n int) string {
		var buff bytes.Buffer
		for i := 0; i < n; i++ {
			buff.WriteString(s)
		}
		return buff.String()
	}
	return v
}

/*
HandleResource make a handler that serves HTTP for static file that use for template system.
For example if you want handle the static files at example.com/statics/ you should call:
	*View.HandleResource("/statics/", "path/to/statics/folder/")
And then in the .tmpl file you can call {{resource "css/index.css"}}, it will returns
"/statics/css/index.css"
*/
func (v *View) HandleResource(prefix, folder string) {
	v.resource = prefix
	http.Handle(prefix, http.StripPrefix(prefix,
		http.FileServer(http.Dir(folder))))
}

// AddFunc add the function to the template system. You call call the function you added in the .tmpl
// files by {{function-name}}. An error return if you add an invalid function.
// Note: this function must be call before Parse.
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

// SetDefault change the set to default. It call Parse if need.
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

// Parse parses the view-set you want to use. You may call Parse for all view-set you have and then
// switching beetwen them by call SetDefault.
func (v *View) Parse(set string) error {
	setFolder := filepath.Join(v.root, set)

	tmpl := template.Must(template.New("layout.tmpl").Funcs(v.funcsMap).
		ParseGlob(filepath.Join(setFolder, "shared", "*.tmpl")))
	vs := NewViewSet()
	vs.diver = tmpl
	//parse page
	setroot, err := os.Open(setFolder)
	defer setroot.Close()

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
			} else {
				return err
			}
		}
	}

	v.set[set] = vs
	v.current = set
	return nil
}

// Load render the template you specific with name and write it to the Writer.W
func (v *View) Load(w io.Writer, pageName string, data interface{}) error {
	p, ok := v.set[v.current].page[pageName]
	if ok {
		return p.ExecuteTemplate(w, "layout.tmpl", data)
	}
	fmt.Fprintf(w, "%#v", data)
	return errors.New("view: cannot load template")
}

// SetLang set the language use with the current template system. The method must be call before Parse.
// After call this method you can use these command in your .tmpl files:
// 	{{lang "filename.lang" "key"}}
// 	{{langset "set" "filename.lang" "key"}}
func (v *View) SetLang(l *locale.Lang) {
	v.funcsMap["lang"] = func(file, key string) string {
		return l.Load(file, key)
	}
	v.funcsMap["langset"] = func(set, file, key string) string {
		return l.LoadSet(set, file, key)
	}
}
