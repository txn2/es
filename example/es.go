package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/txn2/es/v2"
	"go.uber.org/zap"
)

// TestObject
type TestObject struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func main() {
	zapCfg := zap.NewDevelopmentConfig()

	logger, err := zapCfg.Build()
	if err != nil {
		fmt.Printf("Can not build logger: %s\n", err.Error())
		os.Exit(1)
	}

	esClient := es.CreateClient(es.Config{
		Log:           logger,
		ElasticServer: "http://localhost:9200",
	})

	tstObj := TestObject{
		Name:        "test",
		Description: "This is a test",
	}

	code, res, errorResult, err := esClient.PostObj("es_test/_doc/test", tstObj)
	if err != nil {
		logger.Fatal("Error posting test object", zap.Error(err))
		spew.Dump(errorResult)
	}

	if code != 200 {
		logger.Warn("Got non 200 from elastic.", zap.Int("code", code))
		spew.Dump(errorResult)
	} else {
		spew.Dump(res)
	}

}
