package testcmd_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lintnet/lintnet/pkg/controller/testcmd"
)

func TestDataFile_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		data string
		exp  *testcmd.DataFile
	}{
		{
			name: "string",
			data: `"foo.yaml"`,
			exp: &testcmd.DataFile{
				Path:     "foo.yaml",
				FakePath: "foo.yaml",
			},
		},
		{
			name: "non string",
			data: `{
	"path": "foo.yaml",
	"fake_path": ".github/workflows/foo.yaml"
}`,
			exp: &testcmd.DataFile{
				Path:     "foo.yaml",
				FakePath: ".github/workflows/foo.yaml",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dataFile := &testcmd.DataFile{}
			if err := json.Unmarshal([]byte(tt.data), dataFile); err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.exp, dataFile); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
