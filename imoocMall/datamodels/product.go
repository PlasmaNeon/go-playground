package datamodels

type Product struct {
	ID           int64  `json:"id" sql:"id" imooc:"id"`
	ProductName  string `json:"product-name" sql:"product-name" imooc:"product-name"`
	ProductNum   int64  `json:"product-num" sql:"product-num" imooc:"product-num"`
	ProductImage string `json:"product-image" sql:"product-image" imooc:"product-image"`
	ProductUrl   string `json:"product-url" sql:"product-url" imooc:"product-url"`
}
