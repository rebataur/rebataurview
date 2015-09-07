package cmdimpl

import (
	"fmt"
	"os/exec"
	"time"
	"log"
	"strings"
	"database/sql"
)
func dbStarted() bool{
	if db == nil {
		db = getDB()
	}
	var result []byte
	db.QueryRow("SELECT 2*3").Scan(&result)

	if len(result) ==  1{
		return true
	}
	return false

}
func StartPG(dbPath string) {
	log.Println("starting up")
	if dbStarted() == false {
		go func() {
				out, err := exec.Command(strings.Join([]string{dbPath,"\\bin\\pg_ctl"},""), "start", "-D", "D:\\Program Files\\PostgreSQLPortable-9.4\\Data\\data").Output()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(string(out))

				}()


				time.Sleep(10 * time.Second)
	}else{
		log.Println("DB is already running")
	}

}


func StopPG(dbPath string) {

	fmt.Println("Cleaning up")
	go func() {
			out, err := exec.Command(strings.Join([]string{dbPath,"\\bin\\pg_ctl"},""), "stop", "-D", "D:\\Program Files\\PostgreSQLPortable-9.4\\Data\\data", "-W", "-m", "immediate").Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(out))
	}()


	time.Sleep(10 * time.Second)
}

func PGCopyCmd(tableName string, filePath string){
	if db == nil {
		db = getDB()
	}
	copyQuery := fmt.Sprintf("COPY %s FROM '%s' DELIMITERS ',' CSV HEADER",tableName, filePath)
	db.QueryRow(copyQuery)
}
func PGCreateTable(tableName string, colsType []string){
		log.Println("Creating tables")
		dropTableSQL := fmt.Sprintf("DROP TABLE IF  EXISTS %s; \n", tableName)
		createTableSQL := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ( %s ); \n ", tableName, strings.Join(colsType," ,"))

		log.Println(dropTableSQL)
		log.Println(createTableSQL)
		execPrepareStmt(dropTableSQL)
		execPrepareStmt(createTableSQL)



}
func GetMetaData(mdType, mdVal string) []byte {
	if db == nil {
		db = getDB()
	}

	var query string

	if mdType == "schema" {
		query = `

				select array_to_json(array_agg(row_to_json(t))) as result
				from (
 					select table_name from information_schema.tables
						where table_schema = 'public'
				) t ;

				`
	} else if mdType == "table" {
		query = fmt.Sprintf(`
				select array_to_json(array_agg(row_to_json(t))) as result
				from (
 					select column_name,ordinal_position,data_type,numeric_precision  from information_schema.columns
						where table_name = '%s'
				) t ;
			`, mdVal)

	} else if mdType == "count" {
		query = fmt.Sprintf(`
				select row_to_json(t) as result
				from (
 					select count(*) from %s
				) t ;
			`, mdVal)
	} else if mdType == "data" {
		query = fmt.Sprintf(`
				select array_to_json(array_agg(row_to_json(t))) as result
				from (
 					select * from %s limit 10
				) t ;
			`, mdVal)
	} else if mdType == "query"{
		query = fmt.Sprintf(`
				select array_to_json(array_agg(row_to_json(t))) as result
				from (
 					%s
				) t ;
			`, mdVal)
	}

	var result []byte
	db.QueryRow(query).Scan(&result)
	return result
}

func GetFrequencyCount(colName string,tableName string, limit int) []byte {

	query := fmt.Sprintf("select %s, count(%s) as cnt from %s group by state order by cnt desc limit %d",colName,colName,tableName,limit)
	return queryStmt(query)
}

func execPrepareStmt(prepStmt string) {
	if db == nil {
		db = getDB()
	}

	stmt, err := db.Prepare(prepStmt)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec()

	if err != nil {
		log.Fatal(err)
	}
}
func queryStmt(queryStmt string) []byte{
	if db == nil {
		db = getDB()
	}
	var result []byte
	db.QueryRow(getJSONQuery(queryStmt)).Scan(&result)

	rows,_ := db.Query(queryStmt)
	printRows(rows)
	return result;

}


func getJSONQuery(queryStmt string) string{
	jsonQuery := fmt.Sprintf("	select array_to_json(array_agg(row_to_json(t))) as result	from ( %s ) t ;",queryStmt)
	return jsonQuery

}


func printRows(rows *sql.Rows){

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	// Make a slice for the values
	values := make([]interface{}, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		// Print data
		for i, value := range values {
			switch value.(type) {
			case nil:
				fmt.Println(columns[i], ": NULL")

			case []byte:
				fmt.Println(columns[i], ": ", string(value.([]byte)))

			default:
				fmt.Println(columns[i], ": ", value)
			}
		}
		fmt.Println("-----------------------------------")
	}
}
