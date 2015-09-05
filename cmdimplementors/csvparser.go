package cmdimplementors

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)
import (
	"database/sql"
	_ "github.com/lib/pq"
)

type ColumnType struct {
	name   string
	dbType string
}

// func main() {
// 	ReadRepositoryFile("rep", "D:\\uploads\\4.csv")
// 	fmt.Println("========DONE==========")
// }

func ReadRepositoryFile(fileName, filePath string) []string {
	log.Println("Reading repository file")
	repositoryName := strings.Split(fileName,".csv")[0]
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	rawCSVDat, err := reader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	var columns []string
	var dbCols []ColumnType
	var createTableSQL string

	for i, each := range rawCSVDat {

		if i == 0 {
			columns = createTabel(each)
		} else if i == 1 {
			// figure out data types and create table
			createTableSQL, dbCols = getCreateTableQuery(columns, each)
			createTable(repositoryName, createTableSQL)
		}
	}

	log.Println("inserting data");
	for i, each := range rawCSVDat {
		if i > 0 {

			insertData(repositoryName, dbCols, each)
		}

	}

	log.Println("Finished Processing data")
	return columns
}

func replaceRegexMatchesForSpecialChars(src, repl string) string {
	matches := regexp.MustCompile("[`~!@#$%^&*()_+={}|\"';:/?.>,<']")
	rawStr := matches.ReplaceAllLiteralString(src, repl)
	removeBlank := strings.Replace(rawStr, " ", "_", -1)

	removeHyphen := strings.Replace(removeBlank, "-", "_", -1)
	finalStr := strings.ToLower(removeHyphen)
	return finalStr
}

func createTabel(columns []string) []string {
	colsize := len(columns)
	col := make([]string, colsize)
	for j, e := range columns {
		c := replaceRegexMatchesForSpecialChars(e, "")
		col[j] = strings.Join([]string{c, ""}, "")
	}
	return col
}

func checkErr(err error) {
	if err != nil {
		log.Fatal("Error Occurred : ", err)
	}

}

func getCreateTableQuery(columns, colVal []string) (string, []ColumnType) {
	ddMMStr := `(((0[1-9]|[12][0-9]|3[01])([/])(0[13578]|10|12)([/])(\d{4}))|(([0][1-9]|[12][0-9]|30)([/])(0[469]|11)([/])(\d{4}))|((0[1-9]|1[0-9]|2[0-8])([/])(02)([/])(\d{4}))|((29)(\.|-|\/)(02)([/])([02468][048]00))|((29)([/])(02)([/])([13579][26]00))|((29)([/])(02)([/])([0-9][0-9][0][48]))|((29)([/])(02)([/])([0-9][0-9][2468][048]))|((29)([/])(02)([/])([0-9][0-9][13579][26])))`
	ddMM := regexp.MustCompile(ddMMStr)

	mmDDStr := `^((0?[13578]|10|12)(-|\/)(([1-9])|(0[1-9])|([12])([0-9]?)|(3[01]?))(-|\/)((19)([2-9])(\d{1})|(20)([01])(\d{1})|([8901])(\d{1}))|(0?[2469]|11)(-|\/)(([1-9])|(0[1-9])|([12])([0-9]?)|(3[0]?))(-|\/)((19)([2-9])(\d{1})|(20)([01])(\d{1})|([8901])(\d{1})))$`
	mmDD := regexp.MustCompile(mmDDStr)

	genericDate := regexp.MustCompile(`(^\d{1,2}-\w{3,9}-\d{2,4}$)|(^\d{1,2}\\w{3,9}\\d{2,4}$)|(^\d{1,2}/\w{3,9}/\d{2,4}$)`)

	integerMatchesPlusMinus := regexp.MustCompile("^(\\+|-)\\d*[\\d]$")
	integerMatches := regexp.MustCompile("^\\d*[\\d]$")

	decimalMatchesPlusMinus := regexp.MustCompile("^(\\+|-)\\d*\\.\\d*[\\d]$")
	decimalMatches := regexp.MustCompile("^\\d*\\.\\d*[\\d]$")
	pctMatches := regexp.MustCompile("(^\\d*\\.\\d*%$)|(^\\d*%$)")

	dbCols := make([]ColumnType, len(columns))

	for i, each := range colVal {

		if len(each) == 0 {
			dbCols[i].name = each
			dbCols[i].dbType = "varchar"
		} else if ddMM.MatchString(each) {
			dbCols[i].name = each
			dbCols[i].dbType = "date"
		} else if mmDD.MatchString(each) {
			dbCols[i].name = each
			dbCols[i].dbType = "date"
		} else if genericDate.MatchString(each) {
			dbCols[i].name = each
			dbCols[i].dbType = "date"
		} else if integerMatchesPlusMinus.MatchString(each) || integerMatches.MatchString(each) {
			dbCols[i].name = each
			dbCols[i].dbType = "bigint"
		} else if decimalMatchesPlusMinus.MatchString(each) || decimalMatches.MatchString(each) {
			dbCols[i].name = each
			dbCols[i].dbType = "decimal"
		} else if pctMatches.MatchString(each) {
			dbCols[i].name = each
			dbCols[i].dbType = "percentage"
		} else {
			dbCols[i].name = each
			dbCols[i].dbType = "varchar"
		}

	}

	sqlColDef := make([]string, len(columns))
	for i, each := range dbCols {
		if each.dbType == "percentage" {
			columns[i] = strings.Join([]string{columns[i], "_percentage"}, "")
			each.dbType = "decimal"
		}
		sqlColDef[i] = strings.Join([]string{columns[i], " ", each.dbType}, "")
	}

	return strings.Join(sqlColDef, ","), dbCols

}

