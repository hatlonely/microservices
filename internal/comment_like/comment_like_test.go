package comment_like

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDoLike(t *testing.T) {
	Convey("用户在阅读", t, func() {
		ip := "127.0.0.1"
		ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"
		title := "golang 网络框架之 grpc"

		Convey("点一个赞", func() {
			err := DoLike(ip, ua, title)
			So(err, ShouldBeNil)
		})

		Convey("重复点赞", func() {
			err := DoLike(ip, ua, title)
			So(err, ShouldNotBeNil)
		})

		Convey("取消赞", func() {
			err := DoUnlike(ip, ua, title)
			So(err, ShouldBeNil)
			var count int
			db.Model(&Like{}).Where(&Like{Ip:ip, Ua:ua, Title:title}).Count(&count)
			So(count, ShouldEqual, 0)
		})
	})
}

func TestDoUnlike(t *testing.T) {
	//Convey("用户在阅读", t, func() {
	//	ip := "127.0.0.1"
	//	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"
	//	title := "golang 网络框架之 grpc"
	//
	//
	//})
}
