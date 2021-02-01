package client

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/pkg/errors"
)

//Response has
type Response struct {
	Fit    bool
	Link   string
	Serial int64
	Status bool
	//Serial int64
}

//Check is fun
func Check(p string, ts []Response) ([]Response, error) {
	fmt.Println("New Launcher with proxy ", p)
	l := launcher.New().
		Set("proxy-server", "socks5://"+p). // add a flag, here we set a http proxy
		Headless(true).
		Set("blink-settings", "imagesEnabled=false").
		Devtools(false)

	defer l.Cleanup() // remove user-data-dir

	url := l.MustLaunch()

	browser := rod.New().
		ControlURL(url).
		Trace(true).
		SlowMotion(1 * time.Second).
		MustConnect()

	// auth the proxy
	// here we use cli tool "mitmproxy --proxyauth user:pass" as an example
	defer browser.Close()
	var rs []Response
	var page *rod.Page

	for i, link := range ts {
		if i == 0 {
			err := rod.Try(func() {
				page = browser.Timeout(120 * time.Second).MustPage(link.Link)
				//page = browser.Timeout(120 * time.Second).Must

			})
			if err != nil {
				return rs, errors.Wrap(err, fmt.Sprintf("Cannot navigate to %s with %s", link.Link, p))
			}
		} else {
			err := rod.Try(func() {
				page.MustNavigate(link.Link)
				//page = browser.Timeout(120 * time.Second).Must

			})
			if err != nil {
				return rs, errors.Wrap(err, fmt.Sprintf("Cannot navigate to %s with %s", link.Link, p))
			}
		}

		err := page.WaitLoad()
		if err != nil {
			return rs, errors.Wrap(err, fmt.Sprintf("Cannot navigate to %s with %s", link.Link, p))
		}

		//fmt.Println(page.MustElement("#fitanalytics__button").) // print the size of the image
		elems, err := page.Elements("#fitanalytics__button")
		if err != nil {
			return rs, errors.Wrap(err, fmt.Sprintf("Cannot find elements at %s with %s", link.Link, p))
		}

		//fmt.Println(len(elems))
		//utils.Pause() // pause goroutine
		if len(elems) > 0 {
			//return r, nil
			link.Fit = true
			link.Status = true
			fmt.Printf("[INFO] Serial: %d is %t\n", link.Serial, link.Fit)
			rs = append(rs, link)
		} else {
			link.Fit = false
			link.Status = true
			fmt.Printf("[INFO] Serial: %d is %t\n", link.Serial, link.Fit)
			rs = append(rs, link)
		}

	}
	//page := browser.MustPage(link)

	return rs, nil

}
