package controllers

import "html/template"

//ProductController _
type ProductController struct {
	BaseController
}

//ProductList ProductList
func (c *ProductController) ProductList() {
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "product/list.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "product/list-script.html"
	c.Render()
}

//GetProductList GetProductList
func (c *ProductController) GetProductList() {
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "product/list.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "product/list-script.html"
	c.Render()
}
