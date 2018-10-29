package index

import (
	"io/ioutil"
	"sort"
	"strings"
	"sync"
)

type Database map[string]Index

//Index хранит слово и структуру extFiles, относящуюся к нему
type Index struct {
	Word  string
	Files []ExtFiles
}

// ExtFiles хранит имя файла и вес этого файла
//(для каждого слова - количество встреч в данном файле)
type ExtFiles struct {
	Filename string
	Weight   int
}

// ResStruct хранит уже обработанные данные
type ResStruct struct {
	Filename string
	Weight   int
}

// MergeIt копит результаты работы всех тредов в одну базу
func MergeIt(database Database, chWork chan Index, chError chan error, filesCnt int, waiter *sync.WaitGroup) {
	defer waiter.Done()
	successCnt := 0
MergeLoop:
	for {
		select {
		case err := <-chError:
			if err == nil {
				successCnt++
				if successCnt == filesCnt {
					break MergeLoop
				}
			}
		case val := <-chWork:
			_, ok := database[val.Word]
			if ok {
				index := database[val.Word]
				index.Files = append(index.Files, val.Files[0])
				database[val.Word] = index
			} else {
				database[val.Word] = val
			}
		}
	}
}

// ReadIt читает файл и создает базу по нему
func ReadIt(filename string, ch chan Index, chError chan error) {
	myDatabase := make(Database, 1)
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		chError <- err
	} else {
		FileIndexing(myDatabase, fileBytes, filename)
		for _, val := range myDatabase {
			ch <- val
		}
		chError <- nil
	}
}

func cleanWord(in string) string {
	in = strings.ToLower(in)
	in = strings.TrimFunc(in, func(c rune) bool {
		return (c < 48 || c > 57) && (c < 97 || c > 122) && c != 45
	})
	return string(in)
}

// FileIndexing обновляет стркутуру обратного индекса в файле filename
func FileIndexing(arrayIndexes Database, inputContent []byte, filename string) {
	str := string(inputContent)
	words := strings.Fields(str)
	for i := 0; i < len(words); i++ {
		word := cleanWord(words[i])
		_, ok := arrayIndexes[word]
		if !ok {
			newWordIdx := Index{Word: word}
			newFile := ExtFiles{filename, 1}
			newWordIdx.Files = append(newWordIdx.Files, newFile)
			arrayIndexes[word] = newWordIdx
		} else {
			var isExist bool
			for j := 0; j < len(arrayIndexes[word].Files); j++ {
				if arrayIndexes[word].Files[j].Filename == filename {
					arrayIndexes[word].Files[j].Weight++
					isExist = true
					sort.SliceStable(arrayIndexes[word].Files, func(i, j int) bool {
						return arrayIndexes[word].Files[i].Weight > arrayIndexes[word].Files[j].Weight
					})
				}
			}
			if !isExist {
				newFile := ExtFiles{filename, 1}
				tmp := arrayIndexes[word]
				tmp.Files = append(tmp.Files, newFile)
				arrayIndexes[word] = tmp
			}
		}
	}
}

// GetResults обрабатывает базу данных и возвращает структуру результата
func GetResults(arrayIndexes Database, userTokens []string) []ResStruct {
	var resultIdx []ResStruct
	for i := 0; i < len(userTokens); i++ {
		userWord := cleanWord(userTokens[i])
		val, ok := arrayIndexes[userWord]
		if ok {
			for j := 0; j < len(val.Files); j++ {
				isExists := false
				for k := 0; k < len(resultIdx); k++ {
					if resultIdx[k].Filename == val.Files[j].Filename {
						resultIdx[k].Weight += val.Files[j].Weight
						isExists = true
					}
				}
				if !isExists {
					tmp := ResStruct{val.Files[j].Filename, val.Files[j].Weight}
					resultIdx = append(resultIdx, tmp)
				}
			}
		}
	}
	sort.SliceStable(resultIdx, func(i, j int) bool {
		return resultIdx[i].Weight > resultIdx[j].Weight
	})
	return resultIdx
}
