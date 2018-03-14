package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"
	h "txpos/helpers"
	m "txpos/models"

	"github.com/go-playground/form"
)

//SaleController SaleController
type SaleController struct {
	BaseController
}

//Get _
func (c *SaleController) Get() {
	docID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	if strings.Contains(c.Ctx.Request.URL.RequestURI(), "read") {
		c.Data["r"] = "readonly"
	}
	if docID == 0 {
		c.Data["title"] = "ขายสินค้า/วัตถุดิบ"
	} else {
		doc, _ := m.GetSale(int(docID))
		c.Data["m"] = doc
		if !doc.Active {
			c.Data["r"] = "readonly"
		}
		c.Data["RetCount"] = len(doc.SaleSub)
		c.Data["RetPayCount"] = len(doc.Payment)
		c.Data["title"] = "แก้ไขขายสินค้า/วัตถุดิบ : " + doc.DocNo
	}
	c.Data["Channel"] = m.GetAllChannel()
	c.Data["CurrentDate"] = time.Now()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "sale/sale.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["html_head"] = "sale/sale-style.html"
	c.LayoutSections["scripts"] = "sale/sale-script.html"
	c.Render()
}

//Post _
func (c *SaleController) Post() {
	doc := m.Sale{}
	doc.Flag = 4 // นับ
	var RetID int64
	actionUser, _ := m.GetUser(h.GetUser(c.Ctx.Request))
	retJSON := m.RetModel{RetOK: true}
	decoder := form.NewDecoder()
	fmt.Println(c.GetString("Member.ID"))
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

		if retJSON.RetOK && c.GetString("SendDate") != "" {
			if docDate, err := h.ParseDateTime(c.GetString("SendDate"), doc.DocTime); err == nil {
				doc.SendDate = docDate
				doc.SendUser = &actionUser
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
			RetID, parsFormErr = m.CreateSale(doc, actionUser)
			if parsFormErr == nil {
				retJSON.RetOK = true
				retJSON.RetData = "บันทึกสำเร็จ"
			} else {
				retJSON.RetOK = false
				retJSON.RetData = parsFormErr.Error()
			}
		}

		if retJSON.RetOK && doc.ID != 0 {
			_, parsFormErr = m.UpdateSale(doc, actionUser)
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

//CancelSale _
func (c *SaleController) CancelSale() {
	ID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	ret := m.RetModel{}
	dataTemplate := m.RetModel{}
	dataTemplate.ID = ID
	dataTemplate.Title = "กรุณาระบุ หมายเหตุ การยกเลิก"
	dataTemplate.XSRF = c.XSRFToken()
	t, err := template.ParseFiles("views/sale/sale-cancel.html")
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

//UpdateCancelSale _
func (c *SaleController) UpdateCancelSale() {
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
		_, err := m.UpdateCancelSale(ID, remark, actionUser)
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

//Sale _
func (c *SaleController) Sale() {
	ID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	ret := m.RetModel{}
	dataTemplate := m.RetModel{}
	dataTemplate.ID = ID
	dataTemplate.Title = "กรุณาระบุ วันที่รับสินค้า"
	dataTemplate.XSRF = c.XSRFToken()
	t, err := template.ParseFiles("views/sale/sale-sale.html")
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

//UpdateSale _
func (c *SaleController) UpdateSale() {
	actionUser, _ := m.GetUser(h.GetUser(c.Ctx.Request))
	ret := m.RetModel{}
	ID, _ := c.GetInt("ID")
	remark := c.GetString("Remark")
	saleDate := c.GetString("SaleDate")
	saleTime := c.GetString("SaleTime")
	var saleDateTime time.Time
	ret.RetOK = true
	if docDate, err := h.ParseDateTime(saleDate, saleTime); err == nil {
		saleDateTime = docDate
	} else {
		ret.RetOK = false
		ret.RetData = "วันที่/เวลา รับสินค้า ไม่ถูกต้อง"
	}
	if ret.RetOK {
		_, err := m.UpdateSendDate(ID, saleDateTime, saleTime, remark, actionUser)
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

//SaleList SaleList
func (c *SaleController) SaleList() {
	pn := h.NewPaging(0, 10, 0)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["Page"] = 0
	c.Data["PerPage"] = 10
	c.Data["Pages"] = pn
	c.Layout = "layout.html"
	c.TplName = "sale/sale-list.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "sale/sale-list-script.html"
	c.LayoutSections["html_head"] = "sale/sale-list-style.html"
	c.Render()
}

//GetSaleListJSON GetSaleListJSON
func (c *SaleController) GetSaleListJSON() {
	searchTxt := c.Ctx.Request.FormValue("SearchTxt")
	txtDateBegin := c.Ctx.Request.FormValue("TxtDateBegin")
	txtDateEnd := c.Ctx.Request.FormValue("TxtDateEnd")
	saleState := c.Ctx.Request.FormValue("SaleState")
	payState := c.Ctx.Request.FormValue("payState")
	page, _ := strconv.ParseInt(c.Ctx.Request.FormValue("Page"), 10, 64)
	perPage, _ := strconv.ParseInt(c.Ctx.Request.FormValue("PerPage"), 10, 64)
	page = h.PrePaging(page)
	offset := h.CalOffsetPaging(page, perPage)
	ret := m.SaleJSON{}

	num, list, _ := m.GetSaleList(uint(offset), uint(perPage), txtDateBegin, txtDateEnd, saleState, payState, searchTxt)

	pn := h.NewPaging(page, perPage, num)
	ret.SaleList = &list

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	t, _ := template.ParseFiles("views/paging.html")
	var tpl bytes.Buffer
	t.Execute(&tpl, pn)
	ret.Paging = tpl.String()
	ret.Page = uint(num)
	c.Data["json"] = ret
	c.ServeJSON()
}

//Print _
func (c *SaleController) Print() {
	docID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)

	doc, _ := m.GetSale(int(docID))
	com, _ := m.GetComFirst()
	c.Data["m"] = doc
	c.Data["com"] = com
	c.Data["RetCount"] = len(doc.SaleSub)
	c.Data["CurrentDate"] = time.Now()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "sale/sale-invoice.html"
	if c.Ctx.Request.URL.Query().Get("po") == "1" {
		c.TplName = "sale/sale-iv-invoice.html"
		c.Data["title"] = "ใบสั่งขายสินค้า"
	} else {
		c.Data["title"] = "ใบรับสินค้า"
	}
	c.Render()
}
