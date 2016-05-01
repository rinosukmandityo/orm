package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
			if strings.HasSuffix(wfi.Name(), ".orm") {
				files = append(files, wfi)
			}
		}
		return nil
	})

	return
}

func genGo(fi os.FileInfo, out string) error {
	log.Info("Processing " + fi.Name())

	fn := filepath.Join(wd, fi.Name())
	var txts []string
	if bs, e := ioutil.ReadFile(fn); e != nil {
		return errors.New("Open error " + e.Error())
	} else {
		txts = strings.Split(string(bs), "\n")
	}

	var e error
	var pkgName string
	sm := new(StructModel)

	commenting := false
    commentTxt := ""
	for idx, txt := range txts {
		if idx == 0 {
			pkgName, e = getPackage(txt)
			if e != nil {
				return e
			}
			sm.PkgName = pkgName
		}

		if commenting {
            if strings.HasSuffix(txt,"*/"){
                commenting=false
                commentTxt += txt[:len(txt)-2] 
                sm.makeComment(commentTxt)
                commentTxt=""     
            } else {
                commentTxt += txt
            }
		} else {
			if strings.HasPrefix(txt, "struct") {
				writeModel(sm, out)
				sm = new(StructModel)
				sm.PkgName = pkgName
				sm.Name = getStructName(txt)
			} else if strings.HasPrefix(txt,"TableName") {

			} else if strings.HasPrefix(txt,"Get") {

			} else if strings.HasPrefix(txt,"Find") {

			} else if strings.HasPrefix(txt,"/*C:") {
				commenting = true
                if len(commentTxt)>4{
                    commentTxt=txt[4:]
                    if strings.HasPrefix(commentTxt,"*/"){
                        commentTxt=commentTxt[:len(commentTxt)-2]
                        commenting=false
                        sm.makeComment(commentTxt)
                        commentTxt=""
                    }
                }
			} else {

			}
		}
	}

	writeModel(sm, out)

	log.Info("Processing " + fi.Name() + " done")
	return nil
}

func writeModel(sm *StructModel, out string) {
	if sm.Name == "" {
		return
	}

	//path := filepath.Join(out, sm.Name+".go")
	e := sm.Write(out)
	if e != nil {
		log.Error("Writing " + sm.Name + ".go error: " + e.Error())
	} else {
		log.Info("Writing " + sm.Name + ".go success")
	}
}
