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

// Obj represents any structure
type Obj map[string]interface{}

// Result
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
	SeqNo       int                    `json:"_seq_no"`
	PrimaryTerm int                    `json:"_primary_term"`
	Source      map[string]interface{} `json:"_source"`
	Error       string                 `json:"error"`
	Status      int                    `json:"status"`
}

// HitsMeta
type HitsMeta struct {
	Total    int      `json:"total"`
	MaxScore float64  `json:"max_score"`
	Hits     []Result `json:"hits"`
}

// SearchResults
type SearchResults struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits   HitsMeta `json:"hits"`
	Error  string   `json:"error"`
	Status int      `json:"status"`
}

// IndexTemplate
type IndexTemplate struct {
	Name     string
	Template Obj
}
