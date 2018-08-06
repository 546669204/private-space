package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"

	"github.com/tidwall/gjson"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin:      func(r *http.Request) bool { return true },
		HandshakeTimeout: time.Duration(time.Second * 5),
	}
)

type handlerMsgModel struct {
	ID     int64
	Method string
	Data   string
}
type handlerProgressModel struct {
	Name       string
	ID         int64
	Percentage int64
	Import     bool
	Conn       *websocket.Conn
}

func main() {
	router := gin.New()
	router.GET("/", handleWebsocket)
	initSqlite()
	defer sess.Close()
	sess.SetLogging(true)
	router.Run(":31687")

}

func handleWebsocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("cant upgrade connection:", err)
		return
	}
	//onOpen
	for {
		msgType, msgData, err := conn.ReadMessage()
		if err != nil {
			log.Println("cant read message:", err)

			switch err.(type) {
			case *websocket.CloseError:
				//onClose
				return
			default:
				//onError
				return
			}
		}

		// Skip binary messages
		if msgType != websocket.TextMessage {
			continue
		}
		//onMessage
		log.Printf("incoming message: %s\n", msgData)

		var req handlerMsgModel
		json.Unmarshal(msgData, &req)
		if req.Method == "getAllFolder" {
			data, err := getAllFolder()
			if err != nil {
				conn.WriteJSON(gin.H{"id": req.ID, "method": "return", "data": gin.H{"code": 1, "msg": err.Error()}})
				continue
			}
			conn.WriteJSON(gin.H{"id": req.ID, "method": "return", "data": gin.H{"code": 0, "msg": "success", "data": data}})
			continue
		}
		if req.Method == "deleteFolder" {
			err := deleteFolder(gjson.Parse(req.Data).Get("id").Int())
			if err != nil {
				conn.WriteJSON(gin.H{"id": req.ID, "method": "return", "data": gin.H{"code": 1, "msg": err.Error()}})
				continue
			}
			conn.WriteJSON(gin.H{"id": req.ID, "method": "return", "data": gin.H{"code": 0, "msg": "success"}})
			continue
		}
		if req.Method == "createFolder" {
			err := createFolder(gjson.Parse(req.Data).Get("name").String(), gjson.Parse(req.Data).Get("father").Int())
			if err != nil {
				conn.WriteJSON(gin.H{"id": req.ID, "method": "return", "data": gin.H{"code": 1, "msg": err.Error()}})
				continue
			}
			conn.WriteJSON(gin.H{"id": req.ID, "method": "return", "data": gin.H{"code": 0, "msg": "success"}})
			continue
		}
		if req.Method == "getFolderFile" {
			data, err := getFolderFile(gjson.Parse(req.Data).Get("id").Int())
			if err != nil {
				conn.WriteJSON(gin.H{"id": req.ID, "method": "return", "data": gin.H{"code": 1, "msg": err.Error()}})
				continue
			}
			conn.WriteJSON(gin.H{"id": req.ID, "method": "return", "data": gin.H{"code": 0, "msg": "success", "data": data}})
			continue
		}
		if req.Method == "fileImport" {
			files := []string{}
			filesResult := gjson.Parse(req.Data).Get("file").Array()
			for index := 0; index < len(filesResult); index++ {
				files = append(files, filesResult[index].String())
			}
			err := fileImport(files, gjson.Parse(req.Data).Get("father").Int(), conn)
			if err != nil {
				conn.WriteJSON(gin.H{"id": req.ID, "method": "return", "data": gin.H{"code": 1, "msg": err.Error()}})
				continue
			}
			conn.WriteJSON(gin.H{"id": req.ID, "method": "return", "data": gin.H{"code": 0, "msg": "success"}})
			continue
		}
		if req.Method == "fileExport" {
			files := []int64{}
			filesResult := gjson.Parse(req.Data).Get("file").Array()
			for index := 0; index < len(filesResult); index++ {
				files = append(files, filesResult[index].Int())
			}
			err := fileExport(gjson.Parse(req.Data).Get("dstFile").Array()[0].String(), files, conn)
			if err != nil {
				conn.WriteJSON(gin.H{"id": req.ID, "method": "return", "data": gin.H{"code": 1, "msg": err.Error()}})
				continue
			}
			conn.WriteJSON(gin.H{"id": req.ID, "method": "return", "data": gin.H{"code": 0, "msg": "success"}})
			continue
		}
		if req.Method == "fileExportTemp" {
			err := fileExportTemp(gjson.Parse(req.Data).Get("file").Int(), conn)
			if err != nil {
				conn.WriteJSON(gin.H{"id": req.ID, "method": "return", "data": gin.H{"code": 1, "msg": err.Error()}})
				continue
			}
			conn.WriteJSON(gin.H{"id": req.ID, "method": "return", "data": gin.H{"code": 0, "msg": "success"}})
			continue
		}
		if req.Method == "login" {
			h := md5.New()
			h.Write([]byte("xm" + gjson.Parse(req.Data).Get("pass").String() + "521"))

			ok, err := sess.Collection("xm_user").Find().Where(`user = ? and pass = ?`, gjson.Parse(req.Data).Get("user").String(), hex.EncodeToString(h.Sum(nil))[8:24]).Exists()
			if err != nil {
				conn.WriteJSON(gin.H{"id": req.ID, "method": "return", "data": gin.H{"code": 1, "msg": err.Error()}})
				continue
			}
			if ok {
				conn.WriteJSON(gin.H{"id": req.ID, "method": "return", "data": gin.H{"code": 0, "msg": "success"}})
				continue
			} else {
				conn.WriteJSON(gin.H{"id": req.ID, "method": "return", "data": gin.H{"code": 1, "msg": "帐号不存在或密码错误"}})
				continue
			}

		}

	}
}

