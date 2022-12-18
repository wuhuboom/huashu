/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package Library

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangyi/makeAccount/dao/mysql"
	"github.com/wangyi/makeAccount/model"
	"github.com/wangyi/makeAccount/util"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
	"time"
)

//设置话术库
func SetWordsLibraryV2(c *gin.Context) {
	action := c.Query("action")
	if action == "ADD" { //添加话术
		Chinese := c.PostForm("chinese")
		India := c.PostForm("india")
		English := c.PostForm("english")
		password := c.PostForm("password")
		if password != "password" {
			util.JsonWrite(c, -101, nil, "密码不准确")
			return
		}
		if Chinese == "" {
			//fmt.Println(md5str1)
			util.JsonWrite(c, -101, nil, "不能为空")
			return
		}

		md5str1 := fmt.Sprintf("%x", md5.Sum([]byte(Chinese))) //将[]byte转成16进制
		//fmt.Println(md5str1)

		//判断是否重复添加
		err2 := mysql.DB.Where("chinese_md5=?", md5str1).First(&model.WordsArtLibraryV2{}).Error
		if err2 == nil { //找到了重复的数据
			util.JsonWrite(c, -101, nil, "不要重复添加")
			return
		}
		err := mysql.DB.Save(&model.WordsArtLibraryV2{Chinese: Chinese, India: India, English: English, ChineseMd5: md5str1, Created: time.Now().Unix()}).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "添加失败")
			return
		}

		util.JsonWrite(c, 200, nil, "添加成功")
	}

	if action == "GET" { //中文搜索 暂时只支持
		search := c.PostForm("search")

		words := make([]model.WordsArtLibraryV2, 0)
		err := mysql.DB.Where("chinese LIKE ?", "%"+search+"%").Find(&words).Error

		if err != nil {
			util.JsonWrite(c, -101, nil, "没有这条数据")
			return
		}

		if username, isExist := c.GetPostForm("username"); isExist == true {
			for k, va := range words {
				//
				err := mysql.DB.Where("words_id=? and username=?", va.ID, username).First(&model.ThePhraseV2{}).Error
				if err == nil {
					words[k].CollectionStatus = 1 //已收藏
				}
			}
		}

		util.JsonWrite(c, 200, words, "获取成功")
		return
	}

	if action == "DEL" {
		id := c.PostForm("id")
		admin := c.PostForm("password")
		if admin != "password" {
			util.JsonWrite(c, -101, nil, "密码错误")

			return
		}
		err := mysql.DB.Delete(&model.WordsArtLibraryV2{}, id).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "删除失败")
			return
		}
		util.JsonWrite(c, 200, nil, "删除成功")
		return

	}

	if action == "collection" {
		id := c.PostForm("id")
		username := c.PostForm("username")
		if username == "" {
			util.JsonWrite(c, -101, nil, "用户名不可以为空")
			return
		}
		world := model.WordsArtLibraryV2{}
		err := mysql.DB.Where("id=?", id).First(&world).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "收藏的话术不存在")
			return
		}

		//  查看 话术 是否已经添加过了

		err1 := mysql.DB.Where("words_id=? and username=?", world.ID, username).First(&model.ThePhraseV2{}).Error
		if err1 == nil {
			util.JsonWrite(c, -101, nil, "该话术已经收藏过了 ")
			return
		}

		add := model.ThePhraseV2{
			WordsId:  int(world.ID),
			Username: username,
			Created:  time.Now().Unix(),
		}

		err = mysql.DB.Save(&add).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "收藏失败")
			return
		}
		util.JsonWrite(c, 200, nil, "收藏成功")
		return

	}

}

