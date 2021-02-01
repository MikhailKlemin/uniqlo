package crawler

import (
	"bytes"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/MikhailKlemin/uniclo.uk/pkg/analyzer"
	"github.com/MikhailKlemin/uniclo.uk/pkg/config"
	"github.com/MikhailKlemin/uniclo.uk/pkg/product"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
)

var serialRe = regexp.MustCompile(`\d{6}`)

func cleanLink(link string) (out string) {
	index := strings.Index(link, "?")
	if index > 0 {
		link = link[:index]
	}
	return link
}

func getSerial(link string) string {
	//https://www.uniqlo.com/uk/en/size/429159_size.html
	//o := strings.TrimPrefix(link, "https://www.uniqlo.com/uk/en/size/")
	if m := serialRe.FindStringSubmatch(link); len(m) > 0 {
		return m[0]
	}
	return link

}

//CollectProducts collect src pages of producs
func CollectProducts(conf config.DefaultConfig) {

	//db2, err := dbase.NewDB("./asstes/meta.data.db")
	type pair struct {
		value []byte
		link  string
	}

	ch := make(chan pair, 10)
	done := make(chan bool)

	go func(ch <-chan pair) {

		db, err := sql.Open("sqlite3", conf.DB)
		if err != nil {
			log.Fatal(err)
		}

		//db := data.GetDB()

		//defer db.Close()

		um := make(map[string]bool)

		rows, err := db.Query("select serial from products_data")
		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			var serial string
			if err := rows.Scan(&serial); err != nil {
				log.Fatal(err)
			}
			um[serial] = true
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		c := 0
		for {
			c++
			j, more := <-ch
			if more {
				//fmt.Println("received job")
				//fmt.Println("LEN:\t", len(j))
				serial := getSerial(j.link)
				if _, ok := um[serial]; !ok {
					fmt.Println("SERIAL:\t", serial)
					um[serial] = true
					data, gender := product.Parse(j.value)
					_, err = tx.Exec("insert into products_data (serial, link, gender, parsed, updated) values(?,?,?,?,?)",
						serial,
						j.link,
						gender,
						data,
						time.Now())
					if err != nil {
						log.Panicln(err)
					}
				} else {
					_, err = tx.Exec("update products_data set updated = ?, link = ? where serial =?  ", time.Now(), j.link, serial)
					if err != nil {
						log.Panicln(err)
					}
				}
				if c%100 == 0 && c != 0 {
					tx.Commit()
					tx, err = db.Begin()
					if err != nil {
						log.Fatal(err)
					}

				}
			} else {
				fmt.Println("received all jobs")
				done <- true
				//return
				break
			}
		}
		tx.Commit()

	}(ch)

	c := colly.NewCollector(
		colly.AllowedDomains("www.uniqlo.com", "uniqlo.com"),
		colly.CacheDir("/media/mike/WDC4_1/uniclo-cache"),
		colly.URLFilters(
			regexp.MustCompile("https://www.uniqlo.com/uk/en/.*$"),
		),
		colly.UserAgent("Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:83.0) Gecko/20100101 Firefox/83.0"),
		colly.Async(true),
	)

	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 20,
	})
	// Rotate two socks5 proxies
	rp, err := proxy.RoundRobinProxySwitcher(loadProxy()...)
	if err != nil {
		log.Fatal(err)
	}
	c.SetProxyFunc(rp)

	// Print the response
	c.OnResponse(func(r *colly.Response) {
		//log.Printf("Link: %s; Proxy Address: %s\n", r.Request.URL.String(), r.Request.ProxyURL)
		/*	if strings.HasPrefix(r.Request.URL.String(), "https://www.uniqlo.com/uk/en/size/") {
				err := db.PutByte(r.Request.URL.String(), r.Body)
				if err != nil {
					log.Fatal(err)
				}
			}
			//log.Printf("%s\n", bytes.Replace(r.Body, []byte("\n"), nil, -1))
			//fmt.Println(r.Request.URL.String())
		*/

		if strings.HasPrefix(r.Request.URL.String(), "https://www.uniqlo.com/uk/en/product/") {
			ch <- pair{link: r.Request.URL.String(), value: r.Body}
		}
	})

	c.OnHTML(`a[data-seoproducturl^="https://www.uniqlo.com/uk/en/product/"]`, func(e *colly.HTMLElement) {
		link := e.Attr("data-seoproducturl")
		// Print link
		//log.Printf("!!!!Link found: %q -> %s\n", e.Text, link)
		//os.Exit(1)
		e.Request.Visit(link)
	})

	c.OnHTML(`a[href^="https://www.uniqlo.com/uk/en/size/"]`, func(e *colly.HTMLElement) {
		link := e.Attr(`href`)
		//log.Printf("!!!!Link found: %q -> %s @ %s\n", e.Text, link, e.Request.URL.String())
		e.Request.Visit(link)
	})

	c.OnXML("//urlset/url/loc", func(e *colly.XMLElement) {
		courseURL := e.Text
		e.Request.Visit(courseURL)
	})
	// Fetch httpbin.org/ip five times
	c.Visit("https://www.uniqlo.com/uk/en/sitemap_mobile.xml")
	c.Wait()
	fmt.Println("Channel size:\t", len(ch))
	time.Sleep(10 * time.Second)
	fmt.Println("Channel size:\t", len(ch))

	close(ch)
	<-done
}

