package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var baseGo string = `
import(
    "github.com/eaciit/dbox"
    "github.com/eaciit/orm"
)

var _db *orm.DataContext

func SetDb(conn dbox.IConnection)error{
    CloseDb()
    _db = orm.New(conn)
    return nil
}

func CloseDb(){
    if _db!=nil{
        _db.Close()    
    }
}

func DB() *orm.DataContext{
    return _db
}`

type ExistingFunctionList struct {
	Name  string
	Lines []string
}
type ExistingVars struct {
	Name  string
	Lines []string
}

var fieldImports = make(map[string]string)

type StrucMap struct {
	PackageName string
	StructName  string
	TableName   string
	Imports     []ImportStructure
	Fields      []FieldStructure
	Functions   []FunctionStructure
}
type FieldStructure struct {
	FieldName       string
	FieldType       string
	IsDefault       bool
	DefaultValue    string
	IsBson          bool
	BsonName        string
	IsJson          bool
	JsonName        string
	IsReference     bool
	ReferenceStruct string
	IsPK            bool
}

type FunctionStructure struct {
	Name       string
	ParamName  string
	ParamType  string
	FieldName  string
	ReturnType string
}
type ImportStructure struct {
	ImportType string
	ImportUrl  string
}

var wg sync.WaitGroup

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d
	}()
)
var sourceFile, outPath string
var structMap []StrucMap

func init() {
	fieldImports["time"] = "time"
	fieldImports["bson"] = "gopkg.in/mgo.v2/bson"
}

