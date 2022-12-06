// Colly爬虫框架的示例

package xunacg

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

type Result struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// 签到
func SignIn() {

	// 1.创建collector收集器
	c := colly.NewCollector()

	//2.设置gbk编码，可重复访问
	c.DetectCharset = true
	c.AllowURLRevisit = true
	c.Async = false // 异步

	// beginRevist := false

	cookie := "__51vcke__JKGOgqiI9YURpbFt=cef27301-76b4-5350-99a8-e3d5ffd2d908; __51vuft__JKGOgqiI9YURpbFt=1667371952189; bbs_sid=2d6u03b5ntc6kk1ninjb52h9lq; _wish_accesscount_visited=1; __51uvsct__JKGOgqiI9YURpbFt=50; bbs_token=pGWZ67tVueYj2bwqK4q_2FINaZF_2FBTIBF5mcbl2uISItka3rUl; __vtins__JKGOgqiI9YURpbFt=%7B%22sid%22%3A%20%225ce2fad9-6512-5687-a7bd-56a6f890b569%22%2C%20%22vd%22%3A%2010%2C%20%22stt%22%3A%202920451%2C%20%22dr%22%3A%201318%2C%20%22expires%22%3A%201670307356613%2C%20%22ct%22%3A%201670305556613%7D"
	// d := &data{
	// 	uid: 54311,
	// }
	// da, err := json.Marshal(d)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("请求前，准备工作")
		r.Headers.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 Edg/107.0.1418.62")
		// r.Headers.Set(":method", "POST")
		// r.Headers.Set(":path", "/my-free.htm")
		// r.Headers.Set(":scheme", "https")
		r.Headers.Set("origin", "https://www.xunacg.xyz")
		r.Headers.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
		// r.Headers.Add(":authority", "www.xunacg.xyz")
		r.Headers.Set("x-requested-with", "XMLHttpRequest")

		r.Headers.Set("cookie", cookie)
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("请求失败：" + err.Error())
	})
	count := 0
	max := 10
	c.OnResponse(func(response *colly.Response) {

		var result Result
		err := json.Unmarshal(response.Body, &result)
		if err != nil {
			fmt.Println("签到失败:" + string(response.Body))
		} else {
			if count >= max {
				fmt.Printf("已签到: %d 次，已完成", count)
				return
			}
			count += 1
			fmt.Printf("已签到: %d 次", count)
			// 等待下一次
			time.Sleep(30 * time.Second)

			c.PostRaw("https://www.xunacg.xyz/my-free.htm", []byte("uid=54311"))
		}

	})
	fmt.Println("开始签到")
	c.PostRaw("https://www.xunacg.xyz/my-free.htm", []byte("uid=54311"))

}
