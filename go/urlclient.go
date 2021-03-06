package main

import (
	"fmt"
	"net/http"
	"os"
	"io/ioutil"
	"time"
	"runtime"
	"strconv"
	"sync"
//	"stats"
)

var urls = []string{"http://www.python.org",
	"http://www.ibm.com",
	"http://google.com",
	"http://www.bbc.co.uk/",
	"http://www.json.org/",
	"http://www.lemonde.fr/",
	"http://yahoo.com",
	"http://www.amazon.com",
	"http://www.cnn.com",
	"http://www.nytimes.com",
	"http://www.bradmehldau.com",
	"http://bobhancock.org",
	"http://eventlet.net/doc/design_patterns.html",
	"http://golang.org/doc/effective_go.html",
	"http://docs.sun.com/source/816-6698-10/replicat.html",
	"http://rss.cnn.com/rss/cnn_world.rss",
	"http://rss.cnn.com/rss/cnn_us.rss",
	"http://rss.cnn.com/rss/si_topstories.rss",
	"http://rss.cnn.com/rss/money_latest.rss",
	"http://rss.cnn.com/rss/cnn_allpolitics.rss",
	"http://rss.cnn.com/rss/cnn_crime.rss",
	"http://rss.cnn.com/rss/cnn_tech.rss",
	"http://rss.cnn.com/rss/cnn_space.rss",
	"http://rss.cnn.com/rss/cnn_health.rss",
	"http://rss.cnn.com/rss/cnn_showbiz.rss",
	"http://rss.cnn.com/rss/cnn_travel.rss",
	"http://rss.cnn.com/rss/cnn_living.rss",
	"http://rss.cnn.com/rss/cnn_freevideo.rss",
	"http://rss.cnn.com/rss/cnn_mostpopular.rss",
	"http://rss.cnn.com/rss/cnn_latest.rss",
	"http://www.nytimes.com/services/xml/rss/nyt/Business.xml",
	"http://finance.yahoo.com/rss/headline?s=mhp",
	"http://www.ft.com/servicestools/newstracking/rss#world",
	"http://finance.yahoo.com/rss/headline?s=mhp",
	"http://golang.org"}

const MS_DIVISOR = 1e-9

type request struct {
	url    string
	replyc chan string
}

var times []int64

func urlClient(service chan *request, wg *sync.WaitGroup) {
	for req := range service {
		process(req)
	}
	wg.Done()
}


func process(req *request) {
	//start_process := time.Nanoseconds()
	//fmt.Printf("%s,process start,%d\n",req.url,start_process)
	reply := getURL(req.url)
	//return_op := time.Nanoseconds()
	//fmt.Printf("%s,process opGet,%d\n",req.url, time.Nanoseconds() - start_process)
	req.replyc <- reply
	//fmt.Printf("%s,process <- replyt,%d\n",req.url, time.Nanoseconds() - return_op)
	//fmt.Printf("%s,process all opGet,%d\n", req.url, time.Nanoseconds() - start_process)
}


// Return the contents of a URL 
func getURL(url string) string {
	var d time.Duration
	start := d.Nanoseconds()
	//var b[]byte
	//fmt.Printf("Getting %s\n", url)
	r, err := http.Get(url)
	if err != nil {
		fmt.Printf("\n---> Get()-Error: %s\n\n", err)
		os.Exit(1)
	}
	//b, _ = ioutil.ReadAll(r.Body);
	_, _ = ioutil.ReadAll(r.Body)
	r.Body.Close()

	nsecs := (d.Nanoseconds() - start)
	times = append(times, nsecs)
	//r_time := float64(nsecs) / MS_DIVISOR
	//return fmt.Sprintf("%s|%d|%f",url,len(string(b)),r_time)
	//return fmt.Sprintf("%s,%d,%d",url,len(string(b)),nsecs)
	return fmt.Sprintf("%s,getUrl,%d", url, nsecs)
}

func main() {
	// prog maxprocs host connections reps
	usage := "usage: urlclient iterations gomaxprocs"
	if len(os.Args) < 2 {
		fmt.Println(usage)
		return 
	}

	max_procs, err := strconv.Atoi(os.Args[2])
	runtime.GOMAXPROCS(max_procs)

	maxiters, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Atoi: %d", err)
		return
	}
	if maxiters < 1 {
		maxiters = len(urls)
	}
	iterations := maxiters

	max_urls := len(urls)

	service := make(chan *request)
	wg := new(sync.WaitGroup)
	wg.Add(1)

	// setup consumers
	for i := 0; i < max_urls; i++ {
		go urlClient(service, wg)
	}
	//fmt.Println("iterations=",iterations)

	// setup producers
	for i := 0; i < iterations; i++ {
		reqs := make([]request, max_urls)

		//start_inside_loop := time.Nanoseconds()
		for i := 0; i < max_urls; i++ {
			req := &reqs[i]
			req.url = urls[i]
			req.replyc = make(chan string, 1024)
			service <- req
		}
		//for i := n-1; i >= 0; i-- {   // doesn't matter what order
		//	fmt.Println(<-reqs[i].replyc)
		//}
		//fmt.Printf("Inside_loop %d:%d\n", i, time.Nanoseconds()-start_inside_loop)
	}
	close(service)

	// wiat for the client to send a done signal
	wg.Wait()

	var ftimes []float64
	for _, tv := range times {
//		fmt.Printf("%d %f\n", tv, float64(tv))
		ftimes = append(ftimes, float64(tv))
	}

//	fmt.Printf("Min: %f\n", float64(stats.StatsMin(ftimes)) * 1e-9)
//	fmt.Printf("Max: %f\n", float64(stats.StatsMax(ftimes)) * 1e-9)
//	fmt.Printf("Mean: %f\n", float64(stats.StatsSum(ftimes)) / float64(len(times)) * 1e-9 )
//	fmt.Printf("StdDev: %f\n", stats.StatsSampleStandardDeviation(ftimes) * 1e-9)
}
