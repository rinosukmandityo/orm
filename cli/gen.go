package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/eaciit/toolkit"
	//"io/ioutil"
)

const (
	mandatoryLibs string = "github.com/eaciit/dbox,github.com/eaciit/orm"
)

type PackageModel struct {
	Name       string
	ObjectLibs string
	Structs    []StructModel
}

type StructModel struct {
	Name      string
	TableName string
	Libs      string
	Fields    []*FieldModel
	Methods   []*MethodModel
}

type FieldModel struct {
	Name    string
	Type    string
	Default interface{}
	Tag     string
}

const (
	MethodGet  string = "get"
	MethodFind        = "find"
)

type MethodModel struct {
	Type  string
	Field string
}

func libs(ls ...string) string {
	ret := map[string]int{}
	for _, s := range ls {
		libraries := strings.Split(s, ",")
		for _, l := range libraries {
			l = strings.Trim(l, " ")
			if l != "" {
				_, exist := ret[l]
				if !exist {
					ret[l] = 1
				}
			}
		}
	}
	txt := ""
	for k, _ := range ret {
		txt += "\"" + k + "\"\n"
	}
	return txt
}

func (pkg *PackageModel) WriteBase(path string) error {
	filename := filepath.Join(path, "base.go")
	f, e := os.Open(filename)
	if e == nil {
		os.Remove(filename)
	}

	f, e = os.Create(filename)
	if e != nil {
		return toolkit.Errorf("Failed to write %s: %s", "base.go", e.Error())
	}
	defer f.Close()

	b := bufio.NewWriter(f)
	b.WriteString(toolkit.Formatf(baseGo, pkg.Name))
	e = b.Flush()
	if e != nil {
		return toolkit.Errorf("Failed to write base.go: %s", e.Error())
	}

	toolkit.RunCommand("/bin/sh", "-c", "gofmt -w "+filename)
	return nil
}

func (sm *StructModel) Write(pkg *PackageModel, path string) error {
	if pkg.Name == "" || sm.Name == "" {
		return toolkit.Errorf("Both package name and struct name should be defined")
	}
	//return toolkit.Errorf("Fail to write %s.%s : method Write is not yet implemented", pkg.Name, sm.Name)

	//-- write base
	e := pkg.WriteBase(path)
	if e != nil {
		return e
	}

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
		return toolkit.Errorf("Failed to write %s.%s: %s", pkg.Name, sm.Name, e.Error())
	}
	defer f.Close()

	txts := []string{}
	//--- package
	txts = append(txts, "package "+pkg.Name)

	//--- imports
	txts = append(txts,
		toolkit.Sprintf("import (%s)",
			libs(mandatoryLibs, pkg.ObjectLibs, sm.Libs)))

	//--- struct definition
	txts = append(txts, "type "+sm.Name+" struct {\n"+
		"orm.ModelBase `bson:\"-\" json:\"-\"`")
	for _, f := range sm.Fields {
		if f.Type == "" {
			f.Type = "string"
		}
		fieldStr := toolkit.Sprintf("%s %s %s", f.Name, f.Type, f.Tag)
		txts = append(txts, fieldStr)
	}
	txts = append(txts, "}")

	//--- tablename
	pluralNames := strings.ToLower(sm.Name)
	if strings.HasSuffix(pluralNames, "s") {
		pluralNames = pluralNames + "es"
	} else {
		pluralNames = pluralNames + "s"
	}
	tablename := toolkit.Sprintf("func (o *%s) TableName()string{"+
		"return \"%s\"\n"+
		"}", sm.Name, pluralNames)
	txts = append(txts, tablename)

	//--- new
	fieldBuilders := ""
	for _, field := range sm.Fields {
		notEmpty := !toolkit.IsNilOrEmpty(field.Default)
		if notEmpty {
			def := toolkit.Sprintf("%v", field.Default)
			if field.Type == "string" {
				def = "\"" + def + "\""
			}
			fieldBuilders +=
				toolkit.Sprintf("o.%s=%s", field.Name, def) +
					"\n"
		}
	}
	newfn := "func New{0}() *{0}{\n" +
		"o:=new({0})\n" +
		fieldBuilders +
		"return o" +
		"}"
	newfn = toolkit.Formatf(newfn, sm.Name)
	txts = append(txts, newfn)

	//--- find
	tpl := `func {0}Find(filter *dbox.Filter, fields, orders string, limit, skip int) dbox.ICursor {
            config := makeFindConfig(fields, orders, skip, limit)
            if filter != nil {
                config.Set("where", filter)
            }
		    c, _ := DB().Find(new({0}), config)
            return c
        }`
	txts = append(txts, toolkit.Formatf(tpl, sm.Name))

	//--- get
	tpl = `func {0}Get(filter *dbox.Filter, orders string, skip int) (emp *{0}, err error) {
        config := makeFindConfig("", orders, skip, 1)
        if filter != nil {
            config.Set("where", filter)
        }
        c, ecursor := DB().Find(new({0}), config)
        if ecursor != nil {
            return nil, ecursor
        }
        defer c.Close()

        emp = new({0})
        err = c.Fetch(emp, 1, false)
        return emp, err
    }`
	txts = append(txts, toolkit.Formatf(tpl, sm.Name))

	//-- method & get
	for _, method := range sm.Methods {
		txts = append(txts, sm.buildMethod(
			pkg,
			method.Type,
			method.Field))
	}

	b := bufio.NewWriter(f)
	for _, txt := range txts {
		b.WriteString(txt + "\n")
	}

	e = b.Flush()
	if e != nil {
		return toolkit.Errorf("Failed to write %s.%s: %s", pkg.Name, sm.Name, e.Error())
	}

	toolkit.RunCommand("/bin/sh", "-c", "gofmt -w "+filename)
	return nil
}

