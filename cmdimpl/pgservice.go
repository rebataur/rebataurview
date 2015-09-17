package cmdimpl

import (
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

func dbStarted() bool {
	if db == nil {
		db = getDB()
	}
	var result []byte
	db.QueryRow("SELECT 2*3").Scan(&result)

	if len(result) == 1 {
		return true
	}
	return false
}
func StartPG(dbPath string) {
	log.Println("starting up")
	if dbStarted() == false {
		go func() {
			out, err := exec.Command(strings.Join([]string{dbPath, "\\bin\\pg_ctl"}, ""), "start", "-D", "D:\\Program Files\\PostgreSQLPortable-9.4\\Data\\data").Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(out))

		}()

		time.Sleep(10 * time.Second)
	} else {
		log.Println("DB is already running")
	}

}

func StopPG(dbPath string) {

	fmt.Println("Cleaning up")
	go func() {
		out, err := exec.Command(strings.Join([]string{dbPath, "\\bin\\pg_ctl"}, ""), "stop", "-D", "D:\\Program Files\\PostgreSQLPortable-9.4\\Data\\data", "-W", "-m", "immediate").Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(out))
	}()

	time.Sleep(10 * time.Second)
}
func SetupDB() {
	if db == nil {
		db = getDB()
	}
	db.QueryRow(`CREATE TABLE TABLE_METADATA(
	id bigserial primary key not null ,
	table_name varchar unique not null,
	table_desc varchar,
	table_origin varchar,
	creation_date date not null default CURRENT_DATE

	)`)

	db.QueryRow(`create table columns_metadata(
	id bigserial primary key not null,
	column_name varchar,
	ordinal_position int,
	data_type varchar,
	measure boolean default false,
	owning_table bigserial REFERENCES table_metadata(id)
)`)
	db.QueryRow(`create extension pg_trgm`)
	db.QueryRow(`CREATE INDEX trgm_idx ON columns_metadata USING gist (column_name gist_trgm_ops)`)
}

func PGCopyCmd(tableName string, filePath string) {
	if db == nil {
		db = getDB()
	}
	copyQuery := fmt.Sprintf("COPY %s FROM '%s' DELIMITERS ',' CSV HEADER", tableName, filePath)
	db.QueryRow(copyQuery)
}
func PGCreateTable(tableName string, colsType []string) {
	log.Println("Creating tables")
	dropTableSQL := fmt.Sprintf("DROP TABLE IF  EXISTS %s; \n", tableName)
	createTableSQL := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ( %s ); \n ", tableName, strings.Join(colsType, " ,"))

	execPrepareStmt(dropTableSQL)
	execPrepareStmt(createTableSQL)

}
func GetMetaData(mdType, mdVal string) []byte {
	if db == nil {
		db = getDB()
	}

	var query string

	if mdType == "table_schema" {
		query = `

				select array_to_json(array_agg(row_to_json(t))) as result
				from (
 					select table_name from information_schema.tables
						where table_schema = 'public'
				) t ;

				`
	} else if mdType == "columns_schema" {
		query = fmt.Sprintf(`

				select array_to_json(array_agg(row_to_json(t))) as result
				from (
					select column_name,ordinal_position,data_type,numeric_precision from information_schema.columns
						where table_name = '%s'
				) t ;
				`, mdVal)
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
	} else if mdType == "query" {
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

// Create table metadata from the imported table
type tableMetaData struct {
	column_name      string
	ordinal_position int
	data_type        string
	dimension        bool
	measure          bool
	owning_table     int
}

func CreateTableMetaData(tableName string) {
	if db == nil {
		db = getDB()
	}

	// Insert into Table metadata
	var id int
	q := fmt.Sprintf(`select id from table_metadata where table_name = '%s'`, tableName)
	db.QueryRow(q).Scan(&id)
	if id > 0 {
		s := fmt.Sprintf(`delete from columns_metadata where owning_table = %d`, id)
		execPrepareStmt(s)
		s = fmt.Sprintf(`delete from table_metadata where table_name = '%s'`, tableName)
		execPrepareStmt(s)

	}
	insertStmt := fmt.Sprintf(`insert into table_metadata (table_name,table_desc,table_origin) values ('%s','%s','%s') RETURNING id`, tableName, tableName, tableName)
	err := db.QueryRow(insertStmt).Scan(&id)
	// select the PG table to get metadata
	queryStmt := fmt.Sprintf(`select column_name,ordinal_position,data_type  from information_schema.columns
	where table_name = '%s'	`, strings.ToLower(tableName))

	rows, err := db.Query(queryStmt)

	if err != nil {
		log.Fatal("Error occured while create table metadata")
	}
	for rows.Next() {
		tblMeta := new(tableMetaData)
		err := rows.Scan(&tblMeta.column_name, &tblMeta.ordinal_position, &tblMeta.data_type)
		if err == nil {

			var isMeasure bool = findIfColumnIsMeasure(tblMeta, tableName)
			//  insertStmtColums := fmt.Sprintf(,tableName,tableName,tableName)
			insertStmtCol := fmt.Sprintf(`insert into columns_metadata
			(column_name,ordinal_position,data_type,measure,owning_table)
			 values ('%s',%d,'%s',%t,%d)`, tblMeta.column_name, tblMeta.ordinal_position, tblMeta.data_type, isMeasure, id)
			db.QueryRow(insertStmtCol)

		} else {
			log.Println("Error parsing tbl meta", err)
		}
	}
}

// Find whether column is a measure
func findIfColumnIsMeasure(tblMeta *tableMetaData, tableName string) bool {
	switch tblMeta.data_type {
	case "numeric":
		return true
	case "bigint":
		lowerCol := strings.ToLower(tblMeta.column_name)
		if strings.Contains(lowerCol, "id") {
			return false
		} else {
			return true
		}
	case "character varying":
		return false
	case "date":
		return false
	default:
		return false
	}
}

// Execute a prepared statement
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

// Query Statement to fetch data
func queryStmt(queryStmt string) []byte {
	if db == nil {
		db = getDB()
	}
	var result []byte
	err := db.QueryRow(getJSONQuery(queryStmt)).Scan(&result)

	//rows, _ := db.Query(queryStmt)
	//printRows(rows)
	if err != nil {
		result = []byte("there was an error executing query")
	}
	return result

}
func getJSONQuery(queryStmt string) string {
	jsonQuery := fmt.Sprintf("select array_to_json(array_agg(row_to_json(t))) as result	from ( %s ) t ;", queryStmt)
	return jsonQuery

}

// Print Query Output
func printRows(rows *sql.Rows) {

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
