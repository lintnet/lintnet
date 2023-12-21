package lint

import (
	"encoding/json"
	"fmt"

	"github.com/lintnet/lintnet/pkg/encoding"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type Data struct {
	Text     string      `json:"text"`
	Value    interface{} `json:"value"`
	FilePath string      `json:"file_path"`
	FileType string      `json:"file_type"`
	JSON     []byte      `json:"-"`
}

func (c *Controller) parse(filePath string) (*Data, error) {
	unmarshaler, fileType, err := encoding.NewUnmarshaler(filePath)
	if err != nil {
		return nil, logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
			"file_path": filePath,
		})
	}
	b, err := afero.ReadFile(c.fs, filePath)
	if err != nil {
		return nil, fmt.Errorf("read a file: %w", err)
	}
	input, err := unmarshaler.Unmarshal(b)
	if err != nil {
		return nil, fmt.Errorf("decode a file: %w", err)
	}
	data := &Data{
		Text:     string(b),
		FilePath: filePath,
		FileType: fileType,
		Value:    input,
	}
	dataB, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal data as JSON: %w", err)
	}
	data.JSON = dataB
	return data, nil
}
