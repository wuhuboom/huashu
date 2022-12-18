/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/wangyi/makeAccount/model"

	"github.com/spf13/viper"
)

var (
	DB  *gorm.DB
	err error
)

func Init() error {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		viper.GetString("mysql.dbname"),
	)
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		fmt.Println("数据库链接失败", err.Error())
		panic(err)
		return err
	}

	//设置连接池
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
	////////////////////////////////////////////////////////////////////////模型初始化
	model.CheckIsExistModelFish(DB)
	model.CheckIsExistModelFootball(DB)
	model.CheckIsExistModelWordsArtLibrary(DB)
	model.CheckIsExistModelThePhrase(DB)
	model.CheckIsExistModelThePhraseV2(DB)
	model.CheckIsExistModelWordsArtLibraryV2(DB)
	model.CheckIsExistModelHash(DB)
	////////////////////////////////////////////////////////////////////////模型初始化
	return err
}

func Close() {
	defer DB.Close()
}