func main() {
	//	log.Println("Current Dir", wd)
	if len(os.Args) == 1 {
		fmt.Println("Usage ormgen -file=inputfile -out=folderoutpath")
		os.Exit(2)
	} else if len(os.Args) == 2 {
		sourceFile = os.Args[1]
		//		sourceFile = strings.Split(sourceFile, "=")[1]
		outPath = wd + string(os.PathSeparator) + "gen" + string(os.PathSeparator)
	} else if len(os.Args) == 3 {
		sourceFile = os.Args[1]
		outPath = os.Args[2]
	}
	sourceFile = strings.Split(sourceFile, "=")[1]
	if _, err := os.Stat(sourceFile); err != nil {
		log.Println(sourceFile, "Not found; combine with working directory ")
		sourceFile = wd + string(os.PathSeparator) + sourceFile
	}
	log.Println(" INPUT_FILE => ", sourceFile, "; OUTPATH => ", outPath)
	inputLines, err := readLines(sourceFile)
	if err != nil {
		log.Println("Error reading source ORM file")
		os.Exit(2)
	}
	var pkgName, structName, tableName string
	var isStruct, closeStruct bool = false, false
	var fnCount, structCount, fieldCount int = -1, -1, 0
	for i, line := range inputLines {
		line := strings.TrimSpace(line)
		if strings.HasPrefix(line, "package") {
			pkgName = strings.TrimSpace(strings.Split(line, " ")[1])
		}
		if strings.HasPrefix(line, "struct") {
			if strings.Index(line, "/*") > -1 {
				line = strings.TrimSpace(line[0:strings.Index(line, "/*")])
			}
			if len(strings.Split(line, " ")) != 2 {
				log.Println("Invalid struct defined at => ", sourceFile, " line ", (i + 1), " ", line)
				os.Exit(2)
			} else {
				structName = strings.TrimSpace(strings.Split(line, " ")[1])
				if strings.HasSuffix(structName, "s") {
					tableName = strings.ToLower(structName + "es")
				} else {
					tableName = strings.ToLower(structName + "s")
				}
			}
			structCount++
			isStruct = true
			closeStruct = false
			structMap = append(structMap, StrucMap{})
			structMap[structCount] = StrucMap{}
			structMap[structCount].StructName = structName
			structMap[structCount].TableName = tableName
			structMap[structCount].PackageName = pkgName
			structMap[structCount].Imports = append(structMap[structCount].Imports, ImportStructure{".", "github.com/eaciit/orm"})
		} else if isStruct && strings.TrimSpace(line) != "" {
			if strings.Index(line, "/*") > -1 {
				line = strings.TrimSpace(line[0:strings.Index(line, "/*")])
			}
			if strings.HasPrefix(line, "TableName") {
				if len(strings.Split(line, ":")) != 2 {
					log.Println("Invalid TableName defined at => ", sourceFile, " line ", (i + 1), " ", line)
					os.Exit(2)
				} else {
					tableName = strings.TrimSpace(strings.Split(line, ":")[1])
				}
				structMap[structCount].TableName = tableName
			} else if strings.Index(line, "()") > -1 {
				//Function Here
				if strings.Index(line, "GetBy") > -1 || strings.Index(line, "FindBy") > -1 {
					fnCount++
					structMap[structCount].Functions = append(structMap[structCount].Functions, FunctionStructure{})
					fnName := strings.Replace(strings.Replace(line, ")", "", -1), "(", "", -1)
					if strings.Index(line, "GetBy") > -1 {
						//						log.Println("Detected FnName => " + fnName + "; fieldName => " + fnName[strings.Index(fnName, "GetBy")+5:])
						fnParam := strings.TrimSpace(fnName[strings.Index(fnName, "GetBy")+5:])
						if ok, fnFieldName, fnFieldType := fieldInSlice(fnParam, structMap[structCount].Fields); ok {
							structMap[structCount].Functions[fnCount] = FunctionStructure{}
							structMap[structCount].Functions[fnCount].Name = structName + fnName
							structMap[structCount].Functions[fnCount].ParamName = strings.ToLower(fnParam)
							structMap[structCount].Functions[fnCount].FieldName = fnFieldName
							structMap[structCount].Functions[fnCount].ParamType = fnFieldType
							structMap[structCount].Functions[fnCount].ReturnType = "*" + structName
						}
					} else if strings.Index(line, "FindBy") > -1 {
						//						log.Println("Detected FnName => " + fnName + "; fieldName => " + fnName[strings.Index(fnName, "FindBy")+6:])
						fnParam := strings.TrimSpace(fnName[strings.Index(fnName, "FindBy")+6:])
						if ok, fnFieldName, fnFieldType := fieldInSlice(fnParam, structMap[structCount].Fields); ok {
							structMap[structCount].Functions[fnCount] = FunctionStructure{}
							structMap[structCount].Functions[fnCount].Name = structName + fnName
							structMap[structCount].Functions[fnCount].ParamName = strings.ToLower(fnParam)
							structMap[structCount].Functions[fnCount].FieldName = fnFieldName
							structMap[structCount].Functions[fnCount].ParamType = fnFieldType
							structMap[structCount].Functions[fnCount].ReturnType = "dbox.ICursor"
							if !impInSlice("github.com/eaciit/toolkit", structMap[structCount].Imports) {
								structMap[structCount].Imports = append(structMap[structCount].Imports, ImportStructure{"", "github.com/eaciit/toolkit"})
							}
							if !impInSlice("github.com/eaciit/dbox", structMap[structCount].Imports) {
								structMap[structCount].Imports = append(structMap[structCount].Imports, ImportStructure{"", "github.com/eaciit/dbox"})
							}
						}
					}
				} else {
					//Not Supported Yet
				}
			} else {
				fields := strings.Split(line, ":")
				var fieldStru FieldStructure
				for _, field := range fields {
					fieldStru.FieldName = fields[0]
					fieldStru.FieldType = fields[1]
					if strings.Index(field, "default_") > -1 {
						fieldStru.IsDefault = true
						//						fieldStru.FieldName = fields[0]
						//						fieldStru.FieldType = fields[1]
						fieldStru.DefaultValue = strings.TrimSpace(strings.Split(field, "_")[1])
					} else if strings.Index(field, "reference_") > -1 {
						//						fieldStru.FieldName = fields[0]
						//						fieldStru.FieldType = fields[1]
						fieldStru.IsReference = true
						fieldStru.ReferenceStruct = strings.TrimSpace(strings.Split(field, "_")[1])
					} else if strings.Index(field, "primary_key") > -1 {
						//						fieldStru.FieldName = fields[0]
						//						fieldStru.FieldType = fields[1]
						fieldStru.IsPK = true
					} else if strings.Index(field, "bson`") > -1 {
						fieldStru.IsBson = true
						fieldStru.BsonName = strings.TrimSpace(strings.Split(field, "`")[1])
					} else if strings.Index(field, "json`") > -1 {
						fieldStru.IsJson = true
						fieldStru.JsonName = strings.TrimSpace(strings.Split(field, "`")[1])
					} else {
					}
				}
				structMap[structCount].Fields = append(structMap[structCount].Fields, FieldStructure{})
				structMap[structCount].Fields[fieldCount] = fieldStru
				fieldCount++
			}
		} else if isStruct && !closeStruct {
			isStruct = false
			closeStruct = true
			fnCount = 0
			fieldCount = 0
		}
	}
	//	log.Println("++++++++++++++++++++++++++++++++++++++++++++++")
	//	log.Println(" ")
	//	log.Println("")
	for _, maps := range structMap {
		//		fmt.Println("Package : ", maps.PackageName)
		//		fmt.Println("Structure : ", maps.StructName)
		//		fmt.Println("TableName : ", maps.TableName)
		//		for c, fields := range maps.Fields {
		//			fmt.Printf("%d %+v\n", c, fields)
		//		}
		for c, functions := range maps.Functions {
			//			fmt.Printf("%d %+v\n", c, functions)
			pField := strings.Replace(functions.Name, "GetBy", "", -1)
			pField = strings.Replace(pField, "FindBy", "", -1)
			//			fmt.Println("Param Field=> ", pField)
			for _, fields := range maps.Fields {
				if fields.FieldName == pField {
					functions.ParamName = strings.ToLower(fields.FieldName)
					functions.ParamType = fields.FieldType
					//					fmt.Printf("%+v \n", functions)
				}
			}
			maps.Functions[c] = functions
		}
		//		for c, functions := range maps.Functions {
		//			fmt.Printf("%d %+v\n", c, functions)
		//		}
		//		fmt.Println("--------            -------------")
		//		fmt.Println("-------- END STRUCT -------------")
		//		fmt.Println("--------            -------------")
	}
	baseFileName := outPath + "base.go"
	writeBaseFile(pkgName, outPath, baseFileName)
	for _, stMap := range structMap {
		var exFnList, keepFnList []ExistingFunctionList
		var exVarList []ExistingVars
		fileName := outPath + strings.ToLower(stMap.StructName) + ".go"
		log.Println("OUTPUT FILE => ", fileName)
		_, err := os.Stat(outPath)
		if err != nil {
			err = os.MkdirAll(outPath, 0644)
			checkError(err)
		}
		_, err = os.Stat(fileName)
		if err != nil {
			file, err := os.Create(fileName)
			checkError(err)
			defer file.Close()
		} else {
			exFnList, exVarList = readExistingSource(fileName)
			err := os.Remove(fileName)
			checkError(err)
			file, err := os.Create(fileName)
			checkError(err)
			defer file.Close()
		}
		fileOut, err := os.OpenFile(fileName, os.O_RDWR, 0644)
		checkError(err)
		defer fileOut.Close()
		_, err = fileOut.WriteString("package " + stMap.PackageName + "\n")
		_, err = fileOut.WriteString("import (\n")
		for _, imp := range stMap.Imports {
			_, err = fileOut.WriteString(imp.ImportType + " \"" + imp.ImportUrl + "\"\n")
		}
		_, err = fileOut.WriteString(" )\n")
		_, err = fileOut.WriteString("type " + stMap.StructName + " struct {\n")
		_, err = fileOut.WriteString("ModelBase `bson:\"-\",json:\"-\"`\n")
		for _, fields := range stMap.Fields {
			var bsonStr, jsonStr = "", ""
			if fields.IsBson {
				bsonStr = " bson:\"" + fields.BsonName + "\""
			}
			if fields.IsJson {
				jsonStr = "json:\"" + fields.JsonName + "\""
			}
			if len(bsonStr) > 0 && len(jsonStr) > 0 {
				bsonStr = " `" + bsonStr + " , " + jsonStr + " `"
			} else if fields.IsBson || fields.IsJson {
				bsonStr = " `" + bsonStr + jsonStr + " `"
			}
			_, err = fileOut.WriteString(fields.FieldName + " " + fields.FieldType + bsonStr + "\n")
		}
		_, err = fileOut.WriteString("}\n\n")
		for _, exVar := range exVarList {
			for _, line := range exVar.Lines {
				_, err = fileOut.WriteString(line + "\n")
			}
		}
		for _, functions := range stMap.Functions {
			//			fmt.Printf("%+v\n", functions)
			if strings.Index(functions.Name, "GetBy") > 0 {
				_, err = fileOut.WriteString("func " + functions.Name + "(" + functions.ParamName + " " + functions.ParamType + ") " + functions.ReturnType + " {\n")
				_, err = fileOut.WriteString(strings.ToLower(stMap.StructName) + " := new(" + stMap.StructName + ")\n")
				_, err = fileOut.WriteString("DB().GetById(" + strings.ToLower(stMap.StructName) + ", " + functions.ParamName + ")\n")
				_, err = fileOut.WriteString("return " + strings.ToLower(stMap.StructName) + "\n")
				_, err = fileOut.WriteString("}\n")

			} else if strings.Index(functions.Name, "FindBy") > 0 {
				_, err = fileOut.WriteString("func " + functions.Name + "(" + functions.ParamName + " " + functions.ParamType + ", order []string, skip, limit int) " + functions.ReturnType + " {\n")
				_, err = fileOut.WriteString("c, _ := DB().Find(new(" + stMap.StructName + "),\n")
				_, err = fileOut.WriteString("toolkit.M{}.Set(\"where\", []*dbox.Filter{dbox.Eq(\"" + functions.FieldName + "\"," + functions.ParamName + ")}).\n")
				_, err = fileOut.WriteString("Set(\"order\",order).\n")
				_, err = fileOut.WriteString("Set(\"skip\",skip).\n")
				_, err = fileOut.WriteString("Set(\"limit\",limit))\n")
				_, err = fileOut.WriteString("return dbox.NewCursor(c) \n}\n\n")
			}
		}
		fnNew := "func New" + stMap.StructName + "() *" + stMap.StructName + "{\n"
		fnNew = fnNew + "e:= new(" + stMap.StructName + ") \n"
		for _, fields := range stMap.Fields {
			//			fmt.Printf("%+v \n", fields)
			if fields.IsDefault {
				switch fields.FieldType {
				case "string":
					fnNew = fnNew + " e." + fields.FieldName + "=\"" + fields.DefaultValue + "\"\n"
				default:
					fnNew = fnNew + " e." + fields.FieldName + "=" + fields.DefaultValue + "\n"
				}
				stMap.Functions = append(stMap.Functions, FunctionStructure{"New" + stMap.StructName, "", "", "", ""})
			} else if fields.IsPK {
				fnRecId := "func (e *" + stMap.StructName + " ) RecordID() interface{} {\n"
				fnRecId = fnRecId + " return e." + fields.FieldName + " \n }\n\n"
				_, err = fileOut.WriteString(fnRecId)
				stMap.Functions = append(stMap.Functions, FunctionStructure{"RecordID", "", "", "", ""})
			}
		}
		fnNew = fnNew + "return e\n }\n\n"
		_, err = fileOut.WriteString(fnNew)
		fnTblName := "func (e *" + stMap.StructName + ") TableName() string {\n"
		fnTblName = fnTblName + "return \"" + stMap.TableName + "\" \n }\n\n"
		_, err = fileOut.WriteString(fnTblName)
		stMap.Functions = append(stMap.Functions, FunctionStructure{"TableName", "", "", "", ""})

		for _, exFn := range exFnList {
			//			log.Println("Check FN[" + exFn.Name + " ]")
			if !funcInSlice(exFn.Name, stMap.Functions) {
				//				log.Println("Check FN[" + exFn.Name + " ] NOT EXIST(), give to Keeper ")
				keepFnList = append(keepFnList, exFn)
			}
		}
		//Keeper of the lights
		for _, exFn := range keepFnList {
			//			log.Println("Function Keep => ", exFn.Name)
			for _, exFnLine := range exFn.Lines {
				_, err := fileOut.WriteString(exFnLine + "\n")
				checkError(err)
			}
		}

		// save changes
		err = fileOut.Sync()
		checkError(err)
	}
	//	fmt.Println(runtime.GOOS)
	switch runtime.GOOS {
	case "windows":
		err = exec.Command("cmd", "/c", "gofmt", "-w", outPath).Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
			os.Exit(1)
		}
	case "linux":
		cmd := exec.Command("/bin/sh", "-c", "gofmt -w "+outPath)
		err := cmd.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
			os.Exit(1)
		}
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

