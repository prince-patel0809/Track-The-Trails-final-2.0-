package config

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/joho/godotenv"
)

var CLD *cloudinary.Cloudinary

func ConnectCloudinary() {

	godotenv.Load()

	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUD_NAME"),
		os.Getenv("API_KEY"),
		os.Getenv("API_SECRET"),
	)

	if err != nil {
		log.Fatal(err)
	}

	CLD = cld
	log.Println("Cloudinary Connected")
}
