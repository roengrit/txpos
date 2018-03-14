package controllers

import (
	"bytes"
	"html/template"
	"strconv"
	"time"
	"txpos/helpers"
	"txpos/models"

	"github.com/go-playground/form"
)

//ChannelController _
type ChannelController struct {
	BaseController
}

//ChannelList ChannelList
func (c *ChannelController) ChannelList() {
	pn := helpers.NewPaging(0, 10, 0)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["Page"] = 0
	c.Data["PerPage"] = 10
	c.Data["Pages"] = pn
	c.Layout = "layout.html"
	c.TplName = "channel/list.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "channel/list-script.html"
	c.Render()
}

//GetChannelListJSON GetChannelListJSON
func (c *ChannelController) GetChannelListJSON() {
	searchTxt := c.Ctx.Request.FormValue("SearchTxt")

	page, _ := strconv.ParseInt(c.Ctx.Request.FormValue("Page"), 10, 64)
	perPage, _ := strconv.ParseInt(c.Ctx.Request.FormValue("PerPage"), 10, 64)
	page = helpers.PrePaging(page)
	offset := helpers.CalOffsetPaging(page, perPage)
	ret := models.ChannelListJSON{}

	num, list, _ := models.GetChannelList(uint(offset), uint(perPage), searchTxt)

	pn := helpers.NewPaging(page, perPage, num)
	ret.ChannelList = &list

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	t, _ := template.ParseFiles("views/paging.html")
	var tpl bytes.Buffer
	t.Execute(&tpl, pn)
	ret.Paging = tpl.String()
	ret.Page = uint(num)
	c.Data["json"] = ret
	c.ServeJSON()
}

//CreateChannel CreateChannel
func (c *ChannelController) CreateChannel() {
	cateID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	if cateID == 0 {
		c.Data["title"] = "สร้างช่องทางการขาย"
	} else {
		c.Data["title"] = "แก้ไขช่องทางการขาย"
		category, _ := models.GetChannel(int(cateID))
		c.Data["m"] = category
	}
	c.Data["ret"] = models.RetModel{}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "channel/channel.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "channel/channel-script.html"
	c.Render()
}

//UpdateChannel _
func (c *ChannelController) UpdateChannel() {

	var category models.Channel
	decoder := form.NewDecoder()
	err := decoder.Decode(&category, c.Ctx.Request.Form)
	ret := models.RetModel{}
	actionUser, _ := models.GetUser(helpers.GetUser(c.Ctx.Request))

	ret.RetOK = true
	if err != nil {
		ret.RetOK = false
		ret.RetData = err.Error()
	} else if c.GetString("Name") == "" {
		ret.RetOK = false
		ret.RetData = "กรุณาระบุชื่อ"
	}

	if ret.RetOK && category.ID == 0 {
		category.CreatedAt = time.Now()
		category.Creator = &actionUser
		_, err := models.CreateChannel(category)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		} else {
			ret.RetData = "บันทึกสำเร็จ"
		}
	} else if ret.RetOK && category.ID > 0 {
		category.EditedAt = time.Now()
		category.Editor = &actionUser
		err := models.UpdateChannel(category)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		} else {
			ret.RetData = "บันทึกสำเร็จ"
		}
	}
	ret.XSRF = c.XSRFToken()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["json"] = ret
	c.ServeJSON()
}

//DeleteChannel DeleteChannel
func (c *ChannelController) DeleteChannel() {
	ID, _ := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 32)
	ret := models.RetModel{}
	ret.RetOK = true
	err := models.DeleteChannel(int(ID))
	if err != nil {
		ret.RetOK = false
		ret.RetData = err.Error()
	} else {
		ret.RetData = "ลบข้อมูลสำเร็จ"
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["json"] = ret
	c.ServeJSON()
}