func impInSlice(a string, list []ImportStructure) bool {
	for _, field := range list {
		if field.ImportUrl == a {
			return true
		}
	}
	return false
}
func fieldInSlice(a string, list []FieldStructure) (bool, string, string) {
	for _, field := range list {
		if field.FieldName == a {
			if field.IsBson {
				return true, field.BsonName, field.FieldType
			} else {
				return true, strings.ToLower(field.FieldName), field.FieldType
			}
		}
	}
	return false, "", ""
}

func funcInSlice(a string, list []FunctionStructure) bool {
	for _, fn := range list {
		if fn.Name == a {
			return true
		}
	}
	return false
}
func writeBaseFile(pkgName, outPath, fileName string) {
	_, err := os.Stat(outPath)
	if err != nil {
		err = os.MkdirAll(outPath, 0644)
		checkError(err)
	}
	_, err = os.Stat(fileName)
	if err != nil {
		file, err := os.Create(fileName)
		checkError(err)
		defer file.Close()
	} else {
		err := os.Remove(fileName)
		checkError(err)
		file, err := os.Create(fileName)
		checkError(err)
		defer file.Close()
	}
	fileOut, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	checkError(err)
	defer fileOut.Close()
	_, err = fileOut.WriteString("package " + pkgName)
	_, err = fileOut.WriteString(baseGo)
	err = fileOut.Sync()
	checkError(err)
}

