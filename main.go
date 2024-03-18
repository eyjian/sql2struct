// wrote by yijian on 2024/03/03
package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "regexp"
    "strings"
    "time"
    "unicode"
)

const Version string = "0.0.5"

var (
    help    = flag.Bool("h", false, "Display a help message and exit.")
    version = flag.Bool("v", false, "Display version info and exit.")

    sqlFile     = flag.String("sf", "", "SQL file containing \"CREATE TABLE\".")
    packageName = flag.String("package", "", "Package name.")

    tablePrefix = flag.String("tp", "t_", "Prefix of table name.")
    fieldPrefix = flag.String("fp", "f_", "Prefix of field name.")

    withGorm = flag.Bool("gorm", true, "With gorm tag.")
    withJson = flag.Bool("json", true, "With json tag.")
    withDb   = flag.Bool("db", true, "With db tag.")
    withForm = flag.Bool("form", true, "With form tag.")

    withTableNameFunc = flag.Bool("with-tablename-func", false, "Generate a function that takes the table name.")
    jsonWithPrefix    = flag.Bool("json-with-prefix", false, "Whether json tag retains prefix of field name.")
    formWithPrefix    = flag.Bool("form-with-prefix", false, "Whether from tag retains prefix of field name.")

    tags = flag.String("tags", "", "Custom tags, separate multiple tags with commas, example: -tags=\"sql,-xorm\".")
)

// SqlTableField 表字段
type SqlTableField struct {
    RawFieldName  string // 未处理的原始字段名
    FieldName     string // 字段名
    FieldType     string // 字段类型
    FieldComment  string // 字段的注释
    IsUnsigned    bool   // 是否为无符号类型
    IsPrimaryKey  bool   // 是否为主键
    AutoIncrement bool   // 是否为自增字段
}

func (s *SqlTableField) Print() {
    fmt.Printf("FieldName:%s, FieldType:%s, FieldComment:%s, IsUnsigned:%v\n", s.FieldName, s.FieldType, s.FieldComment, s.IsUnsigned)
}

// SqlTable 表
type SqlTable struct {
    RawTableName string           // 原始表名
    TableName    string           // 表名
    TableComment string           // 表的注释
    Fields       []*SqlTableField // 字段列表
}

var sqlTable *SqlTable

func main() {
    flag.Parse()
    if *help {
        showUsage()
        os.Exit(1)
    }
    if *version {
        showVersion()
        os.Exit(1)
    }
    if len(*sqlFile) == 0 {
        fmt.Fprintf(os.Stderr, "Parameter --sf is not set.\n")
        os.Exit(2)
    }

    // 打开文件
    file, err := os.Open(*sqlFile)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Open %s error: %s.\n", *sqlFile, err.Error())
        os.Exit(3)
    }
    defer file.Close()

    // 创建一个扫描器，用于按行读取文件
    scanner := bufio.NewScanner(file)

    // 标记是否在注释块中
    inCommentBlock := false

    // 逐行读取文件
    sqlTable = &SqlTable{}
    for scanner.Scan() {
        // 读取一行文本
        line := scanner.Text()

        // 删除前后空格
        line = strings.TrimSpace(line)

        // 删除尾部的逗号
        line = strings.TrimSuffix(line, ",")

        // 删除前后空格
        line = strings.TrimSpace(line)

        // 过滤掉空行
        if len(line) == 0 {
            continue
        }

        // 删除指定字符
        line = strings.ReplaceAll(line, "`", "")

        // 全部转为小写，简化后续处理
        line = strings.ToLower(line)

        // 结束注释了
        if inCommentBlock {
            if strings.HasPrefix(line, "*/") ||
                strings.HasSuffix(line, "*/") {
                inCommentBlock = false
            }
            continue
        }

        // 进入多行注释
        if strings.HasPrefix(line, "/*") {
            if !strings.HasSuffix(line, "*/") {
                inCommentBlock = true // 如果 /* 同 */ 不在同一行则进入注释块状态中
            }
            continue
        }

        // 在注释中
        if inCommentBlock {
            continue
        }

        // 需要过滤掉的
        if skipLine(line) {
            continue
        }

        // 解析文本行
        if !parseLine(line) {
            os.Exit(4)
        }
    }

    if len(sqlTable.TableName) > 0 && len(sqlTable.Fields) > 0 {
        sqlTable.ToGoStruct()
    }

    os.Exit(0)
}

func showUsage() {
    fmt.Fprintln(os.Stderr, "A tool to convert table creation SQL into Go struct, TAB is not supported in SQL file, only spaces are supported.")
    flag.Usage()
}

func showVersion() {
    fmt.Printf("Version: %s, build at %s\n", Version, time.Now().Format("2006-01-02 15:04:05"))
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
}

