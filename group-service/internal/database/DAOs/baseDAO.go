package DAO

import (
	"group-service/internal/database"
	"group-service/pkg/utils"
)

var logger = utils.GetLogger()

type BaseDAO struct {
	dbClient *database.DBClient `gorm:"-:all" json:"-"`
	model    any
}

func (b *BaseDAO) Create(data interface{}) error {
	result := b.dbClient.Session.Create(data)
	if result.Error != nil {
		return result.Error
	}

	logger.Infof("Model saved: %v", data)
	return nil
}

func (b *BaseDAO) Update(data interface{}) error {
	result := b.dbClient.Session.Updates(data)
	if result.Error != nil {
		return result.Error
	}

	logger.Infof("Model updated: %v", data)
	return nil
}

func (b *BaseDAO) Delete(data interface{}) error {
	result := b.dbClient.Session.Delete(data)
	if result.Error != nil {
		return result.Error
	}

	logger.Infof("Model deleted: %v", data)
	return nil
}

func (b *BaseDAO) FindAll() []any {
	var result []any
	b.dbClient.Session.Find(&result)
	return result
}

// func (b *BaseDAO) FindByID(id string,result interface{}){
// 	b.dbClient.Session.Find(result, id)
// }


func (b *BaseDAO) FindByID(id string, result interface{}, preloadRelations ...string) error {
	// Start with the base query to find the record by ID
	query := b.dbClient.Session.Where("id = ?", id)

	// Dynamically preload the relations if provided
	for _, relation := range preloadRelations {
		query = query.Preload(relation)
	}

	// Fetch the record by ID and preload the relations
	if err := query.First(result).Error; err != nil {
		// Handle error (e.g., record not found or DB error)
		// log.Println("Error finding record by ID:", err)
		return err
	}

	return nil
}