func createTable(tableName, columns string) {
	log.Println("Creating tables")
	dropTableSQL := fmt.Sprintf("DROP TABLE IF  EXISTS %s; \n", tableName)
	createTableSQL := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ( %s ); \n ", tableName, columns)

		// file, err := os.OpenFile("D:\\uploads\\insert.sql", os.O_APPEND|os.O_WRONLY, 0600)
		// if err != nil {
		// 	panic(err)
		// }
		// defer file.Close()
		//
		// if _, err = file.WriteString(createTableSQL); err != nil {
		// 	panic(err)
		// }
  // dbCreateTable("set datestyle = DMY;")
	dbCreateTable(dropTableSQL)
	dbCreateTable(createTableSQL)


}

func insertData(tableName string, dbCols []ColumnType, dataRow []string) {

	dbData := make([]string, len(dbCols))
	for i, each := range dataRow {
		escEach := strings.Replace(each, "'", "''", -1)
		if len(each) == 0 || len(strings.TrimSpace(each)) == 0  {
			dbData[i] = "NULL"
		} else if dbCols[i].dbType == "percentage" {
			dbData[i] = strings.Replace(escEach, "%", "", -1)
		} else if dbCols[i].dbType == "varchar" || dbCols[i].dbType == "date" {

			if dbCols[i].dbType == "date" {
				dbData[i] = strings.Join([]string{"'", dateCorrectionToDDMMYY(escEach), "'"}, "")

			} else {
				dbData[i] = strings.Join([]string{"'", escEach, "'"}, "")
			}
		} else {
			dbData[i] = strings.Join([]string{escEach}, "")
		}

	}

	insertTableSQL := fmt.Sprintf("INSERT INTO %s values( %s ); \n", tableName, strings.Join(dbData, ","))

	// open files r and w

		// file, err := os.OpenFile("D:\\uploads\\insert.sql", os.O_APPEND|os.O_WRONLY, 0600)
		// if err != nil {
		// 	panic(err)
		// }
		// defer file.Close()
		//
		// if _, err = file.WriteString(insertTableSQL); err != nil {
		// 	panic(err)
		// }

	dbInsertData(insertTableSQL)

}

