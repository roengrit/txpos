package models

import (
	"errors"
	"fmt"
	"time"
	"txpos/helpers"

	"github.com/astaxie/beego/orm"
)

//Receive _
type Receive struct {
	ID                   int
	Flag                 int
	Active               bool
	DocType              int
	SaleType             int
	VatType              int
	DocNo                string    `orm:"size(30)"`
	DocDate              time.Time `form:"-" orm:"null"`
	DocTime              string    `orm:"size(6)"`
	DocRefNo             string    `orm:"size(30)"`
	Member               *Member   `orm:"rel(fk)"`
	MemberName           string    `orm:"size(300)"`
	Tel                  string    `orm:"size(100)"`
	Fax                  string    `orm:"size(100)"`
	Email                string    `orm:"size(100)"`
	Address              string    `orm:"size(500)"`
	TaxNo                string    `orm:"size(13)"`
	BranchNo             string    `orm:"size(25)"`
	BranchName           string    `orm:"size(200)"`
	DiscountType         int
	DiscountWord         string  `orm:"size(300)"`
	VatVal               float64 `orm:"digits(12);decimals(2)"`
	TotalDiscount        float64 `orm:"digits(12);decimals(2)"`
	TotalAmount          float64 `orm:"digits(12);decimals(2)"`
	TotalAmountExludeVat float64 `orm:"digits(12);decimals(2)"`
	TotalNetAmount       float64 `orm:"digits(12);decimals(2)"`
	TotalPay             float64 `orm:"digits(12);decimals(2)"`
	CreditDay            int
	CreditDate           time.Time    `form:"-" orm:"type(date)"`
	ReceiveDate          time.Time    `form:"-" orm:"null;type(datetime)"`
	ReceiveTime          string       `orm:"size(6)"`
	ReceiveUser          *User        `orm:"null;rel(fk)"`
	Remark               string       `orm:"size(300)"`
	CancelRemark         string       `orm:"size(300)"`
	ReceiveRemark        string       `orm:"size(300)"`
	Creator              *User        `orm:"rel(fk)"`
	CreatedAt            time.Time    `orm:"auto_now_add;type(datetime)"`
	Editor               *User        `orm:"null;rel(fk)"`
	EditedAt             time.Time    `orm:"null;auto_now;type(datetime)"`
	CancelUser           *User        `orm:"null;rel(fk)"`
	CancelAt             time.Time    `orm:"null;type(datetime)"`
	ReceiveSub           []ReceiveSub `orm:"-"`
	Payment              []Payment    `orm:"-"`
}

//ReceiveSub _
type ReceiveSub struct {
	ID                  int
	Flag                int
	Active              bool
	DocNo               string    `orm:"size(30)"`
	DocDate             time.Time `form:"-" orm:"null"`
	Product             *Product  `orm:"rel(fk)"`
	ProductName         string    `orm:"size(300)"`
	ProductDesc         string    `orm:"size(300)"`
	Unit                *Unit     `orm:"rel(fk)"`
	Qty                 float64   `orm:"digits(12);decimals(2)"`
	RemainQty           float64   `orm:"digits(12);decimals(2)"`
	AverageCost         float64   `orm:"digits(12);decimals(2)"`
	Price               float64   `orm:"digits(12);decimals(2)"`
	VatVal              float64   `orm:"digits(12);decimals(2)"`
	DiscountWord        string    `orm:"size(300)"`
	DiscountVal         float64   `orm:"digits(12);decimals(2)"`
	TotalPrice          float64   `orm:"digits(12);decimals(2)"`
	TotalPriceExludeVat float64   `orm:"digits(12);decimals(2)"`
	Creator             *User     `orm:"rel(fk)"`
	CreatedAt           time.Time `orm:"auto_now_add;type(datetime)"`
}

//ReceiveJSON ReceiveJSON
type ReceiveJSON struct {
	ReceiveList *[]Receive
	Paging      string
	Page        uint
	PerPage     uint
}

