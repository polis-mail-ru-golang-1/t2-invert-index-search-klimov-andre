package index

import (
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

//ResStruct хранит уже обработанные данные
type ResStruct struct {
	Filename string
	Weight   int
}

//FileIndexing обновляет стркутуру обратного индекса в файле filename
func FileIndexing(arrayIndexes map[string]Index, inputContent []byte, filename string) {
	str := string(inputContent)
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

// GetResults обрабатывает базу данных и возвращает структуру результата
func GetResults(arrayIndexes map[string]Index, userTokens []string) []ResStruct {
	var resultIdx []ResStruct
	for i := 0; i < len(userTokens); i++ {
		val, ok := arrayIndexes[userTokens[i]]
		if ok {
			for j := 0; j < len(val.Files); j++ {
				var isExists bool
				for k := 0; k < len(resultIdx); k++ {
					if resultIdx[k].Filename == val.Files[j].Filename {
						resultIdx[k].Weight += val.Files[j].Weight
						isExists = true
					}
				}
				if !isExists {
					tmp := ResStruct{val.Files[j].Filename, val.Files[j].Weight}
					resultIdx = append(resultIdx, tmp)
					sort.SliceStable(resultIdx, func(i, j int) bool { return resultIdx[i].Weight > resultIdx[j].Weight })
				}
			}
		}
	}
	return resultIdx
}
