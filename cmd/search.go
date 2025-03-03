//underscores => variables
//camelCase => functions

package cmd

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"github.com/spf13/cobra"
)

func getDirectoryFiles()([]string, error){
	var files []string

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) (error){
		if _, found := excludedExtensions[filepath.Ext(path)]; found{
			continue
		}
		
		if !d.IsDir(){
			files = append(files, path)
		}
		return nil
	})

	if err != nil{
		return nil, err
	}

	return files, nil
}

var search_term string
var directory string
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search a keyword",
	Long:  `Use this command to search for a keyword within your directory`,
	Run: func(cmd *cobra.Command, args []string) {
	  files, err := getDirectoryFiles()
	  if err != nil{
		fmt.Printf("Error retrieving files from directory: %v", err)
		return
	  }

	  for _, file_path := range(files){
		fmt.Println(file_path)
	  }
	},
}

func init(){
	searchCmd.Flags().StringVarP(&search_term, "term", "t", "", "search term for search command")
	// searchCmd.MarkFlagRequired("term")
	searchCmd.Flags().StringVarP(&directory, "directory", "d", ".", "directory you would like to search in")
}

