package controllers

import "time"

//POSController POSController
type POSController struct {
	BaseController
}

//Get _
func (c *POSController) Get() {
	c.Data["CurrentDate"] = time.Now()
	c.Layout = "layout.html"
	c.TplName = "pos/pos.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["html_head"] = "pos/pos-style.html"
	c.LayoutSections["scripts"] = "pos/pos-script.html"
	c.Render()
}
