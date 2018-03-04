package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

//Provinces Provinces
type Provinces struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func init() {
	orm.RegisterModel(new(Provinces)) // Need to register model in init
}

//GetAllProvince _
func GetAllProvince() (req *[]Provinces) {
	o := orm.NewOrm()
	reqGet := &[]Provinces{}
	o.QueryTable("provinces").RelatedSel().OrderBy("ID").All(reqGet)
	return reqGet
}
