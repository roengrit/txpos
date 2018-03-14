package controllers

//POSController POSController
type POSController struct {
	BaseController
}

//Get _
func (c *POSController) Get() {
	c.Layout = "layout.html"
	c.TplName = "pos/pos.html"
	c.Render()
}
