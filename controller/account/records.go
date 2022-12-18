/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package account

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/makeAccount/dao/mysql"
	"github.com/wangyi/makeAccount/model"
	"github.com/wangyi/makeAccount/util"
	"net/http"
	"strconv"
	"time"
)

func Account(c *gin.Context) {
	action := c.Query("action")

	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		var total int = 0
		Db := mysql.DB
		fish := make([]model.Account, 0)

		if foxAddress, isExist := c.GetQuery("phone"); isExist == true {
			Db = Db.Where("phone LIKE ?", "%"+foxAddress+"%")
		}
		if foxAddress, isExist := c.GetQuery("card"); isExist == true {
			Db = Db.Where("card LIKE ?", "%"+foxAddress+"%")
		}

		Db.Table("accounts").Count(&total)
		Db = Db.Model(&fish).Offset((page - 1) * limit).Limit(limit).Order("created desc")
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

	if action == "ADD" {
		phone := c.Query("phone")
		card := c.Query("card")
		//寻找 重复的手机号 活着 卡.
		err := mysql.DB.Where("phone =? or card=?", phone, card).First(&model.Account{}).Error
		if err == nil {
			util.JsonWrite(c, -101, "nil", "不可以重复添加!")
			return
		}
		err = mysql.DB.Save(&model.Account{Created: time.Now().Unix(), Phone: phone, Card: card}).Error
		if err != nil {
			util.JsonWrite(c, -101, "nil", "添加失败:"+err.Error())
			return
		}
		util.JsonWrite(c, 200, "nil", "添加成功")
	}

	if action == "DEL" {
		id := c.Query("id")
		err := mysql.DB.Where("id=?", id).First(&model.Account{}).Error
		if err != nil {
			util.JsonWrite(c, -101, "nil", "要删除的数据不存在!")
			return
		}
		err = mysql.DB.Delete(&model.Account{}, id).Error
		if err != nil {
			util.JsonWrite(c, -101, "nil", "删除失败!")
			return
		}
		util.JsonWrite(c, 200, "nil", "删除成功")
	}

}
