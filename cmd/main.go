package main

import (
	"flag"
	"fmt"
	"github.com/abadojack/whatlanggo"
	"github.com/tidwall/gjson"
	"github.com/xiaoxuan6/deeplx"
	db "github.com/xiaoxuan6/go-package-db"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	url       string
	desc      string
	languages string
)

func main() {
	flag.StringVar(&url, "url", "", "第三方包地址")
	flag.StringVar(&desc, "desc", "", "第三方包描述")
	flag.StringVar(&languages, "language", "", "第三方包语言")
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

		response, err := http.DefaultClient.Get(fmt.Sprintf("https://ungh.xiaoxuan6.me/repos/%s/%s", split[1], split[2]))
		if err == nil {
			defer response.Body.Close()

			b, _ := ioutil.ReadAll(response.Body)
			desc = gjson.GetBytes(b, "repo.description").String()

			lang := whatlanggo.DetectLang(desc)
			sourceLang := strings.ToLower(lang.Iso6391())
			if sourceLang != "zh" {
				res := deeplx.Translate(desc, sourceLang, "zh")
				println("翻译结果：", res)
				if res.Code == 200 {
					desc = res.Data
				} else {
					response, err = http.DefaultClient.Post("https://xiaoxuan6s-yd-translate.hf.space/api/translate", "application/json", strings.NewReader(`{"text":"`+desc+`"}`))
					if err != nil {
						defer response.Body.Close()
						b, _ = ioutil.ReadAll(response.Body)
						println("有道翻译结果：", string(b))
						if 200 == gjson.GetBytes(b, "code").Int() {
							desc = gjson.GetBytes(b, "data").String()
						}
					}
				}
			}
		}
	}

	if err := db.DB.Where(db.Collect{Url: url}).Attrs(db.Collect{Name: desc, Language: languages}).FirstOrCreate(&db.Collect{}).Error; err != nil {
		fmt.Println(fmt.Sprintf("插入数据失败：%s", err.Error()))
	}

	fmt.Println("插入成功！")
}
