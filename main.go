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
	beego.Router("/service/secure/json/", &c.ServiceController{}, "get:GetXSRF")

	beego.Router("/product/?:id", &c.ProductController{}, "get:CreateProduct;post:UpdateProduct;delete:DeleteProduct")
	beego.Router("/product/list", &c.ProductController{}, "get:ProductList;post:GetProductListJSON")

	beego.Router("/unit/?:id", &c.UnitController{}, "get:CreateUnit;post:UpdateUnit;delete:DeleteUnit")
	beego.Router("/unit/list", &c.UnitController{}, "get:UnitList;post:GetUnitListJSON")

	beego.Router("/category/?:id", &c.CategoryController{}, "get:CreateCategory;post:UpdateCategory;delete:DeleteCategory")
	beego.Router("/category/list", &c.CategoryController{}, "get:CategoryList;post:GetCategoryListJSON")

	beego.AddFuncMap("ThCommaSeperate", h.ThCommaSeperate)
	beego.AddFuncMap("HTMLRowOrder", h.HTMLRowOrder)
	beego.Run()
}
