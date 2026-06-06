package config
import(
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB()  {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        getEnv("DB_USER", "root"),
        getEnv("DB_PASSWORD", "Password123"),
        getEnv("DB_HOST", "localhost"),
        getEnv("DB_PORT", "3306"),
        getEnv("DB_NAME", "db_auth_golang_jwt"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database : ", err)
	}

	DB = db
	fmt.Println("Database connected successfully")
}

func getEnv(key, fallback string) string  {
	if val := os.Getenv(key); val != ""{
		return val
	}

	return fallback
}

