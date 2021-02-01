package scraper

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/MikhailKlemin/uniclo.uk/pkg/sizer/generator"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/pkg/errors"
)

//Response is
type Response struct {
	Serial int64
	Height int
	Weight int
	Shape  int
	Chest  int
	Age    int

	FitMatrixID    int64
	BestFitSize    string
	BestFitPercent int
	NextFitSize    string
	NextFitPercent int
	Cookies        []byte
}

/*
//Params holds parameters for scraping
type Params struct {
	Proxy   string
	Serial  int64
	Gender  int
	Height  int
	Shape   int
	Chest   int
	Age     int
	Weights []int
}
*/

/*
//Begin is
func Begin(p Params) {
	//rotor := rotator.NewRotaingProxy("/media/mike/WDC4_1/Neo/proxy_socks_ip.txt")
	//p := rotor.Get()
	//link := fmt.Sprintf("https://www.uniqlo.com/uk/en/product/%d.html", p.Serial)
	//fmt.Println(link)
	l := launcher.New().
		Set("proxy-server", "socks5://"+p.Proxy). // add a flag, here we set a http proxy
		Headless(false).
		Set("blink-settings", "imagesEnabled=false").
		Devtools(false)

	defer l.Cleanup() // remove user-data-dir
	//l.ProfileDir("/media/mike/WDC4_1/chrome-profiles/" + p)

	url := l.MustLaunch()

	browser := rod.New().
		ControlURL(url).
		Trace(true).
		SlowMotion(1 * time.Second).
		MustConnect()

	defer browser.Close()

}
*/

//Start is
func Start(combos []generator.Combo, proxy string) (rs []Response, err error) {
	if len(combos) == 0 {
		return rs, errors.New("empty combos")
	}

	//link := fmt.Sprintf("https://www.uniqlo.com/uk/en/product/%d.html", p.Serial)
	//fmt.Println(link)
	l := launcher.New().
		Set("proxy-server", "socks5://"+proxy). // add a flag, here we set a http proxy
		Headless(false).
		Set("blink-settings", "imagesEnabled=false").
		Devtools(false)

	defer l.Cleanup() // remove user-data-dir
	//l.ProfileDir("/media/mike/WDC4_1/chrome-profiles/" + p)

	url := l.MustLaunch()

	browser := rod.New().
		ControlURL(url).
		Trace(false).
		SlowMotion(1 * time.Second).
		MustConnect()

	defer browser.Close()

	browser.Timeout(60 * time.Second)
	//first := combos[0]

	resp, page, err := InitWith(browser, combos[0])
	if err != nil {
		browser.Close()
		return nil, err
	}

	if page == nil || resp.BestFitSize == "" {
		return nil, errors.New("empty page or size")
	}

	//gc get cookies
	gc := func(browser *rod.Browser) []byte {
		cookies := browser.MustGetCookies()
		bc, err := json.Marshal(cookies)
		if err != nil {
			log.Fatal(err)
		}

		return bc
	}
	resp.Cookies = append(resp.Cookies, gc(browser)...)

	rs = append(rs, resp)

	for _, p := range combos[1:] {
		resp, err := SwitchWeightHeight(page, p.Height, p.Weight)
		if err != nil {
			return nil, errors.Wrap(err, "Cannot Switch !")
		}
		resp.Age = p.Age
		resp.Chest = p.Chest
		resp.Height = p.Height
		resp.Weight = p.Weight
		resp.Serial = p.Serial
		resp.Shape = p.Shape
		resp.FitMatrixID = int64(p.ID)
		resp.Cookies = append(resp.Cookies, gc(browser)...)
		//fmt.Printf("Done for height %d and weight %d\n", resp.Height, resp.Weight)
		fmt.Print(".")
		rs = append(rs, resp)
	}

	return rs, nil

}

