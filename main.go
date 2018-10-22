package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	. "./index"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("File list doesn't specified")
		fmt.Println("Use:")
		fmt.Println(os.Args[0] + " [file1 file2 ...]")
		os.Exit(0)
	}
	var database map[string]Index
	database = make(map[string]Index, 1)

	for i := 1; i < len(os.Args); i++ {
		filename := os.Args[i]
		fileBytes, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println("Error occured while reading file:")
			fmt.Println(err)
			os.Exit(0)
		}
		FileIndexing(database, fileBytes, filename)
	}

	fmt.Println("Enter your phrase: ")
	var userStr string
	fmt.Scanln(&userStr)

	userTokens := strings.Split(userStr, " ")
	resultIdx := GetResults(database, userTokens)
	for i := 0; i < len(resultIdx); i++ {
		fmt.Println(resultIdx[i].Filename, "; совпадений -", resultIdx[i].Weight)
	}
}
