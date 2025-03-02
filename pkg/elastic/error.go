package elastic

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

var (
	ErrNotFound = errors.New("not found")
	ErrInvalid  = errors.New("invalid")
)

type Err struct {
	Type   string     `json:"type"`
	Reason string     `json:"reason"`
	Causes []ErrCause `json:"root_cause"`
}

type ErrCause struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
	Index  string `json:"index"`
}

type Error struct {
	Err Err `json:"error"`
}

func isErr(res *esapi.Response) error {
	if !res.IsError() {
		return nil
	}

	var err error

	switch res.StatusCode {
	case http.StatusNotFound:
		err = ErrNotFound

	default:
		err = fmt.Errorf("%w, status: %d", ErrInvalid, res.StatusCode)
	}

	if res.Body == nil {
		return err
	}

	var e Error
	if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		return fmt.Errorf("decode response body: %w", err)
	}

	errStrs := []string{e.Err.Type + "/" + e.Err.Reason}

	for _, cause := range e.Err.Causes {
		errStrs = append(errStrs, fmt.Sprintf("%s/%s", cause.Index, cause.Reason))
	}

	return fmt.Errorf("%w: %s", err, strings.Join(errStrs, ", "))
}
