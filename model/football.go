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
)

type Football struct {
	ID      uint   `gorm:"primaryKey;comment:'主键'"`
	Phone   string `gorm:"varchar(225)"`
	Card    string `gorm:"varchar(225)"`
	Created int64
}

func CheckIsExistModelFootball(db *gorm.DB) {
	if db.HasTable(&Football{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Football{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.CreateTable(&Football{})
	}
}
