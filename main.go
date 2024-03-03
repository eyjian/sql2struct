// wrote by yijian on 2024/03/03
package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "regexp"
    "strings"
    "unicode"
)

var (
    help        = flag.Bool("h", false, "Display a help message and exit")
    sql         = flag.String("sql", "", "SQL file containing \"CREATE TABLE\"")
    tablePrefix = flag.String("tp", "t_", "Prefix of table name")
    fieldPrefix = flag.String("fp", "f_", "Prefix of field name")
)

// SqlTableField 表字段
type SqlTableField struct {
    FieldName    string // 字段名
    FieldType    string // 字段类型
    FieldComment string // 字段的注释
    IsUnsigned   bool   // 是否为无符号类型
}

func (s SqlTableField) Print() {
    fmt.Printf("FieldName:%s, FieldType:%s, FieldComment:%s, IsUnsigned:%v\n", s.FieldName, s.FieldType, s.FieldComment, s.IsUnsigned)
}

// SqlTable 表
type SqlTable struct {
    TableComment string          // 表的注释
    TableName    string          // 表名
    Fields       []SqlTableField // 字段列表
}

var sqlTable SqlTable

func main() {
    flag.Parse()
    if *help {
        usage()
        os.Exit(1)
    }
    if len(*sql) == 0 {
        fmt.Printf("Parameter --sql is not set\n")
        usage()
        os.Exit(2)
    }

    // 打开文件
    file, err := os.Open(*sql)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    // 创建一个扫描器，用于按行读取文件
    scanner := bufio.NewScanner(file)

    // 逐行读取文件
    for scanner.Scan() {
        // 读取一行文本
        line := scanner.Text()

        // 删除前后空格
        line = strings.TrimSpace(line)

        // 删除尾部的逗号
        line = strings.TrimSuffix(line, ",")

        // 删除指定字符
        line = strings.ReplaceAll(line, "`", "")

        // 全部转为小写，简化后续处理
        line = strings.ToLower(line)

        if strings.HasPrefix(line, "key") ||
            strings.HasPrefix(line, "index") ||
            strings.HasPrefix(line, "primary") ||
            strings.HasPrefix(line, "unique") ||
            strings.HasPrefix(line, "(") ||
            strings.HasPrefix(line, ")") ||
            strings.HasPrefix(line, "--") ||
            strings.HasPrefix(line, "drop") {
            continue
        }

        // 解析文本行
        if !parseLine(line) {
            break
        }
    }
}

func usage() {
    flag.Usage()
}

func parseLine(line string) bool {
    // 替换多个连续的空格为一个空格
    re := regexp.MustCompile(`\s+`)
    line = re.ReplaceAllString(line, " ")

    // 检查是否包含 "create table"
    if strings.Contains(line, "create table") {
        return parseCreateTable(line)
    } else {
        return parseNonCreateTable(line)
    }
    return false
}

func parseCreateTable(line string) bool {
    // 去除可能存在的"-- 创建"部分
    newLine := strings.Split(line, "--")[0]

    // 创建正则表达式，用于匹配 "CREATE TABLE" 后的表名
    re := regexp.MustCompile(`create\s+table\s+(\w+)`)

    // 在输入字符串中查找匹配项
    matches := re.FindStringSubmatch(newLine)

    // 如果找到匹配项，则输出表名
    if len(matches) > 1 {
        // 得到可能含前缀的表名
        tableName := matches[1]

        // 去掉字符串指定的前缀部分
        if len(*tablePrefix) > 0 {
            tableName = strings.TrimPrefix(tableName, *tablePrefix)
        }

        sqlTable.TableName = toStructName(tableName)
        fmt.Printf("Table name: %s\n", sqlTable.TableName)
        return true
    } else {
        fmt.Printf("No table name found: %s\n", line)
        return false
    }
}

func parseNonCreateTable(line string) bool {
    var sqlTableField SqlTableField
    //fmt.Println(line)

    // 取得字段的注释
    re := regexp.MustCompile(`comment\s+'(.+)'`)
    matches := re.FindStringSubmatch(line)
    if len(matches) > 1 {
        sqlTableField.FieldComment = matches[1]
    }

    // 取得字段名和字段类型
    // 使用正则表达式匹配字符串
    re = regexp.MustCompile(`(\w+)\s+(\w+)`)
    matches = re.FindStringSubmatch(line)
    if len(matches) > 2 {
        sqlTableField.FieldName = toStructName(matches[1])
        sqlTableField.FieldType = matches[2]
    }

    sqlTableField.Print()

    return true
}

// toStructName 将"err_code"转为"ErrCode"
func toStructName(name string) string {
    var result []rune
    for i, v := range name {
        if i == 0 || name[i-1] == '_' {
            result = append(result, unicode.ToUpper(v))
        } else {
            result = append(result, v)
        }
    }

    // 将结果转换为字符串并去掉所有的"_"
    return strings.ReplaceAll(string(result), "_", "")
}
