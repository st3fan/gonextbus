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
