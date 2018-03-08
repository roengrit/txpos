package models

import (
	"errors"
	"time"
	"txpos/helpers"

	"github.com/astaxie/beego/orm"
)

//PickUp _
type PickUp struct {
	ID             int
	Flag           int
	Active         bool
	DocType        int
	DocNo          string    `orm:"size(30)"`
	DocDate        time.Time `form:"-" orm:"null"`
	DocTime        string    `orm:"size(6)"`
	DocRefNo       string    `orm:"size(30)"`
	DiscountType   int
	DiscountWord   string  `orm:"size(300)"`
	TotalDiscount  float64 `orm:"digits(12);decimals(2)"`
	TotalAmount    float64 `orm:"digits(12);decimals(2)"`
	TotalNetAmount float64 `orm:"digits(12);decimals(2)"`
	CreditDay      int
	CreditDate     time.Time   `orm:"type(date)"`
	Remark         string      `orm:"size(300)"`
	CancelRemark   string      `orm:"size(300)"`
	Creator        *User       `orm:"rel(fk)"`
	CreatedAt      time.Time   `orm:"auto_now_add;type(datetime)"`
	Editor         *User       `orm:"null;rel(fk)"`
	EditedAt       time.Time   `orm:"null;auto_now;type(datetime)"`
	CancelUser     *User       `orm:"null;rel(fk)"`
	CancelAt       time.Time   `orm:"null;type(datetime)"`
	PickUpSub      []PickUpSub `orm:"-"`
}

//PickUpSub _
type PickUpSub struct {
	ID          int
	Flag        int
	Active      bool
	DocNo       string    `orm:"size(30)"`
	DocDate     time.Time `form:"-" orm:"null"`
	Product     *Product  `orm:"rel(fk)"`
	Unit        *Unit     `orm:"rel(fk)"`
	BalanceQty  float64   `orm:"digits(12);decimals(2)"`
	Qty         float64   `orm:"digits(12);decimals(2)"`
	DiffQty     float64   `orm:"digits(12);decimals(2)"`
	RemainQty   float64   `orm:"digits(12);decimals(2)"`
	AverageCost float64   `orm:"digits(12);decimals(2)"`
	Price       float64   `orm:"digits(12);decimals(2)"`
	TotalPrice  float64   `orm:"digits(12);decimals(2)"`
	Remark      string    `orm:"size(300)"`
	Creator     *User     `orm:"rel(fk)"`
	CreatedAt   time.Time `orm:"auto_now_add;type(datetime)"`
}

//PickUpJSON PickUpJSON
type PickUpJSON struct {
	PickUpList *[]PickUp
	Paging     string
	Page       uint
	PerPage    uint
}

func init() {
	orm.RegisterModel(new(PickUp), new(PickUpSub)) // Need to register model in init
}

//CreatePickUp _
func CreatePickUp(PickUp PickUp, user User) (retID int64, errRet error) {
	PickUp.DocNo = helpers.GetMaxDoc("pick_up", "PI")
	PickUp.Creator = &user
	PickUp.CreatedAt = time.Now()
	PickUp.CreditDay = 0
	PickUp.CreditDate = time.Now()
	PickUp.Active = true
	var fullDataSub []PickUpSub
	for _, val := range PickUp.PickUpSub {
		if val.Product.ID != 0 {
			Product, _ := GetProduct(val.Product.ID)
			val.CreatedAt = time.Now()
			val.Creator = &user
			val.DocNo = PickUp.DocNo
			val.Flag = PickUp.Flag
			val.BalanceQty = Product.BalanceQty
			val.DiffQty = val.Qty - Product.BalanceQty
			val.DocDate = PickUp.DocDate
			val.AverageCost = val.Price
			fullDataSub = append(fullDataSub, val)
		}
	}
	o := orm.NewOrm()
	o.Begin()
	id, err := o.Insert(&PickUp)
	_, err = o.InsertMulti(len(fullDataSub), fullDataSub)
	if err == nil {
		retID = id
		o.Commit()
	} else {
		o.Rollback()
	}

	errRet = err
	return retID, errRet
}

