package sizer

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/MikhailKlemin/uniclo.uk/pkg/config"
	"github.com/MikhailKlemin/uniclo.uk/pkg/dbase"
	"github.com/MikhailKlemin/uniclo.uk/pkg/rotator"
	"github.com/MikhailKlemin/uniclo.uk/pkg/sizer/generator2"
	"github.com/MikhailKlemin/uniclo.uk/pkg/sizer2/scraper"
)

type myResult struct {
	resp []dbase.Fit
	err  error
}

var rot = rotator.NewRotaingProxy(`/media/mike/WDC4_1/Neo/proxy_socks_ip.txt`)

func worker2(task <-chan []dbase.Fit, results chan<- myResult) {
	for cs := range task {
		//fmt.Println(cs[0].Height)
		temp := make(chan []dbase.Fit)
		go func(cs []dbase.Fit) {
			var resp []dbase.Fit
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
			results <- myResult{[]dbase.Fit{}, errors.New("time out")}
		}
	}

}

//Start starts
func Start(conf config.DefaultConfig) {
	mdb, err := dbase.NewDB(conf)
	if err != nil {
		log.Fatal(err)
	}

	xclusters, err := mdb.GetClusters()
	if err != nil {
		log.Fatal(err)
	}

	defer mdb.Close()
	//mdb.Close()
	xclusters = xclusters[:10]

	for i, c := range xclusters {
		if i > 0 {
			//break
		}
		fmt.Printf("[INFO] Starting %s\n", c)
		xStart2(c, mdb)
	}
}

//xStart2 fresh start
func xStart2(c string, mdb *dbase.DB) {
	opts := generator2.Generate(generator2.Opts{Age: 20, Gender: "Men", Shape: 1, Preference: 4, Chest: 1})

	count := 0
	var mfits []dbase.Fit
	var mfitsbyfive [][]dbase.Fit
	var err error

	for _, o := range opts {
		var f dbase.Fit
		f.Age = o.Age
		f.Chest = o.Chest
		f.Gender = o.Gender
		f.Height = o.Height
		f.Preference = o.Preference
		f.Shape = o.Shape
		f.Weight = o.Weight
		f.Cluster = c
		f.Links, err = mdb.GetLinksToProducstPerCluster(c)
		if err != nil {
			log.Fatal(err)
		}
		exists, err := mdb.FitExists(f)
		if err != nil {
			log.Fatal(err)
		}
		if !exists {
			//fmt.Printf("[INFO] h: %d, w: %d, c: %s not exists\n", f.Height, f.Weight, f.Cluster)

			/*
				fmt.Printf("select id from fit where cluster_id='%s' and "+
					"gender='%s' and weight=%d and height=%d and shape=%d and chest=%d and age=%d and preference=%d\n",
					f.Cluster, f.Gender, f.Weight, f.Height, f.Shape, f.Chest, f.Age, f.Preference,
				)
			*/

			mfits = append(mfits, f)
			count++
			if len(mfits) > 20 {
				mfitsbyfive = append(mfitsbyfive, mfits)
				mfits = []dbase.Fit{}
			}
		}
	}

	if len(mfits) > 0 {
		mfitsbyfive = append(mfitsbyfive, mfits)
	}

	task := make(chan []dbase.Fit, len(mfitsbyfive))
	results := make(chan myResult, len(mfitsbyfive))

	for w := 0; w < 6; w++ {
		go worker2(task, results)
	}

	for _, mfbf := range mfitsbyfive {

		task <- mfbf
	}
	close(task)

	for i := 0; i < len(mfitsbyfive); i++ {
		resp := <-results
		if resp.err != nil {
			log.Println(resp.err)
			continue
		}
		for _, r := range resp.resp {
			/*_, err := db.Exec(q, r.FitMatrixID, r.BestFitSize, r.BestFitPercent, r.NextFitSize, r.NextFitPercent, cluster)
			if err != nil {
				log.Fatal(err)
			}*/
			mdb.InsertFit(r)

		}

	}

}

/*
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

	xclusters = xclusters[:30]

	count := 0
	var mfits []dbase.Fit
	var mfitsbyfive [][]dbase.Fit
	for _, c := range xclusters {
		for _, o := range opts {
			var f dbase.Fit
			f.Age = o.Age
			f.Chest = o.Chest
			f.Gender = o.Gender
			f.Height = o.Height
			f.Preference = o.Preference
			f.Shape = o.Shape
			f.Weight = o.Weight
			f.Cluster = c
			f.Links, err = mdb.GetLinksToProducstPerCluster(c)
			if err != nil {
				log.Fatal(err)
			}
			exists, err := mdb.FitExists(f)
			if err != nil {
				log.Fatal(err)
			}
			if !exists {
				//fmt.Printf("%#v\n", f)
				mfits = append(mfits, f)
				count++
				if len(mfits) > 20 {
					mfitsbyfive = append(mfitsbyfive, mfits)
					mfits = []dbase.Fit{}
				}
			}
		}
	}
	mfitsbyfive = append(mfitsbyfive, mfits)

	task := make(chan []dbase.Fit, len(mfitsbyfive))
	results := make(chan myResult, len(mfitsbyfive))

	for w := 0; w < 1; w++ {
		go worker2(task, results)
	}

	for _, mfbf := range mfitsbyfive {

		task <- mfbf
	}
	close(task)

	for i := 0; i < len(mfitsbyfive); i++ {
		resp := <-results
		if resp.err != nil {
			log.Println(resp.err)
			continue
		}
		for _, r := range resp.resp {
			mdb.InsertFit(r)

		}

	}

}
*/
