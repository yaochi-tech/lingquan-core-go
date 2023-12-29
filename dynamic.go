package lingquan

import "github.com/yaochi-tech/lingquan-core-go/db"

func StartEngine(driver, dsn string) (*db.Engine, error) {
	return db.NewEngine(driver, dsn)
}

func StopEngine(engine *db.Engine) {
	engine.Close()
}
