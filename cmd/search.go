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
	"strings"
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
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	return fileLines, nil
}

//worker function
func searchFile(id int, jobs <-chan Job, results chan <- Result, wg *sync.WaitGroup){
	defer wg.Done()
	//Ensures that jobs are distributed evenly (more or less) among workers
	//Each worker does approx jobCount / workerCount jobs
	for job := range jobs{

		for line_number, line := range job.file_content{
			if strings.Contains(line, search_term){
				results <- Result{
					file_name: job.file_name, 
					line_content: line,
					line_number: line_number,
				}
			}
		}
	}
}

//collect results worker function
func collectResults(results <- chan Result, wg *sync.WaitGroup){
	defer wg.Done()
	for result := range results{
		fmt.Printf("\nFile: %s  Line #%d:  %s", result.file_name, result.line_number, result.line_content)
	}
}


type Job struct{
	file_name string
	file_content []string
}

type Result struct{
	file_name string
	line_content string
	line_number int
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
 
	  var wg sync.WaitGroup
	  jobs := make(chan Job, 10)
	  results := make(chan Result, 10)

	  //Replace with actual worker count
	  workerCount := 5
	  wg.Add(workerCount)

	  //Start workers
	  for w := 1; w <= workerCount; w++{
		go searchFile(w, jobs, results, &wg)
	  }

	  //Start collecting results
	  var resultsWg sync.WaitGroup
	  resultsWg.Add(1) //Only need to add 1 to results wait group since only starting one go routine for collecting results
	  go collectResults(results, &resultsWg)


	  //Distribute jobs
	  for j := 0; j <= workerCount; j++{
		name := file_paths[j]
		content, err := getFileContent(name)
		if err != nil{
			fmt.Printf("Error retrieving file content for %s: %v", name, err)
			continue
		}
		jobs <- Job{file_name: name, file_content: content}
	  }
	  close(jobs)
	  wg.Wait()
	  close(results)

	  //Ensure all results are collected
	  resultsWg.Wait()




	  
	},
}

func init(){
	searchCmd.Flags().StringVarP(&search_term, "term", "t", "", "search term for search command")
	// searchCmd.MarkFlagRequired("term")
	searchCmd.Flags().StringVarP(&directory, "directory", "d", ".", "directory you would like to search in")
}