var ddMMStr, mmDDStr string
var mmDD, ddMM, genericDate *regexp.Regexp
var ismmDD bool = false
func dateCorrectionToDDMMYY(dateInString string) string {
	dateInString = strings.Split(dateInString," ")[0]
	if len(ddMMStr) == 0 {

		ddMMStr = `(((0[1-9]|[12][0-9]|3[01])([/])(0[13578]|10|12)([/])(\d{4}))|(([0][1-9]|[12][0-9]|30)([/])(0[469]|11)([/])(\d{4}))|((0[1-9]|1[0-9]|2[0-8])([/])(02)([/])(\d{4}))|((29)(\.|-|\/)(02)([/])([02468][048]00))|((29)([/])(02)([/])([13579][26]00))|((29)([/])(02)([/])([0-9][0-9][0][48]))|((29)([/])(02)([/])([0-9][0-9][2468][048]))|((29)([/])(02)([/])([0-9][0-9][13579][26])))`
		ddMM = regexp.MustCompile(ddMMStr)

		mmDDStr = `^((0?[13578]|10|12)(-|\/)(([1-9])|(0[1-9])|([12])([0-9]?)|(3[01]?))(-|\/)((19)([2-9])(\d{1})|(20)([01])(\d{1})|([8901])(\d{1}))|(0?[2469]|11)(-|\/)(([1-9])|(0[1-9])|([12])([0-9]?)|(3[0]?))(-|\/)((19)([2-9])(\d{1})|(20)([01])(\d{1})|([8901])(\d{1})))$`
		mmDD = regexp.MustCompile(mmDDStr)
		genericDate = regexp.MustCompile(`(^\d{1,2}-\w{3,9}-\d{2,4}$)|(^\d{1,2}\\w{3,9}\\d{2,4}$)|(^\d{1,2}/\w{3,9}/\d{2,4}$)`)

	}

	if ddMM.MatchString(dateInString) && !ismmDD {

		if strings.Contains(dateInString, "/") {
			splitDate := strings.Split(dateInString, "/")
			inDDMMFormat := strings.Join([]string{splitDate[2], splitDate[1], splitDate[0]}, "/")
			return strings.Replace(inDDMMFormat, "/", "-", -1)

		} else if strings.Contains(dateInString, "-") {
			splitDate := strings.Split(dateInString, "-")
			inDDMMFormat := strings.Join([]string{splitDate[2], splitDate[1], splitDate[0]}, "/")
			return strings.Replace(inDDMMFormat, "/", "-", -1)
		}

	} else if mmDD.MatchString(dateInString) || ismmDD {
		ismmDD = true
		if strings.Contains(dateInString, "/") {
			splitDate := strings.Split(dateInString, "/")
			inDDMMFormat := strings.Join([]string{splitDate[2], splitDate[0], splitDate[1]}, "/")
			return strings.Replace(inDDMMFormat, "/", "-", -1)

		} else if strings.Contains(dateInString, "-") {
			splitDate := strings.Split(dateInString, "-")
			inDDMMFormat := strings.Join([]string{splitDate[2], splitDate[0], splitDate[1]}, "/")
			return strings.Replace(inDDMMFormat, "/", "-", -1)
		}

	}else if genericDate.MatchString(dateInString) {
			// TODO This needs to be looked into
		dateStr1 := strings.Split(dateInString, "-")
		dateStr2 := strings.Split(dateInString, "/")
		dateStr3 := strings.Split(dateInString, "\\")

		var monthIndex int

		if len(dateStr1) == 3 {
			monthIndex = findMonthIndex(dateStr1[1])
			inDDMMFormat := strings.Join([]string{dateStr1[0], strconv.Itoa(monthIndex), dateStr1[2]}, "/")
			return inDDMMFormat
		} else if len(dateStr2) == 3 {
			monthIndex = findMonthIndex(dateStr2[1])
			inDDMMFormat := strings.Join([]string{dateStr2[0], strconv.Itoa(monthIndex), dateStr2[2]}, "/")
			return strings.Replace(inDDMMFormat, "/", "-", -1)
		} else if len(dateStr3) == 3 {
			monthIndex = findMonthIndex(dateStr3[1])
			inDDMMFormat := strings.Join([]string{dateStr3[0], strconv.Itoa(monthIndex), dateStr3[2]}, "/")
			return strings.Replace(inDDMMFormat, "/", "-", -1)
		}

	}

	return "ERROR"
}

var shortMonth []string = []string{"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"}
var longMonth []string = []string{"january", "febuary", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december"}

func findMonthIndex(monthName string) int {
	mName := strings.ToLower(monthName)

	for i, each := range shortMonth {
		if mName == each {
			return i + 1
		}
	}

	for j, eachLongMonth := range longMonth {
		if mName == eachLongMonth {
			return j + 1
		}
	}

	return -1
}

/*
*	DB Functions
 */

func dbCreateTable(createTableSql string) {
	if db == nil {
		db = getDB()
	}

	stmt, err := db.Prepare(createTableSql)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec()

	if err != nil {
		log.Fatal(err)
	}
}
const(

		size = 5
	)
var insertDatArr [size]string

var cnt int = 0
var cnt1 int = 0
func dbInsertData(insertStmts string) {

	if(db==nil){
		db= getDB()
	}

	db.Exec(insertStmts)

	// insertDatArr[cnt] = insertStmts
	// cnt++
	//
	// if cnt == size{
	// 	go func(){
	// 		if(db==nil){
	// 			db =getDB()
	// 		}
	// 	  var tempDatArr [size]string
	// 		  copy(tempDatArr[:],insertDatArr[:])
	// 			for i:=0; i<len(tempDatArr);i++{
	// 				if(len(tempDatArr[i]) > 0){
	//
	// 						db.Exec(tempDatArr[i])
	// 				}
	// 			}
	//
	//
	// 	}()
	// 	cnt = 0
	// 	cnt1++
	// }
	//
	// log.Println("CNT ",cnt1)

}

var db *sql.DB

func getDB() *sql.DB {
	// Connect to db
	db, _ = sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable")
	return db
}
