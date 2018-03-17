package models

import (
	"errors"
	"time"
	"txpos/helpers"

	"github.com/astaxie/beego/orm"
)

//Sale _
type Sale struct {
	ID       int
	Flag     int
	Active   bool
	DocType  int
	SaleType int
	VatType  int
	IsPos    bool
	DocNo                string    `orm:"size(30)"`
	DocDate              time.Time `form:"-" orm:"null"`
	DocTime              string    `orm:"size(6)"`
	DocRefNo             string    `orm:"size(30)"`
	Channel              *Channel  `orm:"rel(fk)"`
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
	ReceiveMoney         float64 `orm:"digits(12);decimals(2)"`
	ChangeMoney          float64 `orm:"digits(12);decimals(2)"`
	CreditDay            int
	CreditDate           time.Time `form:"-" orm:"type(date)"`
	SendDate             time.Time `form:"-" orm:"null;type(datetime)"`
	SendTime             string    `orm:"size(6)"`
	SendUser             *User     `orm:"null;rel(fk)"`
	Remark               string    `orm:"size(300)"`
	CancelRemark         string    `orm:"size(300)"`
	SendRemark           string    `orm:"size(300)"`
	Creator              *User     `orm:"rel(fk)"`
	CreatedAt            time.Time `orm:"auto_now_add;type(datetime)"`
	Editor               *User     `orm:"null;rel(fk)"`
	EditedAt             time.Time `orm:"null;auto_now;type(datetime)"`
	CancelUser           *User     `orm:"null;rel(fk)"`
	CancelAt             time.Time `orm:"null;type(datetime)"`
	SaleSub              []SaleSub `orm:"-"`
	Payment              []Payment `orm:"-"`
}

