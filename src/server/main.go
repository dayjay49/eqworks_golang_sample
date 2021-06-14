package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/dayjay49/ws-product-golang-master/src/server/mycounters"
	"github.com/dayjay49/ws-product-golang-master/src/server/mywatcher"
	"github.com/dayjay49/ws-product-golang-master/src/server/ratelimiter"
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func statsHandler(w http.ResponseWriter, r *http.Request, ci mycounters.CounterInterface) {
	if !isAllowed() {
		w.WriteHeader(429)
		return
	}

	// GO GET EM!!
	ms, err := ci.GetMockStore()
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

func viewHandler(w http.ResponseWriter, r *http.Request, ci mycounters.CounterInterface) {
	data = content[rand.Intn(len(content))]

	// GO GET EM!!
	counter, err := ci.GetUpdatedCounter(data)
	checkErr(err)

	counter.IncrementView()

	err = processRequest(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	// simulate random click call
	if rand.Intn(100) < 50 {
		counter.IncrementClick()
	}
}

func main() {
	// Running counter uploader concurrently (go routines inside the function)
	counterInterface, rateLimitInterface, err := mywatcher.NewMyWatcher(
		&mycounters.CounterConfig{
			CycleDuration: 3 * time.Second,
			InitialContent: data,
		}, 
		&ratelimiter.RateLimitConfig{
			ActiveTokenLimit: 5,
			FixedInterval: 15 * time.Second,
		},
	)

	checkErr(err)
	fmt.Println(counterInterface, "--------------COUNTER-------------")
	fmt.Println(rateLimitInterface, "--------------RATE-LIMITER---------------")
		
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/view/", func (w http.ResponseWriter, r *http.Request) {
		viewHandler(w, r, counterInterface)
	})
	http.HandleFunc("/stats/", func (w http.ResponseWriter, r *http.Request) {
		statsHandler(w, r, counterInterface)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
