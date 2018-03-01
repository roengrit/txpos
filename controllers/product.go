package controllers

import (
	"bytes"
	"html/template"
	"strconv"
	"txpos/helpers"
	"txpos/models"
)

//ProductController _
type ProductController struct {
	BaseController
}

//ProductList ProductList
func (c *ProductController) ProductList() {
	num, list, _ := models.GetProductList(0, uint(10), "", "", "")
	pn := helpers.NewPaging(0, 10, num)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["Page"] = 0
	c.Data["PerPage"] = 10
	c.Data["Pages"] = pn
	c.Data["Products"] = list
	c.Data["ProductCategory"] = models.GetAllProductCategory()
	c.Layout = "layout.html"
	c.TplName = "product/list.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "product/list-script.html"
	c.Render()
}

//GetProductListJSON GetProductListJSON
func (c *ProductController) GetProductListJSON() {
	page, _ := strconv.ParseInt(c.Ctx.Request.FormValue("Page"), 10, 64)
	perPage, _ := strconv.ParseInt(c.Ctx.Request.FormValue("PerPage"), 10, 64)
	if page <= 1 {
		page = 0
	} else {
		page = page - 1
	}

	category := c.Ctx.Request.FormValue("Category")
	status := c.Ctx.Request.FormValue("Status")
	searchTxt := c.Ctx.Request.FormValue("SearchTxt")

	num, list, _ := models.GetProductList(uint(page*perPage), uint(perPage), status, category, searchTxt)
	pn := helpers.NewPaging(page+1, perPage, num)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	ret := models.ProductListJSON{}
	ret.ProductList = &list
	t, _ := template.ParseFiles("views/product/paging.html")
	var tpl bytes.Buffer
	t.Execute(&tpl, pn)
	ret.Paging = tpl.String()
	ret.Page = uint(num)
	c.Data["json"] = ret
	c.ServeJSON()
}
