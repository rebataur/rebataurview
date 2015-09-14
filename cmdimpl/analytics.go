package cmdimpl

import (
	"fmt"
)

func GetFrequencyCount(colName string, tableName string, limit int) []byte {
	query := fmt.Sprintf("select %s, count(%s) as cnt from %s group by %s order by cnt desc limit %d", colName, colName, tableName, colName, limit)
	return queryStmt(query)
}

func GetDimensionsAndMeasures(tableName string){

}
