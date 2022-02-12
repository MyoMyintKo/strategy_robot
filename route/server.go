package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myomyintko/strategy_robot/config"
	"github.com/myomyintko/strategy_robot/controller"
	"github.com/myomyintko/strategy_robot/middleware"
	"github.com/myomyintko/strategy_robot/repository"
	"github.com/myomyintko/strategy_robot/service"
	"gorm.io/gorm"
)

var (
	db *gorm.DB = config.SetupDatabaseConnection()
	// user and auth
	userRepository repository.UserRepository = repository.NewUserRepository(db)
	userService    service.UserService       = service.NewUserService(userRepository)
	userController controller.UserController = controller.NewUserController(userService, jwtService)
	authService    service.AuthService       = service.NewAuthService(userRepository)
	authController controller.AuthController = controller.NewAuthController(authService, jwtService)

	// robot
	robotRepository repository.RobotRepository = repository.NewRobotRepository(db)
	robotService    service.RobotService       = service.NewRobotService(robotRepository)
	robotController controller.RobotController = controller.NewRobotController(robotService, jwtService)

	// bind api
	apiRepository repository.APIRepository = repository.NewAPIRepository(db)
	apiService    service.APIService       = service.NewAPIService(apiRepository)
	apiController controller.APIController = controller.NewAPIController(apiService, jwtService)
	//binance
	binanceRepository repository.BinanceRepository = repository.NewBinanceRepository(db)
	binanceService    service.BinanceService       = service.NewBinanceService(binanceRepository)
	binanceController controller.BinanceController = controller.NewBinanceController(binanceService, jwtService)
	// jwt
	jwtService service.JWTService = service.NewJWTService()
)

func InitRoute() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()
	r.Use(Cors())
	//r.GET("ws",controller.TestKline)

	apiV1Routes := r.Group("/api/v1")

	authRoutes := apiV1Routes.Group("auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := apiV1Routes.Group("users", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
	}

	robotRoutes := apiV1Routes.Group("robots", middleware.AuthorizeJWT(jwtService))
	{
		robotRoutes.GET("/", robotController.FindByUserID)
		robotRoutes.POST("/", robotController.Insert)
		robotRoutes.PUT("/:id", robotController.Update)
		robotRoutes.DELETE("/:id", robotController.Delete)
	}

	binanceRoutes := apiV1Routes.Group("binance", middleware.AuthorizeJWT(jwtService))
	{
		binanceRoutes.GET("/get-bind", apiController.All)
		binanceRoutes.POST("/bind", apiController.Insert)
		binanceRoutes.PUT("/update-bind/:id", apiController.Update)
		binanceRoutes.DELETE("/unbind/:id", apiController.Delete)
		// get symbol
		binanceRoutes.GET("/", binanceController.GetSymbolInfo)
		//get coin
		binanceRoutes.GET("/getCoin", binanceController.GetCrypto)
		// order
		binanceRoutes.POST("/orders", binanceController.CreateOrder)
		binanceRoutes.GET("/orders/:id", binanceController.GetOrder)
		binanceRoutes.GET("/orders", binanceController.ListOrders)
		binanceRoutes.DELETE("/orders/:id", binanceController.CancelOrder)
		binanceRoutes.GET("/openOrders", binanceController.ListOpenOrders)
		binanceRoutes.GET("/wsOrders", binanceController.WsListOrdes)
		// kline
		binanceRoutes.GET("/wsKline", binanceController.WsListKline)
		// account
		binanceRoutes.GET("/account", binanceController.GetAccount)
		// stream
		binanceRoutes.POST("/stream", binanceController.StartUserStream)
		binanceRoutes.PUT("/stream", binanceController.KeepAliveUserStream)
	}
	panic(r.Run(":5000"))
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,Origin,X-Requested-With,Content-Type,Accept")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		c.Next()
	}
}
