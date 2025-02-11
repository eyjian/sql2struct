module github.com/eyjian/sql2struct

// 使用本地的 s2s 开发测试时打开，
// 打开后会导致“go install”执行报错“The go.mod file for the module providing named packages contains one or more replace directives. ”
//replace github.com/eyjian/sql2struct/s2s => ./s2s

go 1.21.5

toolchain go1.23.2

// 使用本地的 s2s 开发测试时打开
//require github.com/eyjian/sql2struct/s2s v0.0.0-00010101000000-000000000000 // indirect

// 使用本地的 s2s 开发测试时注释掉
require github.com/eyjian/sql2struct/s2s v0.0.5
