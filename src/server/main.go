package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/dayjay49/ws-product-golang-master/src/server/mycounters"
	"github.com/dayjay49/ws-product-golang-master/src/server/ratelimiter"
)

var (
	content = []string{"sports", "entertainment", "business", "education"}
	data = content[rand.Intn(len(content))]
)

func welcomeHandler(w http.ResponseWriter, r *http.Request, ri ratelimiter.RateLimitInterface) {
	// Acquire a rate limit token and see if more requests are allowed or not
	_, isAllowed, err := ri.Acquire()
	checkErr(err)

	if !isAllowed {
		w.WriteHeader(429)
		return
	}

	fmt.Fprint(w, "Welcome to EQ Works 😎")
}

func viewHandler(w http.ResponseWriter, r *http.Request, ri ratelimiter.RateLimitInterface, ci mycounters.CounterInterface) {
	// Acquire a rate limit token and see if more requests are allowed or not
	_, isAllowed, err := ri.Acquire()
	checkErr(err)
	
	if !isAllowed {
		w.WriteHeader(429)
		return
	}
	
	data = content[rand.Intn(len(content))]

	// Retrieve counter 
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

func statsHandler(w http.ResponseWriter, r *http.Request, ri ratelimiter.RateLimitInterface, ci mycounters.CounterInterface) {
	// Acquire a rate limit token and see if more requests are allowed or not
	_, isAllowed, err := ri.Acquire()
	checkErr(err)
	
	if !isAllowed {
		w.WriteHeader(429)
		return
	}

	// Retrieve the mock store
	ms, err := ci.GetMockStore()
	checkErr(err)

	// display counters from the mock store
	ms.Lock()
	for eventKey, eventValues := range ms.EventHistory {
		fmt.Fprint(w, eventKey+" -------> "+"{views: "+strconv.Itoa(eventValues["views"])+
			", clicks: "+strconv.Itoa(eventValues["clicks"])+"}\n")
	}
	ms.Unlock()
}

func processRequest(r *http.Request) error {
	time.Sleep(time.Duration(rand.Int31n(50)) * time.Millisecond)
	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// counter running concurrently with the website 
	counterInterface, err := mycounters.NewMyCounter(
		&mycounters.CounterConfig{
			CycleDuration: 5 * time.Second,
			InitialContent: data,
		},
	)
	checkErr(err)

	// rate limiter running concurrently with the website 
	rateLimitInterface, err := ratelimiter.NewMyRateLimiter(
		&ratelimiter.RateLimitConfig{
			FixedInterval: 5 * time.Second,
			ActiveTokenLimit: 5,
		},
	)
	checkErr(err)
	
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		welcomeHandler(w,r, rateLimitInterface)
	})
	http.HandleFunc("/view/", func (w http.ResponseWriter, r *http.Request) {
		viewHandler(w,r, rateLimitInterface, counterInterface)
	})
	http.HandleFunc("/stats/", func (w http.ResponseWriter, r *http.Request) {
		statsHandler(w,r, rateLimitInterface, counterInterface)
	})
	
	log.Println("Listening on :1995...")
	log.Fatal(http.ListenAndServe(":1995", nil))
}
