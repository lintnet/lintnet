package encoding

import (
	"fmt"

	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

type DataFileParser struct {
	fs afero.Fs
}

func NewDataFileParser(fs afero.Fs) *DataFileParser {
	return &DataFileParser{
		fs: fs,
	}
}

func (dp *DataFileParser) Parse(filePath *domain.Path) (*domain.TopLevelArgument, error) {
	unmarshaler, fileType, err := NewUnmarshaler(filePath.Abs)
	if err != nil {
		return nil, slogerr.With(err, "file_path", filePath.Raw) //nolint:wrapcheck
	}
	b, err := afero.ReadFile(dp.fs, filePath.Abs)
	if err != nil {
		return nil, fmt.Errorf("read a file: %w", err)
	}
	input, err := unmarshaler.Unmarshal(b)
	if err != nil {
		return nil, fmt.Errorf("decode a file: %w", err)
	}

	return &domain.TopLevelArgument{
		Data: &domain.Data{
			Text:     string(b),
			FilePath: filePath.Raw,
			FileType: fileType,
			Value:    input,
		},
	}, nil
}
