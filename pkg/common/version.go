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
