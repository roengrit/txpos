package models

import (
	"errors"
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

//CategoryListJSON CategoryListJSON
type CategoryListJSON struct {
	CategoryList *[]Category
	Paging       string
	Page         uint
	PerPage      uint
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

//GetProductCategoryList `select * from x offset $1 limit $2`
func GetProductCategoryList(currentPage, lineSize uint, term string) (num int64, categoryList []Category, err error) {
	o := orm.NewOrm()
	var sql = `SELECT 
					T0.i_d,
					T0.name,
					T0.lock					 
			   FROM category T0	    
			   WHERE (lower(T0.name) like lower(?)) `

	num, _ = o.Raw(sql, "%"+term+"%").QueryRows(&categoryList)
	sql += " order by T0.name  offset ? limit ?"
	_, err = o.Raw(sql, "%"+term+"%", currentPage, lineSize).QueryRows(&categoryList)

	return num, categoryList, err
}

//GetProductCategory GetProductCategory
func GetProductCategory(ID int) (category *Category, errRet error) {
	Category := &Category{}
	o := orm.NewOrm()
	o.QueryTable("category").Filter("ID", ID).RelatedSel().One(Category)
	return Category, errRet
}

//CreateProductCategory _
func CreateProductCategory(Category Category) (ID int64, err error) {
	o := orm.NewOrm()
	o.Begin()
	ID, err = o.Insert(&Category)
	o.Commit()
	return ID, err
}

//UpdateProductCategory _
func UpdateProductCategory(category Category) (errRet error) {
	o := orm.NewOrm()
	getUpdate, _ := GetProductCategory(category.ID)
	if getUpdate.Lock {
		errRet = errors.New("ข้อมูลถูก Lock ไม่สามารถแก้ไขได้")
	}
	if getUpdate == nil {
		errRet = errors.New("ไม่พบข้อมูล")
	} else if errRet == nil {
		category.CreatedAt = getUpdate.CreatedAt
		category.Creator = getUpdate.Creator
		if num, errUpdate := o.Update(&category); errUpdate != nil {
			errRet = errUpdate
			_ = num
		}
	}
	return errRet
}

//DeleteProductCategory DeleteProductCategory
func DeleteProductCategory(ID int) (errRet error) {
	o := orm.NewOrm()
	unitDelete, _ := GetProductCategory(ID)
	if unitDelete.Lock {
		errRet = errors.New("ข้อมูลถูก Lock ไม่สามารถแก้ไขได้")
	}
	if num, errDelete := o.Delete(&Category{ID: ID}); errDelete != nil && errRet == nil {
		errRet = errDelete
		_ = num
	}
	return errRet
}
