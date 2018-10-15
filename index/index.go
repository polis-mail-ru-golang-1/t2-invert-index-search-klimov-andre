package index

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

//Index хранит слово и структуру extFiles, относящуюся к нему
type Index struct {
	Word  string
	Files []ExtFiles
}

//ExtFiles хранит имя файла и вес этого файла
//(для каждого слова - количество встреч в данном файле)
type ExtFiles struct {
	Filename string
	Weight   int
}

//FileIndexing обновляет стркутуру обратного индекса в файле filename
func FileIndexing(arrayIndexes map[string]Index, filename string) error {
	myBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error occured while reading file:")
		fmt.Println(err)
		return err
	} else {
		str := string(myBytes)
		words := strings.Split(str, " ")
		for i := 0; i < len(words); i++ {
			word := words[i]
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
						sort.SliceStable(arrayIndexes[word].Files, func(i, j int) bool { return arrayIndexes[word].Files[i].Weight > arrayIndexes[word].Files[j].Weight })
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
	return nil
}
