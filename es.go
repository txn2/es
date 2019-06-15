// Package es implements a simple Elasticsearch client.
package es

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

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

	defaultElasticServer := false
	defaultLogger := false
	defaultHttpClient := false

	// default elastic server if not provided
	if client.ElasticServer == "" {
		defaultElasticServer = true
		client.ElasticServer = "http://elasticsearch:9200"
	}

	// default zap logger if not provided
	if client.Log == nil {
		defaultLogger = true
		zapCfg := zap.NewDevelopmentConfig()

		logger, err := zapCfg.Build()
		if err != nil {
			fmt.Printf("Can not build logger: %s\n", err.Error())
			os.Exit(1)
		}

		client.Log = logger
	}

	// default HttpClient is none provided
	if client.HttpClient == nil {
		defaultHttpClient = true
		// Http Client Configuration for outbound connections
		netTransport := &http.Transport{
			MaxIdleConnsPerHost: 10,
			DialContext: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
		}

		client.HttpClient = &http.Client{
			Timeout:   time.Second * 60,
			Transport: netTransport,
		}

	}

	client.Log.Info("Created es client",
		zap.Bool("defaultHttpClient", defaultHttpClient),
		zap.Bool("defaultLogger", defaultLogger),
		zap.Bool("defaultElasticServer", defaultElasticServer),
		zap.String("ElasticServer", client.ElasticServer))

	return client
}

// Get uses HTTP Get method to retrieve data from elasticsearch
// returns a byte array that may be used to unmarshal into a
// specific type depending on the returned code.
func (es *Client) Get(url string) (int, []byte, error) {
	return es.req(http.MethodGet, url, []byte{})
}

// Put uses HTTP Put method to send data to elasticsearch
func (es *Client) Put(url string, data []byte) (int, Result, *ErrorResponse, error) {
	return es.reqRes(url, data, http.MethodPut)
}

// Post uses HTTP POST method to send data to elasticsearch
func (es *Client) Post(url string, data []byte) (int, Result, *ErrorResponse, error) {
	return es.reqRes(url, data, http.MethodPost)
}

// PutObj Marshals and object into json and PUTs it to Elasticsearch
func (es *Client) PutObj(url string, dataObj interface{}) (int, Result, *ErrorResponse, error) {
	data, err := json.Marshal(dataObj)
	if err != nil {
		es.Log.Error("Error marshaling object to json.", zap.Error(err))
		return 0, Result{}, nil, err
	}

	return es.Put(url, data)
}

// PostObj Marshals and object into json and POSTs it to Elasticsearch
func (es *Client) PostObj(url string, dataObj interface{}) (int, Result, *ErrorResponse, error) {
	data, err := json.Marshal(dataObj)
	if err != nil {
		es.Log.Error("Error marshaling object to json.", zap.Error(err))
		return 0, Result{}, nil, err
	}

	return es.Post(url, data)
}

// PostObjUnmarshal Unmarshals results to retObj, likely
// a overridden es.SearchResults struct
func (es *Client) PostObjUnmarshal(url string, dataObj interface{}, retObj interface{}) (int, *ErrorResponse, error) {
	data, err := json.Marshal(dataObj)
	if err != nil {
		es.Log.Error("Error marshaling object to json.", zap.Error(err))
		return 0, nil, err
	}

	code, res, err := es.req(http.MethodPost, url, data)
	if err != nil {
		return code, &ErrorResponse{Message: string(res)}, err
	}

	if code != 200 {
		return code, &ErrorResponse{Message: string(res)}, err
	}

	err = json.Unmarshal(res, retObj)
	if err != nil {
		es.Log.Error("Error unmarshaling result object.", zap.Error(err))
		return 0, nil, err
	}

	return code, nil, err
}

// SendEsMapping
func (es *Client) SendEsMapping(mapping IndexTemplate) (int, Result, *ErrorResponse, error) {

	es.Log.Info("Sending template",
		zap.String("type", "SendEsMapping"),
		zap.String("mapping", mapping.Name),
	)

	code, esResult, errorResonse, err := es.PutObj(fmt.Sprintf("_template/%s", mapping.Name), mapping.Template)
	if err != nil {
		es.Log.Error("SendEsMapping got error sending template", zap.Error(err))
		return code, esResult, errorResonse, err
	}

	if code != 200 {
		return code, esResult, errorResonse, errors.New("Error setting up " + mapping.Name + " template, got code " + string(code))
	}

	return code, esResult, nil, err
}

// reqRes
func (es *Client) reqRes(url string, data []byte, method string) (int, Result, *ErrorResponse, error) {
	resObj := Result{}

	code, res, err := es.req(method, url, data)
	if err != nil {
		return code, resObj, nil, err
	}

	if code != 200 {
		return code, resObj, &ErrorResponse{Message: string(res)}, err
	}

	err = json.Unmarshal(res, &resObj)
	if err != nil {
		es.Log.Error("reqRes error unmarshaling result object.",
			zap.ByteString("result", res),
			zap.Error(err),
		)
		resObj.Error = string(res)
		return 0, resObj, nil, err
	}

	return code, resObj, nil, nil
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
