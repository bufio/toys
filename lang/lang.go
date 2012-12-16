package lang

import (
	"bufio"
	"os"
	"path/filepath"
)

type Lang struct {
	root    string
	current string
	Set     map[string]map[string]string
}

func NewLang(root string) {
	l := &Lang{}
	l.root = root
	l.Set = make(map[string]map[string]string)
}

// readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned if there is an error with the
// buffered reader.
func readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func (l *Lang) Parse(set string) error {
	_, ok := l.Set[set]
	if ok {
		l.current = set
		return nil
	}

	setFolder := filepath.Join(l.root, set)

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
			//read file
			f, err := os.Open(filepath.Join(setFolder, file.Name()))
			if err != nil {
				continue
			}

			r := bufio.NewReader(f)
			s, err := readln(r)
			if err != nil {
				continue
			}

			m := make(map[string]string)
			for e == nil {
				fmt.Println(s)
				s, err = readln(r)
			}
		}
	}
}
