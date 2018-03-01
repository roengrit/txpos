package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

//Category _
type Category struct {
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
	orm.RegisterModel(new(Category))
}

//GetAllProductCategory GetAllProductCategory
func GetAllProductCategory() *[]Category {
	ProductCategory := &[]Category{}
	o := orm.NewOrm()
	o.QueryTable("category").RelatedSel().All(ProductCategory)
	return ProductCategory
}
