package scraper

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/MikhailKlemin/uniclo.uk/pkg/dbase"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/pkg/errors"
)

/*
//Response is
type Response struct {
	Height  int
	Weight  int
	Shape   int
	Chest   int
	Age     int
	Cluster string

	FitMatrixID    int64
	BestFitSize    string
	BestFitPercent int
	NextFitSize    string
	NextFitPercent int
	//	Cookies        []byte
}
*/

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

func random(min int, max int) int {
	//fmt.Printf("min %d max %d", min, max)
	return rand.Intn(max-min) + min
}

//Start is
func Start(combos []dbase.Fit, proxy string) (rs []dbase.Fit, err error) {
	if len(combos) == 0 {
		return rs, errors.New("empty combos")
	}

	//link := fmt.Sprintf("https://www.uniqlo.com/uk/en/product/%d.html", p.Serial)
	//fmt.Println(link)
	l := launcher.New().
		Set("proxy-server", "socks5://"+proxy). // add a flag, here we set a http proxy
		Headless(true).
		Set("blink-settings", "imagesEnabled=true").
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
	//page.Timeout()
	/*
	   gender='Men' and weight=87 and height=190 and shape=1 and chest=1 and age=20 and preference=1
	*/
	rs = append(rs, resp)

	if len(combos) == 1 {
		return rs, nil

	}
	for _, p := range combos[1:] {
		resp, err := SwitchWeightHeight(page, p.Height, p.Weight)
		if err != nil {
			return nil, errors.Wrap(err, "Cannot Switch !")
		}
		resp.Gender = p.Gender
		resp.Weight = p.Weight
		resp.Height = p.Height
		resp.Shape = p.Shape
		resp.Chest = p.Chest
		resp.Age = p.Age
		resp.Preference = p.Preference
		resp.Cluster = p.Cluster
		fmt.Print(".")
		rs = append(rs, resp)
	}

	return rs, nil

}

//SwitchWeightHeight is
func SwitchWeightHeight(page *rod.Page, height, width int) (resp dbase.Fit, err error) {

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
func InitWith(browser *rod.Browser, p dbase.Fit) (resp dbase.Fit, page *rod.Page, err error) {

	resp.Age = p.Age
	resp.Gender = p.Gender
	resp.Chest = p.Chest
	resp.Height = p.Height
	resp.Weight = p.Weight
	resp.Shape = p.Shape
	resp.Cluster = p.Cluster
	resp.Preference = p.Preference

	//links := dbase.
	//link := fmt.Sprintf("https://www.uniqlo.com/uk/en/product/%d.html", p.Serial)
	/*
		link := fmt.Sprintf("%s", p.ProdLink)
		fmt.Println(p.ProdLink)
	*/
	link := p.Links[random(0, len(p.Links))]

	err = rod.Try(func() {
		page = browser.MustPage(link)
	})
	if err != nil {
		fmt.Println(link)
		//log.Fatal(err)
		return resp, page, errors.Wrap(err, fmt.Sprintf("Cannot navigate to %s", link))
	}

	page.WaitLoad()
	page.WaitIdle(5 * time.Second)

	xcount := 0
	for {
		xcount++
		elems := page.Timeout(30 * time.Second).MustElements("#fitanalytics__button")
		if len(elems) == 0 {
			log.Println("Looks like no fit button")
			elems = page.MustElements("a.productTile__link")
			if len(elems) > 0 {
				err := rod.Try(func() {
					elems[0].MustClick()
				})
				if err != nil {
					return resp, page, errors.Wrap(err, fmt.Sprintf("Cannot find link to product"))
				}
				page.WaitLoad()
			}
			if xcount < 10 {
				time.Sleep(2 * time.Second)
				continue
			}
		}
		break
	}
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

	//####################################################################################################
	//Setting measurements units
	err = setMeasurements(page, p.Height, p.Weight)
	if err != nil {
		return resp, page, errors.Wrap(err, "Cannot set measurrements")
	}
	//####################################################################################################

	//####################################################################################################
	//Setting Tummy units

	header, err := getHeader(page)
	if err != nil {
		return resp, page, errors.Wrap(err, "Cannot get header")

	}

	fmt.Printf("Current header is %s\n", header)
	if header == "Your tummy shape" {
		err = setTummy(page, p.Shape)
	}
	//####################################################################################################

	//####################################################################################################
	//Setting Chges units

	header, err = getHeader(page)
	if err != nil {
		return resp, page, errors.Wrap(err, "Cannot get header")
	}
	fmt.Printf("Current header is %s\n", header)

	if header == "Your chest shape" {
		//click average shape
		err = setChest(page, p.Chest)
		if err != nil {
			return resp, page, errors.Wrap(err, "Cannot set chest")

		}
	}
	//####################################################################################################

	header, err = getHeader(page)
	if err != nil {
		return resp, page, errors.Wrap(err, "Cannot get header")
	}

	fmt.Printf("Current header is %s\n", header)

	if header == "How old are you?" {
		//click average shape
		err = setAge(page, p.Age)
		if err != nil {
			return resp, page, errors.Wrap(err, "Cannot set chest")

		}
	}

	header, err = getHeader(page)
	if err != nil {
		return resp, page, errors.Wrap(err, "Cannot get header")
	}

	fmt.Printf("Current header is %s\n", header)

	if header == "Fit preference" {
		//click average shape
		fmt.Println("[INFO] Setting Fit!\t", p.Preference)
		err = setFit(page, p.Preference)
		if err != nil {
			return resp, page, errors.Wrap(err, "Cannot set preference")

		}
	}

	time.Sleep(2 * time.Second)
	header, err = getHeader(page)
	if err != nil {
		return resp, page, errors.Wrap(err, "Cannot get header")
	}

	fmt.Printf("Current header is %s\n", header)

	if header == "What do you wear?" {
		//click average shape
		log.Println("Setting up Legs ")
		err = setLegsSize(page)
		if err != nil {
			return resp, page, errors.Wrap(err, "Cannot set chest")

		}
	}

	//fit preference
	//page.MustSearch(".uclw_noUi-base").WaitStable(2 * time.Second)
	//page.MustSearch("#uclw_save_info_button").MustClick()

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

	page.WaitIdle(5 * time.Second)
	_, err = page.Eval(`$("#uclw_close_link").click();`)
	if err != nil {
		fmt.Printf("[WARN] init with error: %#v\n", resp)
		return resp, page, err
	}
	fmt.Printf("[INFO] init OK: %#v\n", resp)

	//page.MustSearch("#fitanalytics__button > span").WaitStable(1 * time.Second)
	//size, _ := page.MustSearch("#fitanalytics__button > span").Text()

	//time.Sleep(10 * time.Second)
	return resp, page, nil //errors.New("no element for perfect size found")
}
