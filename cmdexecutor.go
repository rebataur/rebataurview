package main

import (
	"log"
	"strings"
)

import (
	cmdimpl "github.com/ranjanprj/rebataurview/cmdimpl"
	utils "github.com/ranjanprj/rebataurview/utilities"
)
var initialized bool
var config utils.Config

func initialize(){
  if initialized == false{
    var err error
    config, err := utils.ReadConfig()
    if err != nil {
      log.Fatal("Could not read config.json, check whether it's there")
    }
    log.Println(config.Database.DBPath)

    // Start the DB
    cmdimpl.StartPG(config.Database.DBPath)
    initialized = true
  }

}
func cleanup(){
  if initialized{
    var err error
    config, err := utils.ReadConfig()
    if err != nil {
      log.Fatal("Could not read config.json, check whether it's there")
    }
    log.Println(config.Database.DBPath)

    // Start the DB
    cmdimpl.StartPG(config.Database.DBPath)
    initialized = true
  }

}
func LoadDataIntoPG(filePath string, directCopy bool) {
  initialize()
	// Read the CSV metadata
	repName, cols, colsType, err := cmdimpl.FindColumnNameAndType(filePath)
	if err != nil {
		log.Fatal("Could not get Column name and Type from csv")
	}

	log.Println(repName, cols, colsType, err)

	// Create the table
	cmdimpl.PGCreateTable(repName,colsType)

	// If direct copy then issue the command
	if directCopy {
			cmdimpl.PGCopyCmd(repName,filePath)
	}

	log.Println( strings.Split(string(cmdimpl.GetFrequencyCount("state",repName,4)),"},"))
}
