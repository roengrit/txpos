package controllers

import (
	"bytes"
	"html/template"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"txpos/helpers"
	"txpos/models"

	"github.com/go-playground/form"
	"github.com/google/uuid"
)

//ProductController _
type ProductController struct {
	BaseController
}

//ProductList ProductList
func (c *ProductController) ProductList() {
	pn := helpers.NewPaging(0, 10, 0)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["Page"] = 0
	c.Data["PerPage"] = 10
	c.Data["Pages"] = pn
	c.Data["ProductCategory"] = models.GetAllProductCategory()
	c.Layout = "layout.html"
	c.TplName = "product/list.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "product/list-script.html"
	c.Render()
}

//GetProductListJSON GetProductListJSON
func (c *ProductController) GetProductListJSON() {
	category := c.Ctx.Request.FormValue("Category")
	status := c.Ctx.Request.FormValue("Status")
	searchTxt := c.Ctx.Request.FormValue("SearchTxt")

	page, _ := strconv.ParseInt(c.Ctx.Request.FormValue("Page"), 10, 64)
	perPage, _ := strconv.ParseInt(c.Ctx.Request.FormValue("PerPage"), 10, 64)
	page = helpers.PrePaging(page)
	offset := helpers.CalOffsetPaging(page, perPage)

	num, list, _ := models.GetProductList(uint(offset), uint(perPage), status, category, searchTxt)

	pn := helpers.NewPaging(page, perPage, num)
	ret := models.ProductListJSON{}
	ret.ProductList = &list

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	t, _ := template.ParseFiles("views/product/paging.html")
	var tpl bytes.Buffer
	t.Execute(&tpl, pn)
	ret.Paging = tpl.String()
	ret.Page = uint(num)
	c.Data["json"] = ret
	c.ServeJSON()
}

//CreateProduct CreateProduct
func (c *ProductController) CreateProduct() {
	proID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	if strings.Contains(c.Ctx.Request.URL.RequestURI(), "read") {
		c.Data["r"] = "readonly"
	}
	if proID == 0 {
		c.Data["title"] = "สร้างสินค้า"
	} else {
		c.Data["title"] = "แก้ไขสินค้า"
		pro, _ := models.GetProduct(int(proID))
		if len(pro.ImagePath1) > 0 {
			base64, _ := helpers.File64Encode(pro.ImagePath1)
			pro.ImageBase64 = base64
		}
		c.Data["m"] = pro
	}
	c.Data["ret"] = models.RetModel{}
	c.Data["ProductCategory"] = models.GetAllProductCategory()
	c.Data["Unit"] = models.GetAllProductUnit()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "product/product.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "product/product-script.html"
	c.LayoutSections["html_head"] = "product/product-style.html"
	c.Render()
}

//UpdateProduct _
func (c *ProductController) UpdateProduct() {

	var pro models.Product
	decoder := form.NewDecoder()
	err := decoder.Decode(&pro, c.Ctx.Request.Form)
	ret := models.RetModel{}
	actionUser, _ := models.GetUser(helpers.GetUser(c.Ctx.Request))
	file, header, _ := c.GetFile("ImgProduct")
	isNewImage := false
	if file != nil {
		fileName := header.Filename
		fileName = uuid.New().String() + filepath.Ext(fileName)
		filePathSave := "data/product/" + fileName
		err = c.SaveToFile("ImgProduct", filePathSave)
		if err == nil {
			isNewImage = true
			pro.ImagePath1 = filePathSave
			base64, errBase64 := helpers.File64Encode(filePathSave)
			err = errBase64
			pro.ImageBase64 = base64
		}
	}
	_ = file

	ret.RetOK = true
	if err != nil {
		ret.RetOK = false
		ret.RetData = err.Error()
	} else if c.GetString("Name") == "" {
		ret.RetOK = false
		ret.RetData = "กรุณาระบุชื่อ"
	}

	if ret.RetOK && pro.ID == 0 {
		pro.CreatedAt = time.Now()
		pro.Creator = &actionUser
		ID, err := models.CreateProduct(pro)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		} else {
			pro.ID = int(ID)
			ret.RetData = "บันทึกสำเร็จ"
		}
	} else if ret.RetOK && pro.ID > 0 {
		pro.EditedAt = time.Now()
		pro.Editor = &actionUser
		err := models.UpdateProduct(pro, isNewImage)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		} else {
			ret.RetData = "บันทึกสำเร็จ"
		}
	}
	if pro.ID == 0 {
		c.Data["title"] = "สร้าง สินค้า"
		c.Data["m"] = pro
	} else {
		c.Data["title"] = "แก้ไข สินค้า"
		pro, _ := models.GetProduct(int(pro.ID))
		if len(pro.ImagePath1) > 0 {
			base64, _ := helpers.File64Encode(pro.ImagePath1)
			pro.ImageBase64 = base64
		}
		c.Data["m"] = pro
	}
	c.Data["ret"] = ret
	c.Data["ProductCategory"] = models.GetAllProductCategory()
	c.Data["Unit"] = models.GetAllProductUnit()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "product/product.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "product/product-script.html"
	c.Render()
}

//DeleteProduct _
func (c *ProductController) DeleteProduct() {
	ID, _ := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 32)
	ret := models.RetModel{}
	ret.RetOK = true
	err := models.DeleteProduct(int(ID))
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
