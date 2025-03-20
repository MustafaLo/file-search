//underscores => variables
//camelCase => functions

// worker pool blog -> https://rksurwase.medium.com/efficient-concurrency-in-go-a-deep-dive-into-the-worker-pool-pattern-for-batch-processing-73cac5a5bdca
// closing go routines blog -> https://callistaenterprise.se/blogg/teknik/2019/10/05/go-worker-cancellation/
package cmd

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/MustafaLo/file-search/config"
	"github.com/MustafaLo/file-search/utils"
	"github.com/spf13/cobra"
)

func timeTrack(start time.Time, name string) {
    elapsed := time.Since(start)
    fmt.Printf("\n\n%s took %s", name, elapsed)
}




//worker function
func searchFile(id int, ctx context.Context, jobs <-chan Job, results chan <- Result, counter *int32, wg *sync.WaitGroup){
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
					select {
					case <- ctx.Done():
						return
					case results <- Result{
						file_name: job.file_name, 
						line_content: line,
						line_number: line_number + 1,
					}:
					}
				}	
			}
	}

}

func collectResults(results <-chan Result, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("\n====================== 🔍 SEARCH RESULTS 🔍 ======================")

	collected := make([]Result, 0, limit)
	for result := range results {
		if len(collected) < limit{
			collected = append(collected, result)

			if len(collected) == limit{
				cancel()
			}
		}
	}

	for number, result := range collected{
		 // Trim whitespace
		trimmedContent := strings.TrimSpace(result.line_content)

		// Highlight search term in the result
		highlightedContent := strings.ReplaceAll(trimmedContent, search_term, config.Colors["red"]+search_term+config.Colors["reset"])

		fmt.Printf("\n %d 📂 File: %-20s  📍 Line #%-5d\n   👉 Line:  %s\n", 
			number + 1, result.file_name, result.line_number, highlightedContent)
		fmt.Println("---------------------------------------------------------------")
	}

	fmt.Println("\n================================================================")
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
var search_dir string
var limit int
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search a keyword",
	Long:  `Use this command to search for a keyword within your search_dir`,
	Run: func(cmd *cobra.Command, args []string) {
	  defer timeTrack(time.Now(), "Search")

	  file_paths, err := utils.GetDirectoryFiles(search_dir)
	  if err != nil{
		fmt.Printf("Error retrieving files from search_dir: %v", err)
		return
	  }

	  fmt.Println(file_paths)
	  //Set up waitgroup and cancellation context
	  var wg sync.WaitGroup

	  //Set up results counter
	  var counter int32 = 0

	  //Set up cancellation function (to cancel workers)
	  ctx, cancel := context.WithCancel(context.Background())

	  //Channel to ingest jobs (files in directories)
	  jobs := make(chan Job)

	  //Channel that will recieve results from search as we find matches
	  results := make(chan Result, 100)

	  //Changing workercount will make results appear in batches (faster) since 
	  //multiples workers are processing different files at the same time
	  workerCount := 1
	  wg.Add(workerCount)

	  //Start workers
	  for w := 0; w < workerCount; w++{
		go searchFile(w, ctx, jobs, results, &counter, &wg)
	  }

	  //Start collecting results
	  var resultsWg sync.WaitGroup
	  resultsWg.Add(1) //Only need to add 1 to results wait group since only starting one go routine for collecting results
	  go collectResults(results, cancel, &resultsWg)


	  //Distribute jobs
	  jobCount := len(file_paths)
	  jobloop:
	  for j := 0; j < jobCount; j++{
		select {
		case <- ctx.Done():
			break jobloop
		default:
			name := file_paths[j]
			content, err := utils.GetFileContent(name)
			if err != nil{
				fmt.Printf("Error retrieving file content for %s: %v", name, err)
				continue
			}

			select {
			case jobs <- Job{file_name: name, file_content: content}:
			case <-ctx.Done():
				break jobloop
			}
		}
	  }

	  cancel()
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
	searchCmd.Flags().StringVarP(&search_dir, "search_dir", "d", ".", "search_dir you would like to search in")
	searchCmd.Flags().IntVarP(&limit, "limit", "l", 10, "limit number of results to")
}

