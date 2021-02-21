package main

import (
	"fmt"
	"log"
	"time"

	"github.com/MikhailKlemin/uniclo.uk/pkg/config"
	sizer "github.com/MikhailKlemin/uniclo.uk/pkg/sizer2"

	_ "github.com/mattn/go-sqlite3" //comment
)

func main() {

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//client.Start()
	t := time.Now()

	var conf = config.NewDefaultConfig()

	//crawler.CollectProducts(conf)
	//	crawler.CollectSizes(conf)
	//fitanalytic.CheckForFitLink(conf)
	sizer.Start(conf)
	//os.Exit(1)

	fmt.Printf("Done in %s\n", time.Since(t))

}
