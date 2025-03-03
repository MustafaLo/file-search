package cmd

import(
	"fmt"
    "github.com/spf13/cobra"
    "os"
)

var rootCmd = &cobra.Command{
	Use:   "search",
	Short: "Search CLI",
	Long: `This is a directory search CLI tool`,
  }
  
  func Execute() {
	if err := rootCmd.Execute(); err != nil {
	  fmt.Println(err)
	  os.Exit(1)
	}
  }

  func init(){
	rootCmd.AddCommand(searchCmd)
  }