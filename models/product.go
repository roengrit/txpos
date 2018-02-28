package models

import (
	"time"
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
	Remark      string `orm:"size(100)"`
	FixCost     bool
	Active      bool
	Creator     *User     `orm:"rel(fk)"`
	CreatedAt   time.Time `orm:"auto_now_add;type(datetime)"`
	Editor      *User     `orm:"null;rel(fk)"`
	EditedAt    time.Time `orm:"null;auto_now;type(datetime)"`
}
