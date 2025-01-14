// Copyright © 2019 Aqua Security Software Ltd. <info@aquasec.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aquasecurity/bench-common/check"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestLoadConfig(t *testing.T) {
	here, _ := os.Getwd()
	// cfgDir is defined in root.go
	type TestCase struct {
		version     string
		cfgPath     string
		want        string
		expectError bool
	}

	testCases := []TestCase{
		{
			version:     "2.0.0",
			cfgPath:     fmt.Sprintf("%s/../cfg", here),
			want:        "cfg/2.0.0/definitions.yaml",
			expectError: false,
		},
	}
	for _, tc := range testCases {
		cfgDir = tc.cfgPath
		got, _ := loadConfig(tc.version)
		if tc.expectError {
			assert.True(t, strings.Contains(got, tc.want))
		}
	}
}

func TestRunChecks(t *testing.T) {
	b := getMockBench()
	err := runChecks(b, "Server")
	var write bytes.Buffer
	outputWriter = &write
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}
	assert.NoError(t, err)
}

func TestUpdateControl(t *testing.T) {
	here, _ := os.Getwd()
	// cfgDir is defined in root.go
	type TestCase struct {
		version string
		cfgPath string
		want    string
	}

	testCases := []TestCase{
		{
			version: "2.0.0",
			cfgPath: fmt.Sprintf("%s/../cfg", here),
			want:    "cfg/2.0.0/definitions.yaml",
		},
	}
	for _, tc := range testCases {
		cfgDir = tc.cfgPath
		config, err := loadConfig(tc.version)
		assert.NoError(t, err)
		f, err := os.ReadFile(config)
		assert.NoError(t, err)
		var c check.Controls
		err = yaml.Unmarshal(f, &c)
		assert.NoError(t, err)
		got := updateControlCheck(&c, "DomainController")
		assert.Equal(t, got.Groups[0].Checks[0].Audit.(string), "Get-ADDefaultDomainPasswordPolicy -Current LocalComputer | Select -ExpandProperty PasswordHistoryCount")
	}
}
