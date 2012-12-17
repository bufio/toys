package lang

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"path/filepath"
)

type Lang struct {
	root    string
	current string
	set     map[string]map[string]map[string]string
}

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

func (l *Lang) Parse(set string) error {
	setFolder := filepath.Join(l.root, set)

	setroot, err := os.Open(setFolder)
	if err != nil {
		return errors.New("lang: cannot open language set folder")
	}

	files, err := setroot.Readdir(-1)
	if err != nil {
		return errors.New("lang: cannot list file in language set folder")
	}

	l.set[set] = make(map[string]map[string]string)
	for _, file := range files {
		if !file.IsDir() {
			//read file
			f, err := os.Open(filepath.Join(setFolder, file.Name()))
			if err != nil {
				continue
			}

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

func (l *Lang) Load(file, key string) string {
	return l.LoadSet(l.current, file, key)
}

func (l *Lang) LoadSet(set, file, key string) string {
	return l.set[set][file][key]
}
