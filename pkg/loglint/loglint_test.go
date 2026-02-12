package loglint

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/golangci/plugin-module-register/register"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestPluginExample(t *testing.T) {
	t.Run("require-literal", func(t *testing.T) {
		runWithSettings(t, map[string]any{
			"require-literal": true,
		}, "rules/literal")
	})

	t.Run("lowercase-start", func(t *testing.T) {
		runWithSettings(t, map[string]any{
			"require-lowercase-start": true,
		}, "rules/lowercase")
	})

	t.Run("lowercase-start-fix", func(t *testing.T) {
		runWithSettingsSuggestedFixes(t, map[string]any{
			"require-lowercase-start": true,
		}, "rules/lowercase")
	})

	t.Run("english-only", func(t *testing.T) {
		runWithSettings(t, map[string]any{
			"require-english": true,
		}, "rules/english")
	})

	t.Run("no-special-chars", func(t *testing.T) {
		runWithSettings(t, map[string]any{
			"forbid-special-chars": true,
		}, "rules/special")
	})

	t.Run("sensitive-data", func(t *testing.T) {
		runWithSettings(t, map[string]any{
			"forbid-sensitive-data": true,
			"sensitive-keywords":    []string{"password", "api_key", "token"},
			"sensitive-patterns":    []string{`(?i)secret\s*[:=]`},
		}, "rules/sensitive")
	})
}

func testdataDir(t *testing.T) string {
	t.Helper()

	_, testFilename, _, ok := runtime.Caller(1)
	if !ok {
		require.Fail(t, "unable to get current test filename")
	}

	return filepath.Join(filepath.Dir(testFilename), "testdata")
}

func runWithSettings(t *testing.T, overrides map[string]any, packages ...string) {
	t.Helper()

	newPlugin, err := register.GetPlugin("loglint")
	require.NoError(t, err)

	settings := map[string]any{
		"require-literal":         false,
		"require-lowercase-start": false,
		"require-english":         false,
		"forbid-special-chars":    false,
		"forbid-sensitive-data":   false,
	}

	for key, value := range overrides {
		settings[key] = value
	}

	plugin, err := newPlugin(settings)
	require.NoError(t, err)

	analyzers, err := plugin.BuildAnalyzers()
	require.NoError(t, err)

	analysistest.Run(t, testdataDir(t), analyzers[0], packages...)
}

func runWithSettingsSuggestedFixes(t *testing.T, overrides map[string]any, packages ...string) {
	t.Helper()

	newPlugin, err := register.GetPlugin("loglint")
	require.NoError(t, err)

	settings := map[string]any{
		"require-literal":         false,
		"require-lowercase-start": false,
		"require-english":         false,
		"forbid-special-chars":    false,
		"forbid-sensitive-data":   false,
	}

	for key, value := range overrides {
		settings[key] = value
	}

	plugin, err := newPlugin(settings)
	require.NoError(t, err)

	analyzers, err := plugin.BuildAnalyzers()
	require.NoError(t, err)

	analysistest.RunWithSuggestedFixes(t, testdataDir(t), analyzers[0], packages...)
}
