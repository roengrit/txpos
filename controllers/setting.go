package controllers

import (
	"html/template"
	"time"
	"txpos/helpers"
	"txpos/models"

	"github.com/go-playground/form"
)

//SettingController _
type SettingController struct {
	BaseController
}

//CreateSetting _
func (c *SettingController) CreateSetting() {
	setting, _ := models.GetSettingFirst()
	c.Data["title"] = "ตั้งค่าโปรแกรม"
	if setting.ID != 0 {
		c.Data["m"] = setting
	}
	c.Data["ret"] = models.RetModel{}
	c.Data["Province"] = models.GetAllProvince()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "setting/setting.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "setting/setting-script.html"
	c.Render()
}

//UpdateSetting _
func (c *SettingController) UpdateSetting() {

	var setting models.Setting
	decoder := form.NewDecoder()
	errForm := decoder.Decode(&setting, c.Ctx.Request.Form)
	ret := models.RetModel{}
	actionUser, _ := models.GetUser(helpers.GetUser(c.Ctx.Request))
	ret.RetOK = true
	if errForm != nil {
		ret.RetOK = false
		ret.RetData = "ข้อมูลบางอย่างไม่ครบถ้วน"
	}
	if ret.RetOK && setting.ID == 0 {
		setting.EditedAt = time.Now()
		setting.Editor = &actionUser
		_, err := models.CreateSetting(setting)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		} else {
			ret.RetData = "บันทึกสำเร็จ"
		}
	} else if ret.RetOK && setting.ID > 0 {
		setting.EditedAt = time.Now()
		setting.Editor = &actionUser
		err := models.UpdateSetting(setting)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		} else {
			ret.RetData = "บันทึกสำเร็จ"
		}
	}
	m, _ := models.GetSettingFirst()

	c.Data["title"] = "ตั้งค่าโปรแกรม"
	c.Data["m"] = m
	c.Data["ret"] = ret
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "setting/setting.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "setting/setting-script.html"
	c.Render()
}
