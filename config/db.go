package config

import (
	"ExchangeApp/global"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() {
	//data source name
	dsn := AppConfig.DataBase.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) //&gorm.Config{} 是 GORM 框架的初始化配置结构体指针,传入空结构体（如代码所示）表示采用 GORM 的默认行为

	if err != nil {
		log.Fatalf("Failed to initialize DataBase,got error: %v", err)
	}

	sqlDB, err := db.DB()

	sqlDB.SetMaxIdleConns(AppConfig.DataBase.MaxIdleConns)

	sqlDB.SetMaxOpenConns(AppConfig.DataBase.MaxOpenConns)

	sqlDB.SetConnMaxIdleTime(time.Hour)

	if err != nil {
		log.Fatalf("Failed to configue DataBase,got error: %v", err)
	}
	global.Db = db
	//db 与 sqlDB 的类型与区别
	//两者分别处理不同抽象层级的数据库交互：
	/*
		db (ORM 实例)

		类型：*gorm.DB

		作用：GORM 提供的对象关系映射客户端。**它负责将 Go 代码和结构体翻译为具体的 SQL 语句**，处理数据映射和业务逻辑验证。例如：db.Create(&user) 或 db.AutoMigrate()。
	*/

	/*
		sqlDB (底层连接池)

		类型：*database/sql.DB

		作用：Go 语言原生标准库提供的数据库连接池对象。它负责维护与 MySQL 进程的物理 TCP 连接生命周期。
		建立连接方面的规则,如代码中设置的 SetMaxIdleConns 和 SetMaxOpenConns
	*/
}
