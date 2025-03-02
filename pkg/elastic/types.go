package elastic

import (
	"encoding/json"
)

type _map map[string]any

type ElMap _map

type Bool struct {
	Filter             []_map `json:"filter,omitempty"`
	Must               []_map `json:"must,omitempty"`
	Should             []_map `json:"should,omitempty"`
	MustNot            []_map `json:"must_not,omitempty"`
	MinimumShouldMatch int    `json:"minimum_should_match,omitempty"`
}

type Query struct {
	Bool     Bool `json:"bool,omitempty"`
	MatchAll bool `json:"-"`
}

type Bucket struct {
	Key   any `json:"key"`
	Count int `json:"doc_count"`
}

type Res struct {
	Buckets []Bucket `json:"buckets"`
}

type Uniq struct {
	Value int `json:"value"`
}

type Aggregation struct {
	Res  Res  `json:"res"`
	Uniq Uniq `json:"uniq"`
}

type Highlight struct {
	Content []string `json:"content"`
}

type Hit struct {
	ID        string          `json:"_id"`
	Source    json.RawMessage `json:"_source"`
	Index     string          `json:"_index"`
	Highlight Highlight       `json:"highlight"`
}

type Hits struct {
	Total Total `json:"total"`
	Hits  []Hit `json:"hits"`
}

type Answer struct {
	ScrollID string      `json:"_scroll_id"`
	Hits     Hits        `json:"hits"`
	Aggs     Aggregation `json:"aggregations"`
}

type Total struct {
	Value int `json:"value"`
}

type FetchResponse struct {
	Hits     []Hit
	Total    int
}

type Cause struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type BulkError struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
	Cause  Cause  `json:"caused_by"`
}

type Index struct {
	ID      string    `json:"_id"`
	Result  string    `json:"result"`
	Status  int       `json:"status"`
	BulkErr BulkError `json:"error"`
}

type Items []struct {
	Index Index `json:"index"`
}

type BulkResponse struct {
	Errors bool  `json:"errors"`
	Items  Items `json:"items"`
}
