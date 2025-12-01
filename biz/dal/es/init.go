package es

import (
	"github.com/olivere/elastic/v7"
	"judgeMore/pkg/base/client"
)

var els *elastic.Client

func Init() {
	var err error
	els, err = client.NewEsClient()
	if err != nil {
		panic(err)
	}
}
