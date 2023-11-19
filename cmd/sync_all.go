package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/joho/godotenv"
	db "github.com/xiaoxuan6/go-package-db"
	"io"
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

	var data []db.Collect
	br := bufio.NewReader(bytes.NewReader(b))
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		re := regexpContent(string(a))
		if len(re) < 1 {
			continue
		}

		if strings.HasPrefix(re[2], "github.com") == false {
			continue
		}

		data = append(data, db.Collect{
			Name:     re[3],
			Url:      re[2],
			Language: "Go", // 默认 Go 语言
		})
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
