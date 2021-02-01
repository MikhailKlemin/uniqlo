package misc

import (
	"database/sql"
	"log"

	"github.com/MikhailKlemin/uniclo.uk/pkg/config"
)

//UpdateProducts updates product by adding values to gender and boolean isBra
func UpdateProducts(conf config.DefaultConfig) {
	db, err := sql.Open("sqlite3", conf.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
