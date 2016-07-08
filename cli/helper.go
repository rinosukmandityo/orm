package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/eaciit/toolkit"
)

func isAbsolutePath(s string) bool {
	if isLinux && strings.HasPrefix(s, string(os.PathSeparator)) {
		return true
	}

	if !isLinux && strings.Contains(s, ":\\") {
		return true
	}

	return false
}

func makePath(s string) string {
	if isAbsolutePath(s) {
		return s
	}

	return filepath.Join(wd, s)
}

func getOrms(p string) (files []os.FileInfo, err error) {
	fi, e := os.Stat(p)
	if e != nil {
		err = errors.New("Get ORM file error: " + e.Error())
		return
	}

	if !fi.IsDir() {
		files = append(files, fi)
		return
	}

	filepath.Walk(p, func(path string, wfi os.FileInfo, e error) error {
		if !wfi.IsDir() {
			if strings.HasSuffix(wfi.Name(), "_orm.json") {
				files = append(files, wfi)
			}
		}
		return nil
	})

	return
}

func genGo(fi os.FileInfo, source, out string) error {
	log.Info("Processing " + fi.Name())
	fn := filepath.Join(source, fi.Name())
	var (
		bs []byte
		e  error
	)
	if bs, e = ioutil.ReadFile(fn); e != nil {
		return errors.New("Open error " + e.Error())
	}

	pkg := new(PackageModel)
	e = toolkit.UnjsonFromString(string(bs), pkg)
	if e != nil {
		return errors.New("Unmarshal JSON: " + e.Error())
	}
	for _, sm := range pkg.Structs {
		e = sm.Write(pkg, out)
		if e != nil {
			return errors.New(toolkit.Sprintf("Write model %s: %s", sm.Name, e.Error()))
		}
		log.Info(toolkit.Sprintf("Writing %s.%s", pkg.Name, sm.Name))
	}

	log.Info("Processing " + fi.Name() + " done")
	return nil
}
