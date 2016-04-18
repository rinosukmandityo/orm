package main

import (
    "os"
    "path/filepath"
    "strings"
    "errors"
    "io/ioutil"
)

func isAbsolutePath(s string)bool{
    if isLinux && strings.HasPrefix(s, string(os.PathSeparator)){
        return true
    }
    
    if !isLinux && strings.Contains(s, ":\\"){
        return true
    }
    
    return false
}

func makePath(s string) string{
    if isAbsolutePath(s){
        return s
    }
    
    return filepath.Join(wd, s)
}

func getOrms(p string)(files []os.FileInfo, err error){
    fi, e := os.Stat(p)
    if e!=nil {
        err = errors.New("Get ORM file error: "+e.Error())
        return 
    }
    
    if !fi.IsDir(){
        files = append(files, fi)
        return
    }
    
    filepath.Walk(p,func(path string, wfi os.FileInfo, e error)error{
        if !wfi.IsDir(){
            if strings.HasSuffix(wfi.Name(), ".orm"){
                files = append(files, wfi)
            }
        }
        return nil
    })
    
    return
}

func genGo(fi os.FileInfo, out string)error{
    log.Info("Processing " + fi.Name())
    
    fn := filepath.Join(wd, fi.Name())
    var txts []string
    if bs, e := ioutil.ReadFile(fn); e!=nil {
        return errors.New("Open error " +  e.Error())
    } else {
       txts = strings.Split(string(bs),"\n")
    }
    
    var e error
    var pkgName string
    var sm StructModel

    for idx, txt := range txts{
        if idx==0 {
            pkgName,e=getPackage(txt)
            if e!=nil {
                return e
            }
            if pkgName!=sm.PkgName && sm.PkgName!="" {
                e = sm.Write(out)
                if e!=nil {
                    return e
                }
            }
            sm.PkgName = pkgName
        }
    }
    
    log.Info("Processing " + fi.Name() + " done")
    return nil
}