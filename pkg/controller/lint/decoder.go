package lint

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/lintnet/pkg/encoding"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

func (c *Controller) parse(filePath string) ([]byte, string, error) {
	newDecoder, fileType, err := encoding.GetNewDecoder(filePath)
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
