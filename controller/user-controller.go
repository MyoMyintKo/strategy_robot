package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/myomyintko/strategy_robot/dto"
	"github.com/myomyintko/strategy_robot/helper"
	"github.com/myomyintko/strategy_robot/service"
)

//UserController is a ....
type UserController interface {
	Update(context *gin.Context)
	Profile(context *gin.Context)
}

type userController struct {
	userService service.UserService
	jwtService  service.JWTService
}

func NewUserController(userService service.UserService, jwtService service.JWTService) UserController {
	return &userController{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (c *userController) Update(context *gin.Context) {
	var userUpdateDTO dto.UserUpdateDTO
	errDTO := context.ShouldBind(&userUpdateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	authHeader := context.GetHeader("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	authHeader = splitToken[1]

	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		res := helper.BuildErrorResponse("Token error", errToken.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	id, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		res := helper.BuildErrorResponse("Parse error", err.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}
	userUpdateDTO.ID = id
	user := c.userService.Update(userUpdateDTO)
	res := helper.BuildResponse(true, "OK!", user)
	context.JSON(http.StatusOK, res)
}

func (c *userController) Profile(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	authHeader = splitToken[1]
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		response := helper.BuildErrorResponse("Token Error", errToken.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	user := c.userService.Profile(id)
	fmt.Println(user)
	res := helper.BuildResponse(true, "OK", user)
	context.JSON(http.StatusOK, res)
}
