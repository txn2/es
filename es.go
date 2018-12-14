package es

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin/json"
	"go.uber.org/zap"
)

type EsObj map[string]interface{}

type Result struct {
	Index  string `json:"_index"`
	Type   string `json:"_type"`
	Id     string `json:"_id"`
	Source EsObj  `json:"_source"`
}

type IndexTemplate struct {
	Name     string
	Template EsObj
}

type Config struct {
	Log           *zap.Logger
	HttpClient    *http.Client
	ElasticServer string
}

type Client struct {
	Config
}

// CreateClient returns an Elasticsearch client object
func CreateClient(cfg Config) *Client {
	client := &Client{Config: cfg}

	if client.ElasticServer == "" {
		client.ElasticServer = "http://elasticsearch:9200"
	}

	return client
}

// Get uses HTTP Get method to retrieve data from elasticsearch
func (es *Client) Get(url string) (int, []byte, error) {
	return es.req(http.MethodGet, url, []byte{})
}

// Put uses HTTP Put method to send data to elasticsearch
func (es *Client) Put(url string, data []byte) (int, []byte, error) {
	return es.req(http.MethodPut, url, data)
}

func (es *Client) PutObj(url string, dataObj interface{}) (int, []byte, error) {
	data, err := json.Marshal(dataObj)
	if err != nil {
		es.Log.Error("Error marshaling object to json.", zap.Error(err))
		return 0, nil, err
	}

	return es.req(http.MethodPut, url, data)
}

func (es *Client) req(method string, url string, data []byte) (int, []byte, error) {
	fq := fmt.Sprintf("%s/%s", es.ElasticServer, url)

	req, err := http.NewRequest(method, fq, bytes.NewBuffer(data))
	if err != nil {
		return 0, nil, err
	}

	if method == http.MethodPut {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := es.HttpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, body, nil
}
