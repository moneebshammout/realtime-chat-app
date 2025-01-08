package DAO

import (
	"group-service/internal/database"
	"group-service/pkg/utils"
)

var logger = utils.GetLogger()

type BaseDAO[T any] struct {
	dbClient *database.DBClient `gorm:"-:all" json:"-"`
}

func (b *BaseDAO[T]) FindAll(where ...map[string]interface{}) ([]T, error) {
	var result []T
	query := b.dbClient.Session

	if len(where) > 0 {
		for _, condition := range where {
			query = query.Where(condition)
		}
	}

	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (b *BaseDAO[T]) FindByID(id any, preloadRelations ...string) (T, error) {
	result := *new(T)
	query := b.dbClient.Session.Where("id = ?", id)

	for _, relation := range preloadRelations {
		query = query.Preload(relation)
	}

	if err := query.First(&result).Error; err != nil {
		logger.Errorf("Failed to find by ID: %v", err)
		return result, err
	}

	return result, nil
}

func (b *BaseDAO[T]) Create(data interface{}) error {
	result := b.dbClient.Session.Create(data)
	if result.Error != nil {
		return result.Error
	}

	logger.Infof("Model saved: %v", data)
	return nil
}

func (b *BaseDAO[T]) Update(data interface{}) error {
	result := b.dbClient.Session.Updates(data)
	if result.Error != nil {
		return result.Error
	}

	logger.Infof("Model updated: %v", data)
	return nil
}

func (b *BaseDAO[T]) Delete(ids ...uint) error {
	var model T
	if result := b.dbClient.Session.Where("id IN ?", ids).Delete(&model); result.Error != nil {
		logger.Errorf("Failed to delete records: %v", result.Error)
		return result.Error
	}
	logger.Infof("Models deleted with IDs: %v", ids)
	return nil
}
