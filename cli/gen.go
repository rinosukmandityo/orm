package main 

import (
    "errors"
    "strings"
    "github.com/eaciit/toolkit"
    "path/filepath"
    "os"
    "bufio"
    //"io/ioutil"
)

type StructModel struct{
    PkgName string
    Name string
    TableName string
    Fields []FieldModel
}

type FieldModel struct{
    Name string
    FieldType string
    FieldDefault string
    FieldTag string
}

func (sm *StructModel) Write(path string)error{
    if sm.PkgName=="" || sm.Name=="" {
        return toolkit.Errorf("Both package name and struct name should be defined")
    }
    //return toolkit.Errorf("Fail to write %s.%s : method Write is not yet implemented", sm.PkgName, sm.Name)
    
    filename := filepath.Join(path, strings.ToLower(sm.Name) + ".go")
    //currentCode := ""
    f, e := os.Open(filename)
    if e==nil {
        //bcurrent, _ := ioutil.ReadAll(f)
        //currentCode = string(bcurrent)
        os.Remove(filename)
    }   
    
    f, e = os.Create(filename)
    if e!=nil {
        return toolkit.Errorf("Faild to write %s.%s: %s", sm.PkgName, sm.Name, e.Error())
    }
    defer f.Close()
    
    txts := []string{}
    txts = append(txts,"package " + sm.PkgName)
    txts = append(txts,"/*** OrmGen Auto Generate Code - Start ***/")
    txts = append(txts,"struct " + sm.Name + "{")
    txts = append(txts,"}")
    txts = append(txts,"/*** OrmGen Auto Generate Code - End ***/")
     
    b := bufio.NewWriter(f)
    for _, txt := range txts{
        b.WriteString(txt + "\n")
    }
    
    e = b.Flush()
    if e!=nil {
        return toolkit.Errorf("Faild to write %s.%s: %s", sm.PkgName, sm.Name, e.Error())
    }
    
    toolkit.RunCommand("/bin/sh", "-c", "gofmt -w "+filename)
    return nil
}

func getPackage(txt string)(string,error){
    if !strings.HasPrefix(txt,"package"){
        return "",errors.New("No package has been defined") 
    }
    packages := strings.Split(txt," ")
    if len(packages)<2{
        return "",errors.New("No package has been defined")
    }
    return packages[1],nil
}

func getStructName(s string)string{
    txts := strings.Split(s," ")
    hasStruct := false
    for _,txt:=range txts{
        if !hasStruct && txt=="struct"{
            hasStruct=true
        } else if hasStruct && txt!="" {
            return txt
        }
    }
    return ""
}

func (sm *StructModel) makeGetFn(s string)error{
    return nil
}

func (sm *StructModel) makeFindFn(s string)error{
    return nil
}

func (sm *StructModel) makeComment(s string)error{
    return nil
}