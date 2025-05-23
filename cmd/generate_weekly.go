package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/abadojack/whatlanggo"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	ghttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/goccy/go-yaml"
	"github.com/joho/godotenv"
	"github.com/noelyahan/impexp"
	"github.com/noelyahan/mergi"
	"github.com/sahilm/fuzzy"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/xiaoxuan6/deeplx"
	"image"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"os"
	time2 "package-example/pkg/time"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	path           string
	uri            string
	isDownload     bool
	homepage       string
	img            string
	descriptionVar string
	demoUrl        string
	label          string
	banner         string

	wg            sync.WaitGroup
	gitRepository *git.Repository

	language     = "未知"
	WeeklyGitUrl = "https://github.com/xiaoxuan6/weekly.git"
)

func main() {
	flag.StringVar(&uri, "uri", "", "url 链接地址")
	flag.BoolVar(&isDownload, "is_download", true, "是否下载封面图")
	flag.StringVar(&descriptionVar, "description_var", "", "描述")
	flag.StringVar(&demoUrl, "demo_url", "", "demo 地址")
	flag.StringVar(&label, "label", "pkg", "标签")
	flag.StringVar(&banner, "banner", "", "封面图地址，默认为 github 库中的地址")
	flag.Parse()

	_ = godotenv.Load()

	dir, _ := os.Getwd()
	path = filepath.Join(dir, "/weekly")
	uri = strings.TrimRight(strings.TrimSpace(uri), "/")

	var (
		owner       string
		baseUrl     string
		repository  string
		description string
	)
	u, _ := url2.Parse(uri)
	if strings.Contains(u.Host, "github.com") == true {
		owner, repository, baseUrl = fetchRepositoryAndUrl()

		if len(descriptionVar) < 1 {
			description = fetchDescription(owner, repository, baseUrl)
		} else {
			description = descriptionVar
		}
	} else {
		baseUrl, homepage, repository = uri, uri, uri
		description = descriptionVar

		if len(banner) > 0 {
			homepage = banner
		}
	}

	// --------------------- 去重 START ---------------------
	links := filepath.Join(path, "links.txt")
	if distinct(links) == true {
		logrus.Warning(color.RedString("url [%s] exists", uri))
		return
	}
	// --------------------- 去重 END ---------------------

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
	ch := make(chan bool, 1)             // 是否是创建文件
	ch1 := make(chan bool, 1)            // 内容写入到文件中中是否成功
	mkdocs := make(chan bool, 1)         // mkdocs 是否修改成功
	filenameBeta := make(chan string, 1) // 文件修改时的文件名

	fileCount := getFileCount(root)
	filename := fmt.Sprintf("第%d期（%s）.md", fileCount, time2.NewTime().Date())
	//filename := fmt.Sprintf("%s.md", time2.NewTime().Date())

	go func() {
		defer wg.Done()

		dates := make([]string, 10)
		_ = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			if strings.Contains(info.Name(), ".md") {
				dates = append(dates, info.Name())
			}

			return nil
		})

		newFilenameBeta := filename
		matches := fuzzy.Find(time2.NewTime().Date(), dates)
		if len(matches) > 0 {
			for _, match := range matches {
				filename = match.Str
				break
			}
		}

		if isDownload == false {
			homepage = ""
		}

		if len(homepage) > 0 {
			downloadImage()
		}

		var content string
		if len(descriptionVar) < 1 {
			content = fmt.Sprintf(contentTemplate(), repository, baseUrl, language, description)
		} else if strings.Compare(label, "article") == 0 {
			content = fmt.Sprintf(contentTemplate(), description, baseUrl)
		} else {
			repository = strings.TrimRight(strings.ReplaceAll(repository, "https://", ""), "/")
			content = fmt.Sprintf(contentTemplate(), repository, baseUrl, description)
		}

		newFilename := filepath.Join(root, filename)
		if _, err := os.Stat(newFilename); err != nil {
			newFilename = filepath.Join(root, newFilenameBeta)
			f, errs := os.Create(newFilename)
			if errs != nil {
				ch <- false
				ch1 <- false
				return
			}

			ch <- true
			_, _ = f.WriteString(fmt.Sprintf(`# %s

![view-count](https://count.getloli.com/@xiaoxuan6-weekly-%s)

---%s`, strings.ReplaceAll(filename, ".md", ""), time.Now().Format("20060102"), content))

			ch1 <- true

		} else {

			ch <- false

			f, _ := os.OpenFile(newFilename, os.O_WRONLY|os.O_APPEND, os.ModePerm)
			_, _ = f.WriteString(content)

			ch1 <- true
			filenameBeta <- newFilename
		}

		// ---------------- write links ------------------------
		linkContent := fmt.Sprintf("%s\n", uri)
		if len(homepage) > 0 {
			linkContent = fmt.Sprintf("%s%s\n", linkContent, homepage)
		}
		f, _ := os.OpenFile(links, os.O_WRONLY|os.O_APPEND, os.ModePerm)
		_, _ = f.WriteString(linkContent)

		return
	}()

	go func() {
		defer wg.Done()

		if item := <-ch; item == false {
			mkdocs <- true
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

		mkdocs <- false
		return
	}()

	wg.Wait()

	if result := <-ch1; result != true {
		logrus.Info(color.RedString("创建文件 %s 失败", filename))
		return
	}

	var message string
	commitSlice := make([]string, 2)
	if ok := <-mkdocs; ok {
		filenameCh := <-filenameBeta
		newFilename := strings.ReplaceAll(filenameCh, path, "")
		newFilename = strings.ReplaceAll(newFilename, "/docs", "docs")

		commitSlice = append(commitSlice, newFilename, "links.txt")

		message = fmt.Sprintf("fix: Update %s", filepath.Base(filenameCh))
	} else {
		commitSlice = append(commitSlice, filepath.Join("docs/", year, filename), "mkdocs.yml", "links.txt")

		message = fmt.Sprintf("feat: Add %s", filename)
	}

	logrus.Info("commitSlice", commitSlice)

	if err2 := gitCommit(message, commitSlice...); err2 != nil {
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
	response, err := http.DefaultClient.Get(fmt.Sprintf("https://ungh.xiaoxuan6.me/repos/%s/%s", owner, repo))
	if err != nil {
		return ""
	}
	defer response.Body.Close()

	b, _ := ioutil.ReadAll(response.Body)
	description := gjson.GetBytes(b, "repo.description").String()
	language = gjson.GetBytes(b, "repo.language").String()
	homepage = gjson.GetBytes(b, "repo.homepage").String()

	lang := whatlanggo.DetectLang(description)
	sourceLang := strings.ToLower(lang.Iso6391())
	if sourceLang != "zh" {
		res := deeplx.Translate(description, sourceLang, "zh")
		if res.Code == 200 {
			description = res.Data
		}
	}

	return description
}

func distinct(path string) bool {
	result := false
	r, _ := os.ReadFile(path)
	br := bufio.NewReader(strings.NewReader(string(r)))
	for {
		a, _, errs := br.ReadLine()
		if errs == io.EOF {
			break
		}

		if strings.Compare(string(a), uri) == 0 {
			result = true
			break
		}
	}

	return result
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

func getFileCount(path string) int {
	fileCount := 1
	_ = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if strings.Compare(filepath.Ext(info.Name()), ".md") == 0 {
			fileCount += 1
		}

		return nil
	})

	return fileCount
}

