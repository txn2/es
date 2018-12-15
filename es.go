package es

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

type Obj map[string]interface{}

type Result struct {
	Index      string `json:"_index"`
	Type       string `json:"_type"`
	Id         string `json:"_id"`
	Version    int    `json:"_version"`
	ResultType string `json:"result"`
	Found      bool   `json:"found"`
	Shards     struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	SeqNo       int    `json:"_seq_no"`
	PrimaryTerm int    `json:"_primary_term"`
	Source      Obj    `json:"_source"`
	Error       string `json:"error"`
	Status      int    `json:"status"`
}

type SearchResults struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total    int      `json:"total"`
		MaxScore int      `json:"max_score"`
		Hits     []Result `json:"hits"`
	} `json:"hits"`
	Error  string `json:"error"`
	Status int    `json:"status"`
}

type IndexTemplate struct {
	Name     string
	Template Obj
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
// returns a byte array that may be used to unmarshal into a
// specific type depending on the returned code.
func (es *Client) Get(url string) (int, []byte, error) {
	return es.req(http.MethodGet, url, []byte{})
}

// Put uses HTTP Put method to send data to elasticsearch
func (es *Client) Put(url string, data []byte) (int, Result, error) {
	return es.reqRes(url, data, http.MethodPut)
}

// Post uses HTTP POST method to send data to elasticsearch
func (es *Client) Post(url string, data []byte) (int, Result, error) {
	return es.reqRes(url, data, http.MethodPost)
}

// PutObj Marshals and object into json and PUTs it to Elasticsearch
func (es *Client) PutObj(url string, dataObj interface{}) (int, Result, error) {
	data, err := json.Marshal(dataObj)
	if err != nil {
		es.Log.Error("Error marshaling object to json.", zap.Error(err))
		return 0, Result{}, err
	}

	return es.Put(url, data)
}

// PostObj Marshals and object into json and POSTs it to Elasticsearch
func (es *Client) PostObj(url string, dataObj interface{}) (int, Result, error) {
	data, err := json.Marshal(dataObj)
	if err != nil {
		es.Log.Error("Error marshaling object to json.", zap.Error(err))
		return 0, Result{}, err
	}

	return es.Post(url, data)
}

// PostObjUnmarshal Unmarshals results to retObj, likely
// a overridden es.SearchResults struct
func (es *Client) PostObjUnmarshal(url string, dataObj interface{}, retObj interface{}) (int, error) {
	data, err := json.Marshal(dataObj)
	if err != nil {
		es.Log.Error("Error marshaling object to json.", zap.Error(err))
		return 0, err
	}

	code, res, err := es.req(http.MethodPost, url, data)
	if err != nil {
		return code, err
	}

	err = json.Unmarshal(res, retObj)
	if err != nil {
		es.Log.Error("Error unmarshaling result object.", zap.Error(err))
		return 0, err
	}

	return code, err
}

func (es *Client) reqRes(url string, data []byte, method string) (int, Result, error) {
	resObj := Result{}

	code, res, err := es.req(method, url, data)
	if err != nil {
		return code, resObj, err
	}

	err = json.Unmarshal(res, &resObj)
	if err != nil {
		es.Log.Error("Error unmarshaling result object.", zap.Error(err))
		return 0, resObj, err
	}

	return code, resObj, nil
}

func (es *Client) req(method string, url string, data []byte) (int, []byte, error) {
	fq := fmt.Sprintf("%s/%s", es.ElasticServer, url)

	req, err := http.NewRequest(method, fq, bytes.NewBuffer(data))
	if err != nil {
		return 0, nil, err
	}

	if method == http.MethodPut || method == http.MethodPost {
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
