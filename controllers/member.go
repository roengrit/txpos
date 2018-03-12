package controllers

import (
	"bytes"
	"html/template"
	"strconv"
	"strings"
	"time"
	h "txpos/helpers"
	m "txpos/models"

	"github.com/go-playground/form"
)

//MemberController _
type MemberController struct {
	BaseController
}

//CreateMember _
func (c *MemberController) CreateMember() {
	memID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	if strings.Contains(c.Ctx.Request.URL.RequestURI(), "read") {
		c.Data["r"] = "readonly"
	}
	if memID == 0 {
		c.Data["title"] = "สร้าง สมาชิก"
	} else {
		c.Data["title"] = "แก้ไข สมาชิก"
		mem, _ := m.GetMember(int(memID))
		c.Data["m"] = mem
	}
	c.Data["Province"] = m.GetAllProvince()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "member/mem.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "member/mem-script.html"
	c.Render()
}

//UpdateMember _
func (c *MemberController) UpdateMember() {

	var sub m.Member
	decoder := form.NewDecoder()
	err := decoder.Decode(&sub, c.Ctx.Request.Form)
	ret := m.RetModel{}
	actionUser, _ := m.GetUser(h.GetUser(c.Ctx.Request))

	ret.RetOK = true
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
		_, err := m.CreateMember(sub)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		} else {
			ret.RetData = "บันทึกสำเร็จ"
		}
	} else if ret.RetOK && sub.ID > 0 {
		sub.EditedAt = time.Now()
		sub.Editor = &actionUser
		err := m.UpdateMember(sub)
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

//DeleteMember _
func (c *MemberController) DeleteMember() {
	ID, _ := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 32)
	ret := m.RetModel{}
	ret.RetOK = true
	err := m.DeleteMember(int(ID))
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

//MemberList _
func (c *MemberController) MemberList() {
	pn := h.NewPaging(0, 10, 0)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["Page"] = 0
	c.Data["PerPage"] = 10
	c.Data["Pages"] = pn
	c.Data["title"] = "สมาชิก"
	c.Layout = "layout.html"
	c.TplName = "member/mem-list.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "member/mem-list-script.html"
	c.Render()
}

//GetMemberListJSON GetMemberListJSON
func (c *MemberController) GetMemberListJSON() {
	searchTxt := c.Ctx.Request.FormValue("SearchTxt")
	page, _ := strconv.ParseInt(c.Ctx.Request.FormValue("Page"), 10, 64)
	perPage, _ := strconv.ParseInt(c.Ctx.Request.FormValue("PerPage"), 10, 64)
	page = h.PrePaging(page)
	offset := h.CalOffsetPaging(page, perPage)
	ret := m.MemberJSON{}

	num, list, _ := m.GetMemberListPaging(uint(offset), uint(perPage), searchTxt)

	pn := h.NewPaging(page, perPage, num)
	ret.MemberList = &list

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	t, _ := template.ParseFiles("views/paging.html")
	var tpl bytes.Buffer
	t.Execute(&tpl, pn)
	ret.Paging = tpl.String()
	ret.Page = uint(num)
	c.Data["json"] = ret
	c.ServeJSON()
}

//ListMemberJSON ListMemberJSON
func (c *MemberController) ListMemberJSON() {
	term := strings.TrimSpace(c.GetString("query"))
	ret := m.RetModel{}
	_, list, _ := m.GetMemberListPaging(0, 20, term)
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

//GetMemberJSON GetMemberJSON
func (c *MemberController) GetMemberJSON() {
	id, _ := strconv.Atoi(strings.TrimSpace(c.GetString("id")))
	ret := m.RetModel{}
	mem, err := m.GetMember(id)
	if err == nil {
		ret.RetOK = true
		ret.RetCount = 1
		ret.Data1 = mem
	} else {
		ret.RetOK = false
		ret.RetData = "ไม่พบข้อมูล"
	}
	c.Data["json"] = ret
	c.ServeJSON()
}
