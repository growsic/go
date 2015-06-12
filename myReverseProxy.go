//ReverseProxy with cache and Round-robin
//author: growsic
//2015/6/12

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
)

//this constructor implements RoundTrip method, which is declared in RoundTripper interface
type myTransport struct {
	//this map store cache, whose id is request.URL.Host
	cache map[string][]byte
}

//implements RoundTrip method, which is called automatically in ListenAndServe method
func (t *myTransport) RoundTrip(request *http.Request) (*http.Response, error) {

	//check if the request is in cache map
	if bin, ok := t.cache[request.URL.Host]; ok {
		//initialize
		response := &http.Response{}
		//put response in cache map to response body
		response.Body = ioutil.NopCloser(strings.NewReader(string(bin)))

		fmt.Print("cache called")
		return response, nil
	}

	//if couldn't find cache, send request and get response
	response, err := http.DefaultTransport.RoundTrip(request)

	fmt.Println(response, err)
	//get binary of response
	bin, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return nil, err
	}

	//store binary in cache map
	t.cache[request.URL.Host] = bin
	fmt.Println("key", request.URL)
	fmt.Print("value", bin)
	fmt.Println("cache stored!")
	fmt.Println("map length:", len(t.cache))

	//reinsert response body to response because ReadAll method do something to response
	response.Body = ioutil.NopCloser(strings.NewReader(string(bin)))
	return response, err
}

//return ports by rotation
func returnDerection() func() string {

	i := -1

	ports := []string{
		":9090",
		":9191",
		":9292",
	}

	return func() string {
		i++
		fmt.Println(i % len(ports))
		return ports[i%len(ports)]
	}
}

func main() {
	sourceAddress := ":3000"       //port of ReverseProxy
	derection := returnDerection() //port which ReverseProxy send requst
	ts := &myTransport{make(map[string][]byte)}

	mutex := sync.Mutex{}

	//set derection to director
	director := func(request *http.Request) {
		request.URL.Scheme = "http"
		request.URL.Host = derection()
	}

	//initialize ReverseProxy
	//put args to tell direction and RoundTrip method to ReverseProxy
	proxy := &httputil.ReverseProxy{Director: director, Transport: ts}
	server := http.Server{
		Addr:    sourceAddress,
		Handler: proxy,
	}
	fmt.Println("before listenandserve")
	//launch ReverseProxy server
	server.ListenAndServe()
}