func init() {
	orm.RegisterModel(new(Receive), new(ReceiveSub)) // Need to register model in init
}

//CreateReceive _
func CreateReceive(receive Receive, user User) (retID int64, errRet error) {
	receive.DocNo = helpers.GetMaxDoc("receive", "PO")
	receive.Creator = &user
	receive.CreatedAt = time.Now()
	receive.CreditDay = 0
	receive.CreditDate = time.Now()
	receive.Active = true
	var fullDataSub []ReceiveSub
	var fullDataPayment []Payment
	for _, val := range receive.ReceiveSub {
		if val.Product.ID != 0 {
			val.CreatedAt = time.Now()
			val.Creator = &user
			val.DocNo = receive.DocNo
			val.Flag = receive.Flag
			val.Active = true
			val.DocDate = receive.DocDate
			val.ProductName = val.Product.Name
			fullDataSub = append(fullDataSub, val)
		}
	}
	for _, val := range receive.Payment {
		val.CreatedAt = time.Now()
		val.Creator = &user
		val.DocNo = receive.DocNo
		fullDataPayment = append(fullDataPayment, val)
	}
	o := orm.NewOrm()
	o.Begin()
	id, err := o.Insert(&receive)
	id, err = o.InsertMulti(len(fullDataSub), fullDataSub)
	id, err = o.InsertMulti(len(fullDataPayment), fullDataPayment)
	if err == nil {
		retID = id
		o.Commit()
	} else {
		o.Rollback()
	}
	errRet = err
	return retID, errRet
}

//UpdateReceive _
func UpdateReceive(receive Receive, user User) (retID int64, errRet error) {
	docCheck, _ := GetReceive(receive.ID)
	if docCheck.ID == 0 {
		errRet = errors.New("ไม่พบข้อมูล")
	}
	fmt.Println(receive.MemberName)
	receive.Creator = docCheck.Creator
	receive.CreatedAt = docCheck.CreatedAt
	receive.CreditDay = docCheck.CreditDay
	receive.CreditDate = docCheck.CreditDate
	receive.EditedAt = time.Now()
	receive.Editor = &user
	receive.Active = docCheck.Active
	var fullDataSub []ReceiveSub
	var fullDataPayment []Payment
	for _, val := range receive.ReceiveSub {
		if val.Product.ID != 0 {
			val.CreatedAt = time.Now()
			val.Creator = &user
			val.DocNo = receive.DocNo
			val.Flag = receive.Flag
			val.Active = receive.Active
			val.DocDate = receive.DocDate
			val.ProductName = val.Product.Name
			fullDataSub = append(fullDataSub, val)
		}
	}
	for _, val := range receive.Payment {
		val.CreatedAt = time.Now()
		val.Creator = &user
		val.DocNo = receive.DocNo
		fullDataPayment = append(fullDataPayment, val)
	}
	o := orm.NewOrm()
	o.Begin()
	id, err := o.Update(&receive)
	if err == nil {
		_, err = o.QueryTable("receive_sub").Filter("doc_no", receive.DocNo).Delete()
	}
	if err == nil {
		_, err = o.InsertMulti(len(fullDataSub), fullDataSub)
	}
	if err == nil {
		_, err = o.QueryTable("payment").Filter("doc_no", receive.DocNo).Delete()
	}
	if err == nil {
		_, err = o.InsertMulti(len(fullDataPayment), fullDataPayment)
	}
	if err == nil {
		retID = id
		o.Commit()
	} else {
		o.Rollback()
	}
	errRet = err
	return retID, errRet
}

//GetReceive _
func GetReceive(ID int) (doc *Receive, errRet error) {
	receiveDoc := &Receive{}
	o := orm.NewOrm()
	o.QueryTable("receive").Filter("ID", ID).RelatedSel().One(receiveDoc)
	o.QueryTable("receive_sub").Filter("doc_no", receiveDoc.DocNo).RelatedSel().All(&receiveDoc.ReceiveSub)
	o.QueryTable("payment").Filter("doc_no", receiveDoc.DocNo).OrderBy("ID").RelatedSel().All(&receiveDoc.Payment)
	doc = receiveDoc
	return doc, errRet
}

