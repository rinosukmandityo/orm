package main

import (
    "flag"
    "os"
    "runtime"
    "github.com/eaciit/toolkit"
)

var sourceFile, outPath string
var fieldImports = make(map[string]string)

func init() {
	fieldImports["time"] = "time"
	fieldImports["bson"] = "gopkg.in/mgo.v2/bson"
}

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d
	}()

    isLinux = runtime.GOOS!="windows"
    log, _ = toolkit.NewLog(true,false,"","","")
    flagSource = flag.String("source",".","name or pattern of source file(s). If empty it will be default.orm")
    flagOut = flag.String("out",wd,"Output path. If empty it will be current working directory")
)

func check(e error, fatal bool, pre string){
    if e==nil {
        return
    }
    
    if pre=="" {
        log.Error(pre + " " + e.Error())
    } else {
        log.Error(e.Error())
    }

    if fatal{
        os.Exit(200)
    }
}

func main(){
    flag.Parse()
    source := makePath(*flagSource)
	outPath := makePath(*flagOut)
    
    log.Info(toolkit.Sprintf("Generating *.go files\nSource: %s\nOutput Path: %s", source, outPath))

    fileInfos, e := getOrms(source)
    check(e, true, "")
    
    for _, fi := range fileInfos{
        e := genGo(fi, outPath)
        check(e,true,"Gen-Go")
    }
}