func downloadImage() {
	if len(homepage) < 0 {
		return
	}

	filePath := filepath.Join(path, "/docs/static/images/", time2.NewTime().Date())
	if _, err := os.Stat(filePath); err != nil {
		_ = os.MkdirAll(filePath, os.ModePerm)
	}

	img = fmt.Sprintf("%s/%s.png", filePath, strconv.Itoa(int(time.Now().Unix())))

	response, err := http.DefaultClient.Get(fmt.Sprintf("https://wr.do/api/v1/scraping/screenshot?url=%s&key=8e4c0fac-5526-4d81-a2f0-e534251ea457", strings.TrimRight(homepage, "/")))
	defer response.Body.Close()
	if err != nil {
		return
	}

	if response.StatusCode != 200 {
		return
	}

	f, _ := os.OpenFile(img, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer f.Close()
	b, _ := ioutil.ReadAll(response.Body)
	_, _ = f.Write(b)

	time.Sleep(3)
	watermark(img)
}

func watermark(img string) {
	i, err := mergi.Import(impexp.NewFileImporter(img))
	if err != nil {
		logrus.Error("watermark: ", err.Error())
		return
	}

	watermarkImage, err := mergi.Import(impexp.NewFileImporter("watermark.png"))
	if err != nil {
		logrus.Error("watermark: ", err.Error())
		return
	}

	res, err := mergi.Watermark(watermarkImage, i, image.Pt(0, 0))
	if err != nil {
		logrus.Error("watermark: ", err.Error())
		return
	}

	_ = mergi.Export(impexp.NewFileExporter(res, img))
}

func contentTemplate() (templateBase string) {
	if len(descriptionVar) < 1 {
		templateBase = `
## [%s](%s)
- 所属语言：%s
- 项目说明：%s
`
	} else if strings.Compare(label, "article") == 0 {
		templateBase = `
## [%s](%s)
`
	} else {
		templateBase = `
## [%s](%s)
- 项目说明：%s
`
	}

	if len(homepage) > 0 {
		u := "https://mirror.ghproxy.com/https://raw.githubusercontent.com/xiaoxuan6/weekly/main/docs"
		templateBase = fmt.Sprintf("%s![img](%s/static/images/%s/%s){.img-fluid tag=1}\n", templateBase, u, time2.NewTime().Date(), filepath.Base(img))

		if len(descriptionVar) == 0 {
			templateBase = fmt.Sprintf("%s- 官网地址：[%s](%s)\n", templateBase, homepage, homepage)
		}
	}

	if len(demoUrl) > 0 {
		templateBase = fmt.Sprintf("%s- 相关链接：[Demo](%s)\n", templateBase, demoUrl)
	}

	return
}

func gitCommit(message string, filename ...string) error {
	w, _ := gitRepository.Worktree()

	for _, file := range filename {
		_, err := w.Add(file)
		if err != nil {
			return fmt.Errorf("git add file %s fail: %s", file, err.Error())
		}
	}

	status, _ := w.Status()
	logrus.Info("git status", status.String())

	if _, err := w.Commit(message, &git.CommitOptions{
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
