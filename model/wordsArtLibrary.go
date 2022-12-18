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

type WordsArtLibrary struct {
	ID               uint   `gorm:"primaryKey;comment:'主键'"`
	Chinese          string `gorm:"type:text"` //中文
	English          string `gorm:"type:text"` //英文
	India            string `gorm:"type:text"` //印度
	ChineseMd5       string
	Created          int64
	CollectionStatus int `gorm:"-"` //收藏状态
}

func CheckIsExistModelWordsArtLibrary(db *gorm.DB) {
	if db.HasTable(&WordsArtLibrary{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&WordsArtLibrary{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.CreateTable(&WordsArtLibrary{})
	}
}
