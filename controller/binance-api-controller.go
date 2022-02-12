package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/myomyintko/strategy_robot/dto"
	"github.com/myomyintko/strategy_robot/helper"
	"github.com/myomyintko/strategy_robot/model"
	"github.com/myomyintko/strategy_robot/service"
)

var (
	keys []model.BinanceAPI
	key  model.BinanceAPI
)

type APIController interface {
	All(context *gin.Context)
	FindByID(context *gin.Context)
	FindByUserID(context *gin.Context)
	Insert(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
}

type apiController struct {
	apiService service.APIService
	jwtService service.JWTService
}

func NewAPIController(apiServ service.APIService, jwtServ service.JWTService) APIController {
	return &apiController{
		apiService: apiServ,
		jwtService: jwtServ,
	}
}

func (c *apiController) All(context *gin.Context) {
	keys = c.apiService.All()
	res := helper.BuildResponse(true, "OK", keys)
	context.JSON(http.StatusOK, res)
}

func (c *apiController) FindByID(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	key = c.apiService.FindByID(id)
	fmt.Println(key)
	//if (api == model.BinanceAPIResponse{}) {
	//	res := helper.BuildErrorResponse("Data not found", "No data with given id", helper.EmptyObj{})
	//	context.JSON(http.StatusNotFound, res)
	//} else {
	//	res := helper.BuildResponse(true, "OK", api)
	//	context.JSON(http.StatusOK, res)
	//}
}

func (c *apiController) FindByUserID(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	key = c.apiService.FindByUserID(id)
	fmt.Println(key)
	//if (key == model.BinanceAPIResponse{}) {
	//	res := helper.BuildErrorResponse("Data not found", "No data with given id", helper.EmptyObj{})
	//	context.JSON(http.StatusNotFound, res)
	//} else {
	//	res := helper.BuildResponse(true, "OK", key)
	//	context.JSON(http.StatusOK, res)
	//}
}

func (c *apiController) Insert(context *gin.Context) {
	var apiCreateDTO dto.APICreateDTO
	errDTO := context.ShouldBind(&apiCreateDTO)
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
		response := helper.BuildErrorResponse("Token Error", errToken.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}
	convertedUserID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		response := helper.BuildErrorResponse("Parse Error", err.Error(), helper.EmptyObj{})
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	user := c.apiService.FindByUserID(convertedUserID)
	if user.ID != 0 {
		response := helper.BuildErrorResponse("User already bound keys", "", helper.EmptyObj{})
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	apiCreateDTO.UserID = convertedUserID
	apiCreateDTO.BoundAt = time.Now()
	result := c.apiService.Insert(apiCreateDTO)
	response := helper.BuildResponse(true, "OK", result)
	context.JSON(http.StatusCreated, response)
}

func (c *apiController) Update(context *gin.Context) {
	var apiUpdateDTO dto.APIUpdateDTO
	errDTO := context.ShouldBind(&apiUpdateDTO)
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

	if c.apiService.IsAllowedToEdit(userID, apiUpdateDTO.ID) {
		apiUpdateDTO.UserID = userID
		result := c.apiService.Update(apiUpdateDTO)
		response := helper.BuildResponse(true, "OK", result)
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildErrorResponse("You dont have permission", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}

func (c *apiController) Delete(context *gin.Context) {
	var key model.BinanceAPI
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		response := helper.BuildErrorResponse("Failed tou get id", "No param id were found", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	key.ID = id
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
	if !c.apiService.IsAllowedToEdit(userID, key.ID) {
		response := helper.BuildErrorResponse("You dont have permission", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
		return
	}
	c.apiService.Delete(key)
	res := helper.BuildResponse(true, "Deleted", helper.EmptyObj{})
	context.JSON(http.StatusOK, res)
}

func (c *apiController) getUserIDByToken(token string) (string, error) {
	aToken, err := c.jwtService.ValidateToken(token)
	if err != nil {
		return "", err
	}
	claims := aToken.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id, nil
}
