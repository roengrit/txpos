package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

//Setting _
type Setting struct {
	ID       int
	Lock     bool
	VatValue float64   `orm:"digits(12);decimals(2)"`
	Editor   *User     `orm:"null;rel(fk)"`
	EditedAt time.Time `orm:"null;auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Setting)) //Need to register model in init
}

//CreateSetting _
func CreateSetting(setting Setting) (retID int64, errRet error) {
	o := orm.NewOrm()
	o.Begin()
	id, err := o.Insert(&setting)
	o.Commit()
	if err == nil {
		retID = id
	}
	return retID, err
}

//UpdateSetting _
func UpdateSetting(setting Setting) (errRet error) {
	o := orm.NewOrm()
	getUpdate, _ := GetCom(setting.ID)
	if getUpdate.Lock {
		errRet = errors.New("ข้อมูลถูก Lock ไม่สามารถแก้ไขได้")
	}
	if getUpdate == nil {
		errRet = errors.New("ไม่พบข้อมูล")
	} else if errRet == nil {
		if num, errUpdate := o.Update(&setting); errUpdate != nil {
			errRet = errUpdate
			_ = num
		}
	}
	return errRet
}

//GetSetting _
func GetSetting(ID int) (mem *Setting, errRet error) {
	company := &Setting{}
	o := orm.NewOrm()
	o.QueryTable("setting").Filter("ID", ID).RelatedSel().One(company)
	return company, errRet
}

//GetSettingFirst _
func GetSettingFirst() (sett *Setting, errRet error) {
	setting := &Setting{}
	o := orm.NewOrm()
	o.QueryTable("setting").RelatedSel().One(setting)
	return setting, errRet
}
