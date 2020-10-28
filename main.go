package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/text/language"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"github.com/spf13/viper"
	. "github.com/codingXiang/gogo-i18n"
)

type IndexData struct {
	Title    string
	Content  string
	Login    string
}

type User struct {
	ID       int64  `json:"id" gorm:"primary_key;auto_increase'"`
	Username string `json:"username"`
	Password string `json:""`
}

func test(c *gin.Context) {
	data := new(IndexData)
	data.Title = "首頁"
	data.Content = "我的第一支 gin 專案"
	data.Login = "登入"
	c.HTML(http.StatusOK, "index.html", data)
}

func main() {
	//g := gin.Default()
	//g.LoadHTMLGlob("html/template/*.html")
	////設定靜態資源的讀取
	//g.Static("/assets", "./html/template/assets")
	//g.GET("/login", auth.LoginPage)
	//g.GET("/", test)
	//g.POST("/login", auth.LoginAuth)
	//g.Run(":7777")
	connectDB()
	dsn()
	viperYaml()
	i18n()
}

const (
	USERNAME = "root"
	PASSWORD = "root"
	NETWORK = "tcp"
	SERVER = "127.0.0.1"
	PORT = 8889
	DATABASE = "demo"
)

func connectDB() {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		fmt.Println("開啟 MySQL 連線發生錯誤，原因為：", err)
		return
	}
	if err := db.Ping(); err != nil {
		fmt.Println("資料庫連線錯誤，原因為：", err.Error())
		return
	}
	defer db.Close()
	CreateTable(db)
	InsertUser(db, "FirstTest", "password")
	QueryUser(db, "FirstTest")
}
func CreateTable(db *sql.DB) error {
	sql := `CREATE TABLE IF NOT EXISTS users(
	id INT(4) PRIMARY KEY AUTO_INCREMENT NOT NULL,
        username VARCHAR(64),
        password VARCHAR(64)
	); `

	if _, err := db.Exec(sql); err != nil {
		fmt.Println("建立 Table 發生錯誤:", err)
		return err
	}
	log.Println("建立 Table 成功！")
	return nil
}
func InsertUser(DB *sql.DB, username, password string) error{
	_,err := DB.Exec("insert INTO users(username,password) values(?,?)",username, password)
	if err != nil{
		fmt.Printf("建立使用者失敗，原因是：%v", err)
		return err
	}
	log.Println("建立使用者成功！")
	return nil
}
func QueryUser(db *sql.DB, username string) {
	user := new(User)
	row := db.QueryRow("select * from users where username=?", username)
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		fmt.Printf("映射使用者失敗，原因為：%v\n", err)
		return
	}
	log.Println("查詢使用者成功", *user)
}
func dsn()  {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("使用 gorm 連線 DB 發生錯誤，原因為 " + err.Error())
	}

	if err := db.AutoMigrate(new(User)); err != nil {
		panic("資料庫 Migrate 失敗，原因為 " + err.Error())
	}
	user := &User{
		Username: "SecondTest",
		Password: "password",
	}
	if err := CreateORMUser(db, user); err != nil {
		panic("資料庫 Migrate 失敗，原因為 " + err.Error())
	}
	if user, err := FindUser(db, 1); err == nil {
		log.Println("查詢到 User 為 ", user)
	} else {
		panic("查詢 user 失敗，原因為 " + err.Error())
	}
}
func CreateORMUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}
func FindUser(db *gorm.DB, id int64) (*User, error) {
	user := new(User)
	user.ID = id
	err := db.First(&user).Error
	return user, err
}
func viperYaml()  {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.SetDefault("application.port", 7777)
	err := viper.ReadInConfig()
	if err != nil {
		panic("讀取設定檔出現錯誤，原因為：" + err.Error())
	}
	fmt.Println("application port = " + viper.GetString("application.port"))
}

func i18n()  {
	GGi18n = NewGoGoi18n(language.TraditionalChinese)
	GGi18n.SetFileType("yaml")
	GGi18n.LoadTranslationFile("./i18n",
		language.TraditionalChinese,
		language.English)
	msg := GGi18n.GetMessage("welcome", map[string]interface{}{
		"username": "阿偉",
	})
	fmt.Println(msg)

	GGi18n.SetUseLanguage(language.English)
	msg = GGi18n.GetMessage("welcome", map[string]interface{}{
		"username": "阿偉",
	})
	fmt.Println(msg)
}