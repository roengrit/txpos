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

//UnitController _
type UnitController struct {
	BaseController
}

//UnitList UnitList
func (c *UnitController) UnitList() {
	pn := helpers.NewPaging(0, 10, 0)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["Page"] = 0
	c.Data["PerPage"] = 10
	c.Data["Pages"] = pn
	c.Layout = "layout.html"
	c.TplName = "unit/list.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "unit/list-script.html"
	c.Render()
}

//GetUnitListJSON GetUnitListJSON
func (c *UnitController) GetUnitListJSON() {
	searchTxt := c.Ctx.Request.FormValue("SearchTxt")

	page, _ := strconv.ParseInt(c.Ctx.Request.FormValue("Page"), 10, 64)
	perPage, _ := strconv.ParseInt(c.Ctx.Request.FormValue("PerPage"), 10, 64)
	page = helpers.PrePaging(page)
	offset := helpers.CalOffsetPaging(page, perPage)
	ret := models.UnitListJSON{}

	num, list, _ := models.GetProductUnitList(uint(offset), uint(perPage), searchTxt)

	pn := helpers.NewPaging(page, perPage, num)
	ret.UnitList = &list

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	t, _ := template.ParseFiles("views/paging.html")
	var tpl bytes.Buffer
	t.Execute(&tpl, pn)
	ret.Paging = tpl.String()
	ret.Page = uint(num)
	c.Data["json"] = ret
	c.ServeJSON()
}

//CreateUnit CreateUnit
func (c *UnitController) CreateUnit() {
	unitID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	if unitID == 0 {
		c.Data["title"] = "สร้างหน่วยสินค้า"
	} else {
		c.Data["title"] = "แก้ไขหน่วยสินค้า"
		unit, _ := models.GetProductUnit(int(unitID))
		c.Data["m"] = unit
	}
	c.Data["ret"] = models.RetModel{}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "unit/unit.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "unit/unit-script.html"
	c.Render()
}

//UpdateUnit _
func (c *UnitController) UpdateUnit() {

	var unit models.Unit
	decoder := form.NewDecoder()
	err := decoder.Decode(&unit, c.Ctx.Request.Form)
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

	if ret.RetOK && unit.ID == 0 {
		unit.CreatedAt = time.Now()
		unit.Creator = &actionUser
		_, err := models.CreateProductUnit(unit)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		} else {
			ret.RetData = "บันทึกสำเร็จ"
		}
	} else if ret.RetOK && unit.ID > 0 {
		unit.EditedAt = time.Now()
		unit.Editor = &actionUser
		err := models.UpdateProductUnit(unit)
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

//DeleteUnit DeleteUnit
func (c *UnitController) DeleteUnit() {
	ID, _ := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 32)
	ret := models.RetModel{}
	ret.RetOK = true
	err := models.DeleteProductUnit(int(ID))
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
