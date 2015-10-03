package main

import (
	"log"
	"testing"
	// "strings"
)

import ()

func TestCSVMetaParser(t *testing.T) {

	loadDataIntoPG("D:\\uploads\\minicc.csv", true)
	loadDataIntoPG("D:\\uploads\\Consumer_Complaints.csv", true)
	loadDataIntoPG("D:\\uploads\\Major_Contract_Awards.csv", true)
}

func TestGetTableMetaData(t *testing.T) {
	log.Println("	")
	// log.Println(string(getTableMetadata("cc")))
	// log.Println(string(getColumnMetadata("cc")))
	// res, err := doAnalytics("zip_code", "minicc", "10")
	// if err == nil {
	// 	log.Println("here are the result :", string(res))
	// } else {
	// 	log.Println("There was an error", err)
	// }

}