func readExistingSource(path string) ([]ExistingFunctionList, []ExistingVars) {
	var exFnList []ExistingFunctionList
	var exVarList []ExistingVars
	file, err := os.Open(path)
	if err != nil {
		log.Printf(err.Error())
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	var fnStart /*, varStart*/ bool
	var fnName string
	var fnCount, fnVar int = -1, -1
	for _, line := range lines {
		//		log.Println("FnStart ? ", fnStart)
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "func") {
			fnStart = true
			linet := strings.Replace(line, "func ", "", -1)
			lineParts := strings.Split(linet, " ")
			for _, lp := range lineParts {
				re := regexp.MustCompile("^[A-Za-z_]")
				if re.MatchString(lp) {
					fnName = lp[0:strings.Index(lp, "(")]
					fnCount++
					exFnList = append(exFnList, ExistingFunctionList{})
					exFnList[fnCount].Name = fnName
					exFnList[fnCount].Lines = append(exFnList[fnCount].Lines, line)
					break
				}
			}
		} else if line == "}" && fnStart {
			fnStart = false
			exFnList[fnCount].Lines = append(exFnList[fnCount].Lines, "}")
			//			log.Println("end of func: fnstart ", fnStart)
		} else if fnStart {
			exFnList[fnCount].Lines = append(exFnList[fnCount].Lines, line)
		}
		if strings.HasPrefix(line, "var") && !fnStart {
			fnVar++
			exVarList = append(exVarList, ExistingVars{})
			exVarList[fnVar].Name = strconv.Itoa(fnVar)
			exVarList[fnVar].Lines = append(exVarList[fnVar].Lines, line)
			//			log.Println("Save var[", fnVar, "] => ", line)
		} /*else if line == "}" && varStart {
			varStart = false
			exVarList[fnVar].Lines = append(exVarList[fnVar].Lines, "}")
		} else if varStart {
			exVarList[fnVar].Lines = append(exVarList[fnVar].Lines, line)
		}*/

		//		log.Println("LINE => ", line, " is vars? ", strings.HasPrefix(line, "var"), "; fnStarts? ", fnStart, "; is }?", line == "}")
	}
	return exFnList, exVarList
}
