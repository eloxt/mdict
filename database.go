package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type KeywordIndex struct {
	Keyword     string `gorm:"index"`
	OffsetStart int64
	OffsetEnd   int64
}

type ResourceIndex struct {
	FileName    string `gorm:"index"`
	OffsetStart int64
	OffsetEnd   int64
}

func SetupDatabase(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&KeywordIndex{}, &ResourceIndex{}); err != nil {
		return nil, err
	}

	return db, nil
}

func BatchInsertKeywords(db *gorm.DB, keywords []*KeywordIndex) error {
	return db.CreateInBatches(keywords, 1000).Error
}

func BatchInsertResources(db *gorm.DB, resources []*ResourceIndex) error {
	return db.CreateInBatches(resources, 1000).Error
}

func FindKeyword(db *gorm.DB, word string) (*KeywordIndex, error) {
	var result KeywordIndex
	err := db.First(&result, "keyword = ?", word).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func FindResource(db *gorm.DB, fileName string) (*ResourceIndex, error) {
	var result ResourceIndex
	fileNameWithSlash := "\\" + fileName
	err := db.Where("file_name = ? OR file_name = ?", fileName, fileNameWithSlash).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func FuzzyFindKeywords(db *gorm.DB, prefix string, limit int) ([]string, error) {
	var keywords []string
	pattern := prefix + "%"
	err := db.Model(&KeywordIndex{}).Where("keyword LIKE ?", pattern).Limit(limit).Order("offset_start").Pluck("keyword", &keywords).Error
	if err != nil {
		return nil, err
	}
	return keywords, nil
}
