package controllers

import (
	"github.com/astaxie/beego"
)

const active = "active menu-open"
const parentActive = "active"

//BaseController _
type BaseController struct {
	beego.Controller
}

//Prepare _
func (b *BaseController) Prepare() {

}
