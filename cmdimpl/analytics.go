package cmdimpl

import (
	"fmt"
	"strings"
)

func DoAnalytics(clmn, tbln, fn string, limit int, args []string) ([]byte, error) {
	var lmtStr string
	if limit > 0 {
		lmtStr = fmt.Sprintf(" limit %d", limit)
	}
	clmns := strings.Split(clmn, ",")
	firstCol := clmns[0]
	switch fn {
	case "freq":
		{
			query := fmt.Sprintf("select %s,count(%s) as freq from %s group by %s,product order by freq desc %s", clmn, firstCol, tbln, clmn, lmtStr)
			return queryStmt(query), nil
		}
	default:
		{
			return []byte("No analytics function provided"), nil
		}
	}

}

func AnalyzeTable(tableName string) {
	// query := fmt.Sprintf(`select column_name,ordinal_position,data_type,numeric_precision  from information_schema.columns
	// 	where table_name = '%s'`,tableName)

}
func GetDimensionsAndMeasures(tableName string) {

}
