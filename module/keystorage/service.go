package keystorage

import (
	"errors"
	"gorm.io/gorm"
)

type client struct {
	tableName string
	db        *gorm.DB
	Error     error
}

func NewClient(db *gorm.DB) Client {
	return &client{
		db:    db,
		Error: nil,
	}
}

func (c *client) Get(key string) string {
	c.autoMigrate()
	var result KeyStorage
	err := c.db.Table(c.tableName).Where("key = ?", key).First(&result).Error
	if err != nil {
		c.Error = err
		return ""
	}
	return result.Value
}

func (c *client) Set(key, value string) {
	c.autoMigrate()
	k := KeyStorage{
		Key:   key,
		Value: value,
	}
	c.Get(key)
	if c.Error != nil {
		if errors.Is(c.Error, gorm.ErrRecordNotFound) {
			c.Error = nil
			err := c.db.Table(c.tableName).Create(&k).Error
			if err != nil {
				c.Error = err
				return
			}
			return
		}
		return
	}
	err := c.db.Table(c.tableName).Where("key = ?", key).Updates(&k).Error
	if err != nil {
		c.Error = err
		return
	}
}

func (c *client) Delete(key string) {
	c.autoMigrate()
	err := c.db.Table(c.tableName).Where("key = ?", key).Delete(&KeyStorage{}).Error
	if err != nil {
		c.Error = err
		return
	}
}

func (c *client) Err() error {
	return c.Error
}

func (c *client) SetTableName(name string) {
	if c == nil {
		return
	}
	c.tableName = name
}

func (c *client) autoMigrate() {
	if c == nil {
		return
	}
	defaultTableName := "key_storage"
	if c.tableName == "" {
		c.tableName = defaultTableName
	}
	if c.db.Migrator().HasTable(c.tableName) {
		return
	}
	err := c.db.Table(c.tableName).AutoMigrate(&KeyStorage{})
	if err != nil {
		c.Error = err
		return
	}
}
