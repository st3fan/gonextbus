// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package nextbus

import (
	"strconv"
)

func parseFloat(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return v
	}
	return 0
}

func parseBool(s string) bool {
	v, err := strconv.ParseBool(s)
	if err == nil {
		return v
	}
	return false
}
