package postgresql

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgreSQLConnection() *gorm.DB {
	fmt.Println("Connecting to database ...")

	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	url := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
		host,
		user,
		password,
		dbname,
		port,
	)

	db, err := gorm.Open(
		postgres.Open(url),
		&gorm.Config{
			TranslateError:  true,
			PrepareStmt:     true,
			CreateBatchSize: 500,
			// Logger:         logger.Default.LogMode(logger.Silent),
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"pg_trgm\";")

	generateTable, err := strconv.ParseBool(os.Getenv("GENERATE_TABLE"))
	if generateTable && err == nil {
		autogenerate(db)
	}

	db.Exec(`CREATE INDEX IF NOT EXISTS idx_path_gin_trgm_ops ON files USING gin (regexp_replace(path, '(.*)[\\|/].*', '\1') gin_trgm_ops);`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_name_gin_trgm_ops ON files USING gin (name gin_trgm_ops);`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_content_gin_trgm_ops ON files USING gin (content gin_trgm_ops);`)
	// db.Exec(`CREATE INDEX IF NOT EXISTS idx_file_created_at_tstzrange ON files USING gist (tstzrange(file_created_at, file_created_at, '[]'));`)
	// db.Exec(`CREATE INDEX IF NOT EXISTS idx_file_last_modified_at_tstzrange ON files USING gist (tstzrange(file_last_modified_at, file_last_modified_at, '[]'));`)
	// db.Exec(`CREATE INDEX IF NOT EXISTS idx_file_last_accessed_at_tstzrange ON files USING gist (tstzrange(file_last_accessed_at, file_last_accessed_at, '[]'));`)

	fmt.Println("Connected to database successfully!")

	return db
}
