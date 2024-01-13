package lint

import (
	"fmt"

	"github.com/lintnet/lintnet/pkg/encoding"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

func (c *Controller) parseDataFile(filePath *Path) (*TopLevelArgment, error) {
	parser := &DataFileParser{
		fs: c.fs,
	}
	return parser.Parse(filePath)
}

type Data struct {
	Text     string `json:"text"`
	Value    any    `json:"value"`
	FilePath string `json:"file_path"`
	FileType string `json:"file_type"`
	JSON     []byte `json:"-"`
}

type TopLevelArgment struct {
	Data         *Data          `json:"data,omitempty"`
	CombinedData []*Data        `json:"combined_data,omitempty"`
	Config       map[string]any `json:"config"`
}

type DataFileParser struct {
	fs afero.Fs
}

func (dp *DataFileParser) Parse(filePath *Path) (*TopLevelArgment, error) {
	unmarshaler, fileType, err := encoding.NewUnmarshaler(filePath.Abs)
	if err != nil {
		return nil, logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
			"file_path": filePath.Raw,
		})
	}
	b, err := afero.ReadFile(dp.fs, filePath.Abs)
	if err != nil {
		return nil, fmt.Errorf("read a file: %w", err)
	}
	input, err := unmarshaler.Unmarshal(b)
	if err != nil {
		return nil, fmt.Errorf("decode a file: %w", err)
	}
	return &TopLevelArgment{
		Data: &Data{
			Text:     string(b),
			FilePath: filePath.Raw,
			FileType: fileType,
			Value:    input,
		},
	}, nil
}