/***
  下载  话术
*/
func DownWordWordsV2(c *gin.Context) {
	password := c.Query("password")
	if password != "password" {
		util.JsonWrite(c, -101, nil, "密码错误")
		return
	}
	//查询话术库
	words := make([]model.WordsArtLibraryV2, 0)
	err := mysql.DB.Find(&words).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "生成失败")
		return
	}
	//
	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")
	for i, value := range words {
		o := i + 1
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(o), value.Chinese)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(o), value.English)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(o), value.India)
	}
	f.SetActiveSheet(index)
	today := strconv.FormatFloat(float64(time.Now().Unix()), 'f', 0, 64)

	if err := f.SaveAs("static/" + today + ".xlsx"); err != nil {
		fmt.Println(err)
		util.JsonWrite(c, 200, nil, "生成失败")
		return
	}
	//today := time.Now().Format("2006-01-02")
	//time.Now().Format("2006-01-02 15:04:05")
	util.JsonWrite(c, 200, "static/"+today+".xlsx", "生成成功")
	return

}

/**
上传话术
*/
func UploadWordsV2(c *gin.Context) {

	file, err := c.FormFile("file")
	password := c.PostForm("password")
	if password != "password" {
		util.JsonWrite(c, -101, nil, "密码错误")
		return
	}

	if err != nil {
		util.JsonWrite(c, -101, nil, "上传错误:"+err.Error())
		return
	}
	fileType := strings.Split(file.Filename, ".")
	if fileType[1] != "xlsx" {
		util.JsonWrite(c, -101, nil, "上传格式不对")
		return
	}

	//fmt.Println(file.Filename)
	//f, err2 := excelize.OpenFile("static/Book.xlsx")

	err = c.SaveUploadedFile(file, file.Filename)
	if err != nil {
		util.JsonWrite(c, -101, nil, "上传错误:"+
			err.Error())
		return
	}

	f, err2 := excelize.OpenFile(file.Filename)
	if err2 != nil {
		util.JsonWrite(c, -101, nil, "上传错误:"+
			err2.Error())

		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err1 := f.GetRows("Sheet1")
	if err1 != nil {
		util.JsonWrite(c, -101, nil, "上传错误:"+err1.Error())
		return
	}
	for _, row := range rows {
		//fmt.Println(value)
		add := model.WordsArtLibraryV2{}
		for k, va := range row {
			if k == 0 {
				add.Chinese = va
			}
			if k == 1 {
				add.English = va
			}
			if k == 2 {
				add.India = va
			}
		}
		md5str1 := fmt.Sprintf("%x", md5.Sum([]byte(add.Chinese))) //将[]byte转成16进制
		//判断是否重复添加
		err2 := mysql.DB.Where("chinese_md5=?", md5str1).First(&model.WordsArtLibraryV2{}).Error
		if err2 != nil {
			add.Created = time.Now().Unix()
			err5 := mysql.DB.Save(&add).Error
			if err5 != nil {
				fmt.Println(err5.Error())
			}
		}
	}

	util.JsonWrite(c, 200, nil, "上传成功")

}

/**
  S 设置 和查询 常用语
*/

func SetThePhraseV2(c *gin.Context) {
	action := c.Query("action")

	if action == "GET" { //中文搜索 暂时只支持
		username := c.PostForm("username")
		if username == "" {
			util.JsonWrite(c, -101, nil, "数据为空")
			return
		}

		words := make([]model.ThePhraseV2, 0)
		err := mysql.DB.Where("username =? ", username).Find(&words).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "没有这条数据")
			return
		}
		for k, va := range words {
			pp := model.WordsArtLibraryV2{}
			err := mysql.DB.Where("id=?", va.WordsId).First(&pp).Error
			if err == nil {
				words[k].Chinese = pp.Chinese
				words[k].English = pp.English
				words[k].India = pp.India
			}
		}

		util.JsonWrite(c, 200, words, "获取成功")
		return
	}

	if action == "DEL" {
		id := c.PostForm("id")
		//admin := c.PostForm("password")
		//if admin != "password" {
		//	util.JsonWrite(c, -101, nil, "密码错误")
		//	return
		//}

		err := mysql.DB.Delete(&model.ThePhraseV2{}, id).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "删除失败")
			return
		}
		util.JsonWrite(c, 200, nil, "删除成功")
		return

	}

}
