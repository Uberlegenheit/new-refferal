package postgres

import (
	"errors"
	"gorm.io/gorm"
	"new-refferal/models"
)

func (db *Postgres) GetLastLink() (*models.Link, error) {
	link := new(models.Link)
	result := db.db.Last(link)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &models.Link{
				Code: "",
			}, nil
		}
		return nil, result.Error
	}

	return link, nil
}