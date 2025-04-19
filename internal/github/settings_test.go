package github

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSettings_Template(t *testing.T) {
	yamlTeplate, err := os.ReadFile("../../SETTINGS_TEMPLATE.yaml")
	require.NoError(t, err)

	_, err = DecodeSettings(yamlTeplate)
	require.NoError(t, err)
}
