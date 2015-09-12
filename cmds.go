package main
import(
  // "fmt"
  "sync"
)
import(

	"github.com/spf13/cobra"
)
var result []byte
var mutex = &sync.Mutex{}

var rootCmd = &cobra.Command{Use: "rebataurview"}
var cmdInitPG = &cobra.Command{
  Use:   "initpg",
  Short: "Initialize Postgresql",
  Long: `This will create if not exist, a database in the provided repository path,
      and initialize and start Postgres database.
      `,
  Run: func(cmd *cobra.Command, args []string) {
     loadDataIntoPG(args[0],true)
  },
}
var cmdDescribeTable = &cobra.Command{
  Use:   "describe_table",
  Short: "Describe Table",
  Long: `This describes table`,
  Run: func(cmd *cobra.Command, args []string) {
    mutex.Lock()
    result = describeTable(args[0])
    mutex.Unlock()
  },
}

var cmdDescribeColumns = &cobra.Command{
  Use:   "describe_columns",
  Short: "Describe Column",
  Long: `This describes Columns of table`,
  Run: func(cmd *cobra.Command, args []string) {
    mutex.Lock()
    result = describeColumns(args[0])
    mutex.Unlock()
  },
}
var cmdFrequencyCount = &cobra.Command{
  Use:   "freq",
  Short: "Frequency count",
  Long: `This gets frequency count of column of table`,
  Run: func(cmd *cobra.Command, args []string) {
    mutex.Lock()
    res,err := getColumnFrequency(args[0],args[1],args[2])
    if err !=nil {
      result = []byte("There was an error")
    }else{
      result = res
    }
    mutex.Unlock()
  },
}
func init(){
  rootCmd.AddCommand(cmdInitPG)
  rootCmd.AddCommand(cmdDescribeTable)
  rootCmd.AddCommand(cmdDescribeColumns)
  rootCmd.AddCommand(cmdFrequencyCount)

}
