package main

import (
	"testing"
	"log"
	// "strings"
)

import ()

func TestCSVMetaParser(t *testing.T) {
	loadDataIntoPG("D:\\uploads\\cc.csv", true)
}

func TestGetTableMetaData(t *testing.T){
	log.Println(string(getTableMetadata("cc")))
	log.Println(string(getColumnMetadata("cc")))
	log.Println(string(getColumnFrequency("state","cc",10)))
}
