package shell

import (
	"reflect"
	"regexp"
	"testing"
)

func TestGetOutFilePath(t *testing.T) {
	t.Parallel()

	path, err := GetOutFilePath()
	if err != nil {
		t.Error(err)
	}

	if match, _ := regexp.MatchString(`tf-script-output-.+\.json$`, path); !match {
		t.Error("output file path is not valid")
	}
}

func TestGetErrorFilePath(t *testing.T) {
	t.Parallel()

	path, err := GetErrorFilePath()
	if err != nil {
		t.Error(err)
	}

	if match, _ := regexp.MatchString(`tf-script-error-.+$`, path); !match {
		t.Error("error file path is not valid")
	}
}

func TestReadJSON(t *testing.T) {
	t.Parallel()

	for _, d := range []struct {
		testName string
		path     string
		expected any
		hasErr   bool
	}{
		{
			testName: "missing_file",
			path:     "testdata/missing.json",
			expected: nil,
			hasErr:   true,
		},
		{
			testName: "empty_file",
			path:     "testdata/empty.json",
			expected: nil,
			hasErr:   true,
		},
		{
			testName: "object",
			path:     "testdata/object.json",
			expected: map[string]any{"foo": "bar"},
			hasErr:   false,
		},
		{
			testName: "nested_object",
			path:     "testdata/nested-object.json",
			expected: map[string]any{"foo": map[string]any{"bar": "baz"}},
			hasErr:   false,
		},
		{
			testName: "array",
			path:     "testdata/array.json",
			expected: []any{map[string]any{"foo": "bar"}, map[string]any{"foo": "baz"}},
			hasErr:   false,
		},
	} {
		t.Run(d.testName, func(t *testing.T) {
			t.Parallel()

			raw, err := ReadJSON(d.path)

			if !reflect.DeepEqual(raw, d.expected) {
				t.Errorf("expected %v, got %v", d.expected, raw)
			}

			hasErr := err != nil
			if hasErr != d.hasErr {
				t.Errorf("unexpected error state")
			}
		})
	}
}
