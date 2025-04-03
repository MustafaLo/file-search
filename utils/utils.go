package utils

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"

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