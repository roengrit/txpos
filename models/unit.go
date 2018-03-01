package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

//Unit Unit
type Unit struct {
	ID        int
	Lock      bool
	Name      string `orm:"size(300)"`
	Active    bool
	Creator   *User     `orm:"rel(fk)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	Editor    *User     `orm:"null;rel(fk)"`
	EditedAt  time.Time `orm:"null;auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Unit))
}

//GetAllProductUnit GetAllProductUnit
func GetAllProductUnit() *[]Unit {
	Unit := &[]Unit{}
	o := orm.NewOrm()
	o.QueryTable("unit").RelatedSel().All(Unit)
	return Unit
}

//GetProductUnit _
func GetProductUnit(ID int) (unit *Unit, errRet error) {
	Unit := &Unit{}
	o := orm.NewOrm()
	o.QueryTable("unit").Filter("ID", ID).RelatedSel().One(Unit)
	return Unit, errRet
}
