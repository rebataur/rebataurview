package cmds

import (
	"log"
	"testing"
	// "strings"
	// "fmt"
)

import ()

func TestCSVMetaParser(t *testing.T) {

	LoadDataIntoPG("D:\\uploads\\minicc.csv", true)
	LoadDataIntoPG("D:\\uploads\\Consumer_Complaints.csv", true)
	LoadDataIntoPG("D:\\uploads\\Major_Contract_Awards.csv", true)
	LoadDataIntoPG("D:\\uploads\\banklist.csv", true)
	LoadDataIntoPG("D:\\uploads\\postscndryunivsrvy2013dirinfo.csv", true)
	LoadDataIntoPG("D:\\uploads\\investment.csv", true)
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
