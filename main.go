// Colly爬虫框架的示例

package main

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func getUrl(tag string, pageNum int) string {
	return fmt.Sprintf("http://www.llss.app/wp/tag/%s/page/%d/", tag, pageNum)
}

// 通过询问的方式获取参数（已弃用）
func scanInputs() (int, int, string) {
	var pageStr string
	var currPage, endPage int

	fmt.Print("请输入要采集的页范围，一个值则代表为结束页(默认1,2):")
	fmt.Scanln(&pageStr) // 接收

	var tag string
	fmt.Print("请输入要采集的标签(默认为汉化单行本):")
	fmt.Scanln(&tag) // 接收输入
	if tag == "" {
		tag = "汉化单行本"
	}

	return currPage, endPage, tag
}

/*通过传参的方式获取参数
startPage: 起始页
endPage:结束页
tag:标签
enableAsync:是否开启异步采集
showName:是否显示标题
*/
func getParamsByScanInput() (startPage int, endPage int, tag string, enableAsync bool, showTitle bool) {
	var pageStr string
	flag.StringVar(&pageStr, "page", "1", "页数范围")
	flag.IntVar(&endPage, "end", 0, "结束页")
	flag.StringVar(&tag, "tag", "汉化单行本", "要采集的标签，默认为汉化单行本")
	flag.BoolVar(&enableAsync, "async", true, "是否开启异步采集")
	flag.BoolVar(&showTitle, "title", false, "是否显示标题")
	//转换
	flag.Parse()

	// 解析页数
	if pageStr == "" {
		startPage = 1
	} else {
		tmpArr := strings.Split(pageStr, ",")
		if len(tmpArr) > 1 { // 两个数字（1,2）则是起始页,结束页
			startPage, _ = strconv.Atoi(tmpArr[0])
			endPage, _ = strconv.Atoi(tmpArr[1])
		} else if len(tmpArr) > 0 { // 只输入了一个数字，则是只采集这一页
			startPage, _ = strconv.Atoi(tmpArr[0])

		}
	}

	if endPage <= 0 {
		endPage = startPage
	}
	return
}
func main() {

	// currPage, endPage, tag := scanInputs()
	startPage, endPage, tag, enableAsync, showTitle := getParamsByScanInput()

	// 1.创建collector收集器
	c := colly.NewCollector()

	//2.设置gbk编码，可重复访问
	c.DetectCharset = true
	c.AllowURLRevisit = true
	c.Async = enableAsync // 异步

	//3.clone collector用于内容解析
	contentCollector := c.Clone() //拷贝
	// beginRevist := false

	// 主页中，查找每个子项的地址
	c.OnHTML("article h1 a", func(e *colly.HTMLElement) {
		contentCollector.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		// 主页中的内容不需要
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error %s: %v\n", r.Request.URL, err)
	})
	// 找到的总个数
	count := 0
	// 磁力的正则表达式
	reg := regexp.MustCompile("[A-Za-z0-9]{20,}")
	contentCollector.OnHTML("article", func(e *colly.HTMLElement) {
		count = count + 1
		fmt.Printf("%d:", count)

		if showTitle {
			fmt.Printf("%s\n", e.ChildText(".entry-title"))
		}

		//通过正则表达式找到磁力下载链接
		arr := reg.FindAllString(e.Text, -1) // 查找所有

		for _, value := range arr {

			fmt.Printf("%s\n", value) //也可以在这里把内容保存到文件中
		}

	})

	//正式启动网页访问
	fmt.Printf("开始采集：tag:%s(%d-%d)\n\n", tag, startPage, endPage)
	// 循环采集
	for i := startPage; i <= endPage; i++ {
		fmt.Printf("第%d页\n", i)
		c.Visit(getUrl(tag, i))
	}

	if enableAsync {
		c.Wait()
		contentCollector.Wait()
	}

}
