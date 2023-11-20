package main

import (
	"flag"
	"fmt"
	"github.com/OwO-Network/gdeeplx"
	"github.com/abadojack/whatlanggo"
	"github.com/antchfx/htmlquery"
	db "github.com/xiaoxuan6/go-package-db"
	"os"
	"strings"
)

var (
	url      string
	name     string
	language string
)

func main() {
	flag.StringVar(&url, "url", "", "第三方包地址")
	flag.StringVar(&name, "name", "", "第三方包名")
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

	name = nameDo(url, name)
	url = strings.ReplaceAll(url, "https://", "")
	if err := db.DB.Where(db.Collect{Url: url}).Attrs(db.Collect{Name: name, Language: language}).FirstOrCreate(&db.Collect{}).Error; err != nil {
		fmt.Println(fmt.Sprintf("插入数据失败：%s", err.Error()))
	}

	fmt.Println("插入成功！")
}

func nameDo(url, name string) string {
	if name == "" {
		if strings.HasPrefix(url, "https://") == false {
			url = "https://" + url
		}
		doc, err := htmlquery.LoadURL(url)
		if err == nil {
			a := htmlquery.FindOne(doc, "//*[@id=\"responsive-meta-container\"]/div/p")
			name = strings.TrimSpace(htmlquery.InnerText(a))
		}
	}

	info := whatlanggo.Detect(name)
	lang := info.Lang.String()
	if lang != "" && lang != "Mandarin" {
		result, err := gdeeplx.Translate(name, "", "zh", 0)
		if err == nil {
			res := result.(map[string]interface{})
			name = strings.TrimSpace(res["data"].(string))
		}
	}

	return name
}
