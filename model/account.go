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

type Account struct {
	ID      uint   `gorm:"primaryKey;comment:'主键'"`
	Phone   string `gorm:"varchar(225)"`
	Card    string `gorm:"varchar(225)"`
	Created int64
}

func CheckIsExistModelFish(db *gorm.DB) {
	if db.HasTable(&Account{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Account{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.CreateTable(&Account{})
	}
}
