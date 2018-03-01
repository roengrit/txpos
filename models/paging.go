package models

import (
	"github.com/astaxie/beego/orm"
)

//CountRow `select count(*) from user_profile;`
func CountRow(countSQL string) (uint, error) {
	o := orm.NewOrm()
	var count uint
	err := o.Raw(countSQL).QueryRow(&count)
	if err != nil {
		return 0, err
	}
	return count, err
}
