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

package common

import (
	"emperror.dev/errors"
	"github.com/Masterminds/semver/v3"
)

// ValidateVersionWithConstraint will validate if the given version meets a constraint.
func ValidateVersionWithConstraint(version string, constraint string) (bool, error) {
	versionConstraint, err := semver.NewConstraint(constraint)
	if err != nil {
		return false, errors.WrapIff(err, "could not validate version with constraint, constraint: %v", constraint)
	}

	semVer, err := semver.NewVersion(version)
	if err != nil {
		return false, errors.WrapIff(err, "could not validate version with constraint, version: %v", version)
	}

	isValidVersion, errs := versionConstraint.Validate(semVer)
	if len(errs) != 0 {
		var totalErrors error
		for _, err := range errs {
			totalErrors = errors.Combine(totalErrors, err)
		}

		return false, errors.WrapIf(totalErrors, "could not validate version with constraint")
	}

	return isValidVersion, nil
}
