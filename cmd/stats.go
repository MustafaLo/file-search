package cmd

import "github.com/spf13/cobra"

/*
Stats about directory:
- File structure (including sub directories)
- Number of files
- Largest / Smallest files
- Most common file type
- Most recently modified file
- Least recently modified file
*/


var statsDir string
var statsCmd = &cobra.Command{
	Use: "statistics",
	Short: "find statistics about a directory",
	Long: `Use this command to find different statistics about a directory`,
	Run: func (cmd *cobra.Command, args[] string)  {
		
	},
}

func init(){
	statsCmd.Flags().StringVarP(&statsDir, "directory", "d", ".", "directory you would like to search in")
}