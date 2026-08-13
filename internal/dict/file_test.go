package dict_test

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/dgf/tygo/internal/dict"
)

func writeTempFile(t *testing.T, name string, data string) string {
	t.Helper()

	file, err := os.CreateTemp(t.TempDir(), name)
	if err != nil {
		t.Error(err)
	}

	_, err = file.WriteString(data)
	if err != nil {
		t.Error(err)
	}

	return file.Name()
}

func TestLoadFile_ValidJSON(t *testing.T) {
	t.Parallel()

	name := "test-valid-dict.json"
	words := []string{"foo", "bar"}
	data := fmt.Sprintf(`{"words": ["%v"]}`, strings.Join(words, `", "`))
	filename := writeTempFile(t, name, data)

	dict, err := dict.LoadFile(filename)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(words, dict) {
		t.Errorf("expected valid dict: %q, got: %q", words, dict)
	}
}

func TestLoadFile_FailedAccess(t *testing.T) {
	t.Parallel()

	file, err := os.CreateTemp(t.TempDir(), "test-failed-access.json")
	if err != nil {
		t.Error(err)
	}

	name := file.Name()

	err = os.Remove(name)
	if err != nil {
		t.Error(err)
	}

	words, err := dict.LoadFile(name)
	if err == nil {
		t.Errorf("expected load error, got: %v", words)
	}
}

func TestLoadFile_InvalidJSON(t *testing.T) {
	t.Parallel()

	invalidDictJSON := `{"words": "foo bar"}`
	filename := writeTempFile(t, "test-invalid-dict.json", invalidDictJSON)

	words, err := dict.LoadFile(filename)
	if err == nil {
		t.Errorf("expect marshal error, got: %v", words)
	}
}
