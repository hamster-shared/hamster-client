package chain

import "gorm.io/gorm"

type Chain interface {
	GetDeployParam(appID int, db *gorm.DB) any
}
