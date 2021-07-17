package search

import (
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"os"
	"strings"
)

type Search struct {
	index bleve.Index
}

func NewSearch(path string) *Search {
	var index bleve.Index
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(path, " not exists, create new index db")
		mapping := bleve.NewIndexMapping()
		idx, err := bleve.New(path, mapping)
		if err != nil {
			panic(err)
		}
		index = idx
	} else {
		idx, err := bleve.Open(path)
		if err != nil {
			panic(err)
		}
		index = idx
	}
	return &Search{
		index: index,
	}
}

func (search *Search) doSearch(key string) []map[string]interface{} {
	query := bleve.NewQueryStringQuery(key)
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, _ := search.index.Search(searchRequest)
	return getBleveDocsFromSearchResults(searchResult, search.index)
}

// A naive implementation of search.
// TODO consider search score and resort
func (search *Search) Search(key string) []map[string]interface{} {
	preResult := make([]map[string]interface{}, 0)
	result := make([]map[string]interface{}, 0)
	resultMap := make(map[string]int, 0)

	// split key with space
	keys := strings.Split(key, " ")
	for i := 0; i < len(keys); i++ {
		res := search.doSearch(keys[i])
		for j := 0; j < len(res); j++ {
			resultMap[res[j]["Id"].(string)]++
		}
		preResult = append(preResult, res...)
	}
	// only take docs included all keys result
	for i := 0; i < len(preResult); i++ {
		if resultMap[preResult[i]["Id"].(string)] >= len(keys) {
			result = append(result, preResult[i])
		}
	}
	fmt.Println(resultMap, len(keys))
	if len(result) == 0 {
		return nil
	}
	return result
}

func (search *Search) CreateIndex(content map[string]interface{}) error {
	if content["Id"] == nil {
		return fmt.Errorf("content don't hanve \"Id\" filed.")
	}
	return search.index.Index(content["Id"].(string), content)
}
