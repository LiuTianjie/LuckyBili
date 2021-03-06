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

// Cj Core random process.
func Cj(oid string, page string, fliter bool, condition string) (luckyOne reply) {
	rand.Seed(time.Now().Unix())
	// Req request, only need to change the header.
	req, _ := http.NewRequest("GET", "", nil)
	req.Header.Add("User-Agent", "Mozilla/4.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.106 Safari/537.36")
	url := "https://api.bilibili.com/x/v2/reply/main?&next=" + page + "&type=11&oid=" + oid + "&mode=2&plat=1&_=1624272531135"
	req.URL, _ = urlpkg.Parse(url)
	response, _ := Client.Do(req)
	res, _ := ioutil.ReadAll(response.Body)
	var ChosenPage NewestComment
	if err := json.Unmarshal(res, &ChosenPage); err != nil {
		log.Println("JSON解码错误")
		return
	}
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
func CjManyTimesByFixWorker(oid string, rg int, num int, filter bool, condition string) {
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
					singleTime(&exist, &muList, rg, oid, filter, condition)
				}
				wgList.Done()
			}()
		}
		if i == 9 {
			go func() {
				for k := 0; k < last+avgTaskNum; k++ {
					singleTime(&exist, &muList, rg, oid, filter, condition)
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
func singleTime(exist *sync.Map, muList *sync.Mutex, rg int, oid string, filter bool, condition string) {
	count := 0
	for {
		n := rand.Intn(rg)
		_, ok := exist.Load(n)
		if !ok {
			exist.Store(n, true)
			luckyPerson := Cj(oid, strconv.Itoa(n), filter, condition)
			if filter {
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
			LuckyList = append(LuckyList, luckyPerson)
			muList.Unlock()
			break
		}
	}
}
