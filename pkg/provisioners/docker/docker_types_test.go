package docker

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed testdata/inspect.json
var inspectJSON []byte

func TestParseInspectJSON(t *testing.T) {
	var inspect dockerInspect
	err := json.Unmarshal(inspectJSON, &inspect)
	require.NoError(t, err)
}
