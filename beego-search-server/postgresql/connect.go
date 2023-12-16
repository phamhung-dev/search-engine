package postgresql

import (
	"beego-search-server/models"
	"fmt"
	"os"
	"strconv"

	beeLogger "github.com/beego/bee/v2/logger"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/client/orm/filter/bean"
	_ "github.com/lib/pq"
)

func NewPostgreSQLConnection() orm.Ormer {
	fmt.Println("Connecting to database ...")

	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh", host, user, password, dbname, port)

	orm.RegisterDriver("postgres", orm.DRPostgres)
	if err := orm.RegisterDataBase("default", "postgres", url); err != nil {
		beeLogger.Log.Fatal(err.Error())
	}

	orm.RegisterModel(new(models.User), new(models.File))

	builder := bean.NewDefaultValueFilterChainBuilder(nil, true, true)
	orm.AddGlobalFilterChain(builder.FilterChain)

	orm.Debug = true

	o := orm.NewOrm()

	if _, err := o.Raw(`CREATE EXTENSION IF NOT EXISTS "pg_trgm";`).Exec(); err != nil {
		beeLogger.Log.Fatal(err.Error())
	}

	generateTable, err := strconv.ParseBool(os.Getenv("GENERATE_TABLE"))
	if generateTable && err == nil {
		autogenerate()
	}

	if _, err := o.Raw(`CREATE INDEX IF NOT EXISTS idx_path_gin_trgm_ops ON files USING gin (regexp_replace(path, '(.*)[\\|/].*', '\1') gin_trgm_ops);`).Exec(); err != nil {
		beeLogger.Log.Fatal(err.Error())
	}

	if _, err := o.Raw(`CREATE INDEX IF NOT EXISTS idx_name_gin_trgm_ops ON files USING gin (name gin_trgm_ops);`).Exec(); err != nil {
		beeLogger.Log.Fatal(err.Error())
	}

	if _, err := o.Raw(`CREATE INDEX IF NOT EXISTS idx_content_gin_trgm_ops ON files USING gin (content gin_trgm_ops);`).Exec(); err != nil {
		beeLogger.Log.Fatal(err.Error())
	}

	if _, err := o.Raw(`CREATE INDEX IF NOT EXISTS idx_name_gin_trgm_ops ON files USING gin (name gin_trgm_ops);`).Exec(); err != nil {
		beeLogger.Log.Fatal(err.Error())
	}

	fmt.Println("Connected to database successfully!")

	return o
}
