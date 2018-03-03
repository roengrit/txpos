package models

import (
	"errors"
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

//UnitListJSON UnitListJSON
type UnitListJSON struct {
	UnitList *[]Unit
	Paging   string
	Page     uint
	PerPage  uint
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

//GetProductUnitList `select * from x offset $1 limit $2`
func GetProductUnitList(currentPage, lineSize uint, term string) (num int64, unitList []Unit, err error) {
	o := orm.NewOrm()
	var sql = `SELECT 
					T0.i_d,
					T0.name,
					T0.lock					 
			   FROM unit T0	    
			   WHERE (lower(T0.name) like lower(?)) `

	num, _ = o.Raw(sql, "%"+term+"%").QueryRows(&unitList)
	sql += " order by T0.name  offset ? limit ?"
	_, err = o.Raw(sql, "%"+term+"%", currentPage, lineSize).QueryRows(&unitList)

	return num, unitList, err
}

//GetProductUnit GetProductUnit
func GetProductUnit(ID int) (unit *Unit, errRet error) {
	Unit := &Unit{}
	o := orm.NewOrm()
	o.QueryTable("unit").Filter("ID", ID).RelatedSel().One(Unit)
	return Unit, errRet
}

//CreateProductUnit _
func CreateProductUnit(Unit Unit) (ID int64, err error) {
	o := orm.NewOrm()
	o.Begin()
	ID, err = o.Insert(&Unit)
	o.Commit()
	return ID, err
}

//UpdateProductUnit _
func UpdateProductUnit(unit Unit) (errRet error) {
	o := orm.NewOrm()
	getUpdate, _ := GetProductUnit(unit.ID)
	if getUpdate.Lock {
		errRet = errors.New("ข้อมูลถูก Lock ไม่สามารถแก้ไขได้")
	}
	if getUpdate == nil {
		errRet = errors.New("ไม่พบข้อมูล")
	} else if errRet == nil {
		unit.CreatedAt = getUpdate.CreatedAt
		unit.Creator = getUpdate.Creator
		if num, errUpdate := o.Update(&unit); errUpdate != nil {
			errRet = errUpdate
			_ = num
		}
	}
	return errRet
}

//DeleteProductUnit DeleteProductUnit
func DeleteProductUnit(ID int) (errRet error) {
	o := orm.NewOrm()
	unitDelete, _ := GetProductUnit(ID)
	if unitDelete.Lock {
		errRet = errors.New("ข้อมูลถูก Lock ไม่สามารถแก้ไขได้")
	}
	if num, errDelete := o.Delete(&Unit{ID: ID}); errDelete != nil && errRet == nil {
		errRet = errDelete
		_ = num
	}
	return errRet
}
