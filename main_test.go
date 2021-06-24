package main

//
//func BenchmarkCjManyTimes(b *testing.B) {
//	url := "https://api.bilibili.com/x/v2/reply/main?&next=0&type=11&oid=135459278&mode=2&plat=1&_=1624272531135"
//	request, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		log.Println(err)
//	} else {
//		request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (HTML, like Gecko) Chrome/91.0.4472.106 Safari/537.36")
//	}
//	response, _ := cj.Client.Do(request)
//	result, _ := ioutil.ReadAll(response.Body)
//	var res cj.NewestComment
//	err = json.Unmarshal(result, &res)
//	b.ResetTimer()
//	cj.CjManyTimesByFixWorker("135459278", res.Data.Cursor.Prev, 50, false, "")
//}

//func BenchmarkCj(b *testing.B) {
//	url := "https://api.bilibili.com/x/v2/reply/main?&next=0&type=11&oid=135459278&mode=2&plat=1&_=1624272531135"
//	request, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		log.Println(err)
//	} else {
//		request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (HTML, like Gecko) Chrome/91.0.4472.106 Safari/537.36")
//	}
//	response, _ := cj.Client.Do(request)
//	result, _ := ioutil.ReadAll(response.Body)
//	var res cj.NewestComment
//	err = json.Unmarshal(result, &res)
//	n := rand.Intn(res.Data.Cursor.Prev)
//	b.ResetTimer()
//	cj.Cj("135459278", strconv.Itoa(n), false, "")
//}

// TODO: How to write a currency benchmark test in go with channel?

//func BenchmarkCjWithChannel(b *testing.B) {
//	url := "https://api.bilibili.com/x/v2/reply/main?&next=0&type=11&oid=135459278&mode=2&plat=1&_=1624272531135"
//	request, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		log.Println(err)
//	} else {
//		request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (HTML, like Gecko) Chrome/91.0.4472.106 Safari/537.36")
//	}
//	response, _ := cj.Client.Do(request)
//	result, _ := ioutil.ReadAll(response.Body)
//	var res cj.NewestComment
//	err = json.Unmarshal(result, &res)
//	b.ResetTimer()
//	cj.CJWithChannelByFixWorker("135459278", res.Data.Cursor.Prev, 50, false, "")
//}
