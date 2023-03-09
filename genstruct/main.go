package main

import (
	"e.coding.net/dcoder/micro-shop/server/app/pkg/database/genstruct/util"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"text/template"

	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
)

func genModelFile(db *sqlx.DB, render *template.Template, packageName, tableName string) {
	tableSchema := &[]util.TABLE_SCHEMA{}
	err := db.Select(tableSchema,
		"SELECT COLUMN_NAME, DATA_TYPE,COLUMN_KEY,COLUMN_COMMENT,IS_NULLABLE,COLUMN_TYPE,EXTRA from COLUMNS where "+
			"TABLE_NAME"+"='"+tableName+"' and "+"table_schema = '"+*dbName+"'")

	if err != nil {
		fmt.Println(err)
		return
	}

	if len(*tableSchema) <= 0 {
		fmt.Println(tableName, "tableSchema is null")
		return
	}

	//如果表有前缀，替换掉
	exportTableName := tableName
	if *dbTablePre != "" {
		exportTableName = strings.Replace(exportTableName, *dbTablePre, "", 1)
	}

	fileName := *modelFolder + util.SecondCharUpper(exportTableName) + ".go"

	os.Remove(fileName)
	f, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	model := &util.ModelInfo{
		PackageName:  packageName,
		BDName:       *dbName,
		DBConnection: *dbConnection,
		TableName:    tableName,
		ModelName:    exportTableName,
		TableSchema:  tableSchema,
	}

	if err := render.Execute(f, model); err != nil {
		log.Fatal(err)
	}
	fmt.Println(fileName)
	cmd := exec.Command("goimports", "-w", fileName)
	cmd.Run()
}

var tplFile = flag.String("tplFile", "./model.tpl", "the path of tpl file")
var modelFolder = flag.String("modelFolder", "../model/", "the path for folder of model files")
var genTable = flag.String("genTable", "", "the name of table to be generated")
var dbInstanceName = flag.String("dbInstanceName", "dbhelper.DB", "the name of db instance used in model files")
var dbConnection = flag.String("dbConnection", "", "the name of db connection instance used in model files")
var packageName = flag.String("packageName", "", "packageName")
var dbIP = flag.String("dbIP", "127.0.0.1", "the ip of db host")
var dbPort = flag.Int("dbPort", 3306, "the port of db host")
var dbName = flag.String("dbName", "dbnote", "the name of db")
var dbTablePre = flag.String("dbTablePre", "", "the name of dbpre")
var userName = flag.String("userName", "root", "the user name of db")
var pwd = flag.String("pwd", "123456", "the password of db")

func main() {

	flag.Parse()

	logDir, _ := filepath.Abs(*modelFolder)
	if _, err := os.Stat(logDir); err != nil {
		os.Mkdir(logDir, os.ModePerm)
	}

	data, err := ioutil.ReadFile(*tplFile)
	if nil != err {
		fmt.Println("read tplFile err:", err)
		return
	}

	render := template.Must(template.New("model").
		Funcs(template.FuncMap{
			"FirstCharUpper":       util.FirstCharUpper,
			"TypeConvert":          util.TypeConvert,
			"Tags":                 util.Tags,
			"ExportColumn":         util.ExportColumn,
			"Join":                 util.Join,
			"MakeQuestionMarkList": util.MakeQuestionMarkList,
			"ColumnAndType":        util.ColumnAndType,
			"ColumnWithPostfix":    util.ColumnWithPostfix,
		}).
		Parse(string(data)))

	var tablaNames []string
	sysDB := util.GetDB(*dbIP, *dbPort, "information_schema", *userName, *pwd)

	if len(*genTable) > 0 {
		tablaNames = strings.Split(*genTable, "#")
	} else {
		err = sysDB.Select(&tablaNames,
			"SELECT table_name from tables where table_schema = '"+*dbName+"'")
		if err != nil {
			fmt.Println(err)
		}
	}

	for _, table := range tablaNames {
		genModelFile(sysDB, render, *packageName, table)
	}

	sysDB.Close()
}
