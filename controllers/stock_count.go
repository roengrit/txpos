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

//StockCountController _
type StockCountController struct {
	BaseController
}

//Get _
func (c *StockCountController) Get() {
	docID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	if strings.Contains(c.Ctx.Request.URL.RequestURI(), "read") {
		c.Data["r"] = "readonly"
	}
	if docID == 0 {
		c.Data["title"] = "นับสต๊อคสินค้า/วัตถุดิบ"
		c.Data["temp"] = 1
	} else {
		doc, _ := m.GetStockCount(int(docID))
		c.Data["m"] = doc
		if !doc.Active {
			c.Data["r"] = "readonly"
		}
		c.Data["temp"] = doc.FlagTemp
		c.Data["RetCount"] = len(doc.StockCountSub)
		c.Data["title"] = "แก้ไข นับสต๊อคสินค้า/วัตถุดิบ : " + doc.DocNo
	}
	c.Data["CurrentDate"] = time.Now()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "stock/stock.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["html_head"] = "stock/stock-style.html"
	c.LayoutSections["scripts"] = "stock/stock-script.html"
	c.Render()
}

//StockDiff _
func (c *StockCountController) StockDiff() {
	docID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	doc, _ := m.GetStockCount(int(docID))
	c.Data["m"] = doc
	c.Data["RetCount"] = len(doc.StockCountSub)
	c.Data["title"] = "ผลต่างการนับสต๊อคสินค้า/วัตถุดิบ : " + doc.DocNo
	c.Data["CurrentDate"] = time.Now()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "stock/stock-diff.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["html_head"] = "stock/stock-diff-style.html"
	c.LayoutSections["scripts"] = "stock/stock-diff-script.html"
	c.Render()
}

//Post _
func (c *StockCountController) Post() {
	doc := m.StockCount{}
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
			RetID, parsFormErr = m.CreateStockCount(doc, actionUser)
			if parsFormErr == nil {
				retJSON.RetOK = true
				retJSON.RetData = "บันทึกสำเร็จ"
			} else {
				retJSON.RetOK = false
				retJSON.RetData = parsFormErr.Error()
			}
		}
		if retJSON.RetOK && doc.ID != 0 {
			_, parsFormErr = m.UpdateStockCount(doc, actionUser)
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

//StockList StockList
func (c *StockCountController) StockList() {
	pn := h.NewPaging(0, 10, 0)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["Page"] = 0
	c.Data["PerPage"] = 10
	c.Data["Pages"] = pn
	c.Layout = "layout.html"
	c.TplName = "stock/stock-list.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "stock/stock-list-script.html"
	c.LayoutSections["html_head"] = "stock/stock-list-style.html"
	c.Render()
}

//GetStockListJSON GetStockListJSON
func (c *StockCountController) GetStockListJSON() {
	searchTxt := c.Ctx.Request.FormValue("SearchTxt")
	txtDateBegin := c.Ctx.Request.FormValue("TxtDateBegin")
	txtDateEnd := c.Ctx.Request.FormValue("TxtDateEnd")
	page, _ := strconv.ParseInt(c.Ctx.Request.FormValue("Page"), 10, 64)
	perPage, _ := strconv.ParseInt(c.Ctx.Request.FormValue("PerPage"), 10, 64)
	page = h.PrePaging(page)
	offset := h.CalOffsetPaging(page, perPage)
	ret := m.StockCountJSON{}

	num, list, _ := m.GetStockCountList(uint(offset), uint(perPage), txtDateBegin, txtDateEnd, searchTxt)

	pn := h.NewPaging(page, perPage, num)
	ret.StockCountList = &list

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	t, _ := template.ParseFiles("views/paging.html")
	var tpl bytes.Buffer
	t.Execute(&tpl, pn)
	ret.Paging = tpl.String()
	ret.Page = uint(num)
	c.Data["json"] = ret
	c.ServeJSON()
}

//CancelStockCount _
func (c *StockCountController) CancelStockCount() {
	ID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	ret := m.RetModel{}
	dataTemplate := m.RetModel{}
	dataTemplate.ID = ID
	dataTemplate.Title = "กรุณาระบุ หมายเหตุ การยกเลิก"
	dataTemplate.XSRF = c.XSRFToken()
	t, err := template.ParseFiles("views/stock/stock-cancel.html")
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

//UpdateCancelStockCount _
func (c *StockCountController) UpdateCancelStockCount() {
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
		_, err := m.UpdateCancelStockCount(ID, remark, actionUser)
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

//UpdateActiveStockCount _
func (c *StockCountController) UpdateActiveStockCount() {
	actionUser, _ := m.GetUser(h.GetUser(c.Ctx.Request))
	ret := m.RetModel{}
	ID, _ := c.GetInt("ID")
	ret.RetOK = true
	if ret.RetOK {
		_, err := m.UpdateActiveStockCount(ID, actionUser)
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
