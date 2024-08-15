### sql2struct

一个根据"CREATE TABLE"建表语句生成对应的Go语言结构体的工具，暂只支持 MySQL 表。支持自定义的 tags，通过参数"-tags"指定，多个自定义的 tags 间使用逗号分开。如果自定义的 tags 值以横杠"-"打头，则表示使用去掉字段名前缀作为值，否则使用字段名作为值。

### 开发目的

在 github 中找到一些 sql2struct，但要么是 chrome 插件，要么是在线工具，要么是需要连接 MySQL，不是很方便。本 sql2struct 根据 SQL 文件中的建表语句来生成 Go 的 struct，可集成到 Makefile 等中，方便使用。

### 安装方法

```shell
go install github.com/eyjian/sql2struct@latest
```

执行成功后，在 $GOPATH/bin 目录下可找到 sql2struct：

```shell
# file `go env GOPATH`/bin/sql2struct
/root/go/bin/sql2struct: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, not stripped
```

### 使用示例

```shell
sql2struct % cat example-01.sql
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
sql2struct % 
sql2struct % ./sql2struct -sf=./example-01.sql --package="test"
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

* 建议将 sql2struct 放到 PATH 指定的目录，比如 /usr/bin/ 或 $GOPATH/bin/ 目录下，以便在任何地方都可以直接使用
* 运行成功，程序退出码为 0，否则为非 0，Shell 中可通过"$?”的值来区分
* 结果直接屏幕输出，可重定向到文件中
* 通过重定向，可实现多个 SQL 文件对应一个 Go 代码文件
* 默认不输出取表名函数，可通过参数"--with-tablename-func"开启
* 默认 json 和 form 两种 tag 会去掉字段名的前缀部分，但可通过命令行参数"-json-with-prefix”和"-form-with-prefix"分别控制

### Makefile 中应用示例

```shell
all: sql sql-01 sql-02 sql-03

.PHONY: sql

sql:
	rm -f example.go

sql-01: example-01.sql
	sql2struct -sf=$< -package="main" -with-tablename-func=true >> example.go

sql-02: example-02.sql
	echo "" >> example.go&&sql2struct -sf=$< -with-tablename-func=true -tags="sql" >> example.go

sql-03: example-03.sql
	echo "" >> example.go&&sql2struct -sf=$< -json-with-prefix=true -tags="sql,-xorm" >> example.go
```

### 使用帮助

```shell
% sql2struct -h
A tool to convert table creation SQL into Go struct, TAB is not supported in SQL file, only spaces are supported.
Usage of sql2struct:
  -db
        With db tag. (default true)
  -form
        With form tag. (default true)
  -form-with-prefix
        Whether from tag retains prefix of field name.
  -fp string
        Prefix of field name. (default "f_")
  -gorm
        With gorm tag. (default true)
  -h    Display a help message and exit.
  -json
        With json tag. (default true)
  -json-with-prefix
        Whether json tag retains prefix of field name.
  -package string
        Package name.
  -sf string
        SQL file containing "CREATE TABLE".
  -tags string
        Custom tags, separate multiple tags with commas.
  -tp string
        Prefix of table name. (default "t_")
  -v    Display version info and exit.
  -with-tablename-func
        Generate a function that takes the table name.
```
