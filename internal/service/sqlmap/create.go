package sqlmap

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"diploma-project/internal/model"
)

func (s *Service) Create(ctx context.Context, cmd model.SQLMapCommand) error {
	report, err := createReport(cmd.URL)
	if err != nil {
		return fmt.Errorf("exec sqlmap: %w", err)
	}

	cmd.Report = report

	err = s.r.Create(ctx, cmd)
	if err != nil {
		return fmt.Errorf("create sqlmap report: %w", err)
	}

	return nil
}

func createReport(url string) ([]byte, error) {
	cmd := exec.Command(
		"sqlmap",
		"-u",
		url,
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("sqlmap running: %w", err)
	}

	return out.Bytes(), nil
}
