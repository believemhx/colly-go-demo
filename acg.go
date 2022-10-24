// Colly爬虫框架的示例

package main

import (
	// "encoding/json"
	"flag"
	"fmt"

	// "io/ioutil"
	// "log"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Result struct {
	title          string
	baiduCloudUrl  string
	baiduCloudCode string
}

func getUrl(category string) string {
	if category == "" {
		return fmt.Sprintf("http://acgyyg.ru")
	}
	return fmt.Sprintf("http://acgyyg.ru/category/%s/", category)
}

// 通过询问的方式获取参数（已弃用）
func scanInputs() (int, int, string) {
	var pageStr string
	var currPage, endPage int

	fmt.Print("请输入要采集的页范围，一个值则代表为结束页(默认1,2):")
	fmt.Scanln(&pageStr) // 接收

	var category string
	fmt.Print("请输入要采集的分类:")
	fmt.Scanln(&category) // 接收输入

	return currPage, endPage, category
}

/*通过传参的方式获取参数
startPage: 起始页
endPage:结束页
tag:标签
enableAsync:是否开启异步采集
showName:是否显示标题
*/
func getParamsByScanInput() (startPage int, endPage int, category string, enableAsync bool, showTitle bool) {
	var pageStr string
	flag.StringVar(&pageStr, "page", "1", "页数范围")
	flag.IntVar(&endPage, "end", 0, "结束页")
	flag.StringVar(&category, "tag", "", "要采集的分类")
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
	startPage, endPage, category, enableAsync, showTitle := getParamsByScanInput()

	// 1.创建collector收集器
	c := colly.NewCollector()

	//2.设置gbk编码，可重复访问
	c.DetectCharset = true
	c.AllowURLRevisit = true
	c.Async = enableAsync // 异步

	var cookies = "X_CACHE_KEY=d856b607d979b62e02f9683a5040967e; PHPSESSID=09cq15t65lkrihsc0usk1ebm12; wordpress_logged_in_9203cc2a22839a9fb197313b29830cf4=BelieveYou%7C1667611948%7CnJAITdHBhRLa9qNBNgKKWjqs5HngQDhVAro9MTwA9vv%7C9d7122aba68020d1547c9cddda6792f66273f8e7fcdf2ced6c4bf887657f62a4"

	//3.clone collector用于内容解析
	contentCollector := c.Clone() //拷贝
	// beginRevist := false

	resMap := make(map[int]*Result, 100)

	// 主页中，查找每个子项的地址
	c.OnHTML(".post-list-view .post-thumbnail a", func(e *colly.HTMLElement) {

		contentCollector.Cookies(cookies)

		// if temp == 2 {
		contentCollector.Visit(e.Attr("href"))
		// }

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
	// regStr := "/(https?|http|ftp|file):\/\/[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]/g"

	// regStr := `^https:\/\/pan\.baidu\.com\/s\/(?=.*[a-z])(?=.*\d)(?=.*[A-Z])[a-z-\dA-Z]{23}$`
	// 百度网盘链接匹配（\Q和\E）代表忽略中间文本中的正则语法
	bdLinkRegStr := `\Qhttps://pan.baidu.com/s/\E[\PP]+ `
	bdLinkReg := regexp.MustCompile(bdLinkRegStr)

	bdGetCodeRegStr := `\Q提取码：\E[a-z0-9A-Z]{4}`
	bdGetCodeReg := regexp.MustCompile(bdGetCodeRegStr)

	contentCollector.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Cookie", cookies)
	})

	contentCollector.OnHTML(".container .post-single", func(e *colly.HTMLElement) {

		//通过正则表达式找到百度网盘链接
		baiduCloudUrl := bdLinkReg.FindString(e.Text)
		if baiduCloudUrl == "" {
			return
		}
		var currentRes = new(Result)
		currentRes.title = e.ChildText(".post-title")
		currentRes.baiduCloudUrl = baiduCloudUrl
		currentRes.baiduCloudCode = bdGetCodeReg.FindString(e.Text)

		if !enableAsync {
			fmt.Printf("%d:", count)
		}
		if showTitle {
			fmt.Printf("\n# %s\n", currentRes.title)
		}

		fmt.Printf("%s  %s\n", currentRes.baiduCloudUrl, currentRes.baiduCloudCode) //也可以在这里把内容保存到文件中

		resMap[count] = currentRes

		count = count + 1
	})

	//正式启动网页访问
	fmt.Printf("开始采集：tag:%s(%d-%d)\n\n", category, startPage, endPage)
	// 循环采集
	for i := startPage; i <= endPage; i++ {
		if !enableAsync {
			fmt.Printf("第%d页\n", i)
		}
		c.Visit(getUrl(category))
	}

	if enableAsync {
		c.Wait()
		contentCollector.Wait()
	}

	fmt.Printf("\n\n采集完成，共%d条", count)

}