func createFolder(name string, father int64) error {
	temp_folder := Folder{
		Name:     name,
		FatherID: father,
		IsDir:    true,
	}
	exists, err := sess.Collection("xm_folder").Find().Where(`name = ? and father_id = ?`, name, father).Exists()
	if err != nil {
		return err
	}
	if exists {
		return errors.New("文件夹已存在")
	}
	sess.Collection("xm_folder").Insert(temp_folder)
	return nil
}
func deleteFolder(id int64) error {
	_, err := sess.Update("xm_folder").Where(`id = ?`, id).Set("delete_time = ?", time.Now().Unix()).Exec()
	return err
}
func getAllFolder() ([]Folder, error) {
	var AllFolder []Folder
	err := sess.Collection("xm_folder").Find().Where(`delete_time = 0 and is_dir = 1`).All(&AllFolder)
	return AllFolder, err
}
func getFolderFile(id int64) ([]FolderFileFileHistory, error) {
	var a1 []FolderFileFileHistory
	var a2 []FolderFileFileHistory

	err := sess.
		Select("d.id as folder_id", "*"). // Note how we set an alias for book.id.
		From("xm_folder AS d").
		Join("xm_file_history AS fs").On("d.id = fs.file_id and d.version = fs.file_version").Where("d.father_id = ? and d.delete_time = 0 and d.is_dir = 0", id).
		All(&a1)
	err = sess.Collection("xm_folder").Find().Select("id as folder_id", "*").Where("father_id = ? and delete_time = 0 and is_dir = 1", id).All(&a2)

	return append(a1, a2...), err
}

func fileImport(srcPath []string, father int64, conn *websocket.Conn) error {
	for index := 0; index < len(srcPath); index++ {
		FilePath := "./data/" + randomPath() + "/"
		FileName := randomFileName()
		RealFileName := getFileName(srcPath[index])
		FilePassword := randomPassword()

		exists, err := sess.Collection("xm_folder").Find().Where("father_id = ? and name = ?", father, RealFileName).Exists()
		if err != nil {
			return err
		}
		if exists {
			return errors.New("已存在同名文件")
		}
		insertFolder, err := sess.InsertInto("xm_folder").Values(Folder{
			Name:       RealFileName,
			FatherID:   father,
			CreateTime: time.Now().Unix(),
			UpdateTime: time.Now().Unix(),
			Version:    1,
			IsDir:      false,
		}).Exec()
		if err != nil {
			return err
		}
		FID, _ := insertFolder.LastInsertId()

		os.MkdirAll(FilePath, 0666)
		temp_p := handlerProgressModel{
			Name:       RealFileName,
			ID:         FID,
			Percentage: 0,
			Import:     true,
			Conn:       conn,
		}
		FileInfo, err := encryptFile(srcPath[index], FilePath+FileName, FilePassword, temp_p)
		if err != nil {
			return err
		}

		_, err = sess.Collection("xm_file_history").Insert(FileHistory{
			Path:        FilePath + FileName,
			Password:    FilePassword,
			FileID:      FID,
			FileVersion: 1,
			MD5:         FileInfo.MD5,
			Size:        FileInfo.Size,
			CreateTime:  time.Now().Unix(),
			UpdateTime:  FileInfo.ModTime.Unix(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func fileExport(dstPath string, fileID []int64, conn *websocket.Conn) error {
	var a1 []FolderFileFileHistory
	sess.
		Select("d.id as folder_id", "*"). // Note how we set an alias for book.id.
		From("xm_folder AS d").
		Join("xm_file_history AS fs").On("d.id = fs.file_id and d.version = fs.file_version").Where("d.id in ?", fileID).
		All(&a1)
	for index := 0; index < len(a1); index++ {
		temp_p := handlerProgressModel{
			Name:       a1[index].Name,
			ID:         a1[index].ID,
			Percentage: 0,
			Import:     false,
			Conn:       conn,
		}
		log.Println(temp_p)
		err := decryptFile(a1[index].Path, path.Join(dstPath, a1[index].Name), a1[index].Password, temp_p)
		if err != nil {
			return err
		}
	}
	return nil
}

func fileExportTemp(fileID int64, conn *websocket.Conn) error {
	var a1 []FolderFileFileHistory
	sess.
		Select("d.id as folder_id", "*"). // Note how we set an alias for book.id.
		From("xm_folder AS d").
		Join("xm_file_history AS fs").On("d.id = fs.file_id and d.version = fs.file_version").Where("d.id = ?", fileID).
		All(&a1)
	for index := 0; index < len(a1); index++ {
		temp_p := handlerProgressModel{
			Name:       a1[index].Name,
			ID:         a1[index].ID,
			Percentage: 0,
			Import:     false,
			Conn:       conn,
		}
		err := decryptFile(a1[index].Path, path.Join(os.TempDir(), a1[index].Name), a1[index].Password, temp_p)
		if err != nil {
			return err
		}
		log.Println(path.Join(os.TempDir(), a1[index].Name))

		Open(path.Join(os.TempDir(), a1[index].Name))
	}
	return nil
}

var commands = map[string][]string{
	"windows": []string{"cmd", "/c", "start"},
	"darwin":  []string{"open"},
	"linux":   []string{"xdg-open"},
}

func Open(uri string) error {
	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}
	run = append(run, uri)
	cmd := exec.Command(run[0], run[1:]...)
	return cmd.Start()
}
