package analyzer

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

/*
//GetKey get a key
func GetKey(key string) {
	//key:="423539"
	db2, err := dbase.NewDB("./assets/meta.data.v2.db")
	if err != nil {
		log.Fatal(err)
	}

	b, err := db2.GetByte(key)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%#v\n", b)
	defer db2.Close()

}

//Analyze analyzing sizes
func Analyze() {
	db2, err := dbase.NewDB("./assets/meta.data.v2.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db2.Close()

	iter := db2.NewIter()
	c := 0
	for iter.Next() {
		key := iter.Key()
		fmt.Println(string(key))
		c++
	}

	fmt.Println(c)
}
*/

//Analyze analyzing sizes
func Analyze() {
	re := regexp.MustCompile(`\s+`)
	group := make(map[string][]string)
	db, err := sql.Open("sqlite3", "./assets/uniclo.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	rows, err := db.Query(`select sf.serial, sd.content from size_fitanal sf inner join sizes_data sd  on sd.serial = sf.serial where sf.fitanal = 1;`)
	if err != nil {
		log.Fatal(err)
	}

	c := 0
	for rows.Next() {
		c++
		var keys string
		var valx []byte
		if err := rows.Scan(&keys, &valx); err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("%s - %s\n", key, val)
		//fmt.Println(string(key))
		//size := len(val)
		txt, err := Tokenize(bytes.NewReader(valx))
		if err != nil {
			log.Fatal(err)
		}
		txt = strings.TrimSpace(re.ReplaceAllString(txt, ""))
		fmt.Println(len(txt))
		group[txt] = append(group[txt], keys)
	}
	f, _ := os.Create("./assets/cross.csv")
	defer f.Close()
	w := csv.NewWriter(f)

	for _, v := range group {
		if len(v) > 0 {
			//fmt.Println(strings.Join(v, "\t"))
			w.Write(v)
			//fmt.Println("---------------------------------------------------------------")

		}

	}
	w.Flush()

	fmt.Println(len(group), "\t", c)

}

//Tokenize tokenize HTLML
//This function is used to tokenize HTML to
//represent it clean TEXT
func Tokenize(r io.Reader) (string, error) {
	textTags := []string{
		"td",
		//		"p", "span", "em", "string", "blockquote", "q", "cite",
		//		"h1", "h2", "h3", "h4", "h5", "h6", "pre", "ul", "li", "ol",
		//		"mark", "ins", "del", "small", "i", "b",
	}

	tag := ""
	enter := false
	var text []string
	tokenizer := html.NewTokenizer(r)
	for {
		tt := tokenizer.Next()
		token := tokenizer.Token()

		err := tokenizer.Err()
		if err == io.EOF {
			break
		}

		switch tt {
		case html.ErrorToken:
			//log.Fatal(err)
			return "", errors.Wrap(err, "can't parse token")
		case html.StartTagToken, html.SelfClosingTagToken:
			enter = false

			tag = token.Data
			for _, ttt := range textTags {
				if tag == ttt {
					enter = true
					break
				}
			}
		case html.TextToken:
			if enter {
				data := strings.TrimSpace(token.Data)

				if len(data) > 0 {
					//fmt.Println(data)
					text = append(text, data)
				}
			}
		}
	}

	return strings.Join(text, " "), nil
}
