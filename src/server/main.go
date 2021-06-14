package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	counteruploader "github.com/dayjay49/ws-product-golang-master/src/server/uploadCounters"
)

var (
	content = []string{"sports", "entertainment", "business", "education"}
	data = content[rand.Intn(len(content))]
)

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to EQ Works 😎")
}

func processRequest(r *http.Request) error {
	time.Sleep(time.Duration(rand.Int31n(50)) * time.Millisecond)
	return nil
}

func isAllowed() bool {
	return true
}

// func uploadCounters() error {
// 	return nil
// }

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func statsHandler(w http.ResponseWriter, r *http.Request, uploader counteruploader.CounterUploader) {
	if !isAllowed() {
		w.WriteHeader(429)
		return
	}

	// GO GET EM!!
	ms, err := uploader.GetMockStore()
	checkErr(err)

	// display counters
	ms.Lock()
	for eventKey, eventValues := range ms.EventHistory {
		fmt.Println(eventKey+" -------> "+"{views: "+strconv.Itoa(eventValues["views"])+
		", clicks: "+strconv.Itoa(eventValues["clicks"])+"}")
		
		fmt.Fprint(w, eventKey+" -------> "+"{views: "+strconv.Itoa(eventValues["views"])+
			", clicks: "+strconv.Itoa(eventValues["clicks"])+"}\n")
	}
	ms.Unlock()
}

func viewHandler(w http.ResponseWriter, r *http.Request, uploader counteruploader.CounterUploader) {
	data = content[rand.Intn(len(content))]

	// GO GET EM!!
	c, err := uploader.GetUpdatedCounter(data)
	checkErr(err)

	c.IncrementView()

	err = processRequest(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	// simulate random click call
	if rand.Intn(100) < 50 {
		c.IncrementClick()
	}
}

func main() {
	// Running counter uploader concurrently (go routines inside the function)
	uploader, err := counteruploader.NewMyCounterUploader(&counteruploader.Config{
		CycleDuration: 3 * time.Second,
		InitialContent: data,
	})
	checkErr(err)

	// fmt.Println(uploader, "----------------")

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/view/", func (w http.ResponseWriter, r *http.Request) {
		viewHandler(w, r, uploader)
	})
	http.HandleFunc("/stats/", func (w http.ResponseWriter, r *http.Request) {
		statsHandler(w, r, uploader)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
