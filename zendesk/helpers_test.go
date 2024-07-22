package zendesk

import (
	"io"
	"os"
	"testing"
)

func openJsonFile(t *testing.T, filePathAbs string, fileName string) []byte {
	t.Helper()
	jsonFile, err := os.Open(filePathAbs)

	if err != nil {
		t.Fatalf("Unable to open file %s, error: %s", fileName, err.Error())
	}

	byteValue, _ := io.ReadAll(jsonFile)
	return byteValue
}
