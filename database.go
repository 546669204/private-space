package main

import (
	"log"

	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/sqlite"
)

var sess sqlbuilder.Database
var mysqlsettings = sqlite.ConnectionURL{
	Database: `./data.db`, //数据库名字
}

type Folder struct {
	ID         int64  `db:"id,omitempty"`
	Name       string `db:"name"`
	FatherID   int64  `db:"father_id"`
	CreateTime int64  `db:"create_time"`
	UpdateTime int64  `db:"update_time"`
	DeleteTime int64  `db:"delete_time"`
	Version    int64  `db:"version"`
	IsDir      bool   `db:"is_dir"`
}

type FileHistory struct {
	ID          int64  `db:"id,omitempty"`
	Path        string `db:"path"`
	Password    string `db:"password"`
	CreateTime  int64  `db:"create_time"`
	UpdateTime  int64  `db:"update_time"`
	DeleteTime  int64  `db:"delete_time"`
	FileID      int64  `db:"file_id"`
	FileVersion int64  `db:"file_version"`
	Size        int64  `db:"size"`
	MD5         string `db:"md5"`
}

type FolderFileFileHistory struct {
	ID          int64 `db:"folder_id"`
	Folder      `db:",inline"`
	FileHistory `db:",inline"`
}

// CREATE TABLE xm_folder (
//     id        BIGINT       PRIMARY KEY AUTOINCREMENT,
//     name      VARCHAR (50) NOT NULL,
//     father_id BIGINT,
//     file_id   TEXT
// );

func initSqlite() {
	var err error
	sess, err = sqlite.Open(mysqlsettings)
	if err != nil {
		log.Fatal("Conn ==> ", err)
	}
}
