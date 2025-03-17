//underscores => variables
//camelCase => functions

// worker pool blog -> https://rksurwase.medium.com/efficient-concurrency-in-go-a-deep-dive-into-the-worker-pool-pattern-for-batch-processing-73cac5a5bdca
// closing go routines blog -> https://callistaenterprise.se/blogg/teknik/2019/10/05/go-worker-cancellation/
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/MustafaLo/file-search/config"
	"github.com/spf13/cobra"
)

func timeTrack(start time.Time, name string) {
    elapsed := time.Since(start)
    fmt.Printf("\n\n%s took %s", name, elapsed)
}


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
func searchFile(id int, ctx context.Context, cancel context.CancelFunc, jobs <-chan Job, results chan <- Result, counter *int32, wg *sync.WaitGroup){
	defer wg.Done()
	//Ensures that jobs are distributed evenly (more or less) among workers
	//Each worker does approx jobCount / workerCount jobs
	for job := range jobs{
		select {
		case <- ctx.Done():
			return
		default:
		}
			for line_number, line := range job.file_content{
				if strings.Contains(line, search_term){
					time.Sleep(100 * time.Millisecond)
					newCount := atomic.AddInt32(counter, 1)
					fmt.Println(newCount)
					if newCount >= 20{
						cancel()
						return
					}
					select {
					case results <- Result{
						file_name: job.file_name, 
						line_content: line,
						line_number: line_number + 1,
					}:
					case <-ctx.Done():
						return
					}
			
				}	
			}
	}
}

func collectResults(results <-chan Result, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("\n====================== 🔍 SEARCH RESULTS 🔍 ======================")

	for result := range results {
		trimmedContent := strings.TrimSpace(result.line_content) // Trim whitespace

		// Highlight search term in the result
		highlightedContent := strings.ReplaceAll(trimmedContent, search_term, config.Colors["red"]+search_term+config.Colors["reset"])

		fmt.Printf("\n📂 File: %-20s  📍 Line #%-5d\n   👉 Line:  %s\n", 
			result.file_name, result.line_number, highlightedContent)
		fmt.Println("---------------------------------------------------------------")
	}

	// if count == 0 {
	// 	fmt.Println("\n❌ No results found.")
	// }

	fmt.Println("\n================================================================")
}

// func startConsumer(ingest <- chan Job, jobs chan <- Job, ctx context.Context){
// 	for{
// 		select {
// 		case job := <-ingest:
// 			jobs <- job
// 		case <-ctx.Done():
// 			fmt.Println("Consumer received cancellation signal, closing jobsChan!")
// 			close(jobs)
// 			fmt.Println("Consumer closed jobsChan")
// 			return
// 		}
// 	}
// }



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
	  defer timeTrack(time.Now(), "Search")

	  file_paths, err := getDirectoryFiles()
	  if err != nil{
		fmt.Printf("Error retrieving files from directory: %v", err)
		return
	  }

	  var counter int32 = 0
 
	  //Set up waitgroup and cancellation context
	  var wg sync.WaitGroup

	  ctx, cancel := context.WithCancel(context.Background())
	  defer cancel()

	  jobCount := len(file_paths)

	  //Need this channel as a proxy between jobs and results in order to close jobs channel and stop workers gracefully
	//   ingest := make(chan Job)
	  
	  //Channel to ingest jobs (files in directories)
	  jobs := make(chan Job)

	  //Channel that will recieve results from search as we find matches
	  results := make(chan Result)

	  //Changing workercount will make results appear in batches (faster) since 
	  //multiples workers are processing different files at the same time
	  workerCount := 1
	  wg.Add(workerCount)

	  //Start workers
	  for w := 0; w < workerCount; w++{
		go searchFile(w, ctx, cancel, jobs, results, &counter, &wg)
	  }

	  //Start collecting results
	  var resultsWg sync.WaitGroup
	  resultsWg.Add(1) //Only need to add 1 to results wait group since only starting one go routine for collecting results
	  go collectResults(results, &resultsWg)

	  //Start Consumer
	//   go startConsumer(ingest, jobs, ctx)

	  //Distribute jobs
	  //Distribute jobs into ingest first to act as a proxy (ingest -> jobs -> results)
	  for j := 0; j < jobCount; j++{
		select {
		case <- ctx.Done():
			break
		default:
			name := file_paths[j]
			content, err := getFileContent(name)
			if err != nil{
				fmt.Printf("Error retrieving file content for %s: %v", name, err)
				continue
			}
			// ingest <- Job{file_name: name, file_content: content,}


			select {
			case jobs <- Job{file_name: name, file_content: content}:
			case <-ctx.Done():
				break
			}
			
		}
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
	searchCmd.MarkFlagRequired("term")
	searchCmd.Flags().StringVarP(&directory, "directory", "d", ".", "directory you would like to search in")
}