func (sm *StructModel) buildMethod(
	pkg *PackageModel, methodType string, fields string) string {
	fieldIds := strings.Split(fields, ",")
	fieldNameConcat := ""
	filter := ""
	filtersEq := []string{}
	parmNames := []string{}
	for _, fieldId := range fieldIds {
		fieldId = strings.Trim(fieldId, " ")
		field := sm.Field(fieldId)
		if field != nil {
			fieldNameConcat += field.Name
			parmNames = append(parmNames,
				toolkit.Sprintf("p%s %s", field.Name, field.Type))
			fieldNameFn := strings.ToLower(field.Name)
			if fieldNameFn == "id" {
				fieldNameFn = "_id"
			}
			filtersEq = append(filtersEq,
				toolkit.Formatf(`dbox.Eq("{0}",{1})`,
					fieldNameFn,
					"p"+field.Name))
		}
	}
	if len(filtersEq) == 1 {
		filter = filtersEq[0]
	} else if len(filtersEq) > 1 {
		filter = "dbox.And(" + strings.Join(filtersEq, ",") + ")"
	}

	var tpl string
	if methodType == MethodFind {
		tpl = `
            func {0}FindBy{1}({2},fields string,limit,skip int) dbox.ICursor{
			return {0}Find({3},"","",0,0)
			}`
	} else {
		tpl = `
            func {0}GetBy{1}({2},orders string)(*{0},error){
			    return {0}Get({3},"",0)
			}
            `
	}
	return toolkit.Formatf(tpl,
		sm.Name,                      //--0
		fieldNameConcat,              //--1
		strings.Join(parmNames, ","), //--2
		filter)
}

func (sm *StructModel) Field(fn string) *FieldModel {
	smallfn := strings.ToLower(fn)
	for _, f := range sm.Fields {
		if strings.ToLower(f.Name) == smallfn {
			return f
		}
	}
	return nil
}

func (sm *StructModel) Method(fn string, fields string) *MethodModel {
	smallfn := strings.ToLower(fn)
	smallfields := strings.ToLower(fields)
	for _, f := range sm.Methods {
		if strings.ToLower(f.Type) == smallfn && strings.ToLower(f.Field) == smallfields {
			return f
		}
	}
	return nil
}
