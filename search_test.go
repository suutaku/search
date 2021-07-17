package search

import (
	"github.com/stretchr/testify/assert"
	"github.com/suutaku/search"
	"testing"
)

func TestIndexAndSearch(t *testing.T) {
	sch := search.NewSearch("/tmp/search")
	content := make(map[string]interface{}, 0)
	content["Id"] = "2"
	content["Titile"] = "test john some"
	content["Url"] = "https://www.cotnetwork.com"
	content["Body"] = "include some words, like: hello, i'm ok. Fine"
	sch.CreateIndex(content)

	res10 := sch.Search("i'm")
	assert.NotNil(t, res10, "key with ' was not ok")

	res1 := sch.Search("test")
	assert.NotNil(t, res1)

	res11 := sch.Search("some")
	assert.NotNil(t, res11, "word \"some\" was not ok")

	res2 := sch.Search("test some")
	assert.NotNil(t, res2, "word \"some\" was not ok")

	res3 := sch.Search("include hello")
	assert.NotNil(t, res3)

	res4 := sch.Search("ok fine")
	assert.NotNil(t, res4)

	res5 := sch.Search("ok fine suu")
	assert.Nil(t, res5)

	res6 := sch.Search("test john")
	assert.NotNil(t, res6)
}
