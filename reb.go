package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
	// "github.com/ranjanprj/rebataurview/cmdimpl"

)

func main() {

	var cmdInitPG = &cobra.Command{
		Use:   "initpg",
		Short: "Initialize Postgresql",
		Long: `This will create if not exist, a database in the provided repository path,
        and initialize and start Postgres database.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Print: " + strings.Join(args, " "))
		},
	}

	// cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")

	var rootCmd = &cobra.Command{Use: "reb"}
	rootCmd.AddCommand(cmdInitPG)
	rootCmd.Execute()
}
