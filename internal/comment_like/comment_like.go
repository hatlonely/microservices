package comment_like

import "github.com/jinzhu/gorm"
import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spaolacci/murmur3"
	"strings"
)

type Like struct {
	ID    int    `gorm:"primary_key"`
	Ip    string `gorm:"type:varchar(20);not null;"`
	Ua    string `gorm:"type:varchar(256);not null;"`
	Title string `gorm:"type:varchar(256);not null;"`
	Hash  uint64 `gorm:"unique_index:hash_idx;"`
}

type Comment struct {
	ID       int    `gorm:"primary_key"`
	Ip       string `gorm:"type:varchar(20);not null;"`
	Ua       string `gorm:"type:varchar(256);not null;"`
	Title    string `gorm:"type:varchar(128);not null;index:title_idx"`
	Comment  string `gorm:"type:varchar(1024);not null;"`
	NickName string `gorm:"type:varchar(20);"`
	Mail     string `gorm:"type:varchar(256);"`
}

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "hatlonely:keaiduo1@/hatlonely?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	if !db.HasTable(&Like{}) {
		db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Like{})
	}
	if !db.HasTable(&Comment{}) {
		db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Comment{})
	}
}

func DoLike(ip, ua, title string) error {
	like := &Like{
		Ip:    ip,
		Ua:    ua,
		Title: title,
		Hash: murmur3.Sum64([]byte(strings.Join([]string{ip, ua, title}, "-"))) >> 1,
	}

	var count int
	if err := db.Model(&Like{}).Where(like).Count(&count).Error; err != nil {
		return err
	}

	if err := db.Create(like).Error; err != nil {
		return err
	}

	return nil
}

func DoUnlike(ip, ua, title string) error {
	hash := murmur3.Sum64([]byte(strings.Join([]string{ip, ua, title}, "-"))) >> 1
	if err := db.Model(&Like{}).Where(&Like{Hash:hash}).Delete(&Like{}).Error; err != nil {
		return err
	}

	return nil
}

func ShowLike(ip, ua, title string) (bool, error) {
	return false, nil
}
