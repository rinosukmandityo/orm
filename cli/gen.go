package main

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/eaciit/toolkit"
	//"io/ioutil"
)

type PackageModel struct {
	Name    string
	Structs []StructModel
}

type StructModel struct {
	PkgName   string
	Name      string
	TableName string
	Fields    []FieldModel
	Methods   []MethodModel
}

type FieldModel struct {
	Name         string
	FieldType    string
	FieldDefault string
	FieldTag     string
}

const (
	MethodGet  string = "get"
	MethodFind        = "find"
)

type MethodModel struct {
	MethodType string
	Field      string
}

func (sm *StructModel) Write(path string) error {
	if sm.PkgName == "" || sm.Name == "" {
		return toolkit.Errorf("Both package name and struct name should be defined")
	}
	//return toolkit.Errorf("Fail to write %s.%s : method Write is not yet implemented", sm.PkgName, sm.Name)

	filename := filepath.Join(path, strings.ToLower(sm.Name)+".go")
	//currentCode := ""
	f, e := os.Open(filename)
	if e == nil {
		//bcurrent, _ := ioutil.ReadAll(f)
		//currentCode = string(bcurrent)
		os.Remove(filename)
	}

	f, e = os.Create(filename)
	if e != nil {
		return toolkit.Errorf("Faild to write %s.%s: %s", sm.PkgName, sm.Name, e.Error())
	}
	defer f.Close()

	txts := []string{}
	txts = append(txts, "package "+sm.PkgName)
	txts = append(txts, "/*** OrmGen Auto Generate Code - Start ***/")
	txts = append(txts, "type "+sm.Name+" struct{")
	txts = append(txts, "}")
	txts = append(txts, "/*** OrmGen Auto Generate Code - End ***/")

	b := bufio.NewWriter(f)
	for _, txt := range txts {
		b.WriteString(txt + "\n")
	}

	e = b.Flush()
	if e != nil {
		return toolkit.Errorf("Faild to write %s.%s: %s", sm.PkgName, sm.Name, e.Error())
	}

	toolkit.RunCommand("/bin/sh", "-c", "gofmt -w "+filename)
	return nil
}

func (sm *StructModel) setPkgName(txt string) error {
	if !strings.HasPrefix(txt, "package") {
		return errors.New("No package has been defined")
	}
	packages := strings.Split(txt, " ")
	if len(packages) < 2 {
		return errors.New("No package has been defined")
	}
	sm.PkgName = packages[1]
	return nil
}

func (sm *StructModel) setStructName(s string) error {
	txts := strings.Split(s, " ")
	hasStruct := false
	for _, txt := range txts {
		if !hasStruct && txt == "struct" {
			hasStruct = true
		} else if hasStruct && txt != "" {
			sm.Name = txt
			return nil
		}
	}
	return errors.New("No valid structname found: " + s)
}

func (sm *StructModel) setTableName(s string) error {
	txts := strings.Split(s, ":")
	if len(txts) < 2 {
		return errors.New("Invalid parameter - " + s)
	}
	tables := strings.Split(txts[1], " ")
	sm.TableName = strings.Trim(tables[0], " ")
	return nil
}

func (sm *StructModel) addField(s string) error {
	return nil
}

func (sm *StructModel) makeGetFn(s string) error {
	return nil
}

func (sm *StructModel) makeFindFn(s string) error {
	return nil
}

func (sm *StructModel) makeComment(s string) error {
	return nil
}
