package controller

import (
	// "fmt"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syahyudi09/ChatPintar/App/model"
	"github.com/syahyudi09/ChatPintar/App/usecase"
)

type ChatGroupController struct {
	router    *gin.Engine
	chatGroup usecase.ChatGroupUsecase
}

func (cgc *ChatGroupController) CreateGroup(c *gin.Context) {
	phoneNumber := c.Param("phone_number")

	var input model.InputChatGroupModel

	if err := c.ShouldBindJSON(&input); err != nil || input.GroupName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input or missing group_name"})
		return
	}

	groupID, err := cgc.chatGroup.CreateGroup(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create group: %v", err)})
		return
	}

	adminInput := model.InputUserToGroup{
		PhoneNumber: phoneNumber,
		GroupId:     groupID,
		Role:        model.Admin,
	}

	err = cgc.chatGroup.AddUserToGroup(adminInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add admin to group: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"group_id": groupID, "message": "Group created and admin assigned"})

}

func (cgc *ChatGroupController) AddUserGroup(c *gin.Context) {
	groupId := c.Param("group_id")
	adminNumber := c.Param("phone_number")

	var input model.InputUserToGroup
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	input.AddPhoneNumber = adminNumber
	input.GroupId = groupId
	input.Role = model.Member

	err := cgc.chatGroup.AddUserToGroup(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add user to group: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Member added to group successfully"})
}

func (cgc *ChatGroupController) DeleteUsersGroup(c *gin.Context) {
	var req model.InputDeleteUsersGroup
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	err := cgc.chatGroup.DeleteUserGrou(req.PhoneNumber, req.GroupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func NewChatGroupController(r *gin.Engine, chatGroup usecase.ChatGroupUsecase) *ChatGroupController {
	controller := ChatGroupController{
		router:    r,
		chatGroup: chatGroup,
	}

	group := r.Group("/group")

	group.POST("/create/:phone_number", controller.CreateGroup)
	group.POST("/:group_id/:admin_number", controller.AddUserGroup)
	group.DELETE("/delete", controller.DeleteUsersGroup)

	return &controller
}
