package main

import (
	cj "awesomecj/src/cj-core"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	url := "https://api.bilibili.com/x/v2/reply/main?&next=0&type=11&oid=135459278&mode=2&plat=1&_=1624272531135"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	} else {
		request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (HTML, like Gecko) Chrome/91.0.4472.106 Safari/537.36")
	}
	response, _ := cj.Client.Do(request)
	result, _ := ioutil.ReadAll(response.Body)
	var res cj.NewestComment
	err = json.Unmarshal(result, &res)
	// With out condition.
	//luckyList := cj.CjManyTimesByFixWorker("135459278", res.Data.Cursor.Prev, 5,false,"")
	//log.Println("总人数：",len(luckyList))
	//log.Println("中奖列表",luckyList)
	// With condition.
	//luckyListWithCondition := cj.CjManyTimesByFixWorker("135459278", res.Data.Cursor.Prev, 5,true,"拉低")
	//log.Println("总人数：",len(luckyListWithCondition))
	//log.Println("中奖列表",luckyListWithCondition)
	//log.Println("如果少于你的规定人数，请更改条件再试试~")
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println("关闭请求失败", err)
		}
	}(response.Body)
}
