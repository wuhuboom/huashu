/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type Hash struct {
	ID          uint    `gorm:"primaryKey;comment:'主键'"`
	Username    string  //飞机用户名
	Created     int64   //创建时间
	Address     string  //钱包地址
	GetAmount   float64 //领取金额
	RemarkOne   string  //备注一
	RemarkTwo   string  //备注二
	RemarkThree string  //备注三
	RemarkFour  string  //备注四
	Status      int     `gorm:"default:1"` //1没有反馈  2有反馈
	Updated     int64
}

func CheckIsExistModelHash(db *gorm.DB) {
	if db.HasTable(&Hash{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Hash{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.CreateTable(&Hash{})
	}
}

//添加用户名
func (H *Hash) Add(db *gorm.DB) {
	//判断用户是否存在
	err := db.Where("username=?", H.Username).First(&H).Error
	if err == nil {
		return
	}
	H.Created = time.Now().Unix()
	db.Save(&H)
	return
}
