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
	"github.com/myomyintko/strategy_robot/model"
	"github.com/myomyintko/strategy_robot/service"
)

var (
	robot  model.Robot
	robots []model.Robot
)

type RobotController interface {
	All(context *gin.Context)
	FindByID(context *gin.Context)
	FindByUserID(context *gin.Context)
	Insert(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
}

type robotController struct {
	robotService service.RobotService
	jwtService   service.JWTService
}

func NewRobotController(robotServ service.RobotService, jwtServ service.JWTService) RobotController {
	return &robotController{
		robotService: robotServ,
		jwtService:   jwtServ,
	}
}

func (c *robotController) All(context *gin.Context) {
	robots = c.robotService.All()
	res := helper.BuildResponse(true, "OK", robots)
	context.JSON(http.StatusOK, res)
}

func (c *robotController) FindByID(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	robot = c.robotService.FindByID(id)
	if robot.ID == 0 {
		res := helper.BuildResponse(true, "No robot yet!", helper.EmptyObj{})
		context.JSON(http.StatusOK, res)
		return
	}
	res := helper.BuildResponse(true, "OK", robot)
	context.JSON(http.StatusOK, res)
}

func (c *robotController) FindByUserID(context *gin.Context) {

	authHeader := context.GetHeader("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	authHeader = splitToken[1]

	userID, errToken := c.getUserIDByToken(authHeader)
	if errToken != nil {
		res := helper.BuildErrorResponse("Failed to process request", errToken.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}
	convertedUserID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		response := helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}
	robot = c.robotService.FindByUserID(convertedUserID)

	if robot.ID == 0 {
		res := helper.BuildResponse(true, "No robot yet!", helper.EmptyObj{})
		context.JSON(http.StatusOK, res)
		return
	}
	res := helper.BuildResponse(true, "OK", robot)
	context.JSON(http.StatusOK, res)
}

func (c *robotController) Insert(context *gin.Context) {
	var robotCreateDTO dto.RobotCreateDTO
	errDTO := context.ShouldBind(&robotCreateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}
	authHeader := context.GetHeader("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	authHeader = splitToken[1]
	userID, errToken := c.getUserIDByToken(authHeader)
	if errToken != nil {
		res := helper.BuildErrorResponse("Failed to process request", errToken.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}
	convertedUserID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		res := helper.BuildErrorResponse("Parse Error", err.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}
	if c.robotService.IsUserExistRobot(robotCreateDTO.Symbol, convertedUserID) {
		response := helper.BuildErrorResponse("Failed to process request", "Duplicate Robot", helper.EmptyObj{})
		context.JSON(http.StatusConflict, response)
		return
	}
	robotCreateDTO.UserID = convertedUserID
	result := c.robotService.Insert(robotCreateDTO)
	response := helper.BuildResponse(true, "OK", result)
	context.JSON(http.StatusCreated, response)
}

func (c *robotController) Update(context *gin.Context) {
	var robotUpdateDTO dto.RobotUpdateDTO
	errDTO := context.ShouldBind(&robotUpdateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}
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
	userID, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		response := helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}
	robotID, ParamErr := strconv.ParseUint(context.Param("id"), 10, 64)
	if ParamErr != nil {
		response := helper.BuildErrorResponse("Param error", ParamErr.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}

	if c.robotService.IsAllowedToEdit(userID, robotID) {
		robotUpdateDTO.ID = robotID
		robotUpdateDTO.UserID = userID
		result := c.robotService.Update(robotUpdateDTO)
		response := helper.BuildResponse(true, "OK", result)
		context.JSON(http.StatusOK, response)
		return
	} else {
		response := helper.BuildErrorResponse("You dont have permission", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
		return
	}
}

func (c *robotController) Delete(context *gin.Context) {
	var robot model.Robot
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		response := helper.BuildErrorResponse("Failed tou get id", "No param id were found", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	robot.ID = id
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
	userID, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		response := helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}
	if !c.robotService.IsAllowedToEdit(userID, robot.ID) {
		response := helper.BuildErrorResponse("You dont have permission", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
		return
	}
	c.robotService.Delete(robot)
	res := helper.BuildResponse(true, "Deleted", helper.EmptyObj{})
	context.JSON(http.StatusAccepted, res)
}

func (c *robotController) getUserIDByToken(token string) (string, error) {
	aToken, err := c.jwtService.ValidateToken(token)
	if err != nil {
		return "", err
	}
	claims := aToken.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id, nil
}
