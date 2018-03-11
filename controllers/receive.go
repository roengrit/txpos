package controllers

import (
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
		c.Data["title"] = "รับสินค้า/วัตถุดิบ"
	} else {
		doc, _ := m.GetReceive(int(docID))
		c.Data["m"] = doc
		if !doc.Active {
			c.Data["r"] = "readonly"
		}
		c.Data["RetCount"] = len(doc.ReceiveSub)
		c.Data["RetPayCount"] = len(doc.Payment)
		c.Data["title"] = "แก้ไขรับสินค้า/วัตถุดิบ : " + doc.DocNo
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
