package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"fullstack2025-test/database"
	"fullstack2025-test/models"
	"io"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllMyClient(ctx *fiber.Ctx) error {
	var myClient []models.MyClient
	redisKey := "all_my_client"

	cachedData, err := database.RedisClient.Get(context.Background(), redisKey).Result()
	if err == nil {
		err = json.Unmarshal([]byte(cachedData), &myClient)
		if err == nil {
			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "Get my client from cache successfuly",
				"data":    myClient,
			})
		} else {
			fmt.Println("Error unmarshal:", err)
		}
	}
	result := database.DB.Find(&myClient)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed Get My client",
			"error":   result.Error.Error(),
		})
	}

	jsonData, err := json.Marshal(myClient)
	if err == nil {
		database.RedisClient.Set(context.Background(), redisKey, jsonData, time.Hour)
	} else {
		fmt.Println("Error marshal data:", err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Get My client successfully",
		"data":    myClient,
	})

}

func GetMyClientBySlug(ctx *fiber.Ctx) error {
	var myClient models.MyClient
	slug := ctx.Query("slug")

	redisKey := "slug"
	cachedData, err := database.RedisClient.Get(context.Background(), redisKey).Result()
	if err == nil {
		err = json.Unmarshal([]byte(cachedData), &myClient)
		if err == nil {
			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "Get my client from cache successfuly",
				"data":    myClient,
			})
		} else {
			fmt.Println("Error unmarshal:", err)
		}
	}
	result := database.DB.Where("slug LIKE ?", "%"+slug+"%").First(&myClient)
	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "My Client not found",
			"error":   result.Error.Error(),
		})
	}

	jsonData, err := json.Marshal(myClient)
	if err == nil {
		database.RedisClient.Set(context.Background(), redisKey, jsonData, time.Hour)
	} else {
		fmt.Println("Error marshal data:", err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "GetMy Client successfully",
		"data":    myClient,
	})

}

func CreateMyClient(ctx *fiber.Ctx) error {
	myClient := new(models.MyClient)

	database.RedisClient.Del(context.Background(), "all_my_client")

	if form, err := ctx.MultipartForm(); err == nil {
		if len(form.Value["name"]) > 0 {
			myClient.Name = form.Value["name"][0]
		}
		if len(form.Value["slug"]) > 0 {
			myClient.Slug = form.Value["slug"][0]
		}
		if len(form.Value["client_prefix"]) > 0 {
			myClient.ClientPrefix = form.Value["client_prefix"][0]
		}

		if files := form.File["client_logo"]; len(files) > 0 {
			file := files[0]
			src, err := file.Open()
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to open file",
					"error":   err.Error(),
				})
			}
			defer src.Close()

			filename := fmt.Sprintf("my_clients/%s-%s", uuid.New().String(), file.Filename)

			imageURL, err := uploadToCloudinary(src, filename)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Upload file to Cloudinary is failed",
					"error":   err.Error(),
				})
			}

			myClient.ClientLogo = imageURL
		}
	} else {
		if err := ctx.BodyParser(&myClient); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid data",
				"error":   err.Error(),
			})
		}
	}

	result := database.DB.Create(&myClient)
	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create my client",
			"error":   result.Error.Error(),
		})
	}

	jsonData, err := json.Marshal(myClient)
	if err == nil {
		database.RedisClient.Set(context.Background(), myClient.Slug, jsonData, time.Hour)
	} else {
		fmt.Println("Error marshal data:", err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Create My client succesfully",
		"data":    myClient,
	})
}

func UpdateMyClient(ctx *fiber.Ctx) error {
	var myClient models.MyClient
	slug := ctx.Query("slug")

	result := database.DB.Where("slug =", slug).First(&myClient)
	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "My Client not found",
			"error":   result.Error.Error(),
		})
	}
	if form, err := ctx.MultipartForm(); err == nil {
		if len(form.Value["name"]) > 0 {
			myClient.Name = form.Value["name"][0]
		}
		if len(form.Value["slug"]) > 0 {
			myClient.Slug = form.Value["slug"][0]
		}
		if len(form.Value["client_prefix"]) > 0 {
			myClient.ClientPrefix = form.Value["client_prefix"][0]
		}

		if files := form.File["client_logo"]; len(files) > 0 {
			file := files[0]
			src, err := file.Open()
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to open file",
					"error":   err.Error(),
				})
			}
			defer src.Close()

			filename := fmt.Sprintf("my_client/%s-%s", uuid.New().String(), file.Filename)

			imageURL, err := uploadToCloudinary(src, filename)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Upload file to Cloudinary is failed",
					"error":   err.Error(),
				})
			}

			myClient.ClientLogo = imageURL
		}
	} else {
		if err := ctx.BodyParser(&myClient); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid data",
				"error":   err.Error(),
			})
		}
	}

	err := database.RedisClient.Del(context.Background(), slug).Err()
	if err != nil {
		return fmt.Errorf("failed delete %s: %v", slug, err)
	}

	if err := database.DB.Save(&myClient).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed save data My client",
			"error":   err.Error(),
		})
	}

	jsonData, err := json.Marshal(myClient)
	if err == nil {
		database.RedisClient.Set(context.Background(), myClient.Slug, jsonData, time.Hour)
	} else {
		fmt.Println("Error marshal data:", err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Update My client successfully",
		"data":    myClient,
	})

}

func DeleteMyClient(ctx *fiber.Ctx) error {
	var myClient models.MyClient
	slug := ctx.Query("slug")

	result := database.DB.Where("slug LIKE ?", "%"+slug+"%").First(&myClient)
	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "My Client not found",
			"error":   result.Error.Error(),
		})
	}
	database.DB.Delete(&myClient)
	database.RedisClient.Del(context.Background(), "slug")
	database.RedisClient.Del(context.Background(), "all_my_client")
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "delete My client successfully",
	})
}

func uploadToCloudinary(file io.Reader, filename string) (string, error) {
	ctx := context.Background()
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		return "", err
	}

	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: filename,
	})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
