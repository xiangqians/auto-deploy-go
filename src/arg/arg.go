// arg
// @author xiangqian
// @date 00:11 2022/12/31
package arg

import (
	"flag"
	"log"
	"strings"
)

var Port int
var Db string
var TmpDir string
var BuildEnv string

func Parse() {
	// parse
	flag.IntVar(&Port, "port", 8080, "-port 8080")
	flag.StringVar(&Db, "db", "./data/database.db", "-db ./data/database.db")
	flag.StringVar(&TmpDir, "tmpdir", "./tmp", "-tmpdir ./tmp")
	flag.StringVar(&BuildEnv, "buildenv", "default", "-buildenv default | docker:container")
	flag.Parse()

	// trim
	Db = strings.TrimSpace(Db)
	TmpDir = strings.TrimSpace(TmpDir)
	BuildEnv = strings.TrimSpace(BuildEnv)

	// log
	log.Printf("Port: %v\n", Port)
	log.Printf("Db: %v\n", Db)
	log.Printf("TmpDir: %v\n", TmpDir)
	log.Printf("BuildEnv: %v\n", BuildEnv)

	// -db "C:\Users\xiangqian\Desktop\tmp\auto-deploy\data\database.db" -tmpdir "C:\Users\xiangqian\Desktop\tmp\auto-deploy\tmp"
}
