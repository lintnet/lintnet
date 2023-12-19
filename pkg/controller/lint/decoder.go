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
	Text     string
	JSON     []byte
	FilePath string
	FileType string
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
	inputB, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("marshal input as JSON: %w", err)
	}
	return &Data{
		Text:     string(b),
		FilePath: filePath,
		FileType: fileType,
		JSON:     inputB,
	}, nil
}
