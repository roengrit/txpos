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
	ret := models.ProductListJSON{}

	num, list, _ := models.GetProductList(uint(offset), uint(perPage), status, category, searchTxt)

	pn := helpers.NewPaging(page, perPage, num)
	ret.ProductList = &list

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	t, _ := template.ParseFiles("views/paging.html")
	var tpl bytes.Buffer
	t.Execute(&tpl, pn)
	ret.Paging = tpl.String()
	ret.Page = uint(num)
	c.Data["json"] = ret
	c.ServeJSON()
}

//CreateProduct CreateProduct
func (c *ProductController) CreateProduct() {
	productID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	if productID == 0 {
		c.Data["title"] = "สร้างสินค้า"
		c.Data["ProductCode"] = models.GetMaxItemCode("P")
	} else {
		c.Data["title"] = "แก้ไขสินค้า"
		product, _ := models.GetProduct(int(productID))
		if len(product.ImagePath1) > 0 {
			base64, _ := helpers.File64Encode(product.ImagePath1)
			product.ImageBase64 = base64
		}
		c.Data["m"] = product
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

	var product models.Product
	decoder := form.NewDecoder()
	err := decoder.Decode(&product, c.Ctx.Request.Form)
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
			product.ImagePath1 = filePathSave
			base64, errBase64 := helpers.File64Encode(filePathSave)
			err = errBase64
			product.ImageBase64 = base64
		}
	}
	_ = file

	if product.DeleteImage == 1 {
		isNewImage = true
		product.ImagePath1 = ""
	}

	ret.RetOK = true
	if err != nil {
		ret.RetOK = false
		ret.RetData = err.Error()
	} else if c.GetString("Name") == "" {
		ret.RetOK = false
		ret.RetData = "กรุณาระบุชื่อ"
	}

	if ret.RetOK && product.ID == 0 {
		product.CreatedAt = time.Now()
		product.Creator = &actionUser
		ID, err := models.CreateProduct(product)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		} else {
			product.ID = int(ID)
			ret.RetData = "บันทึกสำเร็จ"
		}
	} else if ret.RetOK && product.ID > 0 {
		product.EditedAt = time.Now()
		product.Editor = &actionUser
		err := models.UpdateProduct(product, isNewImage)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		} else {
			ret.RetData = "บันทึกสำเร็จ"
		}
	}
	if product.ID == 0 {
		c.Data["title"] = "สร้าง สินค้า"
		c.Data["m"] = product
	} else {
		c.Data["title"] = "แก้ไข สินค้า"
		product, _ := models.GetProduct(int(product.ID))
		if len(product.ImagePath1) > 0 {
			base64, _ := helpers.File64Encode(product.ImagePath1)
			product.ImageBase64 = base64
		}
		c.Data["m"] = product
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

//ListProductJSON  _
func (c *ProductController) ListProductJSON() {
	term := strings.TrimSpace(c.GetString("query"))
	ret := models.RetModel{}
	_, list, _ := models.GetProductList(0, 20, "", "", term)
	if len(list) > 0 {
		ret.RetOK = true
		ret.RetCount = int64(len(list))
		ret.ListData = list
	} else {
		ret.RetOK = false
		ret.RetData = "ไม่พบข้อมูล"
	}
	c.Data["json"] = ret
	c.ServeJSON()
}

//GetProductJSON  GetProductJSON
func (c *ProductController) GetProductJSON() {
	ID, _ := strconv.ParseInt(c.GetString("id"), 10, 32)
	ret := models.RetModel{}
	product, err := models.GetProduct(int(ID))
	if err == nil {
		ret.RetOK = true
		product.Unit.Creator = nil
		product.Unit.Editor = nil
		product.Category.Creator = nil
		product.Category.Editor = nil
		product.Creator = nil
		product.Editor = nil
		ret.Data1 = product
	} else {
		ret.RetOK = false
		ret.RetData = "ไม่พบข้อมูล"
	}
	c.Data["json"] = ret
	c.ServeJSON()
}

//GetProductByCodeJSON  GetProductByCodeJSON
func (c *ProductController) GetProductByCodeJSON() {
	ret := models.RetModel{}
	product, err := models.GetProductByCode(c.GetString("code"))
	if err == nil {
		ret.RetOK = true
		ret.Data1 = product
	} else {
		ret.RetOK = false
		ret.RetData = "ไม่พบข้อมูล"
	}
	c.Data["json"] = ret
	c.ServeJSON()
}
