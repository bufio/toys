// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package locale help you do build a internationalize your app. The package require a folder with a
simple struct to work with. Assuming you have the languages folder like this:

	path/to/your/languages
	├── en
	│   ├── index.lang
	│   └── login.lang
	└── vi
	    ├── index.lang
	    └── login.lang

The folder should have some folder names with the language code, these folders call lang-set.
Each lang set should containt some file end with .lang, the file name should represent where the
content mostly use. The .lang file have a basic struct like this:

	key1=value1
	key2=value2

For real example, in the en/index.lang we may have
	hi=Hello!
And in vi/index.lang, we habe:
	hi=Xin chào!
Then you can use the package like this:
	lang := locale.NewLang("path/to/your/languages")
	lang.Parse("en")
	lang.Load("index.lang", "hi") // return "Hello!"
*/
package locale

import (
	"bufio"
	"bytes"
	"github.com/kidstuff/toys/util/errs"
	"os"
	"path/filepath"
)

// Lang manages the locale system
type Lang struct {
	root    string
	current string
	set     map[string]map[string]map[string]string
}

// NewLang returns a new Lang iwth the given language folder
func NewLang(root string) *Lang {
	l := &Lang{}
	l.root = root
	l.set = make(map[string]map[string]map[string]string)
	return l
}

// readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned if there is an error with the
// buffered reader.
func readln(r *bufio.Reader) ([]byte, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return ln, err
}

// SetDefault make a lang-set to be default
func (l *Lang) SetDefault(set string) error {
	_, ok := l.set[set]
	if !ok {
		err := l.Parse(set)
		if err != nil {
			return err
		}
	}
	l.current = set
	return nil
}

// Parse parses the lang-set and cache them
func (l *Lang) Parse(set string) error {
	setFolder := filepath.Join(l.root, set)

	setroot, err := os.Open(setFolder)
	if err != nil {
		return errs.New("lang: cannot open language set folder")
	}
	defer setroot.Close()

	files, err := setroot.Readdir(-1)
	if err != nil {
		return errs.New("lang: cannot list file in language set folder")
	}

	l.set[set] = make(map[string]map[string]string)
	for _, file := range files {
		if !file.IsDir() {
			//read file
			f, err := os.Open(filepath.Join(setFolder, file.Name()))
			if err != nil {
				continue
			}
			defer f.Close()

			r := bufio.NewReader(f)
			s, e := readln(r)
			if e != nil {
				continue
			}

			m := make(map[string]string)
			for e == nil {
				pos := bytes.Index(s, []byte{0x3d})
				m[string(s[:pos])] = string(s[pos+1:])
				s, e = readln(r)
			}
			l.set[set][file.Name()] = m
		}
	}
	l.current = set
	return nil
}

// Load returns a value base on file name and key
func (l *Lang) Load(file, key string) string {
	return l.LoadSet(l.current, file, key)
}

// LoadSet returns a value base on file, set name and key.
// It will return the key if no value exist.
func (l *Lang) LoadSet(set, file, key string) string {
	v, ok := l.set[set][file][key]
	if !ok {
		return key
	}

	return v
}
