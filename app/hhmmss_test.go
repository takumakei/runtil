package app_test

import (
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/takumakei/runtil/app"
)

func TestHHMMSS_Parse(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	time.Local = loc

	tests := []struct {
		S string
		W string
	}{
		{S: "10", W: "1h0m0s"},
		{S: "10+09:00", W: "1h0m0s"},
		{S: "10Z", W: "10h0m0s"},
		{S: "10-08:00", W: "18h0m0s"},
		{S: "10-00:00", W: "10h0m0s"},
		{S: "23:02:03.123456789+09:00", W: "14h2m3.123456789s"},
		{S: "23:02:03.123456789Z", W: "23h2m3.123456789s"},
	}
	for i, test := range tests {
		var v app.HHMMSS
		if err := v.Parse(test.S); err != nil {
			t.Errorf("%d %s", i, err.Error())
		}
		if v.String() != test.W {
			t.Errorf("%d %s != %s", i, v, test.W)
		}
	}
}
