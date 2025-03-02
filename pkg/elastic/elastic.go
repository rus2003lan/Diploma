package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

const (
	timeoutSeconds = 360
	boolCount      = 10
)

type ElClient struct {
	es      *elasticsearch.Client
	index   string
	elType  string
	timeout time.Duration
}

type indexConfig struct {
	Settings map[string]any `json:"settings,omitempty"`
	Mappings map[string]any `json:"mappings"`
}

type StartupIndexConfig struct {
	Index         string `yaml:"index" validate:"required"`
	UpdateMapping bool   `yaml:"updateMapping"`
	CreateMapping bool   `yaml:"createMapping" validate:"required"`
	MappingPath   string `yaml:"mappingPath" validate:"required"`
	Alias         string `yaml:"alias"`
	Shards        int    `yaml:"shards"`
	Replics       int    `yaml:"replics"`
}

type Cfg struct {
	Index       string
	MappingPath string
	Type        string
	Hosts       []string
	Login       string
	Password    string
	MaxRetries  int
	Transport   http.RoundTripper
}

func NewElastic(cfg Cfg) (*ElClient, error) {
	var esCfg elasticsearch.Config
	//nolint:exhaustruct
	if cfg.Transport == nil {
		esCfg = elasticsearch.Config{
			Addresses:            cfg.Hosts,
			Username:             cfg.Login,
			Password:             cfg.Password,
			EnableRetryOnTimeout: true,
			RetryOnStatus:        []int{429, 502, 503, 504},
			MaxRetries:           cfg.MaxRetries,
			Transport:            http.DefaultTransport,
		}
	} else {
		esCfg = elasticsearch.Config{
			Transport: cfg.Transport,
		}
	}

	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("new ec client: %w", err)
	}

	return &ElClient{es: es, index: cfg.Index, elType: cfg.Type, timeout: time.Second * timeoutSeconds}, nil
}

func (e *ElClient) GetDoc(ctx context.Context, id string) (*Hit, error) {
	resp, err := e.es.Get(
		e.index,
		id,
		e.es.Get.WithContext(ctx),
		e.es.Get.WithOpaqueID("TAP"),
	)
	if err != nil {
		return nil, fmt.Errorf("get by id %s: %w", id, err)
	}

	defer resp.Body.Close()

	if err := isErr(resp); err != nil {
		return nil, err
	}

	var hit Hit
	if err := json.NewDecoder(resp.Body).Decode(&hit); err != nil {
		return nil, fmt.Errorf("decode to hit: %w", err)
	}

	return &hit, nil
}

func (e *ElClient) doSearch(q map[string]any, opts []func(*esapi.SearchRequest)) (*Answer, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(q); err != nil {
		return nil, fmt.Errorf("error encoding query: %w", err)
	}

	opts = append(opts, e.es.Search.WithBody(&buf))

	res, err := e.es.Search(opts...)
	if err != nil {
		return nil, fmt.Errorf("search es: %w", err)
	}

	defer res.Body.Close()

	if err := isErr(res); err != nil {
		return nil, err
	}

	var a Answer
	if err := json.NewDecoder(res.Body).Decode(&a); err != nil {
		return nil, fmt.Errorf("decode to es answer: %w", err)
	}

	return &a, nil
}

func (e *ElClient) Search(ctx context.Context) (*FetchResponse, error) {
	opts := []func(*esapi.SearchRequest){
		e.es.Search.WithContext(ctx),
		e.es.Search.WithIndex(e.index),
	}

	// if ctx != nil {
	//	opts = append(opts, e.es.Search.WithContext(ctx))
	//}

	a, err := e.doSearch(_map{"query": _map{"match_all": _map{}}}, opts)
	if err != nil {
		return nil, err
	}

	res := FetchResponse{
		Hits:  a.Hits.Hits,
		Total: a.Hits.Total.Value,
	}

	return &res, nil
}

func (e *ElClient) Update(ctx context.Context, id string, body []byte) error {
	//nolint:exhaustruct
	req := esapi.UpdateRequest{
		Index:      e.index,
		DocumentID: id,
		Body:       bytes.NewReader(body),
		Refresh:    "wait_for",
	}

	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	res, err := req.Do(ctx, e.es)
	if err != nil {
		return fmt.Errorf("requst do: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 || res.IsError() {
		//nolint:err113
		return fmt.Errorf("document %s upsert response: %s with code: %d", id, res.String(), res.StatusCode)
	}

	return nil
}

func (e *ElClient) CreateIndex(
	ctx context.Context,
	cfg *StartupIndexConfig,
) error {
	if cfg.MappingPath == "" || !cfg.CreateMapping {
		return nil
	}

	exists, err := e.checkIndexExists(ctx, cfg)
	if err != nil {
		return fmt.Errorf("check index exists: %w", err)
	}

	if exists {
		return nil
	}

	indexconfig, err := readIndexConfigWithDefaults(cfg)
	if err != nil {
		return fmt.Errorf("read index config with defaults: %w", err)
	}

	var buf bytes.Buffer

	err = json.NewEncoder(&buf).Encode(indexconfig)
	if err != nil {
		return fmt.Errorf("encode index config: %w", err)
	}

	resp, err := e.es.Indices.Create(
		cfg.Index,
		e.es.Indices.Create.WithBody(&buf),
		e.es.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("create index: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if err := isErr(resp); err != nil {
		return fmt.Errorf("error in response: %w", err)
	}

	return nil
}

func (e *ElClient) checkIndexExists(
	ctx context.Context,
	cfg *StartupIndexConfig,
) (bool, error) {
	get, err := e.es.Indices.Get(
		[]string{cfg.Index},
		e.es.Indices.Get.WithContext(ctx),
	)
	if err != nil {
		return false, fmt.Errorf("get index: %w", err)
	}

	if get.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}

func readIndexConfigWithDefaults(cfg *StartupIndexConfig) (*indexConfig, error) {
	f, err := os.OpenFile(cfg.MappingPath, os.O_RDONLY, 0o600) //nolint:gomnd, mnd
	if err != nil {
		return nil, fmt.Errorf("open mapping file %s: %w", cfg.MappingPath, err)
	}

	m := new(indexConfig)

	err = json.NewDecoder(f).Decode(m)
	if err != nil {
		return nil, fmt.Errorf("decode mapping file %s: %w", cfg.MappingPath, err)
	}

	if m.Settings == nil {
		m.Settings = make(map[string]any)
	}

	if cfg.Shards == 0 {
		cfg.Shards = 1
	}

	m.Settings["number_of_shards"] = cfg.Shards

	return m, nil
}

func (e *ElClient) Ping(ctx context.Context) error {
	esResp, err := e.es.Ping(e.es.Ping.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	if err = isErr(esResp); err != nil {
		return fmt.Errorf("error in response: %w", err)
	}

	return nil
}
