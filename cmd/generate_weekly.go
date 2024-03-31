package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/OwO-Network/gdeeplx"
	"github.com/abadojack/whatlanggo"
	"github.com/antchfx/htmlquery"
	"github.com/avast/retry-go"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	ghttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/goccy/go-yaml"
	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v48/github"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	url2 "net/url"
	"os"
	time2 "package-example/pkg/time"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	path string
	uri  string

	wg            sync.WaitGroup
	gitRepository *git.Repository

	WeeklyGitUrl = "https://github.com/xiaoxuan6/weekly.git"
)

func main() {
	flag.StringVar(&uri, "uri", "", "")
	flag.Parse()

	_ = godotenv.Load()

	dir, _ := os.Getwd()
	path = filepath.Join(dir, "/weekly")

	owner, repository, baseUrl := fetchRepositoryAndUrl()

	description := fetchDescription(owner, repository, baseUrl)

	if err := cloneRepository(); err != nil {
		return
	}

	// 创建 年 文件
	year := time2.NewTime().Year()
	root := filepath.Join(path, "/docs/", year)
	if _, err := os.Stat(root); err != nil {
		_ = os.MkdirAll(root, os.ModePerm)
	}

	wg.Add(2)
	ch := make(chan bool, 1)
	ch1 := make(chan bool, 1)
	//fileCount := getFileCount(root)
	//filename := fmt.Sprintf("第%d期（%s）.md", fileCount, time2.NewTime().Date())
	filename := fmt.Sprintf("%s.md", time2.NewTime().Date())

	go func() {
		defer wg.Done()

		newFilename := filepath.Join(root, filename)
		if _, err := os.Stat(newFilename); err != nil {

			f, errs := os.Create(newFilename)
			if errs != nil {
				ch <- false
				ch1 <- false
				return
			}

			ch <- true
			_, _ = f.WriteString(fmt.Sprintf(`# %s

---
- 项目地址：[%s](%s)
- 项目说明：%s
---`, strings.ReplaceAll(filename, ".md", ""), repository, baseUrl, description))

			ch1 <- true

		} else {

			ch <- false

			f, _ := os.OpenFile(newFilename, os.O_APPEND, os.ModePerm)
			_, _ = f.WriteString(fmt.Sprintf(`
- 项目地址：[%s](%s)
- 项目说明：%s
---`, repository, baseUrl, description))

			ch1 <- true
		}

		return
	}()

	go func() {
		defer wg.Done()

		item := <-ch
		if item == false {
			return
		}

		mkdocsFile := filepath.Join(path, "mkdocs.yml")
		b, _ := os.ReadFile(mkdocsFile)

		var jsonData yaml.MapSlice
		_ = yaml.Unmarshal(b, &jsonData)

		jsonMap := jsonData.ToMap()
		navSlice := jsonMap["nav"].([]interface{})

		// 判断 xxxx-xx-xx年刊 是否存在, 有则取，无则初始化
		var targetSlice []interface{}
		target := fmt.Sprintf("%s 年刊", year)
		for _, val := range navSlice {
			navs := val.(map[string]interface{})
			if items, ok := navs[target]; ok {
				targetSlice = items.([]interface{})
				break
			}
		}

		// 新增周刊的目录 eg: 2024-3-31.md: ./2024/2020-3-31.md
		yearnNav := make([]interface{}, 0)
		yearnNav = append(yearnNav, map[string]interface{}{
			filename: fmt.Sprintf("./%s/%s", year, filename),
		})

		// 从 旧 的目录中向新的追加
		if len(targetSlice) > 0 {
			for _, navss := range targetSlice {
				yearnNav = append(yearnNav, navss.(map[string]interface{}))
			}
		}

		// 直接标识为不存在
		status := false
		for _, val := range navSlice {
			navs := val.(map[string]interface{})
			for k, _ := range navs {
				if k == target {
					status = true
					navs[target] = yearnNav
				}
			}
		}

		if status == false {
			navSlice = append(navSlice, map[string]interface{}{
				target: yearnNav,
			})
			jsonMap["nav"] = navSlice
		}

		b, _ = yaml.Marshal(jsonMap)
		_ = os.WriteFile(mkdocsFile, b, os.ModePerm)

		return
	}()

	wg.Wait()

	if result := <-ch1; result != true {
		logrus.Info(color.RedString("创建文件 %s 失败", filename))
		return
	}

	if err2 := gitCommit(filename, filepath.Join("docs/", year, filename), "mkdocs.yml"); err2 != nil {
		fmt.Println(color.RedString(err2.Error()))
		return
	}

	if os.Getenv("APP_ENV") == "prod" {
		logrus.Info(color.GreenString("git push"))
		if err := gitPush(); err != nil {
			logrus.Error(color.RedString(err.Error()))
			return
		}

		_ = os.RemoveAll(path)
	}

	logrus.Info(color.GreenString("ok"))
}

