package fitanalytic

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/MikhailKlemin/uniclo.uk/pkg/client"
	"github.com/MikhailKlemin/uniclo.uk/pkg/config"
	"github.com/MikhailKlemin/uniclo.uk/pkg/rotator"
)

//mySerial has serial and proper link just to be sure it's not suspicious to uniqlo
/*
type mySerial struct {
	serial int64
	link   string
}
*/

func splitByTen(in []client.Response) [][]client.Response {
	var divided [][]client.Response
	chunkSize := 10
	for i := 0; i < len(in); i += chunkSize {
		end := i + chunkSize

		if end > len(in) {
			end = len(in)
		}
		divided = append(divided, in[i:end])
	}
	return divided
}

//CheckForFitLink checks for FitSizeLink
func CheckForFitLink(conf config.DefaultConfig) {
	//db, err := sql.Open("sqlite3", "./assets/uniclo.sqlite")
	db, err := sql.Open("sqlite3", conf.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var tempSerials []client.Response
	rows, err := db.Query("select serial, link from products_data where fit_link is null")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		//var serial int64
		var ms client.Response
		if err := rows.Scan(&ms.Serial, &ms.Link); err != nil {
			log.Fatal(err)
		}
		tempSerials = append(tempSerials, ms)
	}
	rows.Close()

	serials := splitByTen(tempSerials)

	fmt.Println("total serials to check ", len(serials))
	//proxies := rotator.NewRotaingProxy("./proxy_socks_ip.txt")
	proxies := rotator.NewRotaingProxy("/media/mike/WDC4_1/Neo/proxy_socks_ip.txt")

	tasks := make(chan []client.Response, len(serials))
	results := make(chan []client.Response, len(serials))

	for w := 0; w < 10; w++ {
		go func(tasks <-chan []client.Response, results chan<- []client.Response) {
			//var serials []mySerial
			for ts := range tasks {
				//link := serial.link
				count := 0
				var resp []client.Response
				for {
					if count > 10 {
						//log.Println("Max retry reached for ", link)
						log.Println("Max retry reached for ")
						//p.ok = false
						break
					}
					resp, err = client.Check(proxies.Get(), ts)
					if err != nil {
						log.Println(err)
						count++
						continue
					}

					break
				}

				results <- resp
			}
		}(tasks, results)
	}

	for _, s := range serials {
		tasks <- s
	}
	close(tasks)

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	for a := 0; a < len(serials); a++ {
		resps := <-results
		for _, res := range resps {
			if !res.Status {
				continue
			}
			if res.Fit {
				_, err := tx.Exec("update products_data set fit_link=? where serial=?", 1, res.Serial)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				_, err := tx.Exec("update products_data set fit_link=? where serial=?", 0, res.Serial)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		if a%10 == 0 && a != 0 {
			tx.Commit()
			tx, err = db.Begin()
			if err != nil {
				log.Fatal(err)
			}
		}

	}
	tx.Commit()

}
