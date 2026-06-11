package model

import "time"

type ShortURLMap struct {
	ID       uint64    `gorm:"primaryKey;autoIncrement"`
	CreateAt time.Time `gorm:"autoCreateTime"`
	CreateBy string    `gorm:"size:64;not null;default:''"`
	IsDel    uint8     `gorm:"not null;default:0;index"`
	Lurl     string    `gorm:"size:2048"`
	Md5      string    `gorm:"size:32;uniqueIndex"`
	Surl     string    `gorm:"size:11;uniqueIndex"`
}

func (ShortURLMap) TableName() string {
	return "short_url_map"
}

type Sequence struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	Stub      string    `gorm:"size:1;not null;uniqueIndex:idx_uniq_stub"`
	Timestamp time.Time `gorm:"autoUpdateTime"`
}

func (Sequence) TableName() string {
	return "sequence"
}
