package main

import (
	"log"
	"os/exec"
	"strings"
	"strconv"
)

import (
	cmdimpl "github.com/ranjanprj/rebataurview/cmdimpl"
	utils "github.com/ranjanprj/rebataurview/utilities"
)
var initialized bool
var config utils.Config

func init() {
	if initialized == false {
		var err error
		config, err := utils.ReadConfig()
		if err != nil {
			log.Fatal("Could not read config.json, check whether it's there")
		}
		log.Println(config.Database.DBPath)

		// Start the DB
		cmdimpl.StartPG(config.Database.DBPath)
		// Start the Node Window
		startNW(config.NW.NWPath)
		initialized = true
	}

}
func cleanup() {
	if initialized {
		// stop the DB
		cmdimpl.StopPG(config.Database.DBPath)
		initialized = false
	}

}
func getConfig() (utils.Config,error){
	return utils.ReadConfig()

}
func startNW(path string){
	go func(){
		fullPath := strings.Join([]string{path,"nw.exe"},"")
		appPath  := strings.Join([]string{path,"\\rebapp"},"")
		out, err := exec.Command(fullPath,appPath).Output()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(out)

		}();
}
func loadDataIntoPG(filePath string, directCopy bool) {
	// Read the CSV metadata
	tableName, cols, colsType, err := cmdimpl.FindColumnNameAndType(filePath)
	if err != nil {
		log.Fatal("Could not get Column name and Type from csv")
	}

	log.Println(tableName, cols, colsType, err)

	// Create the table
	cmdimpl.PGCreateTable(tableName, colsType)

	// If direct copy then issue the command
	if directCopy {
		cmdimpl.PGCopyCmd(tableName, filePath)
	}


}
func describeTable(tableName string) []byte{
	return cmdimpl.GetMetaData("table_schema",tableName)
}
func describeColumns(tableName string) []byte{
	return cmdimpl.GetMetaData("columns_schema",tableName)
}
func getColumnFrequency(columnName string, tableName string, limit string) ([]byte,error){
	l,err := strconv.Atoi(limit)
	return cmdimpl.GetFrequencyCount(columnName, tableName, l),err
}
