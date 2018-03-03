package controllers

import (
	"html/template"
)

//ServiceController ServiceController
type ServiceController struct {
	BaseController
}

//GetXSRF GetXSRF
func (c *ServiceController) GetXSRF() {
	c.Data["json"] = c.XSRFToken()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.ServeJSON()
}
