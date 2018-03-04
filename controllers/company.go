package controllers

import (
	"html/template"
	"path/filepath"
	"time"
	"txpos/helpers"
	"txpos/models"

	"github.com/go-playground/form"
	"github.com/google/uuid"
)

//CompanyController _
type CompanyController struct {
	BaseController
}

//CreateCom _
func (c *CompanyController) CreateCom() {
	company, _ := models.GetComFirst()
	if company.ID == 0 {
		c.Data["title"] = "ข้อมูลร้าน/บริษัท"
	} else {
		c.Data["title"] = "แก้ไขข้อมูลร้าน/บริษัท"
		if len(company.ImageLogo) > 0 {
			base64, _ := helpers.File64Encode(company.ImageLogo)
			company.ImageBase64 = base64
		}
		c.Data["m"] = company
	}
	c.Data["ret"] = models.RetModel{}
	c.Data["Province"] = models.GetAllProvince()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "company/com.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "company/com-script.html"
	c.Render()
}

//UpdateCom _
func (c *CompanyController) UpdateCom() {

	var sub models.Company
	decoder := form.NewDecoder()
	err := decoder.Decode(&sub, c.Ctx.Request.Form)
	ret := models.RetModel{}
	actionUser, _ := models.GetUser(helpers.GetUser(c.Ctx.Request))

	ret.RetOK = true
	isNewImage := false
	file, header, _ := c.GetFile("ImgLogo")
	if file != nil {
		fileName := header.Filename
		fileName = uuid.New().String() + filepath.Ext(fileName)
		filePathSave := "data/company/" + fileName
		err = c.SaveToFile("ImgLogo", filePathSave)
		if err == nil {
			isNewImage = true
			sub.ImageLogo = filePathSave
			helpers.RemoveContentsExcludeFile("data/company", fileName)
			base64, errBase64 := helpers.File64Encode(filePathSave)
			err = errBase64
			sub.ImageBase64 = base64
		}
	}
	_ = file

	if err != nil {
		ret.RetOK = false
		ret.RetData = err.Error()
	} else if c.GetString("Name") == "" {
		ret.RetOK = false
		ret.RetData = "กรุณาระบุชื่อ"
	}

	if ret.RetOK && sub.ID == 0 {
		sub.CreatedAt = time.Now()
		sub.Creator = &actionUser
		_, err := models.CreateCom(sub)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		} else {
			ret.RetData = "บันทึกสำเร็จ"
		}
	} else if ret.RetOK && sub.ID > 0 {
		sub.EditedAt = time.Now()
		sub.Editor = &actionUser
		err := models.UpdateCom(sub, isNewImage)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		} else {
			ret.RetData = "บันทึกสำเร็จ"
		}
	}
	company, _ := models.GetComFirst()
	if company.ID == 0 {
		c.Data["title"] = "ข้อมูลร้าน/บริษัท"
		c.Data["m"] = sub
	} else {
		c.Data["title"] = "แก้ไขข้อมูลร้าน/บริษัท"
		if len(company.ImageLogo) > 0 {
			base64, _ := helpers.File64Encode(company.ImageLogo)
			company.ImageBase64 = base64
		}
		c.Data["m"] = company
	}
	c.Data["ret"] = ret
	c.Data["Province"] = models.GetAllProvince()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "company/com.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "company/com-script.html"
	c.Render()
}
