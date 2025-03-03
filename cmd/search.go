//underscores => variables
//camelCase => functions

package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

func getDirectoryFiles()(string, error){
	files, err := os.ReadDir(directory)
	if err != nil{
		return "", err
	}

	for _, f := range(files){
		fmt.Println(f.Name())
	}

	return "", nil

}

var search_term string
var directory string
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search a keyword",
	Long:  `Use this command to search for a keyword within your directory`,
	Run: func(cmd *cobra.Command, args []string) {
	  getDirectoryFiles()
	},
}

func init(){
	searchCmd.Flags().StringVarP(&search_term, "term", "t", "", "search term for search command")
	// searchCmd.MarkFlagRequired("term")
	searchCmd.Flags().StringVarP(&directory, "directory", "d", ".", "directory you would like to search in")
}

