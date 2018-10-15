package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	. "./index"
)

type resStruct struct {
	filename string
	weight   int
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("File list doesn't specified")
		fmt.Println("Use:")
		fmt.Println(os.Args[0] + " [file1 file2 ...]")
		os.Exit(0)
	}
	var database map[string]Index
	database = make(map[string]Index, 10)

	for i := 1; i < len(os.Args); i++ {
		err := FileIndexing(database, os.Args[i])
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}

	var userStr string

	fmt.Println("Enter your phrase: ")
	fmt.Scanln(&userStr)

	var resultIdx []resStruct
	userTokens := strings.Split(userStr, " ")
	for i := 0; i < len(userTokens); i++ {
		val, ok := database[userTokens[i]]
		if ok {
			for j := 0; j < len(val.Files); j++ {
				var isExists bool
				for k := 0; k < len(resultIdx); k++ {
					if resultIdx[k].filename == val.Files[j].Filename {
						resultIdx[k].weight += val.Files[j].Weight
						isExists = true
					}
				}
				if !isExists {
					tmp := resStruct{val.Files[j].Filename, val.Files[j].Weight}
					resultIdx = append(resultIdx, tmp)
					sort.SliceStable(resultIdx, func(i, j int) bool { return resultIdx[i].weight > resultIdx[j].weight })
				}
			}
		}
	}
	for i := 0; i < len(resultIdx); i++ {
		fmt.Println(resultIdx[i].filename, "; совпадений -", resultIdx[i].weight)
	}
}
