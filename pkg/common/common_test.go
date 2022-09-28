// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.
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

package common_test

import (
	"testing"

	"github.com/cisco-open/cluster-registry-controller/pkg/common"
)

func TestValidateVersionWithConstraint(t *testing.T) {
	t.Parallel()

	tests := map[string][]struct {
		version     string
		constraint  string
		expectMatch bool
	}{
		"version meets constraint": {
			{
				version:     "1.24",
				constraint:  ">= 1.23",
				expectMatch: true,
			},
			{
				// pre-release version
				version:     "1.24.0-test",
				constraint:  ">= 1.23.0-0",
				expectMatch: true,
			},
		},
		"version does not meet constraint": {
			{
				version:     "1.24",
				constraint:  "<= 1.23",
				expectMatch: false,
			},
			{
				version:     "1.24.0-test",
				constraint:  "<= 1.23",
				expectMatch: false,
			},
			{
				// pre-release version
				version:     "1.24.0-test",
				constraint:  ">= 1.23",
				expectMatch: false,
			},
		},
		"bad version": {
			{
				version:     "test",
				constraint:  "<= 1.23",
				expectMatch: false,
			},
		},
		"bad constraint": {
			{
				version:     "1.24",
				constraint:  "test",
				expectMatch: false,
			},
		},
	}

	for title, test := range tests {
		for _, testcase := range test {
			match, err := common.ValidateVersionWithConstraint(testcase.version, testcase.constraint)
			if match != testcase.expectMatch {
				t.Fatalf("returned unexpected boolean value, testcase: %v, err: %v", title, err)
			}
		}
	}
}
