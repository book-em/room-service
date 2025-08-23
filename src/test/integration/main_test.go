package integration

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dbItem struct {
	Service    string
	Connection *gorm.DB
}

var dBs []dbItem
var services = []string{"room", "user"}

func getConnection(service string) *gorm.DB {

	var connection *gorm.DB

	for _, dbItem := range dBs {
		if service == dbItem.Service {
			connection = dbItem.Connection
			break
		}
	}

	if connection == nil {
		panic(fmt.Sprintf("no %v-database connection found", service))
	}

	return connection
}

func cleanup(service string) {

	var db *gorm.DB = getConnection(service)

	var tables []string
	err := db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tables).Error
	if err != nil {
		log.Fatalf("error fetching table names: %v", err)
	}

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)
		err = db.Exec(query).Error
		if err != nil {
			log.Fatalf("error truncating table %s: %v", table, err)
		}
	}

	log.Printf("%v-database cleaned up successfully", service)
}

func connectToDBs() {

	for _, service := range services {

		s := strings.ToUpper(service)
		host := os.Getenv(fmt.Sprintf("%s_DB_HOST", s))
		port := os.Getenv(fmt.Sprintf("%s_DB_PORT", s))
		user := os.Getenv(fmt.Sprintf("%s_DB_USER", s))
		password := os.Getenv(fmt.Sprintf("%s_DB_PASSWORD", s))
		dbname := os.Getenv(fmt.Sprintf("%s_DB_NAME", s))

		dbURL := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname,
		)

		db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to open %v-database", err)
		}

		dbItem := dbItem{
			Service:    service,
			Connection: db,
		}

		dBs = append(dBs, dbItem)

		log.Printf("connected to %v-database", service)
	}

}

func closeDBConnections() {
	for _, dbItem := range dBs {
		db, err := dbItem.Connection.DB()

		if err != nil {
			log.Printf("error getting %v-database instance: %v", dbItem.Service, err)
			continue
		}

		err = db.Close()
		if err != nil {
			log.Printf("error closing %v-database connection: %v", dbItem.Service, err)
		} else {
			log.Printf("closed connection to %v-database", dbItem.Service)
		}
	}
}

func TestMain(m *testing.M) {

	connectToDBs()

	code := m.Run()

	closeDBConnections()

	os.Exit(code)
}
