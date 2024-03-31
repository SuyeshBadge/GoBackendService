package controller

import (
	"backendService/internals/modules/userModule/services"
	"errors"

	"encoding/json"
	"log"

	"io"

	"github.com/gin-gonic/gin"
)

type User_Controller struct {
	userService *services.User_Service
}

func (uc *User_Controller) GetUser(c *gin.Context) (interface{}, error) {
	id := c.Param("id")
	log.Println(id)

	user, err := uc.userService.GetUserByID(id)
	if err != nil {
		return nil, errors.New("failed to retrieve user")
	}
	return user, nil
}

func NewUserController(userService *services.User_Service) *User_Controller {
	return &User_Controller{
		userService: userService,
	}
}

func (uc *User_Controller) CreateUser(c *gin.Context) {
	body, err := io.ReadAll(io.Reader(c.Request.Body))
	if err != nil {
		// Handle error
		c.JSON(400, gin.H{
			"message": "Failed to read request body",
		})
		return
	}
	defer c.Request.Body.Close()

	// Parse request body into CreateUserData struct
	var createData services.CreateUserData
	if err := json.Unmarshal(body, &createData); err != nil {
		// Handle error
		c.JSON(400, gin.H{
			"message": "Invalid request body",
		})
		return
	}
	log.Println(createData)

	// Pass createData to UserService's CreateUser method
	if err := uc.userService.CreateUser(createData); err != nil {
		// Handle error
		c.JSON(500, gin.H{
			"message": "Failed to create user",
		})
		return
	}
}

func (uc *User_Controller) GetAllUsers(c *gin.Context) {
	users, err := uc.userService.GetUsers()
	if err != nil {

		// Handle error
		c.JSON(500, gin.H{
			"message": "Failed to retrieve users",
		})
		return
	}
	c.JSON(200, users)
}
