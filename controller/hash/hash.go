/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package hash

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangyi/makeAccount/dao/mysql"
	"github.com/wangyi/makeAccount/model"
	"github.com/wangyi/makeAccount/util"
	"github.com/xuri/excelize/v2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetHash(c *gin.Context) {
	action := c.Query("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		var total int = 0
		Db := mysql.DB
		fish := make([]model.Hash, 0)

		if foxAddress, isExist := c.GetQuery("Address"); isExist == true {
			Db = Db.Where("address = ?", foxAddress)
		}
		if username, isExist := c.GetQuery("Username"); isExist == true {
			Db = Db.Where("username = ?", username)
		}

		if username, isExist := c.GetQuery("Status"); isExist == true {
			Db = Db.Where("status = ?", username)
		}

		Db.Table("hashes").Count(&total)
		Db = Db.Model(&fish).Offset((page - 1) * limit).Limit(limit).Order("updated desc")
		if err := Db.Find(&fish).Error; err != nil {
			util.JsonWrite(c, -101, nil, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": fish,
		})
		return
	}

	//上传用户名  POST
	if action == "uploadUsername" {
		username := c.PostForm("Username")
		uArray := strings.Split(username, "\n")
		if len(uArray) > 0 {
			for _, i2 := range uArray {
				if i2 != "" {
					h := model.Hash{Username: i2}
					h.Add(mysql.DB)
				}
			}
		}
		util.JsonWrite(c, 200, "nil", "上传成功")
	}

	//更新数据 GET 请求
	if action == "UPDATE" {
		id := c.Query("id")
		ups := make(map[string]interface{})
		if foxAddress, isExist := c.GetQuery("Address"); isExist == true {
			ups["Address"] = foxAddress
		}
		if foxAddress, isExist := c.GetQuery("GetAmount"); isExist == true {
			GetAmount, _ := strconv.ParseFloat(foxAddress, 64)
			ups["GetAmount"] = GetAmount
		}

		if foxAddress, isExist := c.GetQuery("RemarkOne"); isExist == true {
			ups["RemarkOne"] = foxAddress
		}
		if foxAddress, isExist := c.GetQuery("RemarkTwo"); isExist == true {
			ups["RemarkTwo"] = foxAddress
		}
		if foxAddress, isExist := c.GetQuery("RemarkThree"); isExist == true {
			ups["RemarkThree"] = foxAddress
		}
		if foxAddress, isExist := c.GetQuery("RemarkFour"); isExist == true {
			ups["RemarkFour"] = foxAddress
		}
		if foxAddress, isExist := c.GetQuery("Status"); isExist == true {
			status, _ := strconv.Atoi(foxAddress)
			ups["Status"] = status
		}
		ups["Updated"] = time.Now().Unix()

		mysql.DB.Model(&model.Hash{}).Where("id=?", id).Update(ups)
		util.JsonWrite(c, 200, "nil", "更新成功")
		return
	}

}

//下载 飞机内容
func DownHash(c *gin.Context) {
	//password := c.Query("password")
	//if password != "password" {
	//	util.JsonWrite(c, -101, nil, "密码错误")
	//	return
	//}
	//查询话术库
	words := make([]model.Hash, 0)
	var date string
	date = ""
	if date1, isE := c.GetQuery("date"); isE == true {
		date = date1
	}
	err := mysql.DB.Where("updated  > ?", 0).Find(&words).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "生成失败")
		return
	}
	//
	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")
	f.SetCellValue("Sheet1", "A1", "ID")
	f.SetCellValue("Sheet1", "B1", "飞机号")
	f.SetCellValue("Sheet1", "C1", "钱包地址")
	f.SetCellValue("Sheet1", "D1", "转发金额")
	f.SetCellValue("Sheet1", "E1", "备注一")
	f.SetCellValue("Sheet1", "F1", "备注二")
	f.SetCellValue("Sheet1", "G1", "备注三")
	f.SetCellValue("Sheet1", "H1", "备注四")

	for i, value := range words {
		if date != "" {
			if date != time.Unix(value.Updated, 0).Format("2006-01-02") {
				continue
			}
		}
		o := i + 2
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(o), value.ID)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(o), value.Username)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(o), value.Address)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(o), value.GetAmount)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(o), value.RemarkOne)
		f.SetCellValue("Sheet1", "F"+strconv.Itoa(o), value.RemarkTwo)
		f.SetCellValue("Sheet1", "G"+strconv.Itoa(o), value.RemarkThree)
		f.SetCellValue("Sheet1", "H"+strconv.Itoa(o), value.RemarkFour)

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
