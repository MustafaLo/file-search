//underscores => variables
//camelCase => functions

//worker pool blog -> https://rksurwase.medium.com/efficient-concurrency-in-go-a-deep-dive-into-the-worker-pool-pattern-for-batch-processing-73cac5a5bdca

package cmd

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

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
	var fileLines []string

	readFile, err := os.Open(file_path)
	if err != nil{
		return nil, err
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan(){
		//Ignoring any lines with no text
		if fileScanner.Text() == ""{
			continue
		}
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	return fileLines, nil
}

//worker function
func searchFile(id int, jobs <-chan Job, results chan <- Result, wg *sync.WaitGroup){

}

//collect results worker function
func collectResults(results <- chan Result, wg *sync.WaitGroup){
	
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
	  file_paths, err := getDirectoryFiles()
	  if err != nil{
		fmt.Printf("Error retrieving files from directory: %v", err)
		return
	  }
 
	//   var wg sync.WaitGroup
	//   jobs := make(chan Job, 10)
	//   results := make(chan Result, 10)

	// for _, path := range file_paths{
	// 	fmt.Println(path)
	// }

	fileContent, err := getFileContent(file_paths[1])
	for lineNumber, line := range fileContent{
		fmt.Printf("\nLine #%d %s", lineNumber, line)
	}

	  //Start workers

	  //Start collecting results

	  //Distribute jobs




	  
	},
}

func init(){
	searchCmd.Flags().StringVarP(&search_term, "term", "t", "", "search term for search command")
	// searchCmd.MarkFlagRequired("term")
	searchCmd.Flags().StringVarP(&directory, "directory", "d", ".", "directory you would like to search in")
}

