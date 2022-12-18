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

type ThePhraseV2 struct {
	ID       uint `gorm:"primaryKey;comment:'主键'"`
	WordsId  int
	Username string
	Created  int64
	Chinese  string `gorm:"-"` //中文
	English  string `gorm:"-"` //英文
	India    string `gorm:"-"` //印度
}

func CheckIsExistModelThePhraseV2(db *gorm.DB) {
	if db.HasTable(&ThePhrase{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&ThePhrase{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.CreateTable(&ThePhrase{})
	}
}
