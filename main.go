package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	. "github.com/polis-mail-ru-golang-1/t2-invert-index-search-klimov-andre/index"
)

func main() {

	if len(os.Args) == 1 {
		fmt.Println("File list doesn't specified")
		fmt.Println("Use:")
		fmt.Println(os.Args[0] + " [file1 file2 ...]")
		os.Exit(0)
	}

	database := make(Database)
	chWork := make(chan Index)
	chExit := make(chan error)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go MergeIt(database, chWork, chExit, len(os.Args)-1, wg)

	for i := 1; i < len(os.Args); i++ {
		filename := os.Args[i]
		go ReadIt(filename, chWork, chExit)
	}

	wg.Wait()

	fmt.Println("Enter your phrase: ")

	reader := bufio.NewReader(os.Stdin)
	in, _ := reader.ReadString('\n')
	userTokens := strings.Fields(in)

	resultIdx := GetResults(database, userTokens)
	for i := 0; i < len(resultIdx); i++ {
		fmt.Println(resultIdx[i].Filename, "; совпадений -", resultIdx[i].Weight)
	}
}
