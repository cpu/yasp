package yasp

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigValid(t *testing.T) {
	testCases := []struct {
		name     string
		config   Config
		expected error
	}{
		{
			name: "Invalid WinWidth",
			config: Config{
				WinWidth:  -1,
				WinHeight: 1,
			},
			expected: errBadWidth,
		},
		{
			name: "Invalid WinHeight",
			config: Config{
				WinWidth:  1,
				WinHeight: -1,
			},
			expected: errBadHeight,
		},
		{
			name: "Valid",
			config: Config{
				WinWidth:  1,
				WinHeight: 1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Valid()
			assert.Exactly(t, tc.expected, err,
				"expected Config (%#v) Valid() to give err %v got %v",
				tc.config, tc.expected, err)
		})
	}
}

func TestLoadConfigFile(t *testing.T) {
	testCases := []struct {
		name           string
		data           []byte
		expectError    string
		expectedConfig *Config
	}{
		{
			name:        "Bad YAML",
			data:        []byte("["),
			expectError: "yaml: line 1: did not find expected node content",
		},
		{
			name:        "Invalid config",
			data:        []byte(""),
			expectError: "invalid WinWidth",
		},
		{
			name: "Valid config",
			data: []byte("winwidth: 10\nwinheight: 20\n"),
			expectedConfig: &Config{
				WinWidth:  10,
				WinHeight: 20,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tf, err := ioutil.TempFile("", "yasp.config.*.yml")
			require.Nil(t, err, "unexpected err creating temp file: %v", err)
			defer tf.Close()

			_, err = tf.Write(tc.data)
			require.Nil(t, err, "unexpected err writing temp file data: %v", err)

			c, err := LoadConfigFile(tf.Name())
			if tc.expectError != "" {
				assert.EqualError(t, err, tc.expectError)
			}
			assert.Equal(t, tc.expectedConfig, c)
		})
	}
}
