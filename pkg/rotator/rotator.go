package rotator

import (
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"
)

//Proxies is
type Proxies struct {
	counter int
	proxies []string
	mu      sync.Mutex
}

//NewRotaingProxy is
func NewRotaingProxy(path string) *Proxies {
	var ps Proxies
	ps.counter = 0
	ps.proxies = loadProxy(path)
	shuffle(ps.proxies)
	return &ps
}

//Get is
func (ps *Proxies) Get() string {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if ps.counter >= len(ps.proxies)-1 {
		ps.counter = 0
	} else {
		ps.counter++
	}

	return ps.proxies[ps.counter]
}

func loadProxy(path string) []string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	ps := strings.Split(string(b), "\n")

	var proxies []string
	for _, p := range ps {
		if p != "" {
			proxies = append(proxies, strings.TrimSpace(p))
		}
	}

	return proxies

}

func shuffle(a []string) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
}