func fetchRepositoryAndUrl() (string, string, string) {
	u, _ := url2.Parse(uri)

	result := strings.Split(u.Path, "/")
	if len(result) != 3 {
		return "", u.Path, fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)
	}

	return result[1], result[2], fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)
}

func fetchDescription(owner, repo, uri string) string {
	var description string
	if len(owner) < 1 {
		find := func(doc *html.Node) string {
			a := htmlquery.FindOne(doc, "//*[@id=\"responsive-meta-container\"]/div/p")
			return strings.TrimSpace(htmlquery.InnerText(a))
		}

		doc, _ := htmlquery.LoadURL(uri)
		description = find(doc)

		if len(description) < 1 {

			response, err := get(uri)
			if err != nil {
				return ""
			}

			doc, errs := htmlquery.Parse(strings.NewReader(response))
			if errs == nil {
				description = find(doc)
			}
		}

		if len(description) < 1 {
			return ""
		}
	} else {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		rateLimiter, errs := github_ratelimit.NewRateLimitWaiterClient(tc.Transport)
		if errs != nil {
			panic(errs.Error())
		}

		client := github.NewClient(rateLimiter)
		rep, _, _ := client.Repositories.Get(context.Background(), owner, repo)
		description = *rep.Description
	}

	info := whatlanggo.Detect(description)
	lang := info.Lang.String()
	if lang != "" && lang != "Mandarin" {
		result, err1 := gdeeplx.Translate(description, "", "zh", 0)
		if err1 == nil {
			res := result.(map[string]interface{})
			description = strings.TrimSpace(res["data"].(string))
		}
	}

	return description
}

func cloneRepository() error {
	cloneNum := 0
CLONE:
	if _, err1 := os.Stat(path); err1 == nil {
		if err2 := os.RemoveAll(path); err2 != nil {
			logrus.Error(color.RedString("删除文件失败:", err2.Error()))
			return err2
		}
	}

	rep, err := git.PlainCloneContext(context.Background(), path, false, &git.CloneOptions{
		URL:      WeeklyGitUrl,
		Progress: os.Stdout,
	})

	if err != nil {
		if cloneNum == 2 {
			logrus.Error(color.RedString("clone fail: %s", err.Error()))
			return err
		}
		cloneNum += 1
		logrus.Error(color.RedString("clone fail, retrying..."))
		goto CLONE
	}

	gitRepository = rep
	return nil
}

//func getFileCount(path string) int {
//	fileCount := 1
//	_ = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
//		if info.IsDir() {
//			return nil
//		}
//
//		if strings.Compare(filepath.Ext(info.Name()), ".md") == 0 {
//			fileCount += 1
//		}
//
//		return nil
//	})
//
//	return fileCount
//}

func gitCommit(message string, filename ...string) error {
	w, _ := gitRepository.Worktree()

	for _, file := range filename {
		_, err := w.Add(file)
		if err != nil {
			return fmt.Errorf("git add file %s fail: %s", file, err.Error())
		}
	}

	if _, err := w.Commit(fmt.Sprintf("feat: Add %s", message), &git.CommitOptions{
		Author: &object.Signature{
			Name:  os.Getenv("GITHUB_OWNER"),
			Email: os.Getenv("GITHUB_EMAIL"),
			When:  time2.NewTime().DateTime(),
		},
	}); err != nil {
		return fmt.Errorf("git commit fail: %s", err.Error())
	}

	return nil
}

func gitPush() error {
	if err := gitRepository.Push(&git.PushOptions{
		RemoteName: "origin",
		RemoteURL:  WeeklyGitUrl,
		Auth: &ghttp.BasicAuth{
			Username: os.Getenv("GITHUB_OWNER"),
			Password: os.Getenv("GITHUB_TOKEN"),
		},
		Progress: os.Stdout,
	}); err != nil {
		return fmt.Errorf("git push fail: %s", err.Error())
	}
	return nil
}

func get(uri string) (string, error) {
	var result string
	err := retry.Do(
		func() error {
			response, err := http.Get(uri)
			defer func() {
				_ = response.Body.Close()
			}()

			if err != nil {
				return fmt.Errorf("url get fail: %s", err.Error())
			}

			b, err := io.ReadAll(response.Body)
			if err != nil {
				return fmt.Errorf("read body fail: %s", err.Error())
			}

			result = string(b)
			return nil
		},
		retry.Attempts(3),
		retry.LastErrorOnly(true),
		retry.Delay(time.Second*10),
	)

	return result, err
}