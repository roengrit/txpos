package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

//Member _
type Member struct {
	ID         int
	Lock       bool
	Name       string     `orm:"size(300)"`
	Address    string     `orm:"size(300)"`
	Province   *Provinces `orm:"rel(fk)"`
	PostCode   string     `orm:"size(10)"`
	Contact    string     `orm:"size(255)"`
	Tel        string     `orm:"size(100)"`
	MemberType int
	TaxNo      string `orm:"size(100)"`
	Level      int
	Score      int
	Remark     string    `orm:"size(100)"`
	Creator    *User     `orm:"rel(fk)"`
	CreatedAt  time.Time `orm:"auto_now_add;type(datetime)"`
	Editor     *User     `orm:"null;rel(fk)"`
	EditedAt   time.Time `orm:"null;auto_now;type(datetime)"`
}

//MemberJSON MemberJSON
type MemberJSON struct {
	MemberList *[]Member
	Paging     string
	Page       uint
	PerPage    uint
}

func init() {
	orm.RegisterModel(new(Member)) //Need to register model in init
}

//CreateMember _
func CreateMember(sup Member) (retID int64, errRet error) {
	o := orm.NewOrm()
	o.Begin()
	id, err := o.Insert(&sup)
	o.Commit()
	if err == nil {
		retID = id
	}
	return retID, err
}

//UpdateMember _
func UpdateMember(mem Member) (errRet error) {
	o := orm.NewOrm()
	getUpdate, _ := GetMember(mem.ID)
	if getUpdate.Lock {
		errRet = errors.New("ข้อมูลถูก Lock ไม่สามารถแก้ไขได้")
	}
	if getUpdate == nil {
		errRet = errors.New("ไม่พบข้อมูล")
	} else if errRet == nil {
		mem.CreatedAt = getUpdate.CreatedAt
		mem.Creator = getUpdate.Creator
		if num, errUpdate := o.Update(&mem); errUpdate != nil {
			errRet = errUpdate
			_ = num
		}
	}
	return errRet
}

//GetMember _
func GetMember(ID int) (mem *Member, errRet error) {
	Member := &Member{}
	o := orm.NewOrm()
	o.QueryTable("member").Filter("ID", ID).RelatedSel().One(Member)
	return Member, errRet
}

//GetMemberList _
func GetMemberList(term string, limit int) (mem *[]Member, rowCount int, errRet error) {
	Member := &[]Member{}
	o := orm.NewOrm()
	qs := o.QueryTable("member")
	cond := orm.NewCondition()
	cond1 := cond.Or("Name__icontains", term).
		Or("Tel__icontains", term).
		Or("Contact__icontains", term).
		Or("Remark__icontains", term).
		Or("Address__icontains", term)
	qs.SetCond(cond1).RelatedSel().Limit(limit).All(Member)
	return Member, len(*Member), errRet
}

//GetMemberListPaging `select * from x offset $1 limit $2`
func GetMemberListPaging(currentPage, lineSize uint, term string) (num int64, memberList []Member, err error) {
	o := orm.NewOrm()
	var sql = `SELECT 
					T0.i_d,
					T0.name       ,
					T0.address    ,
					T0.post_code    ,
					T0.contact     ,
					T0.tel         ,
					T0.member_type  
			   FROM member T0	    
			   WHERE (lower(T0.name) like lower(?) or lower(T0.tel) like lower(?) 
			          or lower(T0.contact) like lower(?) or lower(T0.remark) like lower(?) or lower(T0.address) like lower(?) ) `
	num, _ = o.Raw(sql, "%"+term+"%", "%"+term+"%", "%"+term+"%", "%"+term+"%", "%"+term+"%").QueryRows(&memberList)
	sql += " order by T0.name  offset ? limit ?"
	_, err = o.Raw(sql, "%"+term+"%", "%"+term+"%", "%"+term+"%", "%"+term+"%", "%"+term+"%", currentPage, lineSize).QueryRows(&memberList)
	return num, memberList, err
}

//DeleteMember _
func DeleteMember(ID int) (errRet error) {
	o := orm.NewOrm()
	getUpdate, _ := GetMember(ID)
	if getUpdate.Lock {
		errRet = errors.New("ข้อมูลถูก Lock ไม่สามารถแก้ไขได้")
	}
	if num, errDelete := o.Delete(&Member{ID: ID}); errDelete != nil && errRet == nil {
		errRet = errDelete
		_ = num
	}
	return errRet
}
