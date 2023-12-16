package postgresql

import (
	"fmt"

	beeLogger "github.com/beego/bee/v2/logger"
	"github.com/beego/beego/v2/client/orm"
)

func autogenerate() {
	fmt.Println("Autogenerating...")

	name := "default"
	err := orm.RunSyncdb(name, true, true)
	if err != nil {
		beeLogger.Log.Fatal(err.Error())
	}

	fmt.Println("Autogenerate successfully!")
}
