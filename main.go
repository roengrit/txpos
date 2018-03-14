package main

import (
	"fmt"
	c "txpos/controllers"
	h "txpos/helpers"
	_ "txpos/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
)

func init() {
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", "host=localhost port=5432 user=postgres password=P@ssw0rd dbname=txpos sslmode=disable")
}

func main() {

	name := "default"
	force := false                             // Drop table and re-create.
	verbose := true                            // Print log.
	err := orm.RunSyncdb(name, force, verbose) // Error.

	if err != nil {
		fmt.Println(err)
	}

	beego.Router("/", &c.AppController{})

	beego.Router("/company", &c.CompanyController{}, "get:CreateCom;post:UpdateCom")
	beego.Router("/setting", &c.SettingController{}, "get:CreateSetting;post:UpdateSetting")

	beego.Router("/service/secure/json/", &c.ServiceController{}, "get:GetXSRF")

	beego.Router("/product/?:id", &c.ProductController{}, "get:CreateProduct;post:UpdateProduct;delete:DeleteProduct")
	beego.Router("/product/list", &c.ProductController{}, "get:ProductList;post:GetProductListJSON")
	beego.Router("/product/list/json", &c.ProductController{}, "get:ListProductJSON")
	beego.Router("/product/json", &c.ProductController{}, "get:GetProductJSON")
	beego.Router("/product/json/code", &c.ProductController{}, "get:GetProductByCodeJSON")

	beego.Router("/unit/?:id", &c.UnitController{}, "get:CreateUnit;post:UpdateUnit;delete:DeleteUnit")
	beego.Router("/unit/list", &c.UnitController{}, "get:UnitList;post:GetUnitListJSON")

	beego.Router("/category/?:id", &c.CategoryController{}, "get:CreateCategory;post:UpdateCategory;delete:DeleteCategory")
	beego.Router("/category/list", &c.CategoryController{}, "get:CategoryList;post:GetCategoryListJSON")

	beego.Router("/channel/?:id", &c.ChannelController{}, "get:CreateChannel;post:UpdateChannel;delete:DeleteChannel")
	beego.Router("/channel/list", &c.ChannelController{}, "get:ChannelList;post:GetChannelListJSON")

	beego.Router("/member/?:id", &c.MemberController{}, "get:CreateMember;post:UpdateMember;delete:DeleteMember")
	beego.Router("/member/read/?:id", &c.MemberController{}, "get:CreateMember")
	beego.Router("/member/list", &c.MemberController{}, "get:MemberList;post:GetMemberListJSON")
	beego.Router("/member/list/json", &c.MemberController{}, "get:ListMemberJSON")
	beego.Router("/member/json/?:id", &c.MemberController{}, "get:GetMemberJSON")

	beego.Router("/stock", &c.StockCountController{})
	beego.Router("/stock/diff", &c.StockCountController{}, "get:StockDiff")
	beego.Router("/stock/cancel", &c.StockCountController{}, "get:CancelStockCount;post:UpdateCancelStockCount")
	beego.Router("/stock/active", &c.StockCountController{}, "post:UpdateActiveStockCount")
	beego.Router("/stock/list", &c.StockCountController{}, "get:StockList;post:GetStockListJSON")

	beego.Router("/pickup", &c.PickUpController{})
	beego.Router("/pickup/cancel", &c.PickUpController{}, "get:CancelPickUp;post:UpdateCancelPickUp")
	beego.Router("/pickup/active", &c.PickUpController{}, "post:UpdateActivePickUp")
	beego.Router("/pickup/list", &c.PickUpController{}, "get:PickUpList;post:GetPickUpListJSON")

	beego.Router("/receive", &c.ReceiveController{})
	beego.Router("/receive/cancel", &c.ReceiveController{}, "get:CancelReceive;post:UpdateCancelReceive")
	beego.Router("/receive/receive", &c.ReceiveController{}, "get:Receive;post:UpdateReceive")
	beego.Router("/receive/print", &c.ReceiveController{}, "get:Print")
	beego.Router("/receive/list", &c.ReceiveController{}, "get:ReceiveList;post:GetReceiveListJSON")

	beego.AddFuncMap("ThCommaSeperate", h.ThCommaSeperate)
	beego.AddFuncMap("TextThCommaSeperate", h.ThCommaSeperate)
	beego.AddFuncMap("TextThCommaAndPercentSeperate", h.TextThCommaAndPercentSeperate)
	beego.AddFuncMap("HTMLRowOrder", h.HTMLRowOrder)
	beego.AddFuncMap("TextNoPercent", h.TextNoPercent)
	beego.Run()
}
