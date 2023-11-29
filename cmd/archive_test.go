package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ChipArtem/k6/cmd/tests"
	"github.com/ChipArtem/k6/errext/exitcodes"
	"github.com/ChipArtem/k6/lib/fsext"
	"github.com/ChipArtem/k6/lib/testutils"
	"github.com/stretchr/testify/require"
)

func TestArchiveThresholds(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		noThresholds bool
		testFilename string

		wantErr bool
	}{
		{
			name:         "archive should fail with exit status 104 on a malformed threshold expression",
			noThresholds: false,
			testFilename: "testdata/thresholds/malformed_expression.js",
			wantErr:      true,
		},
		{
			name:         "archive should on a malformed threshold expression but --no-thresholds flag set",
			noThresholds: true,
			testFilename: "testdata/thresholds/malformed_expression.js",
			wantErr:      false,
		},
		{
			name:         "run should fail with exit status 104 on a threshold applied to a non existing metric",
			noThresholds: false,
			testFilename: "testdata/thresholds/non_existing_metric.js",
			wantErr:      true,
		},
		{
			name:         "run should succeed on a threshold applied to a non existing metric with the --no-thresholds flag set",
			noThresholds: true,
			testFilename: "testdata/thresholds/non_existing_metric.js",
			wantErr:      false,
		},
		{
			name:         "run should succeed on a threshold applied to a non existing submetric with the --no-thresholds flag set",
			noThresholds: true,
			testFilename: "testdata/thresholds/non_existing_metric.js",
			wantErr:      false,
		},
		{
			name:         "run should fail with exit status 104 on a threshold applying an unsupported aggregation method to a metric",
			noThresholds: false,
			testFilename: "testdata/thresholds/unsupported_aggregation_method.js",
			wantErr:      true,
		},
		{
			name:         "run should succeed on a threshold applying an unsupported aggregation method to a metric with the --no-thresholds flag set",
			noThresholds: true,
			testFilename: "testdata/thresholds/unsupported_aggregation_method.js",
			wantErr:      false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			testScript, err := os.ReadFile(testCase.testFilename) //nolint:forbidigo
			require.NoError(t, err)

			ts := tests.NewGlobalTestState(t)
			require.NoError(t, fsext.WriteFile(ts.FS, filepath.Join(ts.Cwd, testCase.testFilename), testScript, 0o644))
			ts.CmdArgs = []string{"k6", "archive", testCase.testFilename}
			if testCase.noThresholds {
				ts.CmdArgs = append(ts.CmdArgs, "--no-thresholds")
			}

			if testCase.wantErr {
				ts.ExpectedExitCode = int(exitcodes.InvalidConfig)
			}
			newRootCommand(ts.GlobalState).execute()
		})
	}
}

func TestArchiveContainsEnv(t *testing.T) {
	t.Parallel()

	// given some script that will be archived
	fileName := "script.js"
	testScript := []byte(`export default function () {}`)
	ts := tests.NewGlobalTestState(t)
	require.NoError(t, fsext.WriteFile(ts.FS, filepath.Join(ts.Cwd, fileName), testScript, 0o644))

	// when we do archiving and passing the `--env` flags
	ts.CmdArgs = []string{"k6", "--env", "ENV1=lorem", "--env", "ENV2=ipsum", "archive", fileName}

	newRootCommand(ts.GlobalState).execute()
	require.NoError(t, testutils.Untar(t, ts.FS, "archive.tar", "tmp/"))

	data, err := fsext.ReadFile(ts.FS, "tmp/metadata.json")
	require.NoError(t, err)

	metadata := struct {
		Env map[string]string
	}{}

	// then unpacked metadata should contain the environment variables with the proper values
	require.NoError(t, json.Unmarshal(data, &metadata))
	require.Len(t, metadata.Env, 2)

	require.Contains(t, metadata.Env, "ENV1")
	require.Contains(t, metadata.Env, "ENV2")

	require.Equal(t, "lorem", metadata.Env["ENV1"])
	require.Equal(t, "ipsum", metadata.Env["ENV2"])
}

func TestArchiveNotContainsEnv(t *testing.T) {
	t.Parallel()

	// given some script that will be archived
	fileName := "script.js"
	testScript := []byte(`export default function () {}`)
	ts := tests.NewGlobalTestState(t)
	require.NoError(t, fsext.WriteFile(ts.FS, filepath.Join(ts.Cwd, fileName), testScript, 0o644))

	// when we do archiving and passing the `--env` flags altogether with `--exclude-env-vars` flag
	ts.CmdArgs = []string{"k6", "--env", "ENV1=lorem", "--env", "ENV2=ipsum", "archive", "--exclude-env-vars", fileName}

	newRootCommand(ts.GlobalState).execute()
	require.NoError(t, testutils.Untar(t, ts.FS, "archive.tar", "tmp/"))

	data, err := fsext.ReadFile(ts.FS, "tmp/metadata.json")
	require.NoError(t, err)

	metadata := struct {
		Env map[string]string
	}{}

	// then unpacked metadata should not contain any environment variables passed at the moment of archive creation
	require.NoError(t, json.Unmarshal(data, &metadata))
	require.Len(t, metadata.Env, 0)
}
