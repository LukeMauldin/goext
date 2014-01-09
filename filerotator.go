package ext

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type fileRotator struct {
	Dirpath     string
	Basename    string
	Extension   string
	Maxnumfiles int
	Maxsize     int
	written     int
	intwidth    int
	w           *os.File
}

func NewFileRotator(path, baseName, extension string, maxFiles, maxSize int) io.Writer {
	fr := &fileRotator{Dirpath: path, Basename: baseName, Extension: extension, Maxnumfiles: maxFiles, Maxsize: maxSize}

	//Call rotate once to setup files
	fr.Rotate()

	return fr
}

func (fr *fileRotator) Write(p []byte) (n int, er error) {
	if len(p)+fr.written > fr.Maxsize {
		er := fr.Rotate()
		if er != nil {
			return 0, er
		}
	}
	n, er = fr.w.Write(p)
	fr.written += n
	return
}

func (fr *fileRotator) Rotate() error {
	if fr.w != nil {
		fr.w.Close()
	}
	infos, er := ioutil.ReadDir(fr.Dirpath)
	if er != nil {
		if os.IsNotExist(er) {
			os.Mkdir(fr.Dirpath, os.ModeDir)
		} else {
			return er
		}
	}
	for i := len(infos) - 1; i >= 0; i-- {
		info := infos[i]
		fname := info.Name()
		fpath := filepath.Join(fr.Dirpath, fname)
		if strings.HasPrefix(fname, fr.Basename) && strings.HasSuffix(fname, fr.Extension) {
			index, er := strconv.Atoi(fname[len(fr.Basename) : len(fname)-4])
			if er != nil {
				fmt.Println(er, index)
			}
			if index > fr.Maxnumfiles {
				os.Remove(fpath)
			} else {
				newfpath := filepath.Join(fr.Dirpath, fmt.Sprintf("%s%0*d.%s", fr.Basename, fr.intwidth, index+1, fr.Extension))
				os.Rename(fpath, newfpath)
			}
		}
	}
	npath := filepath.Join(fr.Dirpath, fmt.Sprintf("%s%0*d.%s", fr.Basename, fr.intwidth, 1, fr.Extension))
	fr.w, er = os.Create(npath)
	fr.written = 0
	return er
}
