package cmd

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/MustafaLo/file-search/utils"
	"github.com/spf13/cobra"
)

/*
Stats about directory:
- File structure (including sub directories)
- Number of files
- Largest / Smallest files
- Most common file type
- Most recently modified file
- Least recently modified file
*/

func printDirectoryStats(fileStructure string, numFiles int, largestFile string, largestSize int64, 
	smallestFile string, smallestSize int64, commonFileType string, commonFileCount int, 
	mostRecentFile string, mostRecentTime string, leastRecentFile string, leastRecentTime string) {

	// Cool ASCII header
	header := `
=====================================
📂 DIRECTORY STATISTICS 📂
=====================================
`

	// Print all statistics with nice formatting
	fmt.Println(header)

	fmt.Println("📁 File Structure:")
	fmt.Println(fileStructure)
	fmt.Println()
	fmt.Printf("📊 Number of Files: %d\n\n", numFiles)

	fmt.Printf("📌 Largest File:       %s (%d)\n", largestFile, largestSize)
	fmt.Printf("📌 Smallest File:      %s (%d)\n\n", smallestFile, smallestSize)

	fmt.Printf("📂 Most Common File Type:  %s (%d files)\n\n", commonFileType, commonFileCount)

	fmt.Printf("⏳ Most Recently Modified:  %s (%s)\n", mostRecentFile, mostRecentTime)
	fmt.Printf("⌛ Least Recently Modified:  %s (%s)\n", leastRecentFile, leastRecentTime)
}

func getExtremesFileSizes(directory_information []os.FileInfo)(os.FileInfo, os.FileInfo, error){
	min_file, min_size := os.FileInfo(nil), math.Inf(-1)
	max_file, max_size := os.FileInfo(nil), math.Inf(1)

	for _, file_information := range directory_information{
		if file_information.Size() < int64(max_size){
			min_file = file_information
			max_size = float64(file_information.Size())
		}

		if file_information.Size() > int64(min_size){
			max_file = file_information
			min_size = float64(file_information.Size())
		}
	}

	if min_file == nil || max_file == nil{
		return nil, nil, errors.New("error evaluating minimum/maximum file size")
	}

	return min_file, max_file, nil

}

func getModifiedFiles(directory_information []os.FileInfo)(os.FileInfo, os.FileInfo, error){
	newest_file, newest_time := os.FileInfo(nil), time.Now()
	oldest_file, oldest_time := os.FileInfo(nil), time.Time{}

	for _, file_information := range directory_information{
		if file_information.ModTime().After(oldest_time){
			newest_file = file_information
			oldest_time = file_information.ModTime()
		}

		if file_information.ModTime().Before(newest_time){
			oldest_file = file_information
			newest_time = file_information.ModTime()
		}
	}


	if newest_file == nil || oldest_file == nil{
		return nil, nil, errors.New("error evaluating most recent/least recent modified file")
	}

	return newest_file, oldest_file, nil
}

func getMostCommonFileType(directory_information []os.FileInfo)(string, int, error){
	filetype_counts := make(map[string]int)
	max_filetype, min_count := "", 0

	for _, file_information := range directory_information{
		ext := filepath.Ext(file_information.Name())
		filetype_counts[ext] += 1
		if filetype_counts[ext] > min_count{
			max_filetype = ext
			min_count = filetype_counts[ext]
		}

	}

	if max_filetype == ""{
		return "", 0, errors.New("Unable to find most common file type")
	}

	return max_filetype, filetype_counts[max_filetype], nil

}


var stats_dir string
var statsCmd = &cobra.Command{
	Use: "stats",
	Short: "find statistics about a directory",
	Long: `Use this command to find different statistics about a directory`,
	Run: func (cmd *cobra.Command, args[] string)  {
		file_paths, err := utils.GetDirectoryFiles(stats_dir)
		if err != nil{
			fmt.Printf("Error retrieving files from search_dir: %v", err)
			return
		}

		dir_information := make([]os.FileInfo, 0)
		for _, file_path := range file_paths{
			file_information, err := utils.GetFileInformation(file_path)
			if err != nil{
				fmt.Printf("Could not get file information for %s: %v",file_path, err)
			}
			dir_information = append(dir_information, file_information)
		}

		// for _, info := range dir_information{
		// 	fmt.Println(info.Size())
		// 	fmt.Println(info.ModTime())
		// 	fmt.Println()
		// }

		min_file, max_file, err := getExtremesFileSizes(dir_information)
		if err != nil{
			fmt.Printf("Error retrieving minimum/maximum files: %v", err)
		}

		// fmt.Printf("Minfile: %d\n", min_file.Size())
		// fmt.Printf("Maxfile: %d\n", max_file.Size())

		recent_modified_file, oldest_modified_file, err := getModifiedFiles(dir_information)
		if err != nil{
			fmt.Printf("Error retrieving newest/oldest modified files: %v", err)
		}

		// fmt.Printf("Most Recent file: %s\n", recent_modified_file.Name())
		// fmt.Printf("Least Recent file: %s\n", oldest_modified_file.Name())

		directory_structure, err := utils.GetDirectoryStructure(stats_dir)
		if err != nil{
			fmt.Printf("Error retrieving directory structure: %v", err)
		}

		// fmt.Println(directory_structure)

		max_file_type, max_file_type_count, err := getMostCommonFileType(dir_information)
		if err != nil{
			fmt.Printf("Error getting most common file type: %v", err)
		}

		// fmt.Println(max_file_type)

		printDirectoryStats(directory_structure, len(file_paths), 
							max_file.Name(), max_file.Size(), 
							min_file.Name(), min_file.Size(),
							max_file_type, max_file_type_count,
							recent_modified_file.Name(), recent_modified_file.ModTime().Format("Mon, Jan 2, 2006 at 3:04 PM"),
							oldest_modified_file.Name(), oldest_modified_file.ModTime().Format("Mon, Jan 2, 2006 at 3:04 PM"))
		

	},
}

func init(){
	statsCmd.Flags().StringVarP(&stats_dir, "directory", "d", ".", "directory you would like to search in")
}