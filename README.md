# Spider_Go
通过Go语言实现的简单爬虫（内涵段子, 豆瓣电影海报）

## Run Mode
go run spider.go

## Close Mode
exit

## jokeSpider编码问题
由于内涵段子吧 采用 gb2312格式编码， Go语言采用utf-8格式，调用 第三方库 iconv 实现转码功能
第三方库下载方法： go get github.com/djimenez/iconv-go
        安装方法： go install github.com/djimenez/iconv-go
