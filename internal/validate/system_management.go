// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/go-playground/validator/v10"
)

func ValidateRiskScoreID(fl validator.FieldLevel) bool {
	return RegexpSystemManagementUserOrEndpointID.MatchString(fl.Field().String())
}
