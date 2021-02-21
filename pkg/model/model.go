package model

import "time"

/*
import (
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/MikhailKlemin/uniclo.uk/pkg/config"
	bolt "go.etcd.io/bbolt"
)

//DB is
type DB struct {
	db *bolt.DB
}
*/

//Product is
type Product struct {
	Serial  int64
	Link    string
	Gender  string
	Bread   []string
	Images  []string
	Name    string
	Price   string
	Sizes   []Size
	Details string
	Fit     bool
	Cluster string
	Created time.Time
	Updated time.Time
}

//Size holds Size information
type Size struct {
	ID      int64
	SizeID  string `json:"ID"`
	Color   string
	Dim     string
	InStock bool
}

/*


//MatchResult per cluster
type MatchResult struct {
	Cluster string
}

func inttobyte(in interface{}) []byte {
	val, ok := in.(int64)
	if !ok {
		log.Fatal("Error converting int to byte")
	}
	bval := make([]byte, 8)
	binary.LittleEndian.PutUint64(bval, uint64(val))
	return bval

}

//Put place an object
func (db *DB) Put(bucket, key, value []byte) error {
	err := db.db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte("Uniqlo")).Bucket(bucket).Put(key, value)
		if err != nil {
			return fmt.Errorf("could not insert weight: %v", err)
		}
		return nil
	})

	return err
}

func (db *DB) clusterIDX2(p Product) error {
	err := db.db.Update(func(tx *bolt.Tx) error {
		cb := tx.Bucket([]byte("Uniqlo")).Bucket([]byte("Products")).Bucket([]byte("ClusterIDX"))
		res := cb.Get([]byte(p.Cluster))
		if res == nil {
			b, _ := json.Marshal([]int64{p.Serial})
			cb.Put([]byte(p.Cluster), b)
			return nil
		}

		var serials []int64
		err := json.Unmarshal(res, &serials)
		if err != nil {
			return err
		}
		serials = append(serials, p.Serial)
		b, _ := json.Marshal(serials)

		err = cb.Put([]byte(p.Cluster), b)
		return err
	})

	return err

}

func (db *DB) clusterIDX(tx *bolt.Tx, p Product) error {
	cb := tx.Bucket([]byte("Uniqlo")).Bucket([]byte("Products")).Bucket([]byte("ClusterIDX"))
	res := cb.Get([]byte(p.Cluster))
	if res == nil {
		b, _ := json.Marshal([]int64{p.Serial})
		cb.Put([]byte(p.Cluster), b)
		return nil
	}

	var serials []int64
	err := json.Unmarshal(res, &serials)
	if err != nil {
		return err
	}
	serials = append(serials, p.Serial)
	b, _ := json.Marshal(serials)

	err = cb.Put([]byte(p.Cluster), b)
	return err
}

//Close closing DB
func (db *DB) Close() {
	db.db.Close()
}

//Get a single value per key
func (db *DB) Get(bucket, key []byte) ([]byte, error) {
	var result []byte
	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Uniqlo")).Bucket(bucket).Get(key)
		if b != nil {
			result = append(result, b...)
		} else {
			return errors.New("no key")
		}
		return nil
	})
	return result, err
}

//NewDB creates new instance
func NewDB(path string) (*DB, error) {
	var db DB
	var err error
	db.db, err = bolt.Open(path, 0666, nil)
	if err != nil {
		return &db, err
	}

	err = db.db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("Uniqlo"))
		if err != nil {
			return fmt.Errorf("could not create root bucket: %v", err)
		}
		_, err = root.CreateBucketIfNotExists([]byte("Products"))
		if err != nil {
			return fmt.Errorf("could not create Products bucket: %v", err)
		}

		return nil
	})

	return &db, err

}


//Test is
func Test(conf config.DefaultConfig) {
	xdb, err := NewDB(conf.BDB)
	if err != nil {
		log.Fatal(err)
	}

	defer xdb.Close()
	/*serial := inttobyte(433515)

	res, err := xdb.Get([]byte("Products"), serial)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", res)

}

//Transfer from sqlite to bolt for fun
func Transfer(conf config.DefaultConfig) {

	xdb, err := NewDB(conf.BDB)
	if err != nil {
		log.Fatal(err)
	}

	defer xdb.Close()

	db, err := sql.Open("sqlite3", conf.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	q := `SELECT serial, gender, link, parsed, fit_link, cluster_id FROM products_data;`

	rows, err := db.Query(q)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var serialI int64
		var gender string
		var link string
		var fit int
		var cluster string
		var parsed []byte
		var p Product
		rows.Scan(&serialI, &gender, &link, &parsed, &fit, &cluster)
		if err := json.Unmarshal(parsed, &p); err != nil {
			log.Fatal(err)
		}
		if fit == 1 {
			p.Fit = true
		}
		p.Link = link
		p.Cluster = cluster
		p.Serial = serialI
		p.Created = time.Now()
		p.Updated = time.Now()

		b, _ := json.Marshal(p)
		//serial := make([]byte, 8)
		//binary.LittleEndian.PutUint64(serial, uint64(serialI))
		serial := inttobyte(serialI)
		xdb.Put([]byte("Products"), serial, b)
	}

}
*/
