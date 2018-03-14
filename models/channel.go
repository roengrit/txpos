package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

//Channel _
type Channel struct {
	ID        int
	Lock      bool
	Name      string `orm:"size(300)"`
	Active    bool
	Creator   *User     `orm:"rel(fk)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	Editor    *User     `orm:"null;rel(fk)"`
	EditedAt  time.Time `orm:"null;auto_now;type(datetime)"`
}

//ChannelListJSON ChannelListJSON
type ChannelListJSON struct {
	ChannelList *[]Channel
	Paging      string
	Page        uint
	PerPage     uint
}

func init() {
	orm.RegisterModel(new(Channel))
}

//GetAllChannel GetAllChannel
func GetAllChannel() *[]Channel {
	Channel := &[]Channel{}
	o := orm.NewOrm()
	o.QueryTable("channel").RelatedSel().All(Channel)
	return Channel
}

//GetChannelList `select * from x offset $1 limit $2`
func GetChannelList(currentPage, lineSize uint, term string) (num int64, channelList []Channel, err error) {
	o := orm.NewOrm()
	var sql = `SELECT 
					T0.i_d,
					T0.name,
					T0.lock					 
			   FROM channel T0	    
			   WHERE (lower(T0.name) like lower(?)) `

	num, _ = o.Raw(sql, "%"+term+"%").QueryRows(&channelList)
	sql += " order by T0.name  offset ? limit ?"
	_, err = o.Raw(sql, "%"+term+"%", currentPage, lineSize).QueryRows(&channelList)

	return num, channelList, err
}

//GetChannel GetChannel
func GetChannel(ID int) (channel *Channel, errRet error) {
	Channel := &Channel{}
	o := orm.NewOrm()
	o.QueryTable("channel").Filter("ID", ID).RelatedSel().One(Channel)
	return Channel, errRet
}

//CreateChannel _
func CreateChannel(Channel Channel) (ID int64, err error) {
	o := orm.NewOrm()
	o.Begin()
	ID, err = o.Insert(&Channel)
	o.Commit()
	return ID, err
}

//UpdateChannel _
func UpdateChannel(channel Channel) (errRet error) {
	o := orm.NewOrm()
	getUpdate, _ := GetChannel(channel.ID)
	if getUpdate.Lock {
		errRet = errors.New("ข้อมูลถูก Lock ไม่สามารถแก้ไขได้")
	}
	if getUpdate == nil {
		errRet = errors.New("ไม่พบข้อมูล")
	} else if errRet == nil {
		channel.CreatedAt = getUpdate.CreatedAt
		channel.Creator = getUpdate.Creator
		if num, errUpdate := o.Update(&channel); errUpdate != nil {
			errRet = errUpdate
			_ = num
		}
	}
	return errRet
}

//DeleteChannel DeleteChannel
func DeleteChannel(ID int) (errRet error) {
	o := orm.NewOrm()
	unitDelete, _ := GetChannel(ID)
	if unitDelete.Lock {
		errRet = errors.New("ข้อมูลถูก Lock ไม่สามารถแก้ไขได้")
	}
	if num, errDelete := o.Delete(&Channel{ID: ID}); errDelete != nil && errRet == nil {
		errRet = errDelete
		_ = num
	}
	return errRet
}