//SwitchWeightHeight is
func SwitchWeightHeight(page *rod.Page, height, width int) (resp Response, err error) {

	var elem *rod.Element

	err = rod.Try(func() {
		elem = page.MustElement("#fitanalytics__button")
	})

	if err != nil {
		return resp, errors.Wrap(err, fmt.Sprintf("Cannot find element"))
	}

	err = rod.Try(func() {
		elem.MustClick()
	})

	if err != nil {
		return resp, errors.Wrap(err, fmt.Sprintf("Cannot click element"))
	}

	page.WaitLoad()

	editSel := ".uclw_edit_info"
	err = rod.Try(func() {
		elem = page.MustSearch(editSel)

	})
	if err != nil {
		return resp, errors.Wrap(err, "Cannot find edit link")

	}

	elem.MustClick()

	err = rod.Try(func() {
		elem = page.MustSearch(`div.uclw_selector_item:nth-child(2) > div:nth-child(1)`)

	})
	if err != nil {
		return resp, errors.Wrap(err, "Cannot find height link")

	}

	elem.MustClick()

	//enter new height
	//#uclw_form_height
	err = rod.Try(func() {
		elem = page.MustSearch(`#uclw_form_height`)

	})
	if err != nil {
		return resp, errors.Wrap(err, "Cannot find height input field link")

	}
	elem.SelectAllText()

	elem.MustInput(fmt.Sprintf("%d", height))

	//changing width
	//div.uclw_selector_item:nth-child(3) > div:nth-child(1)
	err = rod.Try(func() {
		elem = page.MustSearch(`div.uclw_selector_item:nth-child(3) > div:nth-child(1)`)

	})
	if err != nil {
		return resp, errors.Wrap(err, "Cannot find height link")

	}

	elem.MustClick()

	//enter new height
	//#uclw_form_height
	err = rod.Try(func() {
		elem = page.MustSearch(`#uclw_form_weight`)

	})
	if err != nil {
		return resp, errors.Wrap(err, "Cannot find height input field link")

	}
	elem.SelectAllText()

	elem.MustInput(fmt.Sprintf("%d", width))

	//save button
	//.uclw_button
	err = rod.Try(func() {
		elem = page.MustSearch(`.uclw_button`)

	})
	if err != nil {
		return resp, errors.Wrap(err, "Cannot find save button for save height changes")

	}

	//save button
	//.uclw_button
	err = rod.Try(func() {
		elem = page.MustSearch(`.uclw_button`)

	})
	if err != nil {
		return resp, errors.Wrap(err, "Cannot find save button for save height changes")

	}
	elem.MustClick()

	//***************************************************************************************
	// CLOSING AND REPONENING TO READ SIZE
	//***************************************************************************************

	page.MustSearch(`.uclw_button`).WaitStable(5 * time.Second)
	_, err = page.Eval(`$("#uclw_close_link").click();`)
	if err != nil {
		return resp, err
	}

	err = rod.Try(func() {
		elem = page.MustElement("#fitanalytics__button")
	})

	if err != nil {
		return resp, errors.Wrap(err, fmt.Sprintf("Cannot find element"))
	}

	elem.MustClick()
	page.WaitLoad()
	page.MustSearch(`#primary_label`).WaitStable(time.Second)

	//*Parsing
	resp.BestFitSize, err = page.MustSearch(`#primary_label`).Text()
	if err != nil {
		return resp, errors.Wrap(err, "Cannot find size")

	}

	resp.BestFitSize = strings.TrimSpace(resp.BestFitSize)

	fitP, err := page.MustSearch(`#primary_label`).Attribute(`aria-label`)
	if err != nil {
		return resp, errors.Wrap(err, "Cannot find size")

	}
	fit := regexp.MustCompile(`[^\d]`).ReplaceAllString(*fitP, "")

	resp.BestFitPercent, err = strconv.Atoi(fit)
	if err != nil {
		return resp, errors.Wrap(err, "Cannot convert prim percent")

	}

	//*****************NEXT

	resp.NextFitSize, err = page.MustSearch(`#secondary_label`).Text()
	if err != nil {
		return resp, errors.Wrap(err, "Cannot find size")

	}
	resp.NextFitSize = strings.TrimSpace(resp.NextFitSize)

	fitP, err = page.MustSearch(`#secondary_label`).Attribute(`aria-label`)
	if err != nil {
		return resp, errors.Wrap(err, "Cannot find size for second label")

	}

	fit = regexp.MustCompile(`[^\d]`).ReplaceAllString(*fitP, "")

	resp.NextFitPercent, err = strconv.Atoi(fit)
	if err != nil {
		return resp, errors.Wrap(err, "Cannot convert secondary percent")

	}

	_, err = page.Eval(`$("#uclw_close_link").click();`)
	if err != nil {
		return resp, err
	}

	//***************************************************************************************
	// CLOSING
	//***************************************************************************************

	_, err = page.Eval(`$("#uclw_close_link").click();`)
	if err != nil {
		return resp, err
	}

	page.MustSearch("#fitanalytics__button > span").WaitStable(1 * time.Second)
	//size, _ := page.MustSearch("#fitanalytics__button > span").Text()

	return resp, nil

}

