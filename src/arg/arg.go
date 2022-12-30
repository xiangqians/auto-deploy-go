// arg
// @author xiangqian
// @date 00:11 2022/12/31
package arg

import (
	"flag"
	"log"
)

var Port int
var Db string
var TmpDir string

func Parse() {
	// -port
	flag.IntVar(&Port, "port", 8080, "port")
	// -db
	flag.StringVar(&Db, "db", "./data/database.db", "database")
	// -tmpdir
	flag.StringVar(&TmpDir, "tmpdir", "./tmp", "tmpdir")
	flag.Parse()

	log.Printf("Port: %v\n", Port)
	log.Printf("Db: %v\n", Db)
	log.Printf("TmpDir: %v\n", TmpDir)

	// -db "C:\Users\xiangqian\Desktop\tmp\auto-deploy\data\database.db" -tmpdir "C:\Users\xiangqian\Desktop\tmp\auto-deploy\tmp"
}
