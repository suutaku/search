package search

import (
	"encoding/json"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/document"
	"time"
)

func getBleveDocsFromSearchResults(results *bleve.SearchResult, index bleve.Index) []map[string]interface{} {
	docs := make([]map[string]interface{}, 0)
	for _, val := range results.Hits {
		id := val.ID
		doc1, _ := index.Document(id)
		rv := struct {
			ID     string                 `json:"id"`
			Fields map[string]interface{} `json:"fields"`
		}{
			ID:     id,
			Fields: map[string]interface{}{},
		}
		doc := doc1.(*document.Document)
		for _, field := range doc.Fields {
			var newval interface{}
			switch field := field.(type) {
			case *document.TextField:
				newval = string(field.Value())
			case *document.NumericField:
				n, err := field.Number()
				if err == nil {
					newval = n
				}
			case *document.DateTimeField:
				d, err := field.DateTime()
				if err == nil {
					newval = d.Format(time.RFC3339Nano)
				}
			}
			existing, existed := rv.Fields[field.Name()]
			if existed {
				switch existing := existing.(type) {
				case []interface{}:
					rv.Fields[field.Name()] = append(existing, newval)
				case interface{}:
					arr := make([]interface{}, 2)
					arr[0] = existing
					arr[1] = newval
					rv.Fields[field.Name()] = arr
				}
			} else {
				rv.Fields[field.Name()] = newval
			}
		}

		j2, _ := json.Marshal(rv.Fields)
		tmp := make(map[string]interface{}, 0)
		json.Unmarshal(j2, &tmp)
		docs = append(docs, tmp)
	}
	return docs
}
