package main

import (
	"flag"
	"fmt"

	"github.com/dynamicgo/config"
	"github.com/dynamicgo/slf4go"
	"github.com/go-xorm/xorm"
	"github.com/inwecrypto/neodb"
	_ "github.com/lib/pq"
)

var logger = slf4go.Get("neodb")
var configpath = flag.String("conf", "./neodb.json", "neodb config file")

func main() {

	flag.Parse()

	conf, err := config.NewFromFile(*configpath)

	if err != nil {
		logger.ErrorF("load eth indexer config err , %s", err)
		return
	}

	username := conf.GetString("neodb.username", "xxx")
	password := conf.GetString("neodb.password", "xxx")
	port := conf.GetString("neodb.port", "6543")
	host := conf.GetString("neodb.host", "localhost")
	scheme := conf.GetString("neodb.schema", "postgres")

	engine, err := xorm.NewEngine("postgres", fmt.Sprintf("user=%v password=%v host=%v dbname=%v port=%v sslmode=disable", username, password, host, scheme, port))

	if err != nil {
		logger.ErrorF("create postgres orm engine err , %s", err)
		return
	}

	tables := []interface{}{
		new(neodb.Tx),
		new(neodb.Block),
		new(neodb.UTXO),
		new(neodb.Order),
		new(neodb.Wallet),
	}

	if err := engine.Sync2(tables...); err != nil {
		logger.ErrorF("sync table schema error , %s", err)
		return
	}

}
