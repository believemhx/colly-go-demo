// Colly爬虫框架的示例

package xunacg

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gocolly/colly"
)

type Result struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func SignByUser(currentUser XunAcgUser, c *colly.Collector, max int) {
	url := "https://www.xunacg.xyz/my-free.htm"
	requestData := map[string]string{
		"uid": strconv.Itoa(currentUser.Uid),
	}

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Cookie", currentUser.Cookie)
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 Edg/107.0.1418.62")
		// r.Headers.Set(":method", "POST")
		// r.Headers.Set(":path", "/my-free.htm")
		// r.Headers.Set(":scheme", "https")
		r.Headers.Set("Origin", "https://www.xunacg.xyz")
		r.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
		// r.Headers.Add(":authority", "www.xunacg.xyz")
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")
		r.Headers.Set("Cookie", currentUser.Cookie)
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("请求失败：" + err.Error())
	})

	c.OnResponse(func(response *colly.Response) {

		var result Result
		err := json.Unmarshal(response.Save(), &result)
		if err != nil {
			fmt.Println("签到失败:" + string(response.Body))
		} else if result.Code == "0" {
			currentUser.Count += 1

			fmt.Printf("已签到: %d 次,%s \n ", currentUser.Count, string(response.Body))

			if max > 0 && currentUser.Count >= max {
				fmt.Printf("%s 签到已完，共 %d 次", currentUser.Name, currentUser.Count)
				return
			}

			// 等待下一次
			time.Sleep(40 * time.Second)

			c.Post(url, requestData)
		} else if result.Code == "-1" {
			//fmt.Printf("%s 签到已完，共 %d 次", currentUser.Name, currentUser.Count)
			fmt.Println(result.Message)
		}

	})
	fmt.Printf("开始签到:%s(%d) \n", currentUser.Name, currentUser.Uid)

	// c.PostRaw(url, []byte("uid="+))
	c.Post(url, requestData)

}

// 签到
func StartSign() {

	data, err := GetData()
	if err != nil {
		fmt.Printf("Json error: %s \n", err.Error())
		return
	}

	// 1.创建collector收集器
	c := colly.NewCollector()

	//2.设置gbk编码，可重复访问
	c.DetectCharset = true
	c.AllowURLRevisit = false
	c.Async = false // 异步

	for _, user := range data.Users {
		if !user.Status {
			SignByUser(user, c, 0)
		}
	}
	c.Wait()
}
