package dao

import (
	"log"

	"shortURL/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitMySQL(dsn string) {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("failed to connect mysql: %v", err)
	}
	if err := DB.AutoMigrate(&model.Sequence{}, &model.ShortURLMap{}); err != nil {
		log.Fatalf("failed to migrate tables: %v", err)
	}
	log.Println("mysql connected")
}

func NextSequenceID() (uint64, error) {
	var seq model.Sequence
	if err := DB.Exec("REPLACE INTO sequence (stub) VALUES ('a')").Error; err != nil {
		return 0, err
	}
	DB.Raw("SELECT LAST_INSERT_ID()").Scan(&seq)
	return seq.ID, nil
}

func FindByMd5(md5 string) (*model.ShortURLMap, error) {
	var m model.ShortURLMap
	err := DB.Where("md5 = ? AND is_del = 0", md5).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func FindBySurl(surl string) (*model.ShortURLMap, error) {
	var m model.ShortURLMap
	err := DB.Where("surl = ? AND is_del = 0", surl).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func CreateShortURL(m *model.ShortURLMap) error {
	return DB.Create(m).Error
}

func LoadAllSurls() ([]string, error) {
	var surls []string
	err := DB.Model(&model.ShortURLMap{}).Where("is_del = 0").Pluck("surl", &surls).Error
	return surls, err
}
