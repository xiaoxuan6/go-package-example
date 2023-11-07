package main

import (
	"flag"
	"fmt"
	db "github.com/xiaoxuan6/go-package-db"
	"os"
)

var (
	url  string
	name string
)

func main() {
	flag.StringVar(&url, "url", "", "第三方包地址")
	flag.StringVar(&name, "name", "", "第三方包名")
	flag.Parse()

	if url == "" || name == "" {
		fmt.Println("参数错误")
		return
	}

	db.Init(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	defer db.Close()
	db.AutoMigrate()

	_, err := db.FindByUrl(url)
	if err != nil {
		if err = db.Insert(db.Collect{
			Name: name,
			Url:  url,
		}); err != nil {
			fmt.Println(fmt.Sprintf("插入数据失败：%s", err.Error()))
		}

		fmt.Println("插入数据成功")
		return
	}

	fmt.Println("数据已存在")
}
