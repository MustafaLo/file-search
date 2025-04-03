package utils

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/MustafaLo/file-search/config"
)

func GetFileContent(file_path string)([]string, error){
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

func GetDirectoryFiles(directory string)([]string, error){
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

func GetFileInformation(file_path string)(os.FileInfo, error){
	file_info, err := os.Stat(file_path)
	if err != nil{
		return nil, err
	}

	return file_info, err
}

func GetDirectoryStructure(directory string) (string, error) {
	var structure strings.Builder

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip excluded directories
		if d.IsDir() {
			if _, found := config.ExcludedDirs[d.Name()]; found {
				return filepath.SkipDir
			}
		} else {
			// Skip files with excluded extensions
			ext := filepath.Ext(d.Name())
			if _, found := config.ExcludedExtensions[ext]; found {
				return nil
			}
		}

		// Get relative path and depth
		relPath, err := filepath.Rel(directory, path)
		if err != nil {
			return err
		}
		depth := strings.Count(relPath, string(filepath.Separator))
		indent := strings.Repeat("    ", depth)

		// Add line to structure
		if d.IsDir() {
			structure.WriteString(fmt.Sprintf("%s%s/\n", indent, d.Name()))
		} else {
			structure.WriteString(fmt.Sprintf("%s%s\n", indent, d.Name()))
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return structure.String(), nil
}

// func GetDirectoryStructure(directory string)(string, error){
// 	var directory_structure strings.Builder
// 	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
// 		if d.IsDir(){
// 			directory_structure.WriteString(d.Name() + "/\n")
// 		}
// 		fmt.Println(d.Name())
// 		return nil
// 	})

// 	if err != nil{
// 		return "", err
// 	}

// 	return "", nil
	
// }