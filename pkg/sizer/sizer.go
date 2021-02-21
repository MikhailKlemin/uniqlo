package sizer

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	sql "github.com/jmoiron/sqlx"

	"github.com/MikhailKlemin/uniclo.uk/pkg/config"
	"github.com/MikhailKlemin/uniclo.uk/pkg/dbase"
	"github.com/MikhailKlemin/uniclo.uk/pkg/rotator"
	"github.com/MikhailKlemin/uniclo.uk/pkg/sizer/generator"
	"github.com/MikhailKlemin/uniclo.uk/pkg/sizer/generator2"
	"github.com/MikhailKlemin/uniclo.uk/pkg/sizer/scraper"
)

//SelectUndone selects unporcessed combinations
func SelectUndone(conf config.DefaultConfig, gender, shape, chest, age int, cluster string) (map[int][]generator.Combo, []int) {
	db, err := sql.Open("sqlite3", conf.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var combos = make(map[int][]generator.Combo)
	q := `SELECT id, gender, weight, height, shape, chest, age, preference FROM fit_matrix WHERE gender = ? AND shape = ? AND chest = ? AND age = ? AND preference = 1 AND id NOT IN ( SELECT fit_matrix_id FROM fit_results where cluster_id=? ) order by height,weight;`
	rows, err := db.Queryx(q, gender, shape, chest, age, cluster)
	if err != nil {
		log.Fatal(err)
	}

	var keys []int
	for rows.Next() {
		var combo generator.Combo
		err := rows.StructScan(&combo)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("%#v\n", combo)
		if _, ok := combos[combo.Height]; !ok {
			keys = append(keys, combo.Height)
		}
		combos[combo.Height] = append(combos[combo.Height], combo)

	}
	//fmt.Printf("%#v\n", combos[140])
	sort.Ints(keys)
	return combos, keys
}

func random(min int, max int) int {
	//fmt.Printf("min %d max %d", min, max)
	return rand.Intn(max-min) + min
}

type myResult struct {
	resp []scraper.Response
	err  error
}

var rot = rotator.NewRotaingProxy(`/media/mike/WDC4_1/Neo/proxy_socks_ip.txt`)

func worker(task <-chan []generator.Combo, results chan<- myResult) {
	for cs := range task {
		fmt.Println(cs[0].Height)
		temp := make(chan []scraper.Response)
		go func(cs []generator.Combo) {
			var resp []scraper.Response
			var err error
			count := 0
			for {
				resp, err = scraper.Start(cs, rot.Get())
				if err != nil && count < 10 {
					log.Println(err)
					count++
					time.Sleep(2 * time.Second)
					continue
				}
				break
			}
			temp <- resp
		}(cs)
		select {
		case resp := <-temp:
			results <- myResult{resp, nil}
		case <-time.After(2 * time.Hour):
			results <- myResult{[]scraper.Response{}, errors.New("time out")}
		}
	}

}

//Start2 fresh start
func Start2(conf config.DefaultConfig) {

	mdb, err := dbase.NewDB(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer mdb.Close()

	opts := generator2.Generate(generator2.Opts{Age: 20, Gender: "Men", Shape: 1, Preference: 1, Chest: 1})
	xclusters, err := mdb.GetClusters()
	if err != nil {
		log.Fatal(err)
	}

	xclusters = xclusters[:40]

	count := 0
	var mfits []dbase.Fit
	var mfitsbyfive [][]dbase.Fit

	for _, o := range opts {
		for _, c := range xclusters {
			var f dbase.Fit
			f.Age = o.Age
			f.Chest = o.Chest
			f.Gender = o.Gender
			f.Height = o.Height
			f.Preference = o.Preference
			f.Shape = o.Shape
			f.Weight = o.Weight
			f.Cluster = c
			exists, err := mdb.FitExists(f)
			if err != nil {
				log.Fatal(err)
			}
			if !exists {
				fmt.Printf("%#v\n", f)
				mfits = append(mfits, f)
				count++
				if len(mfits) > 5 {
					mfitsbyfive = append(mfitsbyfive, mfits)
					mfits = []dbase.Fit{}
				}
			}
		}
	}
	mfitsbyfive = append(mfitsbyfive, mfits)

	//fmt.Printf("%#v\n", mfitsbyfive[0])
	for i, mfbf := range mfitsbyfive {
		//fmt.Println(mfbf.)
		for _, mf := range mfbf {
			fmt.Println(i, "\t", mf.Height, "\t", mf.Weight)
		}
	}

}

//Start starts
func Start(conf config.DefaultConfig, cluster string, gender, shape, chest, age int) {
	rand.Seed(time.Now().UnixNano())

	if cluster == "" {
		log.Println("no cluster")
		return
	}
	combos, keys := SelectUndone(conf, gender, shape, chest, age, cluster)
	//serials, err := getSerialsByCluster(conf, cluster)
	productLinks, err := getProdLinksByCluster(conf, cluster)

	db, err := sql.Open("sqlite3", conf.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	q := `INSERT INTO fit_results (fit_matrix_id, best_fit_size, best_fit_percent, next_fit_size, next_fit_percent, cluster_id) VALUES ( ?, ?, ?, ?, ?, ?); `

	//fmt.Println("length:\t", len(keys))

	task := make(chan []generator.Combo, len(keys))
	results := make(chan myResult, len(keys))

	for w := 0; w < 5; w++ {
		go worker(task, results)
	}

	/*
		for _, key := range keys {
			for i := 0; i < len(combos[key]); i++ {
				combos[key][i].Serial = serials[random(0, len(serials)-1)]
			}
		}
	*/

	for _, key := range keys {
		for i := 0; i < len(combos[key]); i++ {
			if len(productLinks) > 1 {
				combos[key][i].ProdLink = productLinks[random(0, len(productLinks)-1)]
			} else {
				combos[key][i].ProdLink = productLinks[0]
			}
		}
	}

	for _, key := range keys {
		task <- combos[key]
	}

	close(task)

	for i := 0; i < len(keys); i++ {
		resp := <-results
		if resp.err != nil {
			log.Println(resp.err)
			continue
		}
		for _, r := range resp.resp {
			_, err := db.Exec(q, r.FitMatrixID, r.BestFitSize, r.BestFitPercent, r.NextFitSize, r.NextFitPercent, cluster)
			if err != nil {
				log.Fatal(err)
			}

		}

	}

}

func getSerialsByCluster(conf config.DefaultConfig, cluster string) (serials []int64, err error) {

	db, err := sql.Open("sqlite3", conf.DB)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer db.Close()
	q := `select serial from products_data where cluster_id = ?;`

	row, err := db.Queryx(q, cluster)
	if err != nil {
		//log.Fatal(err)
		return
	}

	for row.Next() {
		var serial int64
		row.Scan(&serial)
		serials = append(serials, serial)
	}

	return
}

func getProdLinksByCluster(conf config.DefaultConfig, cluster string) (links []string, err error) {

	db, err := sql.Open("sqlite3", conf.DB)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer db.Close()
	q := `select link from products_data where cluster_id = ? and updated>"2021-02-03";`

	row, err := db.Queryx(q, cluster)
	if err != nil {
		//log.Fatal(err)
		return
	}

	for row.Next() {
		var link string
		row.Scan(&link)
		links = append(links, link)
	}

	return
}

/*

//ProcessClusters by stored cookies
func ProcessClusters(conf config.DefaultConfig, cluster string) {
	undone, keys := SelectUndone(conf, 0, 1, 1, 20, cluster)
	serials, err := getSerialsByCluster(conf, cluster)
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("sqlite3", conf.DB)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer db.Close()

	for _, key := range keys {
		clusters := undone[key]
		c := clusters[0]
		q := "select cookies from view_fit_results where weight =? and height = ? and shape = ?  and chest =? and age = ? and gender =? limit 1"
		var b []byte

		err := db.QueryRow(q, c.Weight, c.Height, c.Shape, c.Chest, c.Age, c.Gender).Scan(&b)
		if err != nil {
			log.Fatal(err)
		}
		var cookies []*proto.NetworkCookie

		err = json.Unmarshal(b, &cookies)
		if err != nil {
			log.Fatal(err)
		}
		link := fmt.Sprintf("https://www.uniqlo.com/uk/en/product/%d.html", serials[0])
		scraper.ProcessWithCookies(cookies, rot.Get(), link)
		break
	}

}

*/
