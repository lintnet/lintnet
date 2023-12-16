package lint

import (
	"encoding/csv"
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
		Decode() (interface{}, error)
	}
)

type csvDecoder struct {
	reader *csv.Reader
}

func newCSVDecoder(r io.Reader) decoder {
	return &csvDecoder{
		reader: csv.NewReader(r),
	}
}

func (c *csvDecoder) Decode() (interface{}, error) {
	records, err := c.reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parse a file as CSV: %w", err)
	}
	return records, nil
}

type jsonDecoder struct {
	decoder *json.Decoder
}

func (d *jsonDecoder) Decode() (interface{}, error) {
	var dest interface{}
	if err := d.decoder.Decode(&dest); err != nil {
		return nil, fmt.Errorf("parse a file as JSON: %w", err)
	}
	return dest, nil
}

func newJSONDecoder(r io.Reader) decoder {
	return &jsonDecoder{
		decoder: json.NewDecoder(r),
	}
}

type yamlDecoder struct {
	decoder *yaml.Decoder
}

func (d *yamlDecoder) Decode() (interface{}, error) {
	var dest interface{}
	if err := d.decoder.Decode(&dest); err != nil {
		return nil, fmt.Errorf("parse a file as YAML: %w", err)
	}
	return dest, nil
}

func newYAMLDecoder(r io.Reader) decoder {
	return &yamlDecoder{
		decoder: yaml.NewDecoder(r),
	}
}

func getNewDecoder(fileName string) (NewDecoder, string, error) {
	ext := filepath.Ext(fileName)
	switch ext {
	case ".csv":
		return newCSVDecoder, "csv", nil
	case ".json":
		return newJSONDecoder, "json", nil
	case ".toml":
		return newTOMLDecoder, "toml", nil
	case ".yml", ".yaml":
		return newYAMLDecoder, "yaml", nil
	default:
		return nil, "", errors.New("this format is unsupported")
	}
}

type tomlDecoder struct {
	decoder *toml.Decoder
}

func (d *tomlDecoder) Decode() (interface{}, error) {
	var v interface{}
	_, err := d.decoder.Decode(&v)
	if err != nil {
		return nil, fmt.Errorf("parse a file as TOML: %w", err)
	}
	return v, nil
}

func newTOMLDecoder(r io.Reader) decoder {
	return &tomlDecoder{
		decoder: toml.NewDecoder(r),
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
	input, err := newDecoder(f).Decode()
	if err != nil {
		return nil, "", fmt.Errorf("decode a file: %w", err)
	}
	inputB, err := json.Marshal(input)
	if err != nil {
		return nil, "", fmt.Errorf("marshal input as JSON: %w", err)
	}
	return inputB, fileType, nil
}