//UpdatePickUp _
func UpdatePickUp(PickUp PickUp, user User) (retID int64, errRet error) {
	docCheck, _ := GetPickUp(PickUp.ID)
	if docCheck.ID == 0 {
		errRet = errors.New("ไม่พบข้อมูล")
	}
	PickUp.Creator = docCheck.Creator
	PickUp.CreatedAt = docCheck.CreatedAt
	PickUp.CreditDay = docCheck.CreditDay
	PickUp.CreditDate = docCheck.CreditDate
	PickUp.EditedAt = time.Now()
	PickUp.Editor = &user
	PickUp.Active = docCheck.Active
	var fullDataSub []PickUpSub
	for _, val := range PickUp.PickUpSub {
		if val.Product.ID != 0 {
			Product, _ := GetProduct(val.Product.ID)
			val.CreatedAt = time.Now()
			val.Creator = &user
			val.DocNo = PickUp.DocNo
			val.Flag = PickUp.Flag
			val.BalanceQty = Product.BalanceQty
			val.DiffQty = val.Qty - Product.BalanceQty
			val.AverageCost = val.Price
			val.DocDate = PickUp.DocDate
			fullDataSub = append(fullDataSub, val)
		}
	}
	o := orm.NewOrm()
	o.Begin()
	id, err := o.Update(&PickUp)
	if err == nil {
		_, err = o.QueryTable("pick_up_sub").Filter("doc_no", PickUp.DocNo).Delete()
	}
	if err == nil {
		_, err = o.InsertMulti(len(fullDataSub), fullDataSub)
	}
	if err == nil {
		o.Commit()
		retID = id
	} else {
		o.Rollback()
	}
	errRet = err
	return retID, errRet
}

//GetPickUp _
func GetPickUp(ID int) (doc *PickUp, errRet error) {
	PickUp := &PickUp{}
	o := orm.NewOrm()
	o.QueryTable("pick_up").Filter("ID", ID).RelatedSel().One(PickUp)
	o.QueryTable("pick_up_sub").Filter("doc_no", PickUp.DocNo).RelatedSel().All(&PickUp.PickUpSub)
	doc = PickUp
	return doc, errRet
}

//GetPickUpList GetPickUpList
func GetPickUpList(currentPage, lineSize uint, txtDateBegin, txtDateEnd, term string) (num int64, pickupList []PickUp, err error) {
	o := orm.NewOrm()
	var sql = `SELECT 
					T0.i_d,
					T0.doc_no,
					T0.doc_date,
					T0.remark,
					T0.flag_temp,
					T0.active,
					T0.total_amount
			   FROM pick_up T0	   
			   WHERE (lower(T0.doc_no) like lower(?) or lower(T0.remark) like lower(?)) `
	if txtDateBegin != "" && txtDateEnd != "" {
		sql += " and date( T0.doc_date) between '" + helpers.ParseDateString(txtDateBegin) + "' and '" + helpers.ParseDateString(txtDateEnd) + "'"
	}
	num, _ = o.Raw(sql, "%"+term+"%", "%"+term+"%").QueryRows(&pickupList)
	sql += " order by T0.doc_no  offset ? limit ?"
	_, err = o.Raw(sql, "%"+term+"%", "%"+term+"%", currentPage, lineSize).QueryRows(&pickupList)

	return num, pickupList, err
}

//UpdateCancelPickUp _
func UpdateCancelPickUp(ID int, remark string, user User) (retID int64, errRet error) {
	docCheck := &PickUp{}
	o := orm.NewOrm()
	o.QueryTable("pick_up").Filter("ID", ID).RelatedSel().One(docCheck)
	if docCheck.ID == 0 {
		errRet = errors.New("ไม่พบข้อมูล")
	}
	docCheck.Active = false
	docCheck.CancelRemark = remark
	docCheck.CancelAt = time.Now()
	docCheck.CancelUser = &user
	o.Begin()
	_, err := o.Update(docCheck)
	if err == nil {
		_, err = o.Raw("update pick_up_sub set active = false where doc_no = ?", docCheck.DocNo).Exec()
	}
	if err != nil {
		o.Rollback()
	} else {
		o.Commit()
	}
	errRet = err
	return retID, errRet
}

//UpdateActivePickUp _
func UpdateActivePickUp(ID int, user User) (retID int64, errRet error) {
	o := orm.NewOrm()
	o.Begin()
	orm.Debug = true
	_, err := o.Raw("update pick_up set active = true,flag_temp = 0,editor_id = ?,edited_at = now() where i_d = ?", user.ID, ID).Exec()
	if err != nil {
		o.Rollback()
	} else {
		_, err := o.Raw("update pick_up_sub set active = true where doc_no = (select pick_up.doc_no from pick_up where pick_up.i_d = ? limit 1)", ID).Exec()
		if err != nil {
			o.Rollback()
		} else {
			o.Commit()
		}
	}
	errRet = err
	return retID, errRet
}