//InitWith initialize browser with starting conditions
func InitWith(browser *rod.Browser, p generator.Combo) (resp Response, page *rod.Page, err error) {

	resp.Age = p.Age
	resp.Chest = p.Chest
	resp.Height = p.Height
	resp.Weight = p.Weight
	resp.Serial = p.Serial
	resp.Shape = p.Shape
	resp.FitMatrixID = int64(p.ID)
	//resp.

	link := fmt.Sprintf("https://www.uniqlo.com/uk/en/product/%d.html", p.Serial)
	err = rod.Try(func() {
		page = browser.MustPage(link)
	})
	if err != nil {
		return resp, page, errors.Wrap(err, fmt.Sprintf("Cannot navigate to %s", link))
	}

	err = page.WaitLoad()

	var elem *rod.Element

	err = rod.Try(func() {
		elem = page.Timeout(30 * time.Second).MustElement("#fitanalytics__button")
	})
	if err != nil {
		return resp, page, errors.Wrap(err, fmt.Sprintf("Cannot find elements at %s", link))
	}

	err = rod.Try(func() {
		elem.MustClick()
	})
	if err != nil {
		return resp, page, errors.Wrap(err, "Cannot click fit_anal link")
	}

	page.WaitLoad()

	heightSel := "#uclw_form_height"
	weightSel := "#uclw_form_weight"
	cmSel := "#uclw_height_element > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(1)"
	kgSel := "#uclw_weight_element > div:nth-child(1) > div:nth-child(1) > div:nth-child(2)"

	// setting cms
	err = rod.Try(func() {
		elem = page.MustSearch(cmSel)

	})
	if err != nil {
		return resp, page, errors.Wrap(err, fmt.Sprintf("Cannot find input text at %s", link))

	}
	elem.MustClick()

	// setting kgs
	err = rod.Try(func() {
		elem = page.MustSearch(kgSel)

	})
	if err != nil {

	}
	elem.MustClick()

	// setting height
	err = rod.Try(func() {
		elem = page.MustSearch(heightSel)

	})
	if err != nil {
		return resp, page, errors.Wrap(err, fmt.Sprintf("Cannot find input text at %s", link))

	}

	elem.MustInput(fmt.Sprintf("%d", p.Height))

	// setting weight
	err = rod.Try(func() {
		elem = page.MustSearch(weightSel)

	})
	if err != nil {
		return resp, page, errors.Wrap(err, fmt.Sprintf("Cannot find input text at %s", link))

	}
	elem.MustInput(fmt.Sprintf("%d", p.Weight))

	page.WaitLoad()

	//click next
	page.MustSearch("#uclw_save_info_button").MustClick()
	page.WaitLoad()

	//click average tummy
	page.MustSearch(fmt.Sprintf("#uclw_item_shape_%d", p.Shape)).WaitStable(2 * time.Second)
	page.MustSearch(fmt.Sprintf("#uclw_item_shape_%d", p.Shape)).MustClick()
	//page.WaitLoad()

	//click average shape
	page.MustSearch(fmt.Sprintf("#uclw_item_shape_%d", p.Chest)).WaitStable(2 * time.Second)
	page.MustSearch(fmt.Sprintf("#uclw_item_shape_%d", p.Chest)).MustClick()
	page.WaitLoad()

	//enter age
	page.MustSearch(".uclw_input_text").MustInput(fmt.Sprintf("%d", p.Age))

	//save button
	page.MustSearch("#uclw_save_info_button").MustClick()
	page.WaitLoad()

	//fit preference
	page.MustSearch(".uclw_noUi-base").WaitStable(2 * time.Second)
	page.MustSearch("#uclw_save_info_button").MustClick()

	page.MustSearch("#primary_label").WaitStable(2 * time.Second)
	//size, _ := page.MustSearch("#primary_label").Text()

	resp.BestFitSize, err = page.MustSearch(`#primary_label`).Text()
	if err != nil {
		return resp, page, errors.Wrap(err, "Cannot find size")

	}

	resp.BestFitSize = strings.TrimSpace(resp.BestFitSize)

	fitP, err := page.MustSearch(`#primary_label`).Attribute(`aria-label`)
	if err != nil {
		return resp, page, errors.Wrap(err, "Cannot find size")

	}
	fit := regexp.MustCompile(`[^\d]`).ReplaceAllString(*fitP, "")

	resp.BestFitPercent, err = strconv.Atoi(fit)
	if err != nil {
		return resp, page, errors.Wrap(err, "Cannot convert prim percent")

	}

	//*****************NEXT

	resp.NextFitSize, err = page.MustSearch(`#secondary_label`).Text()
	if err != nil {
		return resp, page, errors.Wrap(err, "Cannot find size")

	}
	resp.NextFitSize = strings.TrimSpace(resp.NextFitSize)

	fitP, err = page.MustSearch(`#secondary_label`).Attribute(`aria-label`)
	if err != nil {
		return resp, page, errors.Wrap(err, "Cannot find size for second label")

	}

	fit = regexp.MustCompile(`[^\d]`).ReplaceAllString(*fitP, "")

	resp.NextFitPercent, err = strconv.Atoi(fit)
	if err != nil {
		return resp, page, errors.Wrap(err, "Cannot convert secondary percent")

	}

	_, err = page.Eval(`$("#uclw_close_link").click();`)
	if err != nil {
		return resp, page, err
	}

	//page.MustSearch("#fitanalytics__button > span").WaitStable(1 * time.Second)
	//size, _ := page.MustSearch("#fitanalytics__button > span").Text()

	//time.Sleep(10 * time.Second)
	return resp, page, nil //errors.New("no element for perfect size found")
}

