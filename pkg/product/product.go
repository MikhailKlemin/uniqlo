package product

import (
	"bytes"
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//Item holds all information
type Item struct {
	URL     string
	Gender  string
	Bread   []string
	Images  []string
	Name    string
	Price   string
	Sizes   []Size
	Details string
}

//Size holds Size information
type Size struct {
	ID      string
	Color   string
	Dim     string
	InStock bool
}

//SizeDetail Holds size details
type SizeDetail struct {
	ID         string `json:"id"`
	Attributes struct {
		Color string `json:"color"`
		Size  string `json:"size"`
	} `json:"attributes"`
	Availability struct {
		Status           string `json:"status"`
		StatusQuantity   string `json:"statusQuantity"`
		InStock          bool   `json:"inStock"`
		Ats              string `json:"ats"`
		InStockDate      string `json:"inStockDate"`
		AvailableForSale bool   `json:"availableForSale"`
		//		PurchaseLevel    string `json:"purchaseLevel"`
		Levels struct {
			INSTOCK      int `json:"IN_STOCK"`
			PREORDER     int `json:"PREORDER"`
			BACKORDER    int `json:"BACKORDER"`
			NOTAVAILABLE int `json:"NOT_AVAILABLE"`
		} `json:"levels"`
		IsAvailable  bool   `json:"isAvailable"`
		InStockMsg   string `json:"inStockMsg"`
		PreOrderMsg  string `json:"preOrderMsg"`
		BackOrderMsg string `json:"backOrderMsg"`
	} `json:"availability"`
	Pricing struct {
		ShowStandardPrice bool    `json:"showStandardPrice"`
		IsPromoPrice      bool    `json:"isPromoPrice"`
		Standard          float64 `json:"standard"`
		FormattedStandard string  `json:"formattedStandard"`
		Sale              float64 `json:"sale"`
		FormattedSale     string  `json:"formattedSale"`
		SalePriceMoney    struct {
		} `json:"salePriceMoney"`
		StandardPriceMoney struct {
		} `json:"standardPriceMoney"`
		PricePercentage string `json:"pricePercentage"`
		Quantities      []struct {
			Unit  string `json:"unit"`
			Value int    `json:"value"`
		} `json:"quantities"`
	} `json:"pricing"`
	Applicablebadges []struct {
		ID      string `json:"id"`
		Value   string `json:"value"`
		CoValue string `json:"coValue"`
		Class   string `json:"class"`
	} `json:"applicablebadges"`
}

//Parse parsing product page
func Parse(b []byte) (m []byte, g string) {
	/*	b, err := ioutil.ReadFile("assets/432025.html")
		if err != nil {
			log.Fatal(err)
		}
	*/
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(b))

	/*
		ecommerceProductObject = {"currencyCode":"GBP","detail":{"actionField":{"list":"direct access"},"products":[{"id":"429159","dimension3":"","metric7":"0","name":"Sweatshirt","price":19.9}]}};
	*/
	var p Item
	if m := regexp.MustCompile(`"name":"(.*?)"`).FindSubmatch(b); len(m) > 0 {
		p.Name = string(m[1])
	}

	if m := regexp.MustCompile(`"price":([\d\.]+)`).FindSubmatch(b); len(m) > 0 {
		p.Price = string(m[1])
	}

	//Parse sizes:
	if m := regexp.MustCompile(`(?s)var\s*pdpVariationsJSON\s*=\s*(.*?)\s*;`).FindSubmatch(b); len(m) > 0 {
		siz := make(map[string]SizeDetail)
		xb := m[1]
		if err := json.Unmarshal(xb, &siz); err != nil {
			log.Fatal(err)
		}
		for kk, key := range siz {
			//fmt.Printf("%s\t%s:%s\t%s\n", kk, key.Attributes.Color, key.Attributes.Size, key.Availability.Status)
			size := Size{
				ID:    kk,
				Color: key.Attributes.Color,
				Dim:   key.Attributes.Size,
			}
			if key.Availability.Status == "IN_STOCK" {
				size.InStock = true
			}
			p.Sizes = append(p.Sizes, size)
		}
	}

	//Images

	doc.Find(`img.pdp__verticalSliderImg`).Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr(`src`)
		index := strings.Index(href, "?")
		href = href[:index]
		p.Images = append(p.Images, href)

	})

	//BreadCrumbNavigation
	doc.Find(`.breadCrumb__link`).Each(func(_ int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		text = strings.ToLower(text)
		text = strings.Title(text)
		p.Bread = append(p.Bread, text)
	})
	if len(p.Bread) > 0 {
		p.Gender = p.Bread[0]
	}

	p.Details = doc.Find(`p.deliverySection__text`).Text()
	p.Details = strings.ReplaceAll(p.Details, "\n", " ")
	p.Details = regexp.MustCompile(`\s+`).ReplaceAllString(p.Details, " ")
	p.Details = strings.TrimSpace(p.Details)

	//fmt.Printf("%#v\n", p.Bread)
	b, _ = json.MarshalIndent(p, "", "    ")
	//fmt.Println(string(b))
	return b, p.Gender
}