//CollectSizes collects sizes with SQLite backup
func CollectSizes(conf config.DefaultConfig) {

	//db2, err := dbase.NewDB("./asstes/meta.data.db")
	type pair struct {
		value []byte
		link  string
	}

	ch := make(chan pair, 10)
	done := make(chan bool)

	go func(ch <-chan pair) {
		db, err := sql.Open("sqlite3", conf.DB)
		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()

		um := make(map[string]bool)

		rows, err := db.Query("select serial from products_data where cluster_id is not null")
		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			var serial string
			if err := rows.Scan(&serial); err != nil {
				log.Fatal(err)
			}
			um[serial] = true
		}
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		c := 0
		for {
			c++
			j, more := <-ch
			if more {
				//fmt.Println("received job")
				//fmt.Println("LEN:\t", len(j))
				serial := getSerial(j.link)
				if _, ok := um[serial]; !ok {
					um[serial] = true
					tokenized, err := analyzer.Tokenize(bytes.NewReader(j.value))
					if err != nil {
						log.Fatal(err)
					}

					sha1 := sha1.Sum([]byte(tokenized))
					//fmt.Println("Serial:", serial)
					fmt.Printf("update products_data set cluster_id=%s where serial=%s\n", fmt.Sprintf("%x", sha1), serial)
					_, err = tx.Exec("update products_data set cluster_id=? where serial=? ",
						fmt.Sprintf("%x", sha1),
						serial)

					if err != nil {
						log.Panicln(err)
					}

				}
				if c%100 == 0 && c != 0 {
					tx.Commit()
					tx, err = db.Begin()
					if err != nil {
						log.Fatal(err)
					}

				}
			} else {
				fmt.Println("received all jobs")
				done <- true
				//return
				break
			}
		}
		tx.Commit()

	}(ch)

	c := colly.NewCollector(
		colly.AllowedDomains("www.uniqlo.com", "uniqlo.com"),
		colly.CacheDir("/media/mike/WDC4_1/uniclo-cache"),
		colly.URLFilters(
			regexp.MustCompile("https://www.uniqlo.com/uk/en/.*$"),
		),
		colly.UserAgent("Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:83.0) Gecko/20100101 Firefox/83.0"),
		colly.Async(true),
	)

	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 15,
	})
	// Rotate two socks5 proxies
	rp, err := proxy.RoundRobinProxySwitcher(loadProxy()...)
	if err != nil {
		log.Fatal(err)
	}
	c.SetProxyFunc(rp)

	// Print the response
	c.OnResponse(func(r *colly.Response) {
		//log.Printf("Link: %s; Proxy Address: %s\n", r.Request.URL.String(), r.Request.ProxyURL)
		/*	if strings.HasPrefix(r.Request.URL.String(), "https://www.uniqlo.com/uk/en/size/") {
				err := db.PutByte(r.Request.URL.String(), r.Body)
				if err != nil {
					log.Fatal(err)
				}
			}
			//log.Printf("%s\n", bytes.Replace(r.Body, []byte("\n"), nil, -1))
			//fmt.Println(r.Request.URL.String())
		*/

		if strings.HasPrefix(r.Request.URL.String(), "https://www.uniqlo.com/uk/en/size/") {
			ch <- pair{link: r.Request.URL.String(), value: r.Body}
		}
	})

	c.OnHTML(`a[data-seoproducturl^="https://www.uniqlo.com/uk/en/product/"]`, func(e *colly.HTMLElement) {
		link := e.Attr("data-seoproducturl")
		// Print link
		//log.Printf("!!!!Link found: %q -> %s\n", e.Text, link)
		//os.Exit(1)
		e.Request.Visit(link)
	})

	c.OnHTML(`a[href^="https://www.uniqlo.com/uk/en/size/"]`, func(e *colly.HTMLElement) {
		link := e.Attr(`href`)
		//log.Printf("!!!!Link found: %q -> %s @ %s\n", e.Text, link, e.Request.URL.String())
		e.Request.Visit(link)
	})

	c.OnXML("//urlset/url/loc", func(e *colly.XMLElement) {
		courseURL := e.Text
		e.Request.Visit(courseURL)
	})
	// Fetch httpbin.org/ip five times
	c.Visit("https://www.uniqlo.com/uk/en/sitemap_mobile.xml")
	c.Wait()
	fmt.Println("Channel size:\t", len(ch))
	time.Sleep(10 * time.Second)
	fmt.Println("Channel size:\t", len(ch))

	close(ch)
	<-done
}

func loadProxy() []string {
	b, err := ioutil.ReadFile("/media/mike/WDC4_1/Neo/proxy_socks_ip.txt")
	if err != nil {
		log.Fatal(err)
	}

	ps := strings.Split(string(b), "\n")

	var proxies []string
	for _, p := range ps {
		if p != "" {
			proxies = append(proxies, "socks5://"+strings.TrimSpace(p))
		}
	}

	return proxies

}
