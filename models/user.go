package models

import (
	"gorm.io/gorm"
)

// GORM 通过将 Go 结构体（Go structs） 映射到数据库表来简化数据库交互。
// 结构体 (Struct) -> 表 (Table)
// 结构体字段 (Field) -> 列 (Column)
// 默认情况下，GORM 将结构体名称转换为 snake_case 并为表名加上复数形式
// User -> users  GormUserName->gorm_user_names
// GORM 自动将结构体字段名称转换为 snake_case 作为数据库中的列名
// Username -> usernames

//GORM提供了一个预定义的结构体，名为gorm.Model
/*
	// gorm.Model 的定义
	type Model struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time		 //CreatedAt ：在创建记录时自动设置为当前时间。
	UpdatedAt time.Time		 //UpdatedAt：每当记录更新时，自动更新为当前时间
	DeletedAt gorm.DeletedAt `gorm:"index"`
							 //DeletedAt：用于软删除（将记录标记为已删除，而实际上并未从数据库中删除）。
	}
*/

//将其嵌入在结构体中:
//在结构体中嵌入 gorm.Model ，以便自动包含这些字段。
//这对于在不同模型之间保持一致性并利用GORM内置的约定非常有用

type User struct {
	//GORM 默认将名为 ID 的字段视为表的主键
	gorm.Model
	Uasername string `gorm:"unique"`
	Password  string // 列名:`password`数据表的列名使用的是 struct 字段名的蛇形命名(Snake Case)
}
