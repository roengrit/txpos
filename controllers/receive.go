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

//ReceiveController ReceiveController
type ReceiveController struct {
	BaseController
}

//Get _
func (c *ReceiveController) Get() {
	docID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	if strings.Contains(c.Ctx.Request.URL.RequestURI(), "read") {
		c.Data["r"] = "readonly"
	}
	if docID == 0 {
		c.Data["title"] = "ซื้อสินค้า/วัตถุดิบ"
	} else {
		doc, _ := m.GetReceive(int(docID))
		c.Data["m"] = doc
		if !doc.Active {
			c.Data["r"] = "readonly"
		}
		c.Data["RetCount"] = len(doc.ReceiveSub)
		c.Data["RetPayCount"] = len(doc.Payment)
		c.Data["title"] = "แก้ไขซื้อสินค้า/วัตถุดิบ : " + doc.DocNo
	}
	c.Data["CurrentDate"] = time.Now()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "receive/receive.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["html_head"] = "receive/receive-style.html"
	c.LayoutSections["scripts"] = "receive/receive-script.html"
	c.Render()
}

//Post _
func (c *ReceiveController) Post() {
	doc := m.Receive{}
	doc.Flag = 4 // นับ
	var RetID int64
	actionUser, _ := m.GetUser(h.GetUser(c.Ctx.Request))
	retJSON := m.RetModel{RetOK: true}
	decoder := form.NewDecoder()
	parsFormErr := decoder.Decode(&doc, c.Ctx.Request.Form)
	if parsFormErr == nil {

		if doc.Member.ID == 0 {
			doc.Member.ID = 2 // OTC
		}

		var dataPayment []m.Payment
		for index, val := range doc.Payment {
			if c.GetString("Payment["+strconv.Itoa(index)+"].PayDate") != "" {
				if docDate, err := h.ParseDateTime(c.GetString("Payment["+strconv.Itoa(index)+"].PayDate"), "00:00"); err == nil {
					val.PayDate = docDate

				} else {
					retJSON.RetOK = false
					retJSON.RetData = "วันที่/เวลา การชำระเงินไม่ถูกต้อง"
					break
				}
			} else {
				val.PayDate = time.Time{}
			}
			val.CreatedAt = time.Now()
			val.Creator = &actionUser
			dataPayment = append(dataPayment, val)
		}
		doc.Payment = dataPayment

		if docDate, err := h.ParseDateTime(c.GetString("DocDate"), doc.DocTime); err == nil {
			doc.DocDate = docDate
		} else {
			retJSON.RetOK = false
			retJSON.RetData = "วันที่/เวลา ไม่ถูกต้อง"
		}

		if retJSON.RetOK && c.GetString("ReceiveDate") != "" {
			if docDate, err := h.ParseDateTime(c.GetString("ReceiveDate"), doc.DocTime); err == nil {
				doc.ReceiveDate = docDate
				doc.ReceiveUser = &actionUser
			} else {
				retJSON.RetOK = false
				retJSON.RetData = "วันที่/เวลา ไม่ถูกต้อง"
			}
		}

		if retJSON.RetOK {
			if docDate, err := h.ParseDateTime(c.GetString("CreditDate"), "00:00"); err == nil {
				doc.CreditDate = docDate
			} else {
				retJSON.RetOK = false
				retJSON.RetData = "วันที่ เครดิต ไม่ถูกต้อง"
			}
		}

		if retJSON.RetOK && doc.ID == 0 {
			RetID, parsFormErr = m.CreateReceive(doc, actionUser)
			if parsFormErr == nil {
				retJSON.RetOK = true
				retJSON.RetData = "บันทึกสำเร็จ"
			} else {
				retJSON.RetOK = false
				retJSON.RetData = parsFormErr.Error()
			}
		}

		if retJSON.RetOK && doc.ID != 0 {
			_, parsFormErr = m.UpdateReceive(doc, actionUser)
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

//CancelReceive _
func (c *ReceiveController) CancelReceive() {
	ID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	ret := m.RetModel{}
	dataTemplate := m.RetModel{}
	dataTemplate.ID = ID
	dataTemplate.Title = "กรุณาระบุ หมายเหตุ การยกเลิก"
	dataTemplate.XSRF = c.XSRFToken()
	t, err := template.ParseFiles("views/receive/receive-cancel.html")
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

//UpdateCancelReceive _
func (c *ReceiveController) UpdateCancelReceive() {
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
		_, err := m.UpdateCancelReceive(ID, remark, actionUser)
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

//Receive _
func (c *ReceiveController) Receive() {
	ID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	ret := m.RetModel{}
	dataTemplate := m.RetModel{}
	dataTemplate.ID = ID
	dataTemplate.Title = "กรุณาระบุ วันที่รับสินค้า"
	dataTemplate.XSRF = c.XSRFToken()
	t, err := template.ParseFiles("views/receive/receive-receive.html")
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

//UpdateReceive _
func (c *ReceiveController) UpdateReceive() {
	actionUser, _ := m.GetUser(h.GetUser(c.Ctx.Request))
	ret := m.RetModel{}
	ID, _ := c.GetInt("ID")
	remark := c.GetString("Remark")
	receiveDate := c.GetString("ReceiveDate")
	receiveTime := c.GetString("ReceiveTime")
	var receiveDateTime time.Time
	ret.RetOK = true
	if docDate, err := h.ParseDateTime(receiveDate, receiveTime); err == nil {
		receiveDateTime = docDate
	} else {
		ret.RetOK = false
		ret.RetData = "วันที่/เวลา รับสินค้า ไม่ถูกต้อง"
	}
	if ret.RetOK {
		_, err := m.UpdateReceiveDate(ID, receiveDateTime, receiveTime, remark, actionUser)
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

//ReceiveList ReceiveList
func (c *ReceiveController) ReceiveList() {
	pn := h.NewPaging(0, 10, 0)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["Page"] = 0
	c.Data["PerPage"] = 10
	c.Data["Pages"] = pn
	c.Layout = "layout.html"
	c.TplName = "receive/receive-list.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "receive/receive-list-script.html"
	c.LayoutSections["html_head"] = "receive/receive-list-style.html"
	c.Render()
}

//GetReceiveListJSON GetReceiveListJSON
func (c *ReceiveController) GetReceiveListJSON() {
	searchTxt := c.Ctx.Request.FormValue("SearchTxt")
	txtDateBegin := c.Ctx.Request.FormValue("TxtDateBegin")
	txtDateEnd := c.Ctx.Request.FormValue("TxtDateEnd")
	receiveState := c.Ctx.Request.FormValue("ReceiveState")
	payState := c.Ctx.Request.FormValue("payState")
	page, _ := strconv.ParseInt(c.Ctx.Request.FormValue("Page"), 10, 64)
	perPage, _ := strconv.ParseInt(c.Ctx.Request.FormValue("PerPage"), 10, 64)
	page = h.PrePaging(page)
	offset := h.CalOffsetPaging(page, perPage)
	ret := m.ReceiveJSON{}

	num, list, _ := m.GetReceiveList(uint(offset), uint(perPage), txtDateBegin, txtDateEnd, receiveState, payState, searchTxt)

	pn := h.NewPaging(page, perPage, num)
	ret.ReceiveList = &list

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	t, _ := template.ParseFiles("views/paging.html")
	var tpl bytes.Buffer
	t.Execute(&tpl, pn)
	ret.Paging = tpl.String()
	ret.Page = uint(num)
	c.Data["json"] = ret
	c.ServeJSON()
}
