package lint

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"gopkg.in/yaml.v3"
)

type (
	NewDecoder func(io.Reader) decoder
	decoder    interface {
		Decode(dest interface{}) error
	}
)

func getNewDecoder(fileName string) (NewDecoder, string, error) {
	ext := filepath.Ext(fileName)
	switch ext {
	case ".json":
		return func(r io.Reader) decoder {
			return json.NewDecoder(r)
		}, "json", nil
	case ".yml", ".yaml":
		return func(r io.Reader) decoder {
			return yaml.NewDecoder(r)
		}, "yaml", nil
	case ".toml":
		return func(r io.Reader) decoder {
			return &tomlDecoder{
				decoder: toml.NewDecoder(r),
			}
		}, "toml", nil
	default:
		return nil, "", errors.New("lintnet supports linting only JSON or YAML")
	}
}

type tomlDecoder struct {
	decoder *toml.Decoder
}

func (d *tomlDecoder) Decode(v interface{}) error {
	_, err := d.decoder.Decode(v)
	return err //nolint:wrapcheck
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
