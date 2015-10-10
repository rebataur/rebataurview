package cmds

import (
	"log"
	"os/exec"
	"strings"
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

		// Start the DB
		cmdimpl.StartPG(config.Database.DBPath)
		cmdimpl.SetupDB()
		// Start the Node Window
		//startNW(config.NW.NWPath)
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
func GetConfig() (utils.Config, error) {
	return utils.ReadConfig()

}
func startNW(path string) {
	go func() {
		fullPath := strings.Join([]string{path, "nw.exe"}, "")
		appPath := strings.Join([]string{path, "\\rebapp"}, "")
		out, err := exec.Command(fullPath, appPath).Output()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(out)

	}()
}
func LoadDataIntoPG(filePath string, directCopy bool) {
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
	}else{

	}

	// Create table metadata
	cmdimpl.CreateTableMetaData(tableName)

}
func describeTable() []byte {
	return cmdimpl.GetMetaData("table_schema", "")
}
func describeColumns(tableName string) []byte {
	return cmdimpl.GetMetaData("columns_schema", tableName)
}
func doAnalytics(clmn, tbln, fn string, limit int, args []string) ([]byte, error) {
	return cmdimpl.DoAnalytics(clmn, tbln, fn, limit, args)
}