//ProcessWithCookies loads cookies up, nvagates to specific page and returns result
func ProcessWithCookies(cookies []*proto.NetworkCookie, proxy, link string) (Response, error) {
	l := launcher.New().
		Set("proxy-server", "socks5://"+proxy). // add a flag, here we set a http proxy
		Headless(false).
		Set("blink-settings", "imagesEnabled=false").
		Devtools(false)

	defer l.Cleanup() // remove user-data-dir
	//l.ProfileDir("/media/mike/WDC4_1/chrome-profiles/" + p)

	url := l.MustLaunch()

	browser := rod.New().
		ControlURL(url).
		Trace(true).
		SlowMotion(1 * time.Second).
		MustConnect()

	//browser.MustIncognito()

	browser.MustSetCookies(cookies)
	cookies2 := browser.MustGetCookies()
	fmt.Println("Total cookies\t", len(cookies2), "\t", len(cookies))
	for i, c := range cookies {
		fmt.Println(i, "\t", c.Name)
	}

	for i, c := range cookies2 {
		fmt.Println(i, "\t", c.Name)
	}

	defer browser.Close()

	page := browser.MustPage(link)
	page.WaitLoad()
	page.WaitIdle(30 * time.Second)
	var resp Response

	var elem *rod.Element
	var err error

	err = rod.Try(func() {
		elem = page.MustElement("#fitanalytics__button")
	})

	if err != nil {
		return resp, errors.Wrap(err, fmt.Sprintf("Cannot find element"))
	}

	elem.MustClick()
	page.WaitLoad()
	page.MustSearch(`#primary_label`).WaitStable(time.Second)

	//*Parsing
	resp.BestFitSize, err = page.MustSearch(`#primary_label`).Text()
	if err != nil {
		return resp, errors.Wrap(err, "Cannot find size")

	}

	resp.BestFitSize = strings.TrimSpace(resp.BestFitSize)

	fitP, err := page.MustSearch(`#primary_label`).Attribute(`aria-label`)
	if err != nil {
		return resp, errors.Wrap(err, "Cannot find size")

	}
	fit := regexp.MustCompile(`[^\d]`).ReplaceAllString(*fitP, "")

	resp.BestFitPercent, err = strconv.Atoi(fit)
	if err != nil {
		return resp, errors.Wrap(err, "Cannot convert prim percent")

	}

	//*****************NEXT

	resp.NextFitSize, err = page.MustSearch(`#secondary_label`).Text()
	if err != nil {
		return resp, errors.Wrap(err, "Cannot find size")

	}
	resp.NextFitSize = strings.TrimSpace(resp.NextFitSize)

	fitP, err = page.MustSearch(`#secondary_label`).Attribute(`aria-label`)
	if err != nil {
		return resp, errors.Wrap(err, "Cannot find size for second label")

	}

	fit = regexp.MustCompile(`[^\d]`).ReplaceAllString(*fitP, "")

	resp.NextFitPercent, err = strconv.Atoi(fit)
	if err != nil {
		return resp, errors.Wrap(err, "Cannot convert secondary percent")

	}

	_, err = page.Eval(`$("#uclw_close_link").click();`)
	if err != nil {
		return resp, err
	}

	//***************************************************************************************
	// CLOSING
	//***************************************************************************************

	_, err = page.Eval(`$("#uclw_close_link").click();`)
	if err != nil {
		return resp, err
	}

	page.MustSearch("#fitanalytics__button > span").WaitStable(1 * time.Second)
	//size, _ := page.MustSearch("#fitanalytics__button > span").Text()
	return resp, nil
}