//SaleSub _
type SaleSub struct {
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

//SaleJSON SaleJSON
type SaleJSON struct {
	SaleList *[]Sale
	Paging   string
	Page     uint
	PerPage  uint
}

func init() {
	orm.RegisterModel(new(Sale), new(SaleSub)) // Need to register model in init
}

//CreateSale _
func CreateSale(sale Sale, user User) (retID int64, errRet error) {
	sale.DocNo = helpers.GetMaxDoc("sale", "IV")
	sale.Creator = &user
	sale.CreatedAt = time.Now()
	sale.CreditDay = 0
	sale.CreditDate = time.Now()
	sale.Active = true
	var fullDataSub []SaleSub
	var fullDataPayment []Payment
	for _, val := range sale.SaleSub {
		if val.Product.ID != 0 {
			val.CreatedAt = time.Now()
			val.Creator = &user
			val.DocNo = sale.DocNo
			val.Flag = sale.Flag
			val.Active = true
			val.DocDate = sale.DocDate
			val.ProductName = val.Product.Name
			fullDataSub = append(fullDataSub, val)
		}
	}
	for _, val := range sale.Payment {
		val.CreatedAt = time.Now()
		val.Creator = &user
		val.DocNo = sale.DocNo
		fullDataPayment = append(fullDataPayment, val)
	}
	o := orm.NewOrm()
	o.Begin()
	id, err := o.Insert(&sale)
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

//UpdateSale _
func UpdateSale(sale Sale, user User) (retID int64, errRet error) {
	docCheck, _ := GetSale(sale.ID)
	if docCheck.ID == 0 {
		errRet = errors.New("ไม่พบข้อมูล")
	}
	sale.Creator = docCheck.Creator
	sale.CreatedAt = docCheck.CreatedAt
	sale.CreditDay = docCheck.CreditDay
	sale.CreditDate = docCheck.CreditDate
	sale.EditedAt = time.Now()
	sale.Editor = &user
	sale.Active = docCheck.Active
	var fullDataSub []SaleSub
	var fullDataPayment []Payment
	for _, val := range sale.SaleSub {
		if val.Product.ID != 0 {
			val.CreatedAt = time.Now()
			val.Creator = &user
			val.DocNo = sale.DocNo
			val.Flag = sale.Flag
			val.Active = sale.Active
			val.DocDate = sale.DocDate
			val.ProductName = val.Product.Name
			fullDataSub = append(fullDataSub, val)
		}
	}
	for _, val := range sale.Payment {
		val.CreatedAt = time.Now()
		val.Creator = &user
		val.DocNo = sale.DocNo
		fullDataPayment = append(fullDataPayment, val)
	}
	o := orm.NewOrm()
	o.Begin()
	id, err := o.Update(&sale)
	if err == nil {
		_, err = o.QueryTable("sale_sub").Filter("doc_no", sale.DocNo).Delete()
	}
	if err == nil {
		_, err = o.InsertMulti(len(fullDataSub), fullDataSub)
	}
	if err == nil {
		_, err = o.QueryTable("payment").Filter("doc_no", sale.DocNo).Delete()
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

//GetSale _
func GetSale(ID int) (doc *Sale, errRet error) {
	saleDoc := &Sale{}
	o := orm.NewOrm()
	o.QueryTable("sale").Filter("ID", ID).RelatedSel().One(saleDoc)
	o.QueryTable("sale_sub").Filter("doc_no", saleDoc.DocNo).RelatedSel().All(&saleDoc.SaleSub)
	o.QueryTable("payment").Filter("doc_no", saleDoc.DocNo).OrderBy("ID").RelatedSel().All(&saleDoc.Payment)
	doc = saleDoc
	return doc, errRet
}

//UpdateCancelSale _
func UpdateCancelSale(ID int, remark string, user User) (retID int64, errRet error) {
	docCheck := &Sale{}
	o := orm.NewOrm()
	o.QueryTable("sale").Filter("ID", ID).RelatedSel().One(docCheck)
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
		_, err = o.Raw("update sale_sub set active = false where doc_no = ?", docCheck.DocNo).Exec()
	}
	if err != nil {
		o.Rollback()
	} else {
		o.Commit()
	}
	errRet = err
	return retID, errRet
}

//UpdateSendDate _
func UpdateSendDate(ID int, saleDate time.Time, saleTime, remark string, user User) (retID int64, errRet error) {
	docCheck := &Sale{}
	o := orm.NewOrm()
	o.QueryTable("sale").Filter("ID", ID).RelatedSel().One(docCheck)
	if docCheck.ID == 0 {
		errRet = errors.New("ไม่พบข้อมูล")
	}
	docCheck.SendDate = saleDate
	docCheck.SendTime = saleTime
	docCheck.SendRemark = remark
	docCheck.SendUser = &user
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

//GetSaleList GetSaleList
func GetSaleList(currentPage, lineSize uint, txtDateBegin, txtDateEnd, saleState, payState, term string) (num int64, saleList []Sale, err error) {
	o := orm.NewOrm()
	var sql = `SELECT 
					T0.i_d,
					T0.doc_no,
					T0.doc_date,
					T0.member_name,
					T0.sale_date,
					T0.remark,				 
					T0.active,
					T0.total_amount,
					(select sum(T1.total_amount) from payment T1 where T1.doc_no = T0.doc_no) as total_pay
			   FROM sale T0 	   
			   WHERE (lower(T0.doc_no) like lower(?) or lower(T0.remark) like lower(?) or lower(T0.member_name) like lower(?)) `
	if txtDateBegin != "" && txtDateEnd != "" {
		sql += " and date( T0.doc_date) between '" + helpers.ParseDateString(txtDateBegin) + "' and '" + helpers.ParseDateString(txtDateEnd) + "'"
	}
	if saleState == "1" {
		sql += " and T0.sale_date is not null and T0.active"
	}
	if saleState == "2" {
		sql += " and T0.sale_date is null and T0.active"
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
	num, _ = o.Raw(sql, "%"+term+"%", "%"+term+"%", "%"+term+"%").QueryRows(&saleList)
	sql += " order by T0.doc_no  offset ? limit ?"
	_, err = o.Raw(sql, "%"+term+"%", "%"+term+"%", "%"+term+"%", currentPage, lineSize).QueryRows(&saleList)

	return num, saleList, err
}
