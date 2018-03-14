package models

import (
	"errors"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

//Product Product
type Product struct {
	ID            int
	Lock          bool
	Code          string    `orm:"size(13);unique"`
	Name          string    `orm:"size(300)"`
	Barcode       string    `orm:"size(13)"`
	AverageCost   float64   `orm:"digits(12);decimals(2)"`
	BalanceQty    float64   `orm:"digits(12);decimals(2)"`
	PurchasePrice float64   `orm:"digits(12);decimals(2)"`
	SalePrice     float64   `orm:"digits(12);decimals(2)"`
	Unit          *Unit     `orm:"rel(fk)"`
	Category      *Category `orm:"rel(fk)"`
	ProductType   int
	ImagePath1    string `orm:"size(300)"`
	ImagePath2    string `orm:"size(300)"`
	ImagePath3    string `orm:"size(300)"`
	ImageBase64   string `orm:"-"`
	DeleteImage   int    `orm:"-"`
	Remark        string `orm:"size(100)"`
	FixCost       bool
	Active        bool
	Creator       *User     `orm:"rel(fk)"`
	CreatedAt     time.Time `orm:"auto_now_add;type(datetime)"`
	Editor        *User     `orm:"null;rel(fk)"`
	EditedAt      time.Time `orm:"null;auto_now;type(datetime)"`
}

//ProductList ProductList
type ProductList struct {
	ID           int
	Lock         bool
	Code         string
	Name         string
	Barcode      string
	AverageCost  float64
	BalanceQty   float64
	CategoryName string
	UnitID       int
	UnitName     string
	Active       bool
}

//ProductListJSON ProductListJSON
type ProductListJSON struct {
	ProductList *[]ProductList
	Paging      string
	Page        uint
	PerPage     uint
}

func init() {
	orm.RegisterModel(new(Product))
}

//GetProductList `select * from x offset $1 limit $2`
func GetProductList(currentPage, lineSize uint, statusTerm string, categoryTerm, term string) (num int64, productList []ProductList, err error) {
	o := orm.NewOrm()
	var sql = `SELECT 
					T0.i_d,
					T0.code,
					T0.name,
					T0.lock,
					T0.active,
					T0.balance_qty,
					T0.average_cost,
					T0.barcode,
					T1.i_d as category_i_d,
					T1.name as category_name,
					T2.i_d as unit_i_d,
					T2.name as unit_name
			   FROM product T0	
			   JOIN category T1 ON T0.category_id = T1.i_d	
			   JOIN unit T2 ON T0.unit_id = T2.i_d    
			   WHERE (lower(T0.name) like lower(?) or lower(T0.barcode) like lower(?) or lower(T0.code) like lower(?) ) `

	if statusTerm == "0" {
		sql += " and NOT T0.Active"
	}
	if statusTerm == "1" {
		sql += " and T0.Active"
	}
	if categoryTerm != "" {
		sql += " and T1.i_d = (" + categoryTerm + ")"
	}

	num, _ = o.Raw(sql, "%"+term+"%", "%"+term+"%", "%"+term+"%").QueryRows(&productList)
	sql += " order by T0.name  offset ? limit ?"
	_, err = o.Raw(sql, "%"+term+"%", "%"+term+"%", "%"+term+"%", currentPage, lineSize).QueryRows(&productList)

	return num, productList, err
}

//GetProduct _
func GetProduct(ID int) (pro *Product, errRet error) {
	Product := &Product{}
	o := orm.NewOrm()
	o.QueryTable("product").Filter("ID", ID).RelatedSel().One(Product)
	return Product, errRet
}

//GetProductByCode _
func GetProductByCode(code string) (pro *Product, errRet error) {
	Product := &Product{}
	o := orm.NewOrm()
	o.QueryTable("product").Filter("Code", code).RelatedSel().One(Product)
	return Product, errRet
}

//CreateProduct _
func CreateProduct(Product Product) (ID int64, err error) {
	o := orm.NewOrm()
	o.Begin()
	ID, err = o.Insert(&Product)
	o.Commit()
	if err != nil {
		if strings.Contains(err.Error(), "product_code_key") && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			err = errors.New("รหัสสินค้าซ้ำ")
		}
	}
	return ID, err
}

//UpdateProduct _
func UpdateProduct(pro Product, isNewImage bool) (errRet error) {
	o := orm.NewOrm()
	getUpdate, _ := GetProduct(pro.ID)
	if getUpdate.Lock {
		errRet = errors.New("ข้อมูลถูก Lock ไม่สามารถแก้ไขได้")
	}
	if getUpdate == nil {
		errRet = errors.New("ไม่พบข้อมูล")
	} else if errRet == nil {
		if !isNewImage {
			pro.ImagePath1 = getUpdate.ImagePath1
		}
		pro.BalanceQty = getUpdate.BalanceQty
		if !pro.FixCost {
			pro.AverageCost = getUpdate.AverageCost
		}
		pro.CreatedAt = getUpdate.CreatedAt
		pro.Creator = getUpdate.Creator
		if num, errUpdate := o.Update(&pro); errUpdate != nil {
			errRet = errUpdate
			_ = num
		}
	}
	return errRet
}

//DeleteProduct _
func DeleteProduct(ID int) (errRet error) {
	o := orm.NewOrm()
	getUpdate, _ := GetProduct(ID)
	if getUpdate.Lock {
		errRet = errors.New("ข้อมูลถูก Lock ไม่สามารถแก้ไขได้")
	}
	if num, errDelete := o.Delete(&Product{ID: ID}); errDelete != nil && errRet == nil {
		errRet = errDelete
		_ = num
	}
	return errRet
}

//GetMaxItemCode _
func GetMaxItemCode(format string) (code string) {

	var lists []orm.ParamsList
	var sql = `SELECT 
				 concat( '` + format + `' ,
					  LPAD((COALESCE( MAX ( SUBSTRING ( code FROM '[0-9]{5,}$' )), '0' ) :: NUMERIC + 1)  :: text, 5, '0')
					  ) AS  code 
				 FROM
					 product
				 WHERE
				     code LIKE '` + format + `%'`
	o := orm.NewOrm()
	num, err := o.Raw(sql).ValuesList(&lists)
	if err == nil && num > 0 {
		max := lists[0]
		maxVal := max[0].(string)
		code = maxVal
	} else {
		code = ""
	}
	return code
}
