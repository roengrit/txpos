package controllers

//AppController _
type AppController struct {
	BaseController
}

//Get Home
func (c *AppController) Get() {
	c.Data["title"] = "หน้าหลัก"
	c.Data["base_user_display"] = "Admin"
	c.Layout = "layout.html"
	c.TplName = "main/index.html"
	c.Render()
}
