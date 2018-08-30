package client

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	defaultMsg = "Error"
)

func TestInternalError(t *testing.T) {
	Convey("Test internal error", t, func() {
		So(InternalError("Error"), ShouldImplement, (*error)(nil))
	})
}


func TestNotFound(t *testing.T) {
	Convey("Test not found error", t, func() {
		So(NotFound("Error"), ShouldImplement, (*error)(nil))
	})
}

func Test_messOrDefaultReturnDefaultError(t *testing.T) {
	Convey("Test mess or default", t, func() {
		Convey("Returns default message", func() {
			So(messOrDefault(defaultMsg), ShouldEqual, defaultMsg)
		})

		Convey("Returns given custom message", func() {
			msg := "An error here"
			So(messOrDefault(defaultMsg, msg), ShouldEqual, msg)
		})
	})
}
