package es

// Obj represents any structure
type Obj map[string]interface{}

// ErrorCause
type ErrorCause struct {
	RootCause    []ErrorCause `json:"root_cause"`
	Type         string       `json:"type"`
	Reason       string       `json:"reason"`
	ResourceType string       `json:"resource.type"`
	ResourceId   string       `json:"resource.id"`
	IndexUuid    string       `json:"index_uuid"`
	Index        string       `json:"index"`
}

// ErrorResponse
type ErrorResponse struct {
	Message string
}

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
