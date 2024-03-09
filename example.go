// Package main
// Generated by sql2struct at 2024-03-09 16:16:59
package main

import "time"

// Products Generated by sql2struct at 2024-03-09 16:16:59
type Products struct {
    Id uint32 `gorm:"column:f_id" json:"id" db:"f_id" form:"id"` // 商品id
    Name string `gorm:"column:f_name" json:"name" db:"f_name" form:"name"` // 商品名称
    Description string `gorm:"column:f_description" json:"description" db:"f_description" form:"description"`
    Price float64 `gorm:"column:f_price" json:"price" db:"f_price" form:"price"`
    Weight float32 `gorm:"column:f_weight" json:"weight" db:"f_weight" form:"weight"` // 商品重量（kg）
    Quantity uint32 `gorm:"column:f_quantity" json:"quantity" db:"f_quantity" form:"quantity"` // 商品库存数量
    IsActive int32 `gorm:"column:f_is_active" json:"is_active" db:"f_is_active" form:"is_active"` // 商品是否激活（0 - 未激活，1 - 激活）
    Rating float64 `gorm:"column:f_rating" json:"rating" db:"f_rating" form:"rating"` // 商品评分
    CreatedAt time.Time `gorm:"column:f_created_at" json:"created_at" db:"f_created_at" form:"created_at"` // 商品创建时间
    UpdatedAt time.Time `gorm:"column:f_updated_at" json:"updated_at" db:"f_updated_at" form:"updated_at"` // 商品更新时间
}

func (p *Products) TableName() string {
    return "t_products"
}

// Goods 商品表
// Generated by sql2struct at 2024-03-09 16:16:59
type Goods struct {
    Id uint32 `gorm:"column:id" json:"id" db:"id" form:"id" sql:"id"` // 商品id
    Name string `gorm:"column:name" json:"name" db:"name" form:"name" sql:"name"` // 商品名称
    Description string `gorm:"column:description" json:"description" db:"description" form:"description" sql:"description"` // 商品描述
    Price float64 `gorm:"column:price" json:"price" db:"price" form:"price" sql:"price"` // 商品价格
    Weight float32 `gorm:"column:weight" json:"weight" db:"weight" form:"weight" sql:"weight"` // 商品重量（kg）
    Quantity uint32 `gorm:"column:quantity" json:"quantity" db:"quantity" form:"quantity" sql:"quantity"` // 商品库存数量
    IsActive int32 `gorm:"column:is_active" json:"is_active" db:"is_active" form:"is_active" sql:"is_active"` // 商品是否激活（0 - 未激活，1 - 激活）
    Rating float64 `gorm:"column:rating" json:"rating" db:"rating" form:"rating" sql:"rating"` // 商品评分
    CreatedAt time.Time `gorm:"column:created_at" json:"created_at" db:"created_at" form:"created_at" sql:"created_at"` // 商品创建时间
    UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at" db:"updated_at" form:"updated_at" sql:"updated_at"` // 商品更新时间
}

func (g *Goods) TableName() string {
    return "goods"
}

// Commodity Generated by sql2struct at 2024-03-09 16:16:59
type Commodity struct {
    Id uint32 `gorm:"column:f_id" json:"f_id" db:"f_id" form:"id" sql:"f_id" xorm:"id"` // 商品id
    Name string `gorm:"column:f_name" json:"f_name" db:"f_name" form:"name" sql:"f_name" xorm:"name"` // 商品名称
    Description string `gorm:"column:f_description" json:"f_description" db:"f_description" form:"description" sql:"f_description" xorm:"description"`
    Price float64 `gorm:"column:f_price" json:"f_price" db:"f_price" form:"price" sql:"f_price" xorm:"price"`
    Weight float32 `gorm:"column:f_weight" json:"f_weight" db:"f_weight" form:"weight" sql:"f_weight" xorm:"weight"` // 商品重量（kg）
    Quantity uint32 `gorm:"column:f_quantity" json:"f_quantity" db:"f_quantity" form:"quantity" sql:"f_quantity" xorm:"quantity"` // 商品库存数量
    IsActive int32 `gorm:"column:f_is_active" json:"f_is_active" db:"f_is_active" form:"is_active" sql:"f_is_active" xorm:"is_active"` // 商品是否激活（0 - 未激活，1 - 激活）
    Rating float64 `gorm:"column:f_rating" json:"f_rating" db:"f_rating" form:"rating" sql:"f_rating" xorm:"rating"` // 商品评分
    CreatedAt time.Time `gorm:"column:f_created_at" json:"f_created_at" db:"f_created_at" form:"created_at" sql:"f_created_at" xorm:"created_at"` // 商品创建时间
    UpdatedAt time.Time `gorm:"column:f_updated_at" json:"f_updated_at" db:"f_updated_at" form:"updated_at" sql:"f_updated_at" xorm:"updated_at"` // 商品更新时间
}
