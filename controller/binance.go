package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/myomyintko/strategy_robot/config"
	"github.com/myomyintko/strategy_robot/dto"
	"github.com/myomyintko/strategy_robot/helper"
	"github.com/myomyintko/strategy_robot/repository"
	"github.com/myomyintko/strategy_robot/service"
)

var client *binance.Client

func getBindKey(userId uint64) {
	user := service.NewAPIService(repository.NewAPIRepository(config.SetupDatabaseConnection()))
	res := user.FindByUserID(userId)
	apiKey := res.APIKey
	secretKey := res.SecretKey
	client = binance.NewClient(apiKey, secretKey)
}

type BinanceController interface {
	StartUserStream(context *gin.Context)
	KeepAliveUserStream(context *gin.Context)
	GetSymbolInfo(context *gin.Context)
	GetCrypto(context *gin.Context)
	CreateOrder(context *gin.Context)
	GetOrder(context *gin.Context)
	CancelOrder(context *gin.Context)
	ListOpenOrders(context *gin.Context)
	WsListOrdes(context *gin.Context)
	ListOrders(context *gin.Context)
	WsListKline(context *gin.Context)
	GetAccount(context *gin.Context)
}

type binanceController struct {
	binanceService service.BinanceService
	jwtService     service.JWTService
}

func NewBinanceController(binSer service.BinanceService, jwtSer service.JWTService) BinanceController {
	return &binanceController{
		binanceService: binSer,
		jwtService:     jwtSer,
	}
}

