package comment_like

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDoLike(t *testing.T) {
	Convey("Given 用户在阅读", t, func() {
		ip := "127.0.0.1"
		ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"
		title := "golang 网络框架之 grpc"

		Convey("When 点一个赞", func() {
			err := DoLike(ip, ua, title)
			So(err, ShouldBeNil)
			Convey("Then 数据库里面应该有一条记录", func() {
				var count int
				db.Model(&Like{}).Where(&Like{Ip: ip, Ua: ua, Title: title}).Count(&count)
				So(count, ShouldEqual, 1)
			})
		})

		Convey("When 重复点赞", func() {
			err := DoLike(ip, ua, title)
			So(err, ShouldNotBeNil)
			Convey("Then 数据库里面应该还是一条记录", func() {
				var count int
				db.Model(&Like{}).Where(&Like{Ip: ip, Ua: ua, Title: title}).Count(&count)
				So(count, ShouldEqual, 1)
			})
		})

		Convey("Finally 删除记录", func() {
			db.Where(&Like{Ip: ip, Ua: ua, Title: title}).Delete(Like{})
		})
	})
}

func TestDoUnlike(t *testing.T) {
	Convey("Given 用户在阅读", t, func() {
		ip := "127.0.0.1"
		ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"
		title := "golang 网络框架之 grpc"

		Convey("When 先点个赞，然后取消", func() {
			var err error
			err = DoLike(ip, ua, title)
			So(err, ShouldBeNil)
			err = DoUnlike(ip, ua, title)
			So(err, ShouldBeNil)
			Convey("Then 数据库里面没有记录了", func() {
				var count int
				db.Model(&Like{}).Where(&Like{Ip: ip, Ua: ua, Title: title}).Count(&count)
				So(count, ShouldEqual, 0)
			})
		})
	})
}

func TestShowLike(t *testing.T) {
	Convey("Given 用户在阅读", t, func() {
		ip := "127.0.0.1"
		ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"
		title := "golang 网络框架之 grpc"

		Convey("When 数据库里面没有点赞信息", func() {
			isLike, err := ShowLike(ip, ua, title)
			So(err, ShouldBeNil)
			Convey("Then 显示没有点赞", func() {
				So(isLike, ShouldBeFalse)
			})
		})

		Convey("When 先点个赞", func() {
			err := DoLike(ip, ua, title)
			So(err, ShouldBeNil)
			isLike, err := ShowLike(ip, ua, title)
			So(err, ShouldBeNil)
			Convey("Then 显示点赞了", func() {
				So(isLike, ShouldBeTrue)
			})
		})

		Convey("Finally 删除记录", func() {
			db.Where(&Like{Ip: ip, Ua: ua, Title: title}).Delete(Like{})
		})
	})
}

func TestCountLike(t *testing.T) {
	Convey("Given 用户在阅读", t, func() {
		ip := "127.0.0.1"
		ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"
		title := "golang 网络框架之 grpc"

		Convey("When 数据库里面没有点赞信息", func() {
			count, err := CountLike(title)
			So(err, ShouldBeNil)
			Convey("Then 点赞次数为0", func() {
				So(count, ShouldEqual, 0)
			})
		})

		Convey("When 先点个赞", func() {
			err := DoLike(ip, ua, title)
			So(err, ShouldBeNil)
			count, err := CountLike(title)
			So(err, ShouldBeNil)
			Convey("Then 点赞次数为1", func() {
				So(count, ShouldEqual, 1)
			})
		})

		Convey("When 再点个赞", func() {
			err := DoLike("192.168.0.2", ua, title)
			So(err, ShouldBeNil)
			count, err := CountLike(title)
			So(err, ShouldBeNil)
			Convey("Then 点赞次数为1", func() {
				So(count, ShouldEqual, 2)
			})
		})

		Convey("Finally 删除记录", func() {
			db.Where(&Like{Title: title}).Delete(Like{})
		})
	})
}

func TestDoComment(t *testing.T) {
	Convey("Given 用户在阅读", t, func() {
		ip := "127.0.0.1"
		ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"
		title := "golang 网络框架之 grpc"
		content := "写得不错"
		nickname := "sonic"
		mail := "sonic@foxmail.com"

		Convey("When 用户评论", func() {
			err := DoComment(ip, ua, title, content, nickname, mail)
			So(err, ShouldBeNil)
			Convey("Then 有了一条评论", func() {
				var comments []Comment
				db.Model(&Comment{}).Where(&Comment{Title: title}).Find(&comments)
				So(len(comments), ShouldEqual, 1)
				So(comments[0].Ip, ShouldEqual, ip)
				So(comments[0].Ua, ShouldEqual, ua)
				So(comments[0].Title, ShouldEqual, title)
				So(comments[0].Content, ShouldEqual, content)
				So(comments[0].Nickname, ShouldEqual, nickname)
				So(comments[0].Mail, ShouldEqual, mail)
			})
		})

		Convey("Finally 删除记录", func() {
			db.Where(&Comment{Title: title}).Delete(Comment{})
		})
	})
}

func TestShowComment(t *testing.T) {
	Convey("Given 用户在阅读", t, func() {
		ip := "127.0.0.1"
		ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"
		title := "golang 网络框架之 grpc"
		content := "写得不错"
		nickname := "sonic"
		mail := "sonic@foxmail.com"

		Convey("When 用户评论两次，然后显示评论", func() {
			var err error
			err = DoComment(ip, ua, title, content, nickname, mail)
			So(err, ShouldBeNil)
			err = DoComment(ip, ua, title, content, nickname, mail)
			So(err, ShouldBeNil)
			Convey("Then 显示两条评论", func() {
				comments, err := ShowComment(title)
				So(err, ShouldBeNil)
				So(len(*comments), ShouldEqual, 2)
				So((*comments)[0].Ip, ShouldEqual, ip)
				So((*comments)[0].Ua, ShouldEqual, ua)
				So((*comments)[0].Title, ShouldEqual, title)
				So((*comments)[0].Content, ShouldEqual, content)
				So((*comments)[0].Nickname, ShouldEqual, nickname)
				So((*comments)[0].Mail, ShouldEqual, mail)
			})
		})

		Convey("Finally 删除记录", func() {
			db.Where(&Comment{Title: title}).Delete(Comment{})
		})
	})
}
