package dbase

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/MikhailKlemin/uniclo.uk/pkg/config"
)

//Product is
type Product struct {
	Serial  int64
	Cluster string
	Link    string
	Gender  string
	Bread   []string
	Images  []string
	Name    string
	Price   string
	Sizes   []Size
	Details string
	Fit     bool
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

//Fit is
type Fit struct {
	Cluster        string
	Links          []string
	Gender         string //Men
	Weight         int
	Height         int
	Shape          int
	Chest          int
	Age            int
	Preference     int
	BestFitSize    string
	BestFitPercent int
	NextFitSize    string
	NextFitPercent int
}

//DB is
type DB struct {
	db *sql.DB
}

//NewDB creates new instance
func NewDB(conf config.DefaultConfig) (*DB, error) {
	var db DB
	var err error
	db.db, err = sql.Open("sqlite3", conf.DB)
	if err != nil {
		log.Fatal(err)
	}
	return &db, err

}

//Close closes DB
func (db *DB) Close() {
	db.db.Close()
}

//FitExists -- check if the fit already exists
func (db *DB) FitExists(f Fit) (bool, error) {
	q := "select id from fit where cluster_id=? and gender=? and weight=? and height=? and shape=? and chest=? and age=? and preference=?"
	var id int64

	err := db.db.QueryRow(q, f.Cluster, f.Gender, f.Weight, f.Height, f.Shape, f.Chest, f.Age, f.Preference).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	/**/
	return true, nil
}

//GetClusters select clusters ordered by how many products per it
func (db *DB) GetClusters() ([]string, error) {
	var out []string
	q := "SELECT cluster_id  FROM products_data WHERE fit_link = 1 and gender = 'Men' GROUP BY cluster_id ORDER BY count(cluster_id) DESC;"
	rows, err := db.db.Query(q)
	if err != nil {
		return out, err
	}

	for rows.Next() {
		var c string
		rows.Scan(&c)
		if c != "" {
			out = append(out, c)
		}
	}
	return out, nil
}

//GetLinksToProducstPerCluster select links to Products per cluster
func (db *DB) GetLinksToProducstPerCluster(cluster string) ([]string, error) {
	var out []string
	q := "SELECT link from products_data where cluster_id = ? and fit_link = 1;"
	rows, err := db.db.Query(q, cluster)
	if err != nil {
		return out, err
	}

	for rows.Next() {
		var c string
		rows.Scan(&c)
		if c != "" {
			out = append(out, c)
		}
	}
	return out, nil
}

//InsertFit is
func (db *DB) InsertFit(f Fit) {
	q := `INSERT INTO fit (
		gender,
		weight,
		height,
		shape,
		chest,
		age,
		preference,
		best_fit_size,
		best_fit_percent,
		next_fit_size,
		next_fit_percent,
		cluster_id
	)
	VALUES (
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?
	);
`
	_, err := db.db.Exec(q, f.Gender, f.Weight,
		f.Height, f.Shape, f.Chest, f.Age,
		f.Preference, f.BestFitSize, f.BestFitPercent,
		f.NextFitSize, f.NextFitPercent, f.Cluster)
	if err != nil {
		log.Println(err)
		fmt.Printf("[WARN] NOT inserted %s for h:%d  w:%d\n", f.Cluster, f.Height, f.Weight)

	} else {
		log.Printf("[INFO] inserted %s for h:%d  w:%d\n", f.Cluster, f.Height, f.Weight)
	}
}
