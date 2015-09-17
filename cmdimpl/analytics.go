package cmdimpl

import (
	"fmt"
)

func GetFrequencyCount(colName string, tableName string, limit int) []byte {
	query := fmt.Sprintf("select %s, count(%s) as cnt from %s group by %s order by cnt desc limit %d", colName, colName, tableName, colName, limit)
	return queryStmt(query)
}

func AnalyzeTable(tableName string) {
	// query := fmt.Sprintf(`select column_name,ordinal_position,data_type,numeric_precision  from information_schema.columns
	// 	where table_name = '%s'`,tableName)

}
func GetDimensionsAndMeasures(tableName string) {

}
