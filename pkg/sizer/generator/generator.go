package generator

import (
	"database/sql"
	"fmt"
	"log"
	"sort"

	"github.com/MikhailKlemin/uniclo.uk/pkg/config"
)

//Combo has item sizes
type Combo struct {
	//fit_matrix table
	ID         int    `db:"id"`
	Gender     int    `db:"gender"`
	Weight     int    `db:"weight"`
	Height     int    `db:"height"`
	Shape      int    `db:"shape"`
	Chest      int    `db:"chest"`
	Age        int    `db:"age"`
	Preference int    `db:"preference"`
	Serial     int64  `db:"serial"`
	Cookie     []byte `db:"cookie"`

	//fit_results table
	BestFitSize    string
	BestFitPercent int
	NextFitSize    string
	NextFitPercent int
	ClusterID      string `db:"cluster"`
}

//Generate is
func Generate() []Combo {
	//	var sets [][][]int
	//	sest = append(sets)
	sm := make(map[int][]int)
	sm[140] = makeRange(40, 65)
	sm[145] = makeRange(43, 65)
	sm[155] = makeRange(40, 120)
	sm[160] = makeRange(44, 121)
	sm[165] = makeRange(45, 123)
	sm[165] = makeRange(45, 123)
	sm[165] = makeRange(45, 123)
	sm[165] = makeRange(45, 123)
	sm[170] = makeRange(50, 125)
	sm[175] = makeRange(50, 128)
	sm[180] = makeRange(50, 130)
	sm[185] = makeRange(55, 130)
	sm[190] = makeRange(55, 130)
	sm[195] = makeRange(60, 135)
	sm[200] = makeRange(65, 135)
	sm[205] = makeRange(70, 132)
	sm[210] = makeRange(80, 132)

	var shapes = []int{1, 2, 3}
	var chests = []int{1, 2, 3}
	var ages = []int{20, 22, 24}
	const preference = 1
	var combos []Combo

	for height, weights := range sm {
		for _, weight := range weights {
			for _, shape := range shapes {
				for _, chest := range chests {
					for _, age := range ages {
						var c Combo
						c.Gender = 0
						c.Weight = weight
						c.Height = height
						c.Age = age
						c.Chest = chest
						c.Preference = preference
						c.Shape = shape
						combos = append(combos, c)
					}
				}
			}
		}
	}

	sort.Slice(combos, func(i, j int) bool {
		return combos[i].Height < combos[j].Height
	})

	fmt.Println(len(combos))
	fmt.Printf("%#v\n", combos[0])
	return combos

}

//GenerateAndPopulate generates and populates db with zero
func GenerateAndPopulate(conf config.DefaultConfig) {
	//	var sets [][][]int
	//	sest = append(sets)
	sm := make(map[int][]int)
	sm[140] = makeRange(40, 65)
	sm[145] = makeRange(43, 65)
	sm[155] = makeRange(40, 120)
	sm[160] = makeRange(44, 121)
	sm[165] = makeRange(45, 123)
	sm[165] = makeRange(45, 123)
	sm[165] = makeRange(45, 123)
	sm[165] = makeRange(45, 123)
	sm[170] = makeRange(50, 125)
	sm[175] = makeRange(50, 128)
	sm[180] = makeRange(50, 130)
	sm[185] = makeRange(55, 130)
	sm[190] = makeRange(55, 130)
	sm[195] = makeRange(60, 135)
	sm[200] = makeRange(65, 135)
	sm[205] = makeRange(70, 132)
	sm[210] = makeRange(80, 132)

	var shapes = []int{1, 2, 3}
	var chests = []int{1, 2, 3}
	var ages = []int{20, 22, 24}
	const preference = 1
	var combos []Combo

	for height, weights := range sm {
		for _, weight := range weights {
			for _, shape := range shapes {
				for _, chest := range chests {
					for _, age := range ages {
						var c Combo
						c.Gender = 0
						c.Weight = weight
						c.Height = height
						c.Age = age
						c.Chest = chest
						c.Preference = preference
						c.Shape = shape
						combos = append(combos, c)
					}
				}
			}
		}
	}

	sort.Slice(combos, func(i, j int) bool {
		return combos[i].Height < combos[j].Height
	})

	fmt.Println(len(combos))
	fmt.Printf("%#v\n", combos[0])
	db, err := sql.Open("sqlite3", conf.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	defer tx.Commit()

	q := `INSERT INTO fit_matrix (
		gender,
		weight,
		height,
		shape,
		chest,
		age,
		preference
	)
	VALUES (
		?,
		?,
		?,
		?,
		?,
		?,
		?
	);
`
	for _, combo := range combos {
		_, err := tx.Exec(q, combo.Gender, combo.Weight,
			combo.Height, combo.Shape, combo.Chest,
			combo.Age, combo.Preference)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
