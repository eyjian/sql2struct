### sql2struct

一个根据"CREATE TABLE"建表语句生成对应的Go语言结构体的工具，暂只支持 MySQL 表。

### 使用示例

```shell
yijian@MacBook-Pro-16 sql2struct % cat example-01.sql 
DROP TABLE t_products;
CREATE TABLE t_products (
                          f_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '商品ID',
                          f_name VARCHAR(255) NOT NULL COMMENT '商品名称',
                          f_description TEXT,
                          f_price DECIMAL(10, 2) NOT NULL,
                          f_weight FLOAT NOT NULL COMMENT '商品重量（kg）',
                          f_quantity SMALLINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '商品库存数量',
                          f_is_active TINYINT(1) NOT NULL DEFAULT 1 COMMENT '商品是否激活（0 - 未激活，1 - 激活）',
                          f_rating DOUBLE COMMENT '商品评分',
                          f_created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '商品创建时间',
                          f_updated_at DATETIME ON UPDATE CURRENT_TIMESTAMP COMMENT '商品更新时间',
                          UNIQUE INDEX idx_name_at (f_name),
                          INDEX idx_created_at (f_created_at),
                          KEY idx_updated_at (f_updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';
yijian@MacBook-Pro-16 sql2struct % 
yijian@MacBook-Pro-16 sql2struct % ./sql2struct -sf=./example-01.sql --package="test"
// Package test
// Generated by sql2struct at 2024-03-03 14:36:28
package test

// Products Generated by sql2struct at 2024-03-03 14:36:28
type Products struct {
    Id uint32 `gorm:"column:f_id" json:"Id" db:"f_id" form:"Id"` // 商品id
    Name string `gorm:"column:f_name" json:"Name" db:"f_name" form:"Name"` // 商品名称
    Description string `gorm:"column:f_description" json:"Description" db:"f_description" form:"Description"`
    Price float64 `gorm:"column:f_price" json:"Price" db:"f_price" form:"Price"`
    Weight float32 `gorm:"column:f_weight" json:"Weight" db:"f_weight" form:"Weight"` // 商品重量（kg）
    Quantity uint32 `gorm:"column:f_quantity" json:"Quantity" db:"f_quantity" form:"Quantity"` // 商品库存数量
    IsActive int32 `gorm:"column:f_is_active" json:"IsActive" db:"f_is_active" form:"IsActive"` // 商品是否激活（0 - 未激活，1 - 激活）
    Rating float64 `gorm:"column:f_rating" json:"Rating" db:"f_rating" form:"Rating"` // 商品评分
    CreatedAt time.Time `gorm:"column:f_created_at" json:"CreatedAt" db:"f_created_at" form:"CreatedAt"` // 商品创建时间
    UpdatedAt time.Time `gorm:"column:f_updated_at" json:"UpdatedAt" db:"f_updated_at" form:"UpdatedAt"` // 商品更新时间
}
```

### 使用约束

* sql 中的分割须为空格，而不能是 TAB
* 命令行参数"--sf"指定的 sql 文件只能包含一个"create table"建表语句，不指定同一个 sql 文件含多个建表语句，但大写或者小写不影响
* 生成的时为排版的，需要自行格式化
* 生成的 Go 结构体中，字段名、类型、注释等信息都是从 sql 语句中解析出来的，如果 sql 语句中的字段名、类型、注释等信息不规范，生成的 Go 结构体也会不规范

### 使用提示

* 运行成功，程序退出码为 0，否则为非 0，Shell 中可通过"$?”的值来区分
* 结果直接屏幕输出，可重定向到文件中
* 通过重定向，可实现多个 SQL 文件对应一个 Go 代码文件
