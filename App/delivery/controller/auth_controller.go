package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syahyudi09/ChatPintar/App/helper"
	"github.com/syahyudi09/ChatPintar/App/model"
	"github.com/syahyudi09/ChatPintar/App/usecase"
)

type AuthController struct {
	router  *gin.Engine
	auth usecase.AuthUsecase
}

func (ac *AuthController) Register(c *gin.Context) {
	var input model.UserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		fmt.Println(err)
		errorMessage := gin.H{"errors": "FAILED_TO_PROCESS_REGISTER_REQUEST"}
		response := helper.APIResponse("FAILED_TO_REGISTER_USER", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	exists, err := ac.auth.PhoneNumberExits(input.PhoneNumber)
	if err != nil {
		fmt.Println(err)
		errorMessage := gin.H{"errors": "FAILED_TO_CHECK_PHONENUMBER_EXISTENCE"}
		response := helper.APIResponse("Login Failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	statusCode := http.StatusOK
	message := "CREATE_SUCCESSFULY"

	if exists {
		statusCode = http.StatusConflict
		message = "PHONE_NUMBER_ALREADY_EXISTS"
	} else {
		// Jika pendaftaran berhasil, respons akan menampilkan pesan sukses dan data pengguna yang terdafta
		if err := ac.auth.Register(input); err != nil {
			fmt.Println(err)
			response := helper.APIResponse("FAILED_TO_REGISTER_USER", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}

	// Berhasil melakukan register
	response := helper.APIResponse(message, statusCode, "success", input)
	c.JSON(http.StatusOK, response)
}

func NewAuthController(r *gin.Engine, usecase usecase.AuthUsecase) *AuthController {
	controller := AuthController{
		router:  r,
		auth: usecase,
	}

	auth := r.Group("/auth")

	auth.POST("/register", controller.Register)

	return &controller
} 