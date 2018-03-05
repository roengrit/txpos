package controllers

import (
	"bytes"
	"html/template"
	"strconv"
	"strings"
	"time"
	h "txpos/helpers"
	m "txpos/models"

	"github.com/go-playground/form"
)

//PickUpController PickUpController
type PickUpController struct {
	BaseController
}

//Get _
func (c *PickUpController) Get() {
	docID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	if strings.Contains(c.Ctx.Request.URL.RequestURI(), "read") {
		c.Data["r"] = "readonly"
	}
	if docID == 0 {
		c.Data["title"] = "เบิกสินค้า/วัตถุดิบ"
		c.Data["temp"] = 1
	} else {
		doc, _ := m.GetPickUp(int(docID))
		c.Data["m"] = doc
		if !doc.Active {
			c.Data["r"] = "readonly"
		}
		c.Data["temp"] = doc.FlagTemp
		c.Data["RetCount"] = len(doc.PickUpSub)
		c.Data["title"] = "แก้ไขเบิกสินค้า/วัตถุดิบ : " + doc.DocNo
	}
	c.Data["CurrentDate"] = time.Now()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "pickup/pickup.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["html_head"] = "pickup/pickup-style.html"
	c.LayoutSections["scripts"] = "pickup/pickup-script.html"
	c.Render()
}

//Post _
func (c *PickUpController) Post() {
	doc := m.PickUp{}
	doc.Flag = 4 // นับ
	var RetID int64
	actionUser, _ := m.GetUser(h.GetUser(c.Ctx.Request))
	retJSON := m.RetModel{RetOK: true}
	decoder := form.NewDecoder()
	parsFormErr := decoder.Decode(&doc, c.Ctx.Request.Form)
	if parsFormErr == nil {
		if docDate, err := h.ParseDateTime(c.GetString("DocDate"), doc.DocTime); err == nil {
			doc.DocDate = docDate
		} else {
			retJSON.RetOK = false
			retJSON.RetData = "มีข้อมูลบางอย่างไม่ครบถ้วน"
		}
		if retJSON.RetOK && doc.ID == 0 {
			RetID, parsFormErr = m.CreatePickUp(doc, actionUser)
			if parsFormErr == nil {
				retJSON.RetOK = true
				retJSON.RetData = "บันทึกสำเร็จ"
			} else {
				retJSON.RetOK = false
				retJSON.RetData = parsFormErr.Error()
			}
		}
		if retJSON.RetOK && doc.ID != 0 {
			_, parsFormErr = m.UpdatePickUp(doc, actionUser)
			if parsFormErr == nil {
				retJSON.RetOK = true
				retJSON.RetData = "บันทึกสำเร็จ"
				RetID = int64(doc.ID)
			} else {
				retJSON.RetOK = false
				retJSON.RetData = parsFormErr.Error()
			}
		}
		doc.ID = int(RetID)
		retJSON.ID = int64(doc.ID)
	} else {
		retJSON.RetOK = false
		retJSON.RetData = "มีข้อมูลบางอย่างไม่ครบถ้วน"
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["json"] = retJSON
	c.ServeJSON()
}

//PickUpList PickUpList
func (c *PickUpController) PickUpList() {
	pn := h.NewPaging(0, 10, 0)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["Page"] = 0
	c.Data["PerPage"] = 10
	c.Data["Pages"] = pn
	c.Layout = "layout.html"
	c.TplName = "pickup/pickup-list.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "pickup/pickup-list-script.html"
	c.LayoutSections["html_head"] = "pickup/pickup-list-style.html"
	c.Render()
}

//GetPickUpListJSON GetPickUpListJSON
func (c *PickUpController) GetPickUpListJSON() {
	searchTxt := c.Ctx.Request.FormValue("SearchTxt")
	txtDateBegin := c.Ctx.Request.FormValue("TxtDateBegin")
	txtDateEnd := c.Ctx.Request.FormValue("TxtDateEnd")
	page, _ := strconv.ParseInt(c.Ctx.Request.FormValue("Page"), 10, 64)
	perPage, _ := strconv.ParseInt(c.Ctx.Request.FormValue("PerPage"), 10, 64)
	page = h.PrePaging(page)
	offset := h.CalOffsetPaging(page, perPage)
	ret := m.PickUpJSON{}

	num, list, _ := m.GetPickUpList(uint(offset), uint(perPage), txtDateBegin, txtDateEnd, searchTxt)

	pn := h.NewPaging(page, perPage, num)
	ret.PickUpList = &list

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	t, _ := template.ParseFiles("views/paging.html")
	var tpl bytes.Buffer
	t.Execute(&tpl, pn)
	ret.Paging = tpl.String()
	ret.Page = uint(num)
	c.Data["json"] = ret
	c.ServeJSON()
}

//CancelPickUp _
func (c *PickUpController) CancelPickUp() {
	ID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	ret := m.RetModel{}
	dataTemplate := m.RetModel{}
	dataTemplate.ID = ID
	dataTemplate.Title = "กรุณาระบุ หมายเหตุ การยกเลิก"
	dataTemplate.XSRF = c.XSRFToken()
	t, err := template.ParseFiles("views/pickup/pickup-cancel.html")
	var tpl bytes.Buffer
	if err = t.Execute(&tpl, dataTemplate); err != nil {
		ret.RetOK = err != nil
		ret.RetData = err.Error()
	} else {
		ret.RetOK = true
		ret.RetData = tpl.String()
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["json"] = ret
	c.ServeJSON()
}

//UpdateCancelPickUp _
func (c *PickUpController) UpdateCancelPickUp() {
	actionUser, _ := m.GetUser(h.GetUser(c.Ctx.Request))
	ret := m.RetModel{}
	ID, _ := c.GetInt("ID")
	remark := c.GetString("Remark")
	ret.RetOK = true
	if remark == "" {
		ret.RetOK = false
		ret.RetData = "กรุณาระบุหมายเหตุ"
	}
	if ret.RetOK {
		_, err := m.UpdateCancelPickUp(ID, remark, actionUser)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		}
	}
	ret.XSRF = c.XSRFToken()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["json"] = ret
	c.ServeJSON()
}

//UpdateActivePickUp _
func (c *PickUpController) UpdateActivePickUp() {
	actionUser, _ := m.GetUser(h.GetUser(c.Ctx.Request))
	ret := m.RetModel{}
	ID, _ := c.GetInt("ID")
	ret.RetOK = true
	if ret.RetOK {
		_, err := m.UpdateActivePickUp(ID, actionUser)
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
