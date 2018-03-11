package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

//Payment _
type Payment struct {
	ID          int
	Lock        bool
	DocNo       string    `orm:"size(30)"`
	TotalAmount float64   `orm:"digits(12);decimals(2)"`
	Remark      string    `orm:"size(300)"`
	Number      string    `orm:"size(300)"`
	PayDate     time.Time `form:"-" orm:"null;type(datetime)"`
	Creator     *User     `orm:"rel(fk)"`
	CreatedAt   time.Time `orm:"auto_now_add;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Payment)) // Need to register model in init
}

//GetPayment _
func GetPayment(docNo string) (doc *[]Payment) {
	Payment := &[]Payment{}
	o := orm.NewOrm()
	o.QueryTable("payment").Filter("doc_no", docNo).RelatedSel().All(&Payment)
	doc = Payment
	return doc
}
