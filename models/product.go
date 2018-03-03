package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

//Product Product
type Product struct {
	ID          int
	Lock        bool
	Name        string    `orm:"size(300)"`
	Barcode     string    `orm:"size(13)"`
	AverageCost float64   `orm:"digits(12);decimals(2)"`
	BalanceQty  float64   `orm:"digits(12);decimals(2)"`
	SalePrice   float64   `orm:"digits(12);decimals(2)"`
	Unit        *Unit     `orm:"rel(fk)"`
	Category    *Category `orm:"rel(fk)"`
	ProductType int
	ImagePath1  string `orm:"size(300)"`
	ImagePath2  string `orm:"size(300)"`
	ImagePath3  string `orm:"size(300)"`
	ImageBase64 string `orm:"-"`
	DeleteImage int    `orm:"-"`
	Remark      string `orm:"size(100)"`
	FixCost     bool
	Active      bool
	Creator     *User     `orm:"rel(fk)"`
	CreatedAt   time.Time `orm:"auto_now_add;type(datetime)"`
	Editor      *User     `orm:"null;rel(fk)"`
	EditedAt    time.Time `orm:"null;auto_now;type(datetime)"`
}

//ProductList ProductList
type ProductList struct {
	ID           int
	Lock         bool
	Name         string
	Barcode      string
	AverageCost  float64
	BalanceQty   float64
	CategoryName string
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
					T0.name,
					T0.lock,
					T0.active,
					T0.balance_qty,
					T0.average_cost,
					T0.barcode,
					T1.i_d as category_id,
					T1.name as category_name
			   FROM product T0	
			   JOIN category T1 ON T0.category_id = T1.i_d	    
			   WHERE (lower(T0.name) like lower(?) or lower(T0.barcode) like lower(?) ) `

	if statusTerm == "0" {
		sql += " and NOT T0.Active"
	}
	if statusTerm == "1" {
		sql += " and T0.Active"
	}
	if categoryTerm != "" {
		sql += " and T1.i_d = (" + categoryTerm + ")"
	}

	num, _ = o.Raw(sql, "%"+term+"%", "%"+term+"%").QueryRows(&productList)
	sql += " order by T0.name  offset ? limit ?"
	_, err = o.Raw(sql, "%"+term+"%", "%"+term+"%", currentPage, lineSize).QueryRows(&productList)

	return num, productList, err
}

//GetProduct _
func GetProduct(ID int) (pro *Product, errRet error) {
	Product := &Product{}
	o := orm.NewOrm()
	o.QueryTable("product").Filter("ID", ID).RelatedSel().One(Product)
	return Product, errRet
}

//CreateProduct _
func CreateProduct(Product Product) (ID int64, err error) {
	o := orm.NewOrm()
	o.Begin()
	ID, err = o.Insert(&Product)
	o.Commit()
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
