package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"imoocMall/backend/web/controllers"
	"imoocMall/common"
	"imoocMall/repositories"
	"imoocMall/services"
	"log"
)

func main() {
	// create iris instance
	app := iris.New()
	// set error level, show errors under mvc mode.
	app.Logger().SetLevel("debug")
	// register template
	template := iris.HTML("./backend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	app.HandleDir("/assets", "./backend/web/assets")
	// error redirection

	//app.OnAnyErrorCode(func(ctx iris.Context) {
	//	ctx.ViewData("message", ctx.Values().GetStringDefault("message", "wrong page"))
	//	ctx.ViewLayout("")
	//	ctx.View("shared/error.html")
	//})

	// Controller
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//5.注册控制器
	productRepository := repositories.NewProductManager("product", db)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))

	orderRepository := repositories.NewOrderMangerRepository("order", db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))
	// Start
	app.Run(
		iris.Addr("localhost:18080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
