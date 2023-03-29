package db

import (
	"bulk_issuance/config"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

type DBFileUpload struct {
	gorm.Model
	Filename     string
	TotalRecords int
	UserID       string
	Date         string
}

type DBFilesUpload struct {
	gorm.Model
	Filename string
	Headers  string
	RowData  []byte
}

func Init() {
	var e error
	dbPath := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Config.Database.Host, config.Config.Database.Port,
		config.Config.Database.User, config.Config.Database.Password, config.Config.Database.DBName,
	)
	log.Printf("Using db %s", dbPath)
	db, e = gorm.Open("postgres", dbPath)
	if e != nil {
		log.Printf("Error %+v", e)
		panic("failed to connect to database")
	}
	db.AutoMigrate(&DBFileUpload{})
	db.AutoMigrate(&DBFilesUpload{})
}

func CreateDBFileUpload(data *DBFileUpload) error {
	if result := db.Create(&data); result.Error != nil {
		log.Printf("%v", data.Filename)
	}
	log.Printf("%v", data.Filename)
	return nil
}

func GetDBFilesUpload(id int) *DBFilesUpload {
	filesUpload := &DBFilesUpload{}
	if result := db.First(&filesUpload, "id=?", id); result.Error != nil {
		log.Fatal("Error : %v", result.Error)
	}
	return filesUpload
}

func CreateDBFilesUpload(data *DBFilesUpload) error {
	if result := db.Create(&data); result.Error != nil {
		log.Printf("%v", data.Filename)
	}
	log.Printf("%v", data.Filename)
	return nil
}

func GetAllUploadedFilesData() ([]DBFileUpload, error) {
	var files []DBFileUpload
	if err := db.Find(&files).Error; err != nil {
		log.Fatal(err)
		return nil, err
	}
	return files, nil
}
