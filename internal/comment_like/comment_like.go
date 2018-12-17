package comment_like

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spaolacci/murmur3"
)

type View struct {
	ID        int    `gorm:"primary_key"`
	Ip        string `gorm:"type:varchar(20);not null;index:ip_idx"`
	Ua        string `gorm:"type:varchar(256);not null;"`
	Title     string `gorm:"type:varchar(128);not null;index:title_idx"`
	Hash      uint64
	CreatedAt time.Time
}

type Like struct {
	ID        int    `gorm:"primary_key"`
	Ip        string `gorm:"type:varchar(20);not null;index:ip_idx"`
	Ua        string `gorm:"type:varchar(256);not null;"`
	Title     string `gorm:"type:varchar(128);not null;index:title_idx"`
	Hash      uint64 `gorm:"unique_index:hash_idx;"`
	CreatedAt time.Time
}

type Comment struct {
	ID        int    `gorm:"primary_key"`
	Ip        string `gorm:"type:varchar(20);not null;"`
	Ua        string `gorm:"type:varchar(256);not null;"`
	Title     string `gorm:"type:varchar(128);not null;index:title_idx"`
	Content   string `gorm:"type:varchar(1024);not null;"`
	Nickname  string `gorm:"type:varchar(64);"`
	Mail      string `gorm:"type:varchar(256);"`
	CreatedAt time.Time
}

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "hatlonely:keaiduo1@/hatlonely?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	if !db.HasTable(&View{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&View{}).Error; err != nil {
			panic(err)
		}
	}
	if !db.HasTable(&Like{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Like{}).Error; err != nil {
			panic(err)
		}
	}
	if !db.HasTable(&Comment{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Comment{}).Error; err != nil {
			panic(err)
		}
	}
}

func DoView(ip, ua, title string) error {
	view := &View{
		Ip:        ip,
		Ua:        ua,
		Title:     title,
		Hash:      murmur3.Sum64([]byte(strings.Join([]string{ip, ua, title}, "-"))) >> 1,
		CreatedAt: time.Now(),
	}

	if err := db.Create(view).Error; err != nil {
		return err
	}

	return nil
}

func CountView(title string) (int64, error) {
	var count int64
	err := db.Model(&View{}).Where(&View{Title: title}).Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func DoLike(ip, ua, title string) error {
	hash := murmur3.Sum64([]byte(strings.Join([]string{ip, ua, title}, "-"))) >> 1
	like := &Like{
		Ip:        ip,
		Ua:        ua,
		Title:     title,
		Hash:      hash,
		CreatedAt: time.Now(),
	}

	var count int
	if err := db.Model(&Like{}).Where(&Like{Hash: hash}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		if err := db.Create(like).Error; err != nil {
			return err
		}
	}

	return nil
}

func DoUnlike(ip, ua, title string) error {
	hash := murmur3.Sum64([]byte(strings.Join([]string{ip, ua, title}, "-"))) >> 1
	if err := db.Where(&Like{Hash: hash}).Delete(Like{}).Error; err != nil {
		return err
	}

	return nil
}

func ShowLike(ip, ua, title string) (bool, error) {
	var count int
	err := db.Model(&Like{}).Where(&Like{Ip: ip, Ua: ua, Title: title}).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count != 0, nil
}

func CountLike(title string) (int64, error) {
	var count int64
	err := db.Model(&Like{}).Where(&Like{Title: title}).Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func DoComment(ip, ua, title, content, nickname, mail string) error {
	cmt := &Comment{
		Ip:        ip,
		Ua:        ua,
		Title:     title,
		Content:   content,
		Nickname:  nickname,
		Mail:      mail,
		CreatedAt: time.Now(),
	}

	if err := db.Create(cmt).Error; err != nil {
		return err
	}

	return nil
}

func ShowComment(title string) (*[]Comment, error) {
	var comments []Comment

	if err := db.Model(&Comment{}).Where(&Comment{Title: title}).Find(&comments).Error; err != nil {
		return nil, err
	}

	return &comments, nil
}
