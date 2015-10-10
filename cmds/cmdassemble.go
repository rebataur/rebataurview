package cmds

import (
	// "fmt"
	"strings"
	"sync"
)
import (
	"github.com/spf13/cobra"
)

var Result []byte
var mutex = &sync.Mutex{}
var fn, clmn, tbln string
var limit int
var rootCmd = &cobra.Command{Use: "rebataurview"}
var cmdInitPG = &cobra.Command{
	Use:   "initpg",
	Short: "Initialize Postgresql",
	Long: `This will create if not exist, a database in the provided repository path,
      and initialize and start Postgres database.
      `,
	Run: func(cmd *cobra.Command, args []string) {
		LoadDataIntoPG(args[0], true)
	},
}
var cmdDescribeTable = &cobra.Command{
	Use:   "describe_table",
	Short: "Describe Table",
	Long:  `This describes table`,
	Run: func(cmd *cobra.Command, args []string) {
		mutex.Lock()
		Result = describeTable()
		mutex.Unlock()
	},
}

var cmdDescribeColumns = &cobra.Command{
	Use:   "describe_columns",
	Short: "Describe Column",
	Long:  `This describes Columns of table`,
	Run: func(cmd *cobra.Command, args []string) {
		mutex.Lock()
		Result = describeColumns(args[0])
		mutex.Unlock()
	},
}
var cmdAnalyze = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze",
	Long:  `This is the mother analytics command`,
	Run: func(cmd *cobra.Command, args []string) {
		mutex.Lock()
		res, err := doAnalytics(clmn, tbln, fn, limit, args)
		if err != nil {
			Result = []byte("There was an error")
		} else {
			Result = res
		}
		mutex.Unlock()
	},
}

func init() {
	rootCmd.AddCommand(cmdInitPG)
	rootCmd.AddCommand(cmdDescribeTable)
	rootCmd.AddCommand(cmdDescribeColumns)
	cmdAnalyze.Flags().StringVarP(&clmn, "colname", "c", "", "column name to be analysed")
	cmdAnalyze.Flags().StringVarP(&fn, "funcname", "f", "", "function name to be applied")
	cmdAnalyze.Flags().StringVarP(&tbln, "tablename", "t", "", "name of the table")
	cmdAnalyze.Flags().IntVarP(&limit, "limit", "l", 1000, "limits the number of records to be analyzed to")
	rootCmd.AddCommand(cmdAnalyze)
}
func SetAndExecuteCmd(cmd string) {
	rootCmd.SetArgs(strings.Split(cmd, " "))
	rootCmd.Execute()
}
