package main

import (
	"flag"
	"fmt"
	db "github.com/xiaoxuan6/go-package-db"
	"os"
	"strings"
)

var (
	url  string
	desc string
)

func main() {
	flag.StringVar(&url, "url", "", "第三方包地址")
	flag.StringVar(&desc, "desc", "", "第三方包描述")
	flag.StringVar(&language, "language", "", "第三方包语言")
	flag.Parse()

	if url == "" {
		fmt.Println("参数错误")
		return
	}

	// 本地测试调用
	//_ = godotenv.Load()
	db.Init(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	defer db.Close()
	db.AutoMigrate()

	url = strings.ReplaceAll(url, "https://", "")
	if len(desc) < 1 {
		split := strings.Split(url, "/")
		desc = fetchDescription(split[1], split[2], "")
	}

	if err := db.DB.Where(db.Collect{Url: url}).Attrs(db.Collect{Name: desc, Language: language}).FirstOrCreate(&db.Collect{}).Error; err != nil {
		fmt.Println(fmt.Sprintf("插入数据失败：%s", err.Error()))
	}

	fmt.Println("插入成功！")
}
