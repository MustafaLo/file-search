//underscores => variables
//camelCase => functions

//worker pool blog -> https://rksurwase.medium.com/efficient-concurrency-in-go-a-deep-dive-into-the-worker-pool-pattern-for-batch-processing-73cac5a5bdca

package cmd

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"github.com/MustafaLo/file-search/config"
	"github.com/spf13/cobra"
)


func getDirectoryFiles()([]string, error){
	var files []string

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) (error){
		if d.IsDir(){
			if _, found := config.ExcludedDirs[d.Name()]; found {
				return filepath.SkipDir 
			}
			return nil 
		}

		ext := filepath.Ext(d.Name())


		if _, found := config.ExcludedExtensions[ext]; found{
			return nil
		}
		
		files = append(files, path)
		return nil
	})

	if err != nil{
		return nil, err
	}

	return files, nil
}

func getFileContent(file_path string)([]string, error){
	return []string{}, nil
}


type Job struct{
	file_name string
	file_content []string
}

type Result struct{
	file_name string
	line_content string
	line_number int
	search_term string
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

	  

	  
	},
}

func init(){
	searchCmd.Flags().StringVarP(&search_term, "term", "t", "", "search term for search command")
	// searchCmd.MarkFlagRequired("term")
	searchCmd.Flags().StringVarP(&directory, "directory", "d", ".", "directory you would like to search in")
}

