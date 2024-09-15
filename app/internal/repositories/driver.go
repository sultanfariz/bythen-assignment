package repositories

import (
	"fmt"

	post "app/internal/repositories/post"
	user "app/internal/repositories/user"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
}

func InitDB(c DBConfig) *gorm.DB {
	config := map[string]string{
		"DB_USERNAME": c.Username,
		"DB_PASSWORD": c.Password,
		"DB_HOST":     c.Host,
		"DB_PORT":     c.Port,
		"DB_NAME":     c.Name,
	}

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", config["DB_USERNAME"], config["DB_PASSWORD"], config["DB_HOST"], config["DB_PORT"], config["DB_NAME"])

	DB, err := gorm.Open(mysql.Open(connection), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(
		&user.User{},
		&post.Post{},
	)
	return DB
}
