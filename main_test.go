package main

import (
	"encoding/json"
	"log"
	"testing"

	"upper.io/db.v3/sqlite"
)

func Test_main(t *testing.T) {

	var err error
	sess, err = sqlite.Open(sqlite.ConnectionURL{
		Database: `./data.db`, //数据库名字
	})
	if err != nil {
		log.Fatal("Conn ==> ", err)
	}
	//defer sess.Close()
	sess.SetLogging(true)

	log.Println(sess.Name(), sess.ConnectionURL())
	//getAllFolder()
	log.Println(sess.Collections())

	data, _ := getFolderFile(1)
	sbyte, _ := json.Marshal(data)
	log.Println(string(sbyte))
	sess.Close()
}

func Test_rand(t *testing.T) {
	log.Println((131072 * 100 / 854072))
}
