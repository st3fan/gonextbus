// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package nextbus

import "testing"

func Test_FetchRouteList(t *testing.T) {
	_, err := FetchRouteList("ttc")
	if err != nil {
		t.Error("Expected err to be nil, got: ", err.Error())
	}
}

func Test_FetchRouteConfig(t *testing.T) {
	_, err := FetchRouteConfig("ttc", "501")
	if err != nil {
		t.Error("Expected err to be nil, got: ", err.Error())
	}
}