//UpdateCancelReceive _
func UpdateCancelReceive(ID int, remark string, user User) (retID int64, errRet error) {
	docCheck := &Receive{}
	o := orm.NewOrm()
	o.QueryTable("receive").Filter("ID", ID).RelatedSel().One(docCheck)
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
		_, err = o.Raw("update receive_sub set active = false where doc_no = ?", docCheck.DocNo).Exec()
	}
	if err != nil {
		o.Rollback()
	} else {
		o.Commit()
	}
	errRet = err
	return retID, errRet
}

//UpdateReceiveDate _
func UpdateReceiveDate(ID int, receiveDate time.Time, receiveTime, remark string, user User) (retID int64, errRet error) {
	docCheck := &Receive{}
	o := orm.NewOrm()
	o.QueryTable("receive").Filter("ID", ID).RelatedSel().One(docCheck)
	if docCheck.ID == 0 {
		errRet = errors.New("ไม่พบข้อมูล")
	}
	docCheck.ReceiveDate = receiveDate
	docCheck.ReceiveTime = receiveTime
	docCheck.ReceiveRemark = remark
	docCheck.ReceiveUser = &user
	o.Begin()
	_, err := o.Update(docCheck)
	if err != nil {
		o.Rollback()
	} else {
		o.Commit()
	}
	errRet = err
	return retID, errRet
}

//GetReceiveList GetReceiveList
func GetReceiveList(currentPage, lineSize uint, txtDateBegin, txtDateEnd, receiveState, payState, term string) (num int64, receiveList []Receive, err error) {
	o := orm.NewOrm()
	var sql = `SELECT 
					T0.i_d,
					T0.doc_no,
					T0.doc_date,
					T0.member_name,
					T0.receive_date,
					T0.remark,				 
					T0.active,
					T0.total_amount,
					(select sum(T1.total_amount) from payment T1 where T1.doc_no = T0.doc_no) as total_pay
			   FROM receive T0 	   
			   WHERE (lower(T0.doc_no) like lower(?) or lower(T0.remark) like lower(?) or lower(T0.member_name) like lower(?)) `
	if txtDateBegin != "" && txtDateEnd != "" {
		sql += " and date( T0.doc_date) between '" + helpers.ParseDateString(txtDateBegin) + "' and '" + helpers.ParseDateString(txtDateEnd) + "'"
	}
	if receiveState == "1" {
		sql += " and T0.receive_date is not null and T0.active"
	}
	if receiveState == "2" {
		sql += " and T0.receive_date is null and T0.active"
	}
	if payState == "1" {
		sql += " and (T0.total_amount - (select sum(T1.total_amount) from payment T1 where T1.doc_no = T0.doc_no) ) <= 0 and T0.active"
	}
	if payState == "2" {
		sql += ` and (T0.total_amount - (select sum(T1.total_amount) from payment T1 where T1.doc_no = T0.doc_no)) > 0 
				 and (T0.total_amount - (select sum(T1.total_amount) from payment T1 where T1.doc_no = T0.doc_no)) <> T0.total_amount
				 and T0.active`
	}
	if payState == "3" {
		sql += " and (T0.total_amount - (select sum(T1.total_amount) from payment T1 where T1.doc_no = T0.doc_no)) = T0.total_amount and T0.active"
	}
	num, _ = o.Raw(sql, "%"+term+"%", "%"+term+"%", "%"+term+"%").QueryRows(&receiveList)
	sql += " order by T0.doc_no  offset ? limit ?"
	_, err = o.Raw(sql, "%"+term+"%", "%"+term+"%", "%"+term+"%", currentPage, lineSize).QueryRows(&receiveList)

	return num, receiveList, err
}
