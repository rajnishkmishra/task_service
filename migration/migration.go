package migration

import (
	"log"

	"bitbucket.org/task_service/models"
	"bitbucket.org/task_service/utils"
)

type Service struct {
	mySqlClient *utils.MySQLConn
}

func NewMigrationService(mySqlClient *utils.MySQLConn) Service {
	return Service{mySqlClient: mySqlClient}
}

func (m *Service) migrate(entity interface{}) (err error) {
	defer utils.Recovery()
	if db := m.mySqlClient.DB; db != nil {
		if err = db.AutoMigrate(entity); err != nil {
			log.Fatal("migration failed")
		}
	}
	return err
}

func (m *Service) InitMigration() {
	log.Printf("goroutine::DB table migration started...")
	dbTables := map[string]interface{}{
		"tasks": &models.Task{},
		"users": &models.User{},
	}

	for tableName, table := range dbTables {
		err := m.migrate(table)
		if err != nil {
			log.Fatalf("migration failed for table: %v", tableName)
		}
	}
	log.Printf("goroutine::DB table migration concluded")
}
