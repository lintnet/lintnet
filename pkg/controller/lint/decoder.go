package lint

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"gopkg.in/yaml.v3"
)

func getNewDecoder(fileName string) (NewDecoder, string, error) {
	switch {
	case strings.HasSuffix(fileName, ".json"):
		return func(r io.Reader) decoder {
			return json.NewDecoder(r)
		}, "json", nil
	case strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml"):
		return func(r io.Reader) decoder {
			return yaml.NewDecoder(r)
		}, "yaml", nil
	default:
		return nil, "", errors.New("lintnet supports linting only JSON or YAML")
	}
}

func (c *Controller) parse(filePath string) ([]byte, string, error) {
	newDecoder, fileType, err := getNewDecoder(filePath)
	if err != nil {
		return nil, "", logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
			"file_path": filePath,
		})
	}
	f, err := c.fs.Open(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("open a file: %w", err)
	}
	defer f.Close()
	var input interface{}
	if err := newDecoder(f).Decode(&input); err != nil {
		return nil, "", fmt.Errorf("decode a file: %w", err)
	}
	inputB, err := json.Marshal(input)
	if err != nil {
		return nil, "", fmt.Errorf("marshal input as JSON: %w", err)
	}
	return inputB, fileType, nil
}
