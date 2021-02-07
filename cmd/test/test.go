package main

import (
	"fmt"
	"log"
	"time"

	sql "github.com/jmoiron/sqlx"

	"github.com/MikhailKlemin/uniclo.uk/pkg/config"
	"github.com/MikhailKlemin/uniclo.uk/pkg/sizer"

	_ "github.com/mattn/go-sqlite3" //comment
)

func main() {

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//client.Start()
	t := time.Now()
	//crawler.CollectSizes()
	//product.Parse()
	//crawler.CollectProducts()
	//crawler.CheckForFitLink()
	//analyzer.Analyze()

	//os.Exit(1)
	//dbase.ExportToCSV("")
	//db,err:=
	/*
		db, err := dbase.NewDB("./assets/meta.data.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		err = db.ExportToCSV("./assets/sample.csv")
		if err != nil {
			log.Fatal(err)
		}
	*/

	/*page, err := sizer.InitWith(p, "https://www.uniqlo.com/uk/en/product/men-two-way-single-breasted-coat-425428.html?dwvar_425428_color=COL69&dwvar_425428_size=SMA002", "65", "175")
	if err != nil {
		log.Println(err)
	}

	page, err = sizer.SwitchWeightHeight("65", "175", page)
	if err != nil {
		log.Println(err)
	}*/

	/*
		rotor := rotator.NewRotaingProxy("./proxy_socks_ip.txt")
		p := rotor.Get()

		err := sizer.Start(p, "https://www.uniqlo.com/uk/en/product/men-two-way-single-breasted-coat-425428.html?dwvar_425428_color=COL69&dwvar_425428_size=SMA002")
		if err != nil {
			log.Println(err)
		}

		fmt.Printf("Done in %s\n", time.Since(t))
	*/

	var conf = config.NewDefaultConfig()

	/*crawler.CollectProducts(conf)
	crawler.CollectSizes(conf)
	fitanalytic.CheckForFitLink(conf)
	os.Exit(1)
	*/
	//sizer.Generate()
	//generator.GenerateAndPopulate(conf)
	//sizer.GetSerialsCluster(conf)
	//sizer.SelectUndone(conf, 1, 1, 20)
	//sizer.Begin(conf)
	/*count := 0
	for {
		if count > 20 {
			break
		}
		sizer.Begin(conf)
		count++
		fmt.Println(count)
	}
	*/

	/*

		//sizer.SelectUndone(conf)
	*/
	//sizer.ProcessClusters(conf, "514c23a07c754f6adfd3e5c5701cf90311564100")
	var (
		gender = 0
		shape  = 1
		chest  = 1
		//age    = 30 //20, 30, 40, 50, 60, 70
	)

	var clusters []string
	db, err := sql.Open("sqlite3", conf.DB)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT cluster_id FROM products_data WHERE fit_link = 1 and gender =? GROUP BY cluster_id  ORDER BY count(cluster_id) DESC", "Men")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var cluster string
		rows.Scan(&cluster)
		clusters = append(clusters, cluster)
	}
	rows.Close()
	db.Close()

	for _, cluster := range clusters {
		fmt.Println("Starting ", cluster)
		sizer.Start(conf, cluster, gender, shape, chest, conf.Age)
	}
	fmt.Printf("Done in %s\n", time.Since(t))

}

/*
https://widget.fitanalytics.com/widget/productload?callback=jQuery34107007078046216283_1606223038749&ids%5B%5D=uniqlo-429159COL02&ids%5B%5D=uniqlo-429159COL09&ids%5B%5D=uniqlo-429159COL32&ids%5B%5D=uniqlo-429159COL69&ids%5B%5D=uniqlo-429159COL05&ids%5B%5D=uniqlo-429159COL11&ids%5B%5D=uniqlo-429159COL38&ids%5B%5D=uniqlo-429159COL45&ids%5B%5D=uniqlo-429159COL57&ids%5B%5D=uniqlo-429159COL30&shopCountry=GB&shopLanguage=en&userLanguage=en&sid=GYQwm_SvOgx3ppzwiDPbbV6mTEjOPkFW&_=1606223038750

*/
