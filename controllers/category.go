package controllers

import (
	"bytes"
	"html/template"
	"strconv"
	"time"
	"txpos/helpers"
	"txpos/models"

	"github.com/go-playground/form"
)

//CategoryController _
type CategoryController struct {
	BaseController
}

//CategoryList CategoryList
func (c *CategoryController) CategoryList() {
	pn := helpers.NewPaging(0, 10, 0)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["Page"] = 0
	c.Data["PerPage"] = 10
	c.Data["Pages"] = pn
	c.Layout = "layout.html"
	c.TplName = "category/list.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "category/list-script.html"
	c.Render()
}

//GetCategoryListJSON GetCategoryListJSON
func (c *CategoryController) GetCategoryListJSON() {
	searchTxt := c.Ctx.Request.FormValue("SearchTxt")

	page, _ := strconv.ParseInt(c.Ctx.Request.FormValue("Page"), 10, 64)
	perPage, _ := strconv.ParseInt(c.Ctx.Request.FormValue("PerPage"), 10, 64)
	page = helpers.PrePaging(page)
	offset := helpers.CalOffsetPaging(page, perPage)
	ret := models.CategoryListJSON{}

	num, list, _ := models.GetProductCategoryList(uint(offset), uint(perPage), searchTxt)

	pn := helpers.NewPaging(page, perPage, num)
	ret.CategoryList = &list

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	t, _ := template.ParseFiles("views/paging.html")
	var tpl bytes.Buffer
	t.Execute(&tpl, pn)
	ret.Paging = tpl.String()
	ret.Page = uint(num)
	c.Data["json"] = ret
	c.ServeJSON()
}

//CreateCategory CreateCategory
func (c *CategoryController) CreateCategory() {
	cateID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	if cateID == 0 {
		c.Data["title"] = "สร้างหมวดสินค้า"
	} else {
		c.Data["title"] = "แก้ไขหมวดสินค้า"
		category, _ := models.GetProductCategory(int(cateID))
		c.Data["m"] = category
	}
	c.Data["ret"] = models.RetModel{}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "category/cate.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "category/cate-script.html"
	c.Render()
}

//UpdateCategory _
func (c *CategoryController) UpdateCategory() {

	var category models.Category
	decoder := form.NewDecoder()
	err := decoder.Decode(&category, c.Ctx.Request.Form)
	ret := models.RetModel{}
	actionUser, _ := models.GetUser(helpers.GetUser(c.Ctx.Request))

	ret.RetOK = true
	if err != nil {
		ret.RetOK = false
		ret.RetData = err.Error()
	} else if c.GetString("Name") == "" {
		ret.RetOK = false
		ret.RetData = "กรุณาระบุชื่อ"
	}

	if ret.RetOK && category.ID == 0 {
		category.CreatedAt = time.Now()
		category.Creator = &actionUser
		_, err := models.CreateProductCategory(category)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		} else {
			ret.RetData = "บันทึกสำเร็จ"
		}
	} else if ret.RetOK && category.ID > 0 {
		category.EditedAt = time.Now()
		category.Editor = &actionUser
		err := models.UpdateProductCategory(category)
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

//DeleteCategory DeleteCategory
func (c *CategoryController) DeleteCategory() {
	ID, _ := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 32)
	ret := models.RetModel{}
	ret.RetOK = true
	err := models.DeleteProductCategory(int(ID))
	if err != nil {
		ret.RetOK = false
		ret.RetData = err.Error()
	} else {
		ret.RetData = "ลบข้อมูลสำเร็จ"
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["json"] = ret
	c.ServeJSON()
}
