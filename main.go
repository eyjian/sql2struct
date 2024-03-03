// wrote by yijian on 2024/03/03
package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "regexp"
    "strings"
)

var (
    help        = flag.Bool("h", false, "Display a help message and exit")
    sql         = flag.String("sql", "", "SQL file containing \"CREATE TABLE\"")
    tablePrefix = flag.String("tp", "t_", "Prefix of table name")
    fieldPrefix = flag.String("fp", "f_", "Prefix of field name")
)

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

        // 删除指定字符
        line = strings.ReplaceAll(line, "`", "")

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
    // 全部转为小写，简化后续处理
    lowerLine := strings.ToLower(line)

    // 替换多个连续的空格为一个空格
    re := regexp.MustCompile(`\s+`)
    lowerLine = re.ReplaceAllString(lowerLine, " ")

    // 检查是否包含 "create table"
    if strings.Contains(lowerLine, "create table") {
        return parseCreateTable(lowerLine)
    } else {
        return parseNonCreateTable(lowerLine)
    }
    return false
}

func parseCreateTable(line string) bool {
    // 去除可能存在的"-- 创建"部分
    line = strings.Split(line, "--")[0]

    // 创建正则表达式，用于匹配 "CREATE TABLE" 后的表名
    re := regexp.MustCompile(`create\s+table\s+(\w+)`)

    // 在输入字符串中查找匹配项
    matches := re.FindStringSubmatch(line)

    // 如果找到匹配项，则输出表名
    if len(matches) > 1 {
        fmt.Printf("Table name: %s\n", matches[1])
        return true
    } else {
        fmt.Printf("No table name found: %s\n", line)
        return false
    }
}

func parseNonCreateTable(line string) bool {
    return true
}
