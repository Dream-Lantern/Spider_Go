package main

import "fmt"
import "net/http"
import "log"
import "io/ioutil"
import "regexp"
import "os"
import iconv "github.com/djimenez/iconv-go"

/*
	明确目标  内涵段子
	第1页  https://www.neihanba.com/dz/index.html
	第2页  https://www.neihanba.com/dz/list_2.html
	第n页 https://www.neihanba.com/dz/list_n.html

	1. 首先进入某页的页码主页，----> 取出每个段子链接地址
	 https://www.neihanba.com 拼接一个段子的完整url路径
	得到每个段子路径的正则表达式  `<h4> <a href="(.*?)"`
	https://www.neihanba.com + /dz/1092886.html

	2. 进入每个段子的首页，得到段子的标题和内容
	标题的正则
	`<h1>(.*?)</h1>`

	内容的正则
	`<td><p>(?s:(.*?))</p></td>`
*/

var URL string = "https://www.neihanba.com"

type Spider struct {
	Page       int //当前爬虫已经爬取到了第几页
	Url        string
	regDZ      string
	regTitle   string
	regContent string
}

// 初始化 段子，标题，内容的 正则
func (this *Spider) init_reg() {
	// 段子page 正则
	this.regDZ = `<h4> <a href="(.*?)"`
	// 标题title 正则
	this.regTitle = `<h1>(.*?)</h1>`
	// 内容content 正则
	this.regContent = `<td><p>(?s:(.*?))</p></td>`
}

// 爬取一个某页的菜单页码
func (this *Spider) Spider_one_page() {
	fmt.Println("正在爬取 ", this.Page, " 页")
	if this.Page == 1 {
		this.Url = URL + "/dz/index.html"
	} else {
		this.Url = fmt.Sprintf("https://www.neihanba.com/dz/list_%d.html", this.Page)
	}

	Data, StatCode := this.HttpGet(this.Url)
	if StatCode != 200 {
		fmt.Println("HttpGet err")
		return
	}

	this.init_reg()

	// 获取页码内 所有段子
	DZ := regexp.MustCompile(this.regDZ)
	dzData := DZ.FindAllStringSubmatch(Data, -1)
	// get title
	TITLE := regexp.MustCompile(this.regTitle)
	// get content
	CONTENT := regexp.MustCompile(this.regContent)

	for _, test := range dzData {
		this.Url = URL + test[1]
		titleHTML, titleCode := this.HttpGet(this.Url)
		if titleCode != 200 {
			log.Println("title err")
			return
		}
		ttData := TITLE.FindAllStringSubmatch(titleHTML, -1)
		contentData := CONTENT.FindAllStringSubmatch(titleHTML, -1)

		fileName := fmt.Sprintf("./joke-%d.txt", this.Page)
		for _, text := range ttData {
			retTitle := this.writeFile(fileName, text[1])
			retTitle = this.writeFile(fileName, "\n\n")
			if retTitle == false {
				log.Println("write title err")
			}
		}
		for _, cont := range contentData {
			retCont := this.writeFile(fileName, cont[1])
			retCont = this.writeFile(fileName, "\n\n")
			if retCont == false {
				log.Println("write content err")
			}
		}
	}
}

// 写文件
func (this *Spider) writeFile(fileName string, data string) bool {
	fp, errOpen := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if errOpen != nil {
		log.Println("open file err")
		return false
	}
	defer fp.Close()
	_, errWrite := fp.WriteString(data)
	if errWrite != nil {
		log.Println("write err")
		return false
	}
	return true
}

//请求一个页码将页码中的全部源码content
func (this *Spider) HttpGet(url string) (content string, statusCode int) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		content = ""
		statusCode = -100
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		statusCode = resp.StatusCode
		content = ""
		return
	}
	// gb2312 --> utf-8
	out := make([]byte, len(data))
	out = out[:]

	iconv.Convert(data, out, "gb2312", "utf-8")

	content = string(out)
	statusCode = resp.StatusCode

	return
}

func (this *Spider) DoWork() {
	fmt.Println("Spider begin to  work")
	this.Page = 1

	var cmd string

	for {
		fmt.Println("请输入任意键爬取下一页，输入exit退出")
		fmt.Scanf("%s", &cmd)
		if cmd == "exit" {
			fmt.Println("exit")
			break
		}
		//需要爬取下一页
		this.Spider_one_page()

		this.Page++
	}
}

func main() {
	sp := new(Spider)
	sp.DoWork()
}
