package nmap

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWithTarget(t *testing.T) {
	testHostname := "testhost"

	options := &Options{}
	target := WithTarget(testHostname)

	target(options)

	Convey("Target Hostname", t, func() {
		Convey("Options should contain target hostname", func() {
			So(options.Target, ShouldResemble, testHostname)
		})
	})
}

func TestWithFlags(t *testing.T) {
	testFlags := []string{
		"v",
		"A",
	}
	expectedFlags := []string{
		"-v",
		"-A",
	}

	options := &Options{}
	target := WithFlags(testFlags)

	target(options)

	Convey("Target Flags", t, func() {
		Convey("Options should contain target flags", func() {
			So(options.Flags, ShouldResemble, expectedFlags)
		})
	})
}