func parseCreateTable(line string) bool {
    // 使用字符串分割函数
    parts := strings.Split(line, "--")
    if len(parts) > 1 {
        sqlTable.TableComment = strings.TrimSpace(parts[1])
    }

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
        sqlTable.RawTableName = tableName

        // 去掉字符串指定的前缀部分
        if len(*tablePrefix) > 0 {
            tableName = strings.TrimPrefix(tableName, *tablePrefix)
        }

        sqlTable.TableName = toStructName(tableName)
        //fmt.Printf("Table name: %s\n", sqlTable.TableName)
        return true
    } else {
        fmt.Fprintf(os.Stderr, "No table name found: %s.\n", line)
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

    // 检查是否为主键
    isPrimaryKey := false
    pattern := `\s*primary\s+key\s*`
    matched, err := regexp.MatchString(pattern, line)
    if err == nil && matched {
        isPrimaryKey = true
    }

    // 检查是否为自增主键
    autoIncrement := false
    if strings.Contains(line, "auto_increment") {
        autoIncrement = true
    }

    // 取得字段名和字段类型
    // 使用正则表达式匹配字符串
    re = regexp.MustCompile(`(\w+)\s+(\w+)`)
    matches = re.FindStringSubmatch(line)
    if len(matches) > 2 {
        sqlTableField.RawFieldName = matches[1]
        sqlTableField.FieldType = matches[2]

        // 删除字段前缀
        if len(*fieldPrefix) > 0 {
            sqlTableField.FieldName = strings.TrimPrefix(sqlTableField.RawFieldName, *fieldPrefix)
        } else {
            sqlTableField.FieldName = sqlTableField.RawFieldName
        }
        sqlTableField.FieldName = toStructName(sqlTableField.FieldName)

        // 是否存在 unsigned
        sqlTableField.IsUnsigned = strings.Contains(line, " unsigned ")

        //sqlTableField.Print()
        sqlTableField.IsPrimaryKey = isPrimaryKey
        sqlTableField.AutoIncrement = autoIncrement
        sqlTable.Fields = append(sqlTable.Fields, &sqlTableField)
    }

    // 解析主键
    pattern = `primary\s+key\s*\((.*?)\)`
    re = regexp.MustCompile(pattern)
    matches = re.FindStringSubmatch(line)
    if len(matches) > 0 {
        fieldName := matches[1]
        sqlTable.TagPrimaryKey(fieldName)
    }

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

func (s *SqlTable) ToGoStruct() {
    if len(*packageName) > 0 {
        fmt.Printf("// Package %s\n", *packageName)
        fmt.Printf("// Generated by sql2struct at %s\n", time.Now().Format("2006-01-02 15:04:05"))
        fmt.Printf("package %s\n\n", *packageName)

        if s.HaveTimeMember() {
            fmt.Printf("import \"time\"\n\n")
        }
    }

    // 存在表注释
    if len(s.TableComment) > 0 {
        fmt.Printf("// %s %s\n", s.TableName, s.TableComment)
        fmt.Printf("// Generated by sql2struct at %s\n", time.Now().Format("2006-01-02 15:04:05"))
    } else {
        fmt.Printf("// %s Generated by sql2struct-%s at %s\n", s.TableName, Version, time.Now().Format("2006-01-02 15:04:05"))
    }

    // 输出表名
    fmt.Printf("type %s struct {\n", s.TableName)

    // 输出字段
    for _, field := range s.Fields {
        tag := getTag(field)
        goType := mysqlType2GoType(field)

        // 过滤掉
        if goType == "any" {
            continue
        }

        // 前导 4 个空格
        if len(field.FieldComment) == 0 {
            fmt.Printf("    %s %s%s\n", field.FieldName, goType, tag)
        } else {
            fmt.Printf("    %s %s%s // %s\n", field.FieldName, goType, tag, field.FieldComment)
        }
    }

    // 输出结束
    fmt.Printf("}\n")

    // 输出表名函数
    if *withTableNameFunc {
        firstChar := []rune(s.TableName)[0]
        firstChar = unicode.ToLower(firstChar)
        fmt.Printf("\nfunc (%c *%s) TableName() string {\n", firstChar, s.TableName)
        fmt.Printf("    return \"%s\"\n", s.RawTableName)
        // 输出结束
        fmt.Printf("}\n")
    }
}

// TagPrimaryKey 标记名为 fieldName 的字段为主键
func (s *SqlTable) TagPrimaryKey(fieldName string) {
    for _, field := range s.Fields {
        if field.RawFieldName == fieldName {
            field.IsPrimaryKey = true
            break
        }
    }
}

func (s *SqlTable) HaveTimeMember() bool {
    for _, field := range s.Fields {
        goType := mysqlType2GoType(field)
        if goType == "time.Time" {
            return true
        }
    }

    return false
}

func getTag(field *SqlTableField) string {
    var fieldName string
    var tag string
    rawFieldName := field.RawFieldName

    // 有字段前缀
    if len(*fieldPrefix) > 0 {
        fieldName = strings.TrimPrefix(rawFieldName, *fieldPrefix)
    }

    if *withGorm {
        tag = getGormTag(rawFieldName, field.IsPrimaryKey, field.AutoIncrement)
    }
    if *withJson {
        tag = getJsonTag(tag, rawFieldName, fieldName)
    }
    if *withDb {
        tag = getDbTag(tag, rawFieldName)
    }
    if *withForm {
        tag = getFormTag(tag, rawFieldName, fieldName)
    }
    if len(*tags) > 0 {
        tag = getCustomTags(tag, rawFieldName, fieldName)
    }
    if len(tag) > 0 {
        tag = " `" + tag + "`"
    }
    return tag
}

func mysqlType2GoType(field *SqlTableField) string {
    switch field.FieldType {
    case "tinyint", "smallint", "mediumint", "int", "integer":
        if !field.IsUnsigned {
            return "int32"
        } else {
            return "uint32"
        }
    case "bigint":
        if !field.IsUnsigned {
            return "int64"
        } else {
            return "uint64"
        }
    case "float":
        return "float32"
    case "double", "decimal":
        return "float64"
    case "char", "varchar", "tinytext", "text", "mediumtext", "longtext":
        return "string"
    case "date", "datetime", "timestamp", "time":
        return "time.Time"
    case "tinyblob", "blob", "mediumblob", "longblob", "binary", "varbinary":
        return "[]byte"
    case "bit":
        return "bool"
    case "enum", "set":
        return "string"
    default:
        return "any"
    }
}

func skipLine(line string) bool {
    return strings.HasPrefix(line, "key") ||
        strings.HasPrefix(line, "index") ||
        strings.HasPrefix(line, "unique") ||
        strings.HasPrefix(line, "(") ||
        strings.HasPrefix(line, ")") ||
        strings.HasPrefix(line, "--") ||
        strings.HasPrefix(line, "drop") ||
        strings.HasPrefix(line, "partition") ||
        strings.Contains(line, "engine=") ||
        strings.Contains(line, "auto_increment=") ||
        strings.Contains(line, "charset=") ||
        strings.Contains(line, "partition ")
}

func getGormTag(rawFieldName string, isPrimaryKey, autoIncrement bool) string {
    if !isPrimaryKey {
        if !autoIncrement {
            return fmt.Sprintf("gorm:\"column:%s\"", rawFieldName)
        } else {
            return fmt.Sprintf("gorm:\"column:%s;autoIncrement\"", rawFieldName)
        }
    } else {
        if !autoIncrement {
            return fmt.Sprintf("gorm:\"column:%s;primaryKey\"", rawFieldName)
        } else {
            return fmt.Sprintf("gorm:\"column:%s;primaryKey;autoIncrement\"", rawFieldName)
        }
    }
}

func getJsonTag(tag, rawFieldName, fieldName string) string {
    var newTag string
    if len(tag) == 0 {
        if *jsonWithPrefix {
            newTag = fmt.Sprintf("json:\"%s\"", rawFieldName)
        } else {
            newTag = fmt.Sprintf("json:\"%s\"", fieldName)
        }
    } else {
        if *jsonWithPrefix {
            newTag = fmt.Sprintf("%s json:\"%s\"", tag, rawFieldName)
        } else {
            newTag = fmt.Sprintf("%s json:\"%s\"", tag, fieldName)
        }
    }

    return newTag
}

func getDbTag(tag, rawFieldName string) string {
    var newTag string
    if len(tag) == 0 {
        newTag = fmt.Sprintf("db:\"%s\"", rawFieldName)
    } else {
        newTag = fmt.Sprintf("%s db:\"%s\"", tag, rawFieldName)
    }

    return newTag
}

func getFormTag(tag, rawFieldName, fieldName string) string {
    var newTag string

    if len(tag) == 0 {
        if *formWithPrefix {
            newTag = fmt.Sprintf("form:\"%s\"", rawFieldName)
        } else {
            newTag = fmt.Sprintf("form:\"%s\"", fieldName)
        }
    } else {
        if *formWithPrefix {
            newTag = fmt.Sprintf("%s form:\"%s\"", tag, rawFieldName)
        } else {
            newTag = fmt.Sprintf("%s form:\"%s\"", tag, fieldName)
        }
    }

    return newTag
}

func getCustomTags(tag, rawFieldName, fieldName string) string {
    var newTag string = tag
    customTags := strings.Split(*tags, ",")

    for _, customTag := range customTags {
        // 以"-"打头的表示使用去掉字段头的作为 tag
        useFieldName := strings.HasPrefix(customTag, "-")
        if useFieldName {
            customTag = strings.TrimPrefix(customTag, "-")
        }

        if len(newTag) == 0 {
            if !useFieldName {
                newTag = fmt.Sprintf("%s:\"%s\"", customTag, rawFieldName)
            } else {
                newTag = fmt.Sprintf("%s:\"%s\"", customTag, fieldName)
            }
        } else {
            if !useFieldName {
                newTag = fmt.Sprintf("%s %s:\"%s\"", newTag, customTag, rawFieldName)
            } else {
                newTag = fmt.Sprintf("%s %s:\"%s\"", newTag, customTag, fieldName)
            }
        }
    }

    return newTag
}
