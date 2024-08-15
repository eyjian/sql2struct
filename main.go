// wrote by yijian on 2024/03/03
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)
import (
	"github.com/eyjian/sql2struct/s2s"
)

const Version string = "0.0.9"

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

	customTags = flag.String("custom-tags", "", "Custom tags, separate multiple tags with commas, example: -tags=\"sql,-xorm,ent,reform\".")
)

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

	sqlTable := s2s.NewSqlTable()
	sqlTable.Version = Version
	sqlTable.PackageName = *packageName
	sqlTable.TablePrefix = *tablePrefix
	sqlTable.FieldPrefix = *fieldPrefix
	sqlTable.WithGorm = *withGorm
	sqlTable.WithJson = *withJson
	sqlTable.WithDb = *withDb
	sqlTable.WithForm = *withForm
	sqlTable.WithTableNameFunc = *withTableNameFunc
	sqlTable.JsonWithPrefix = *jsonWithPrefix
	sqlTable.FormWithPrefix = *formWithPrefix
	sqlTable.CustomTags = *customTags

	scanner := bufio.NewScanner(file) // 创建一个扫描器，用于按行读取文件
	structStr, err := sqlTable.Sql2Struct(scanner)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse %s error: %s.\n", *sqlFile, err.Error())
	} else {
		fmt.Fprintf(os.Stdout, "%s\n", structStr)
	}
}

func showUsage() {
	fmt.Fprintln(os.Stderr, "A tool to convert table creation SQL into Go struct, TAB is not supported in SQL file, only spaces are supported.")
	flag.Usage()
}

func showVersion() {
	fmt.Printf("Version: %s, build at %s\n", Version, time.Now().Format("2006-01-02 15:04:05"))
}