func (c *binanceController) StartUserStream(ctx *gin.Context) {
	var bindStreamDTO dto.BindStreamDTO
	authHeader := ctx.GetHeader("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	authHeader = splitToken[1]
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		response := helper.BuildErrorResponse("Token Error", errToken.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userID, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		response := helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	getBindKey(userID)
	res, err := client.NewStartUserStreamService().Do(context.Background())
	if err != nil {
		response := helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	if res == "" {
		response := helper.BuildErrorResponse("Parameter 'listenKey' was empty.", "listenKey", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	bindStreamDTO.UserID = userID
	bindStreamDTO.StreamKey = res
	bindStreamDTO.StreamedAt = time.Now()
	api := service.NewAPIService(repository.NewAPIRepository(config.SetupDatabaseConnection()))
	result := api.BindStream(bindStreamDTO)
	response := helper.BuildResponse(true, "Stream Successfully", result)
	ctx.JSON(http.StatusOK, response)
}

func (c *binanceController) KeepAliveUserStream(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	authHeader = splitToken[1]
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		response := helper.BuildErrorResponse("Token Error", errToken.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userID, UserIdErr := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if UserIdErr != nil {
		response := helper.BuildErrorResponse("error", UserIdErr.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	getBindKey(userID)

	streamKey := c.getStreamKey(userID)
	if streamKey == "" {
		response := helper.BuildErrorResponse("Stream Key was Empty", "", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	stream := client.NewKeepaliveUserStreamService().ListenKey(streamKey).Do(context.Background())
	response := helper.BuildResponse(true, "User Stream Keep Alive Success", stream)
	ctx.JSON(http.StatusBadRequest, response)
}

func (c *binanceController) GetSymbolInfo(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	authHeader = splitToken[1]
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		response := helper.BuildErrorResponse("Token Error", errToken.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userID, UserIDErr := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if UserIDErr != nil {
		response := helper.BuildErrorResponse("error", UserIDErr.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	getBindKey(userID)
	var symbol = ctx.Query("symbol")
	if symbol == "" {
		response := helper.BuildErrorResponse("Symbol was empty", "Param was error", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
	}
	res, err := client.NewDepthService().Symbol(symbol).
		Do(context.Background())
	if err != nil {
		response := helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.BuildResponse(true, "Symbol Info of "+symbol, res)
	ctx.JSON(http.StatusOK, response)
}

func (c *binanceController) GetCrypto(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	authHeader = splitToken[1]
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		response := helper.BuildErrorResponse("Token Error", errToken.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userID, UserIDErr := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if UserIDErr != nil {
		response := helper.BuildErrorResponse("error", UserIDErr.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	getBindKey(userID)
	crypto, err := client.NewGetDepositAddressService().Coin("BTC").Do(context.Background())
	if err != nil {
		response := helper.BuildErrorResponse("NewListDepositsService error", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	fmt.Println(crypto)
	response := helper.BuildResponse(true, "Deposit Address for "+crypto.Coin, crypto)
	ctx.JSON(http.StatusBadRequest, response)
}

func (c *binanceController) CreateOrder(ctx *gin.Context) {
	var orderCreateDTO dto.CreateOrderDTO
	errDTO := ctx.ShouldBind(&orderCreateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
	}
	authHeader := ctx.GetHeader("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	authHeader = splitToken[1]
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		response := helper.BuildErrorResponse("Token Error", errToken.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userID, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		response := helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	getBindKey(userID)

	robotID, _ := strconv.ParseUint(ctx.Query("robot"), 10, 64)
	if CheckRobotByUser(robotID, userID) == "" {
		response := helper.BuildErrorResponse("error", "Invalid user or no robot", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	symbol := CheckRobotByUser(robotID, userID)
	sideType := ctx.Query("type")
	if sideType == "" {
		response := helper.BuildErrorResponse("Type is empty", "there is no param with type", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	order, err := client.NewCreateOrderService().Symbol(symbol).
		Side(binance.SideType(sideType)).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(orderCreateDTO.Quantity).
		Price(orderCreateDTO.Price).Do(context.Background())
	if err != nil {
		response := helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	orderCreateDTO.OrderId = order.OrderID
	orderCreateDTO.ClientOrderId = order.ClientOrderID
	orderCreateDTO.RobotID = robotID

	result := c.binanceService.Insert(orderCreateDTO)
	response := helper.BuildResponse(true, "Order was created successful", result)
	ctx.JSON(http.StatusCreated, response)
}

func (c *binanceController) GetOrder(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	authHeader = splitToken[1]
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		response := helper.BuildErrorResponse("Token Error", errToken.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userID, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		response := helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	getBindKey(userID)

	robotID, err := strconv.ParseUint(ctx.Query("robot"), 10, 64)
	if err != nil {
		response := helper.BuildErrorResponse("ParseUint error", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	if CheckRobotByUser(robotID, userID) == "" {
		response := helper.BuildErrorResponse("error", "Invalid user or no robot", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	symbol := CheckRobotByUser(robotID, userID)
	orderId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response := helper.BuildErrorResponse("ParseInt error", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	order, err := client.NewGetOrderService().Symbol(symbol).
		OrderID(orderId).Do(context.Background())
	if err != nil {
		response := helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.BuildResponse(true, "Order", order)
	ctx.JSON(http.StatusOK, response)
}

func (c *binanceController) CancelOrder(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	authHeader = splitToken[1]
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		response := helper.BuildErrorResponse("Token Error", errToken.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userID, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		response := helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	getBindKey(userID)

	robotID, _ := strconv.ParseUint(ctx.Query("robot"), 10, 64)
	if CheckRobotByUser(robotID, userID) == "" {
		response := helper.BuildErrorResponse("error", "Invalid user or no robot", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	symbol := CheckRobotByUser(robotID, userID)
	orderId, IdError := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if IdError != nil {
		response := helper.BuildErrorResponse("Invalid order id", IdError.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	_, binanceErr := client.NewCancelOrderService().Symbol(symbol).
		OrderID(orderId).Do(context.Background())
	if binanceErr != nil {
		response := helper.BuildErrorResponse("There is no order", binanceErr.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.BuildResponse(true, "Order cancel successful", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}

func (c *binanceController) ListOpenOrders(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	authHeader = splitToken[1]
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		response := helper.BuildErrorResponse("Token Error", errToken.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userID, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		response := helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	getBindKey(userID)

	robotID, _ := strconv.ParseUint(ctx.Query("robot"), 10, 64)
	if CheckRobotByUser(robotID, userID) == "" {
		response := helper.BuildErrorResponse("error", "Invalid user or no robot", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	symbol := CheckRobotByUser(robotID, userID)
	openOrders, err := client.NewListOpenOrdersService().Symbol(symbol).
		Do(context.Background())
	if err != nil {
		response := helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.BuildResponse(true, "Open orders", openOrders)
	ctx.JSON(http.StatusOK, response)
}

func (c *binanceController) ListOrders(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	authHeader = splitToken[1]
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		response := helper.BuildErrorResponse("Token Error", errToken.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userID, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		response := helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	getBindKey(userID)

	robotID, _ := strconv.ParseUint(ctx.Query("robot"), 10, 64)
	if CheckRobotByUser(robotID, userID) == "" {
		response := helper.BuildErrorResponse("error", "Invalid user or no robot", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	symbol := CheckRobotByUser(robotID, userID)
	orders, err := client.NewListOrdersService().Symbol(symbol).
		Do(context.Background())
	if err != nil {
		response := helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.BuildResponse(true, "List Orders", orders)
	ctx.JSON(http.StatusOK, response)
}

func (c *binanceController) WsListKline(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	authHeader = splitToken[1]
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		response := helper.BuildErrorResponse("Token Error", errToken.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userID, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		response := helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	getBindKey(userID)
	robotID, _ := strconv.ParseUint(ctx.Query("robot"), 10, 64)
	if CheckRobotByUser(robotID, userID) == "" {
		response := helper.BuildErrorResponse("error", "Invalid user or no robot", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	symbol := CheckRobotByUser(robotID, userID)

	var interval = ctx.Query("interval")

	wsKlineHandler := func(event *binance.WsKlineEvent) {
		fmt.Println(event)
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneC, _, err := binance.WsKlineServe(symbol, interval, wsKlineHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC
}

//WebSocket

//var upgrader = websocket.Upgrader{
//	ReadBufferSize:  4096,
//	WriteBufferSize: 4096,
//}
//
//func getKline(conn *websocket.Conn, symbol, interval string) {
//	wsKlineHandler := func(event *binance.WsKlineEvent) {
//		fmt.Println(event)
//		err := conn.WriteJSON(event)
//		if err != nil {
//			log.Println(err)
//			return
//		}
//		return
//	}
//	errHandler := func(err error) {
//		fmt.Println(err)
//	}
//	doneC, _, err := binance.WsKlineServe(symbol, interval, wsKlineHandler, errHandler)
//	if err != nil {
//		fmt.Println(err)
//	}
//	<-doneC
//
//}
//
//func reader(conn *websocket.Conn) {
//	for {
//		messageType, payload, err := conn.ReadMessage()
//		if err != nil {
//			log.Println(err)
//			return
//		}
//
//		fmt.Println("message from client: ", string(payload))
//		if err := conn.WriteMessage(messageType, payload); err != nil {
//			log.Println(err)
//			return
//		}
//	}
//}
//
//func TestKline(ctx *gin.Context) {
//	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
//
//	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	fmt.Println("New client is connected")
//
//	getKline(conn, "MATICUSDT", "15m")
//
//	reader(conn)
//}

func (c *binanceController) GetAccount(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	authHeader = splitToken[1]
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		response := helper.BuildErrorResponse("Token Error", errToken.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userID, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		response := helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	getBindKey(userID)

	res, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		response := helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.BuildResponse(true, "Account info", res)
	ctx.JSON(http.StatusOK, response)
}
func (c *binanceController) WsListOrdes(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	authHeader = splitToken[1]
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		response := helper.BuildErrorResponse("Token Error", errToken.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userID, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		response := helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	getBindKey(userID)

	robotID, _ := strconv.ParseUint(ctx.Query("robot"), 10, 64)
	if CheckRobotByUser(robotID, userID) == "" {
		response := helper.BuildErrorResponse("error", "Invalid user or no robot", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	symbol := CheckRobotByUser(robotID, userID)
	trades, err := client.NewAggTradesService().
		Symbol(symbol).StartTime(1508673256594).EndTime(1508673256595).
		Do(context.Background())
	if err != nil {
		response := helper.BuildErrorResponse("Ws list orders error", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.BuildResponse(true, "WS List Orders", trades)
	ctx.JSON(http.StatusOK, response)
}

func CheckRobotByUser(robotID, userId uint64) string {
	robot := service.NewRobotService(repository.NewRobotRepository(config.SetupDatabaseConnection()))
	res := robot.FindByID(robotID)
	robotSymbol := ""
	if res.UserID != userId {
		return robotSymbol
	}
	robotSymbol = res.Symbol
	return robotSymbol
}

func (c *binanceController) getStreamKey(userId uint64) string {
	key := service.NewAPIService(repository.NewAPIRepository(config.SetupDatabaseConnection()))
	res := key.FindByUserID(userId)
	if res.StreamKey == "" {
		return ""
	}
	return res.StreamKey
}
