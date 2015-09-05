package cmdimplementors

import (
	"fmt"
	"os/exec"
	"time"
	"log"
)

func startPG(command string) {
	log.Println("starting up")
	go func() {
			out, err := exec.Command("D:\\Program Files\\PostgreSQLPortable-9.4\\App\\PgSQL\\bin\\pg_ctl", "start", "-D", "D:\\Program Files\\PostgreSQLPortable-9.4\\Data\\data").Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(out))

	}()

	time.Sleep(10 * time.Second)
	return
}

func stopPG(command string) {

	fmt.Println("Cleaning up")
	go func() {
			out, err := exec.Command("D:\\Program Files\\PostgreSQLPortable-9.4\\App\\PgSQL\\bin\\pg_ctl", "stop", "-D", "D:\\Program Files\\PostgreSQLPortable-9.4\\Data\\data", "-W", "-m", "immediate").Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(out))
	}()


	time.Sleep(10 * time.Second)
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
