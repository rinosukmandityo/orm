#ormgen
ormgen is a cli command to automatically generate GO files for model based on a text file

##usage
```
ormgen -file=pathtofile -out=outputfolder
```

pathtofile will define path of file to be read. If no parameter is being defined, by default it will find for default.orm on working directory

ouputfolder will define location of generated GO file. If no parameter is being defined, by default it will be working directory

File created will have name convention into xxx.go where xxx is struct name converted into lower case

##sample of default.orm
```
/* This is a remark */
/* C:Commented remark - any commented remak started with C: will be copied over to generated code and eliminiate its C: part*/
struct Employee     /*Create employee.go*/
TableName:employeTables /*Tablename on orm will be employeTables ... if no tablename define default is employees (plural name of struct in lower case) */
ID:string
Title:string
Enable:bool:default_true    /*Field enable, type is bool, default value when New is true*/
GetByID()          /*Will generate GetByID(id string)*Employee */
FindByTitle()       /*Will generate FindByTitle(title string)[]*Employee */

/*C:Department this is a commented remark and should be copied over to code for documentation purpose*/
struct Department
ID:string
Title:string
Enable:bool:defaut_true
OwnerID:string:reference_Employee /*Field EmployeeID, is a reference to Employee. ormgen should automatically created func (d *Department) Owner()*Employee */
```

Generated file should overwrite existing file and should reserve any changes that has been made outside any definition create within .orm file 