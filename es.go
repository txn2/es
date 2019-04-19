/*
   Copyright 2019 TXN2

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package es

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

// Config
type Config struct {
	Log           *zap.Logger
	HttpClient    *http.Client
	ElasticServer string
}

// Client
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

// SendEsMapping
func (es *Client) SendEsMapping(mapping IndexTemplate) (int, Result, error) {

	es.Log.Info("Sending template",
		zap.String("type", "SendEsMapping"),
		zap.String("mapping", mapping.Name),
	)

	code, esResult, err := es.PutObj(fmt.Sprintf("_template/%s", mapping.Name), mapping.Template)
	if err != nil {
		es.Log.Error("Got error sending template", zap.Error(err))
		return code, esResult, err
	}

	if code != 200 {
		es.Log.Error("Got code", zap.Int("code", code), zap.String("EsResult", esResult.ResultType))
		return code, esResult, errors.New("Error setting up " + mapping.Name + " template, got code " + string(code))
	}

	return code, esResult, err
}

// reqRes
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

// req
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
