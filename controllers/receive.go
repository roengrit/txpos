package controllers

import (
	"html/template"
	"strconv"
	"strings"
	"time"
	m "txpos/models"
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
		doc, _ := m.GetPickUp(int(docID))
		c.Data["m"] = doc
		if !doc.Active {
			c.Data["r"] = "readonly"
		}
		c.Data["RetCount"] = len(doc.PickUpSub)
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
