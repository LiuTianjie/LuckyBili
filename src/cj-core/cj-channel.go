package cj_core

import (
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var replyChan = make(chan reply)
var done = make(chan bool)

func CJWithChannelByFixWorker(oid string, rg int, num int, filter bool, condition string) {
	var wgList sync.WaitGroup = sync.WaitGroup{}
	go WriteLuckyListWithChannel()
	wgList.Add(10)
	// Allocate jobs.
	avgTaskNum := num / 10
	last := num % 10
	for i := 0; i < 10; i++ {
		if i < 9 {
			go func() {
				for j := 0; j < avgTaskNum; j++ {
					SingleWithChannel(oid, rg, filter, condition)
				}
				wgList.Done()
			}()
		}
		if i == 9 {
			go func() {
				for k := 0; k < last+avgTaskNum; k++ {
					SingleWithChannel(oid, rg, filter, condition)
				}
				wgList.Done()
			}()
		}
		// Each goroutine should have a gap.
		time.Sleep(time.Microsecond * 500)
	}
	wgList.Wait()
	close(replyChan)
	if _, ok := <-done; ok {
		close(done)
		return
	}
}

func SingleWithChannel(oid string, rg int, filter bool, condition string) {
	count := 0
	exist := make(map[int]bool)
	for {
		n := rand.Intn(rg)
		if !exist[n] {
			exist[n] = true
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
			replyChan <- luckyPerson
			break
		}
	}
}

func WriteLuckyListWithChannel() {
	for {
		if luckyPerson, ok := <-replyChan; ok {
			LuckyList = append(LuckyList, luckyPerson)
		} else {
			// When the replyChan is closed, notify main goroutine to continue to work.
			done <- true
			return
		}
	}
}
