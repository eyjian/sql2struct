// Package s2s
// Wrote by yijian on 2024/03/03
package s2s

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
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
	IsJsonField   bool   // 是否为json字段
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

	Version     string
	PackageName string
	TablePrefix string
	FieldPrefix string

	WithGorm bool // 是否生成 Gorm tag
	WithJson bool // 是否生成 Json tag
	WithDb   bool // 是否生成 Db tag
	WithForm bool // 是否生成 Form tag

	WithTableNameFunc bool // 是否生成 TableName 方法
	JsonWithPrefix    bool // 生成的 Json tag 时是否添加前缀
	FormWithPrefix    bool // 生成的 Form tag 时是否添加前缀

	CustomTags string // 定制的 tags

	PointerType     bool // 是否映射为指针类型（含 time.Time 字段）
	TimePointerType bool // 是否 time.Time 字段映射为指针类型
}

func NewSqlTable() *SqlTable {
	return &SqlTable{}
}

func (s *SqlTable) Sql2Struct(scanner *bufio.Scanner) (string, error) {
	// 标记是否在注释块中
	inCommentBlock := false

	// 逐行读取文件
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
		err := s.parseLine(line)
		if err != nil {
			return "", err
		}
	}

	if len(s.TableName) > 0 && len(s.Fields) > 0 {
		return s.toGoStruct(), nil
	}

	return "", fmt.Errorf("no table name or fields")
}

func (s *SqlTable) parseLine(line string) error {
	// 替换多个连续的空格为一个空格
	re := regexp.MustCompile(`\s+`)
	line = re.ReplaceAllString(line, " ")

	// 检查是否包含 "create table"
	if strings.Contains(line, "create table") {
		return s.parseCreateTable(line)
	} else {
		return s.parseNonCreateTable(line)
	}
}

func (s *SqlTable) parseCreateTable(line string) error {
	// 使用字符串分割函数
	parts := strings.Split(line, "--")
	if len(parts) > 1 {
		s.TableComment = strings.TrimSpace(parts[1])
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
		s.RawTableName = tableName

		// 去掉字符串指定的前缀部分
		if len(s.TablePrefix) > 0 {
			tableName = strings.TrimPrefix(tableName, s.TablePrefix)
		}

		s.TableName = s.toStructName(tableName)
		//fmt.Printf("Table name: %s\n", sqlTable.TableName)
		return nil
	} else {
		return fmt.Errorf("no table name found: %s.\n", line)
	}
}

func (s *SqlTable) parseNonCreateTable(line string) error {
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
	sqlTableField.AutoIncrement = false
	if strings.Contains(line, "auto_increment") {
		sqlTableField.AutoIncrement = true
	}

	// 取得字段名和字段类型
	// 使用正则表达式匹配字符串
	re = regexp.MustCompile(`(\w+)\s+(\w+)`)
	matches = re.FindStringSubmatch(line)
	if len(matches) > 2 {
		sqlTableField.RawFieldName = matches[1]
		sqlTableField.FieldType = matches[2]

		sqlTableField.IsJsonField = false
		if sqlTableField.FieldType == "json" {
			sqlTableField.IsJsonField = true
		}

		// 删除字段前缀
		if len(s.FieldPrefix) > 0 {
			sqlTableField.FieldName = strings.TrimPrefix(sqlTableField.RawFieldName, s.FieldPrefix)
		} else {
			sqlTableField.FieldName = sqlTableField.RawFieldName
		}
		sqlTableField.FieldName = s.toStructName(sqlTableField.FieldName)

		// 是否存在 unsigned
		sqlTableField.IsUnsigned = strings.Contains(line, " unsigned ")

		//sqlTableField.Print()
		sqlTableField.IsPrimaryKey = isPrimaryKey
		s.Fields = append(s.Fields, &sqlTableField)
	}

	// 解析主键
	pattern = `primary\s+key\s*\((.*?)\)`
	re = regexp.MustCompile(pattern)
	matches = re.FindStringSubmatch(line)
	if len(matches) > 0 {
		fieldName := matches[1]
		s.tagPrimaryKey(fieldName)
	}

	return nil
}

