// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package nextbus

import (
	"testing"
)

func Test_parseFloat(t *testing.T) {
	if v := parseFloat("3.14"); v != 3.14 {
		t.Error("Expected 3.14, got: ", v)
	}
	if v := parseFloat("cheese"); v != 0 {
		t.Error("Expected 0, got: ", v)
	}
}

func Test_parseBool(t *testing.T) {
	if v := parseBool("true"); v != true {
		t.Error("Expected true, got: ", v)
	}
	if v := parseBool("false"); v != false {
		t.Error("Expected false, got: ", v)
	}
	if v := parseBool("cheese"); v != false {
		t.Error("Expected false, got: ", v)
	}
}
