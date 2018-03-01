package models

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
)

//User User
type User struct {
	ID        int
	Lock      bool
	Username  string `orm:"size(50)"`
	Password  string `orm:"size(255)"`
	Name      string `orm:"size(500)"`
	Tel       string `orm:"size(255)"`
	Facebook  string `orm:"size(255)"`
	Line      string `orm:"size(255)"`
	Role      *Role  `orm:"rel(fk)"`
	Active    bool
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	EditedAt  time.Time `orm:"null;auto_now;type(datetime)"`
}

//Role Role
type Role struct {
	ID        int
	Lock      bool
	Name      string    `orm:"size(225)"`
	Creator   *User     `orm:"rel(fk)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	Editor    *User     `orm:"null;rel(fk)"`
	EditedAt  time.Time `orm:"null;auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(User), new(Role))
}

//Login Login
func Login(username, password string) (ok bool, errRet string) {
	o := orm.NewOrm()
	user := User{Username: username}
	err := o.Read(&user, "Username")
	if err == orm.ErrNoRows {
		errRet = "รหัสผู้ใช้/รหัสผ่านไม่ถูกต้อง"
	} else {
		if errCript := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); errCript != nil {
			errRet = "รหัสผู้ใช้/รหัสผ่านไม่ถูกต้อง"
		} else {
			ok = true
		}
	}
	return ok, errRet
}

//GetUser GetUser
func GetUser(username string) (userRet User, errRet error) {
	o := orm.NewOrm()
	user := User{Username: username}
	errRet = o.Read(&user, "Username")
	if errRet == orm.ErrNoRows {
		errRet = errors.New("ไม่พบผู้ใช้งานนี้ในระบบ")
	}
	fmt.Println("DEV")
	user.ID = 1
	user.Name = "Admin"
	return user, errRet
}

//GetUserByID GetUserByID
func GetUserByID(ID int) (user *User, errRet error) {
	o := orm.NewOrm()
	userGet := &User{}
	o.QueryTable("users").Filter("ID", ID).RelatedSel().One(userGet)
	if nil != userGet {
		userGet.Password = ""
	} else {
		errRet = errors.New("ไม่พบผู้ใช้งานนี้ในระบบ")
	}
	return userGet, errRet
}

//GetUserByUserName _
func GetUserByUserName(username string) (user *User, errRet error) {
	userGet := &User{}
	o := orm.NewOrm()
	o.QueryTable("users").Filter("Username", username).RelatedSel().One(userGet)
	if nil != userGet {
		userGet.Password = ""
	} else {
		errRet = errors.New("ไม่พบผู้ใช้งานนี้ในระบบ")
	}
	return userGet, errRet
}

//ForgetPass ForgetPass
func ForgetPass(username, newpass string) (ok bool, errRet error) {
	o := orm.NewOrm()
	user := User{Username: username}
	errRet = o.Read(&user, "Username")
	if errRet == orm.ErrNoRows {
		errRet = errors.New("ไม่พบผู้ใช้งานนี้ในระบบ")
	} else {
		if hash, err := bcrypt.GenerateFromPassword([]byte(newpass), bcrypt.DefaultCost); err != nil {
			errRet = err
		} else {
			user.Password = string(hash)
			if _, errUpdate := o.Update(&user); errUpdate != nil {
				errRet = errUpdate
			} else {
				ok = true
			}
		}
	}
	return ok, errRet
}

//ChangePass ChangePass
func ChangePass(username, newpass string) (ok bool, errRet error) {
	o := orm.NewOrm()
	user := User{Username: username}
	errRet = o.Read(&user, "Username")
	if errRet == orm.ErrNoRows {
		errRet = errors.New("ไม่พบผู้ใช้งานนี้ในระบบ")
	} else {
		if hash, err := bcrypt.GenerateFromPassword([]byte(newpass), bcrypt.DefaultCost); err != nil {
			errRet = err
		} else {
			user.Password = string(hash)
			if num, errUpdate := o.Update(&user); errUpdate != nil {
				errRet = errUpdate
				_ = num
			} else {
				ok = true
			}
		}
	}
	return ok, errRet
}

//RandStringRunes _
func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
