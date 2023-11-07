package main

import (
	"fmt"
	"github.com/joho/godotenv"
	db "github.com/xiaoxuan6/go-package-db"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func main() {
	_ = godotenv.Load()

	b, err := ioutil.ReadFile("README.md")
	if err != nil {
		fmt.Println(fmt.Sprintf("读取文件失败：%s", err.Error()))
		return
	}

	newContent := string(b)
	replacements := []string{"# Go 开源第三方包收集和使用示例", "|分支名|包名|描述|", "|:---|:---|:---|"}
	for _, replaceOld := range replacements {
		newContent = strings.ReplaceAll(newContent, replaceOld, ``)
	}
	newContent = strings.Trim(newContent, "\n")
	contents := strings.Split(newContent, "\n")

	var data []db.Collect
	for _, val := range contents {
		regexpStr := regexpContent(val)
		if regexpStr != nil {
			data = append(data, db.Collect{
				Name: regexpStr[3],
				Url:  regexpStr[2],
			})
		}
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

	_ = db.DeleteAll()
	if err = db.Insert(data...); err != nil {
		fmt.Println("数据插入失败：" + err.Error())
		return
	}

	fmt.Println("同步成功")
}

func regexpContent(val string) []string {
	re := regexp.MustCompile(`\|(.*?)\|(.*?)\|(.*?)\|`)
	matchers := re.FindStringSubmatch(val)
	if len(matchers) < 1 {
		return nil
	}
	return matchers
}
