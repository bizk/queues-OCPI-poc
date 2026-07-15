package configs

import "os"

const DB_NAME = "poc"

func GetDBConnection() string {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	return uri
}
