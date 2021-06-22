package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	urlpkg "net/url"
	"strconv"
	"sync"
	"time"
)

var client http.Client
var data map[string]interface{}

type newestComment struct {
	Data commentData `json:"data"`
}

type commentData struct {
	Cursor  cursor  `json:"cursor"`
	Replies []reply `json:"replies"`
}

type cursor struct {
	AllCount int `json:"all_count"`
	Prev     int `json:"prev"`
	Next     int `json:"next"`
}

type reply struct {
	Floor   int     `json:"floor"`
	Member  member  `json:"member"`
	Content content `json:"content"`
}
type member struct {
	Mid    string `json:"mid"`
	Uname  string `json:"uname"`
	Avatar string `json:"avatar"`
}

type content struct {
	Message string `json:"message"`
}

func main() {
	url := "https://api.bilibili.com/x/v2/reply/main?&next=0&type=11&oid=135459278&mode=2&plat=1&_=1624272531135"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	} else {
		request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (HTML, like Gecko) Chrome/91.0.4472.106 Safari/537.36")
	}
	response, _ := client.Do(request)
	result, _ := ioutil.ReadAll(response.Body)
	var res newestComment
	err = json.Unmarshal(result, &res)
	luckyList := cjManyTimes("135459278", res.Data.Cursor.Prev, 20)
	log.Println("中奖列表为：")
	for _,l:=range luckyList{
		log.Println(l)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println("关闭请求失败", err)
		}
	}(response.Body)
}

func cj(oid string, page string) (luckyOne reply) {
	//生成随机数
	rand.Seed(time.Now().Unix())
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("User-Agent", "Mozilla/4.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.106 Safari/537.36")
	url := "https://api.bilibili.com/x/v2/reply/main?&next=" + page + "&type=11&oid=" + oid + "&mode=2&plat=1&_=1624272531135"
	req.URL, _ = urlpkg.Parse(url)
	response, _ := client.Do(req)
	res, _ := ioutil.ReadAll(response.Body)
	var ChosenOne newestComment
	err = json.Unmarshal(res, &ChosenOne)
	luckyPoint := rand.Intn(len(ChosenOne.Data.Replies))
	luckyOne = ChosenOne.Data.Replies[luckyPoint]
	return
}

func cjManyTimes(oid string, rg int, num int) (luckyList []reply) {
	var exist sync.Map
	var muList sync.Mutex = sync.Mutex{}
	var wgList sync.WaitGroup = sync.WaitGroup{}
	wgList.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			// TODO:代理池分发任务，改造cj函数
			for {
				n := rand.Intn(rg)
				_, ok := exist.Load(n)
				if !ok {
					exist.Store(n, true)
					luckyPerson := cj(oid, strconv.Itoa(n))
					muList.Lock()
					luckyList = append(luckyList, luckyPerson)
					muList.Unlock()
					wgList.Done()
					break
				}
			}
		}()
	}
	wgList.Wait()
	return
}

func cjWithCondition(oid string, condition string, rg int, num int) (luckyList []reply) {
	return
}
