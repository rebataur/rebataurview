package main

import (
	"testing"
	"log"
	// "strings"
)

import ()

func TestCSVMetaParser(t *testing.T) {
	loadDataIntoPG("D:\\uploads\\minicc.csv", true)
	loadDataIntoPG("D:\\uploads\\Consumer_Complaints.csv", true)
}

func TestGetTableMetaData(t *testing.T){
	// log.Println(string(getTableMetadata("cc")))
	// log.Println(string(getColumnMetadata("cc")))
	res,err := getColumnFrequency("zip_code","minicc","10")
	if err == nil{
		log.Println("here are the result :" ,string(res))
	}else{
		log.Println("There was an error",err)
	}

}