// toStructName 将"err_code"转为"ErrCode"
func (s *SqlTable) toStructName(name string) string {
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

func (s *SqlTable) toGoStruct() string {
	var sb strings.Builder

	if len(s.PackageName) > 0 {
		sb.WriteString(fmt.Sprintf("// Package %s\n", s.PackageName))
		if s.Version != "" {
			sb.WriteString(fmt.Sprintf("// Generated by sql2struct-%s at %s\n", s.Version, time.Now().Format("2006-01-02 15:04:05")))
		} else {
			sb.WriteString(fmt.Sprintf("// Generated by sql2struct at %s\n", time.Now().Format("2006-01-02 15:04:05")))
		}
		sb.WriteString(fmt.Sprintf("package %s\n\n", s.PackageName))

		if s.haveTimeMember() {
			sb.WriteString("import \"time\"\n\n")
		}
	}

	// 存在表注释
	if len(s.TableComment) > 0 {
		sb.WriteString(fmt.Sprintf("// %s %s\n", s.TableName, s.TableComment))
		sb.WriteString(fmt.Sprintf("// Generated by sql2struct at %s\n", time.Now().Format("2006-01-02 15:04:05")))
	} else {
		sb.WriteString(fmt.Sprintf("// %s Generated by sql2struct at %s\n", s.TableName, time.Now().Format("2006-01-02 15:04:05")))
	}

	// 输出表名
	sb.WriteString(fmt.Sprintf("type %s struct {\n", s.TableName))

	// 输出字段
	for _, field := range s.Fields {
		tag := s.getTag(field)
		goType := mysqlType2GoType(field)
		if s.PointerType {
			goType = "*" + goType
		} else if s.TimePointerType && goType == "time.Time" {
			goType = "*" + goType
		}

		// 过滤掉
		if goType == "any" {
			continue
		}

		// 前导 4 个空格
		if len(field.FieldComment) == 0 {
			sb.WriteString(fmt.Sprintf("    %s %s%s\n", field.FieldName, goType, tag))
		} else {
			sb.WriteString(fmt.Sprintf("    %s %s%s // %s\n", field.FieldName, goType, tag, field.FieldComment))
		}
	}

	// 输出结束
	sb.WriteString("}\n")

	// 输出表名函数
	if s.WithTableNameFunc {
		firstChar := []rune(s.TableName)[0]
		firstChar = unicode.ToLower(firstChar)
		sb.WriteString(fmt.Sprintf("\nfunc (%c *%s) TableName() string {\n", firstChar, s.TableName))
		sb.WriteString(fmt.Sprintf("    return \"%s\"\n", s.RawTableName))
		// 输出结束
		sb.WriteString("}\n")
	}

	return sb.String()
}

// tagPrimaryKey 标记名为 fieldName 的字段为主键
func (s *SqlTable) tagPrimaryKey(fieldName string) {
	for _, field := range s.Fields {
		if field.RawFieldName == fieldName {
			field.IsPrimaryKey = true
			break
		}
	}
}

func (s *SqlTable) haveTimeMember() bool {
	for _, field := range s.Fields {
		goType := mysqlType2GoType(field)
		if goType == "time.Time" {
			return true
		}
	}

	return false
}

func (s *SqlTable) getTag(field *SqlTableField) string {
	var fieldName string
	var tag string
	rawFieldName := field.RawFieldName

	// 有字段前缀
	if len(s.FieldPrefix) > 0 {
		fieldName = strings.TrimPrefix(rawFieldName, s.FieldPrefix)
	} else {
		fieldName = rawFieldName
	}

	if s.WithGorm {
		tag = s.getGormTag(rawFieldName, field.IsPrimaryKey, field.AutoIncrement)
	}
	if s.WithJson {
		tag = s.getJsonTag(tag, rawFieldName, fieldName)
	}
	if s.WithDb {
		tag = s.getDbTag(tag, rawFieldName)
	}
	if s.WithForm {
		tag = s.getFormTag(tag, rawFieldName, fieldName)
	}
	if len(s.CustomTags) > 0 {
		tag = s.getCustomTags(tag, rawFieldName, fieldName)
	}
	if len(tag) > 0 {
		tag = " `" + tag + "`"
	}
	return tag
}

func (s *SqlTable) getGormTag(rawFieldName string, isPrimaryKey, autoIncrement bool) string {
	tags := []string{fmt.Sprintf("column:%s", rawFieldName)}

	if isPrimaryKey {
		tags = append(tags, "primaryKey")
	}
	if autoIncrement {
		tags = append(tags, "autoIncrement")
	}

	return fmt.Sprintf("gorm:\"%s\"", strings.Join(tags, ";"))
}

func (s *SqlTable) getJsonTag(tag, rawFieldName, fieldName string) string {
	jsonFieldName := fieldName

	if s.JsonWithPrefix {
		jsonFieldName = rawFieldName
	}
	if len(tag) == 0 {
		return fmt.Sprintf("json:\"%s\"", jsonFieldName)
	}

	return fmt.Sprintf("%s json:\"%s\"", tag, jsonFieldName)
}

func (s *SqlTable) getDbTag(tag, rawFieldName string) string {
	if len(tag) == 0 {
		return fmt.Sprintf("db:\"%s\"", rawFieldName)
	}

	return fmt.Sprintf("%s db:\"%s\"", tag, rawFieldName)
}

func (s *SqlTable) getFormTag(tag, rawFieldName, fieldName string) string {
	formFieldName := fieldName

	if s.FormWithPrefix {
		formFieldName = rawFieldName
	}
	if len(tag) == 0 {
		return fmt.Sprintf("form:\"%s\"", formFieldName)
	}

	return fmt.Sprintf("%s form:\"%s\"", tag, formFieldName)
}

func (s *SqlTable) getCustomTags(tag, rawFieldName, fieldName string) string {
	customTags := strings.Split(s.CustomTags, ",")
	newTags := make([]string, 0, len(customTags)+1)

	if len(tag) > 0 {
		newTags = append(newTags, tag)
	}
	for _, customTag := range customTags {
		useFieldName := strings.HasPrefix(customTag, "-")
		if useFieldName {
			customTag = strings.TrimPrefix(customTag, "-")
		}

		tagName := rawFieldName
		if useFieldName {
			tagName = fieldName
		}

		newTags = append(newTags, fmt.Sprintf("%s:\"%s\"", customTag, tagName))
	}

	return strings.Join(newTags, " ")
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
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext", "json":
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