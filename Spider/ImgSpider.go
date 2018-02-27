package main

import "fmt"
import "net/http"
import "log"
import "io/ioutil"
import "regexp"
import "os"

/*
	明确目标 豆瓣排名前250名 电影图片
	第1页  https://movie.douban.com/top250?start=0&filter=
	第2页  https://movie.douban.com/top250?start=25&filter=
	第n页  https://movie.douban.com/top250?start=50&filter=

	图片的格式：<img width="100" alt="指环王2：双塔奇兵" src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p909265336.jpg" class="">
	图片的正则: `src="(.*?)" class="">`

	片名的格式：<img width="100" alt="
	片名的正则：`<img width="100" alt="(.*?)"`
*/

var URL string = "https://movie.douban.com/top250?start="

type Spider struct {
	Page    int
	Url     string
	regName string
	regImg  string
}

// 初始化 正则表达式
func (this *Spider) init_reg() {
	this.regName = `<img width="100" alt="(.*?)"`
	this.regImg = `src="(.*?)" class="">`
}

//爬取一个某页的菜单页码
func (this *Spider) Spider_one_page() {
	fmt.Println("正在爬取 ", this.Page, " 页")
	if this.Page == 1 {
		this.Url = fmt.Sprintf("https://movie.douban.com/top250?start=%d&filter=", 0)
	} else {
		this.Url = fmt.Sprintf("https://movie.douban.com/top250?start=%d&filter=", (this.Page-1)*25)
	}

	Data, StatCode := this.HttpGet(this.Url)
	if StatCode != 200 {
		log.Println("HttpGet err")
		return
	}

	// 改变 工作目录
	dirName := fmt.Sprintf("./Page-%d", this.Page)
	dirErr := os.Mkdir(dirName, os.ModePerm)
	if dirErr != nil {
		log.Println("mkdir err")
		return
	}
	dirErr = os.Chdir(dirName)
	if dirErr != nil {
		log.Println("chdir err")
		return
	}

	// 初始化 正则表达式
	this.init_reg()

	// 获取页码内 所有电影图片和名字
	Name := regexp.MustCompile(this.regName)
	nameData := Name.FindAllStringSubmatch(Data, -1)
	// get imgData
	Img := regexp.MustCompile(this.regImg)
	imgData := Img.FindAllStringSubmatch(Data, -1)
	// 创建 存放电影名词的数组
	fileArr := make([]string, 0)
	suffix := ".jpg"

	for _, value := range nameData {
		fileName := "./" + value[1] + suffix
		fileErr := this.writeFile(fileName, "")
		if fileErr == false {
			log.Println("create file err")
		}
		// 将fileArr数组 填满 文件名字
		fileArr = append(fileArr, fileName)
	}
	for index, imgUrl := range imgData {
		binData, imgCode := this.HttpGet(imgUrl[1])
		if imgCode != 200 {
			log.Println("get ImgUrl err")
			return
		}
		writeErr := this.writeFile(fileArr[index], binData)
		if writeErr == false {
			log.Println("write file err")
			return
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

	content = string(data)
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
		dirFail := os.Chdir("../")
		if dirFail != nil {
			log.Println("chdir err")
			return
		}

		this.Page++
	}
}

func main() {
	sp := new(Spider)
	sp.DoWork()
}
