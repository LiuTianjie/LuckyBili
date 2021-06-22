package cj_core

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	urlpkg "net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Client http.Client

type NewestComment struct {
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

// Cj Core random process.
func Cj(oid string, page string, fliter bool, condition string) (luckyOne reply) {
	rand.Seed(time.Now().Unix())
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("User-Agent", "Mozilla/4.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.106 Safari/537.36")
	url := "https://api.bilibili.com/x/v2/reply/main?&next=" + page + "&type=11&oid=" + oid + "&mode=2&plat=1&_=1624272531135"
	req.URL, _ = urlpkg.Parse(url)
	response, _ := Client.Do(req)
	res, _ := ioutil.ReadAll(response.Body)
	var ChosenPage NewestComment
	err = json.Unmarshal(res, &ChosenPage)
	// Without condition, we just need to randomly chose one.
	if !fliter {
		luckyPoint := rand.Intn(len(ChosenPage.Data.Replies))
		luckyOne = ChosenPage.Data.Replies[luckyPoint]
		return
	} else {
		// If there is some condition, to speed up, we should do a search first.
		var waitList []reply
		for _, p := range ChosenPage.Data.Replies {
			if strings.Contains(p.Content.Message, condition) {
				waitList = append(waitList, p)
			}
		}
		if len(waitList) != 0 {
			luckyPoint := rand.Intn(len(waitList))
			luckyOne = waitList[luckyPoint]
			return
		} else {
			return
		}
	}
}

// CjManyTimesByFixWorker Use 10 workers to request, each worker is sync.
// 6.22 test, max number can be 300.
func CjManyTimesByFixWorker(oid string, rg int, num int, fliter bool, condition string) (luckyList []reply) {
	var wgList sync.WaitGroup = sync.WaitGroup{}
	wgList.Add(10)
	// Allocate jobs.
	avgTaskNum := num / 10
	last := num % 10
	var exist sync.Map
	var muList sync.Mutex = sync.Mutex{}
	for i := 0; i < 10; i++ {
		if i < 9 {
			go func() {
				for j := 0; j < avgTaskNum; j++ {
					singleTime(&exist, &muList, rg, oid, &luckyList, fliter, condition)
				}
				wgList.Done()
			}()
		}
		if i == 9 {
			go func() {
				for k := 0; k < last+avgTaskNum; k++ {
					singleTime(&exist, &muList, rg, oid, &luckyList, fliter, condition)
				}
				wgList.Done()
			}()
		}
		// Each goroutine should have a gap.
		time.Sleep(time.Microsecond * 500)
	}
	wgList.Wait()
	return
}

// Single Time write the sync map and array.
func singleTime(exist *sync.Map, muList *sync.Mutex, rg int, oid string, luckyList *[]reply, fliter bool, condition string) {
	count := 0
	for {
		n := rand.Intn(rg)
		_, ok := exist.Load(n)
		if !ok {
			exist.Store(n, true)
			luckyPerson := Cj(oid, strconv.Itoa(n), fliter, condition)
			if fliter {
				if luckyPerson.Member.Mid == "" {
					count += 1
					// Count is the max retry times.
					if count < 10 {
						continue
					} else {
						break
					}
				}
			}
			muList.Lock()
			*luckyList = append(*luckyList, luckyPerson)
			muList.Unlock()
			break
		}
	}
}
