package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"userVault/config"
	"userVault/internal/domain"
	"userVault/internal/port"
	"userVault/internal/utils"

	aesCipher "userVault/internal/adapters/cipher/aes"
	httpHandler "userVault/internal/adapters/handler/http/v1"
	zapLogger "userVault/internal/adapters/logger/zapLogger"
	authMiddleware "userVault/internal/adapters/middleware/auth"
	mysqlRepo "userVault/internal/adapters/repository/mysql"
	ginRouter "userVault/internal/adapters/router/gin"
	jwtToken "userVault/internal/adapters/token/jwt"
	userUsecase "userVault/internal/usecase/user"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println("failed to load env:", err)
		return
	}
	configPath := os.Getenv("CONFIG_FILE_PATH")
	configName := os.Getenv("CONFIG_FILE_NAME")
	configType := os.Getenv("CONFIG_FILE_TYPE")

	appConfig, err := config.StartConfig(configPath, config.File{
		Name: configName,
		Ext:  configType,
	})

	if err != nil {
		log.Println("failed to load config:", err)
		return
	}

	logger, err := initLogger(appConfig.GetLogger())
	if err != nil {
		log.Println("failed to initialize logger:", err)
		return
	}

	db, err := initDatabase(appConfig)
	if err != nil {
		log.Println("failed to connect to database:", err)
		return
	}
	db.AutoMigrate()

	userService := userUsecase.New(logger, db, appConfig.GetAppName(), appConfig.GetUser())
	services := domain.List{User: userService}

	tokenService, err := initTokenManager()
	if err != nil {
		log.Println("failed to setup token manager:", err)
		return
	}

	router := setupRouter(appConfig, logger, tokenService, services)
	port := appConfig.GetAppPort()
	logger.Infow(context.Background(), "Starting server", "port", port)
	logger.Sync(context.Background())

	if err := router.StartServer(port); err != nil {
		logger.Errorw(context.Background(), "Server stopped with error", "port", port, "error", err)
		logger.Sync(context.Background())
		return
	}

	logger.Infow(context.Background(), "Server stopped", "port", port)
	logger.Sync(context.Background())
}

func initLogger(conf config.Logger) (port.Logger, error) {
	loggerConf := zapLogger.Config{
		Level:           conf.GetLoggerLevel(),
		Encoding:        conf.GetLoggerEncodingMethod(),
		EncodingCaller:  conf.GetLoggerEncodingCaller(),
		OutputPath:      conf.GetLoggerPath(),
		ErrorOutputPath: conf.GetLoggerErrorPath(),
	}
	return zapLogger.New(loggerConf)
}

func initDatabase(conf config.App) (port.RepositoryMySQL, error) {
	cipherKey := os.Getenv("CIPHER_CRYPTO_KEY")
	cipher := aesCipher.New(cipherKey)

	hostEnc, portEnc, userEnc, passEnc, dbName, prefix := conf.GetStoreDatabaseProperties()

	host, err := cipher.Decrypt(hostEnc)
	if err != nil {
		return nil, err
	}
	portVal, err := cipher.Decrypt(portEnc)
	if err != nil {
		return nil, err
	}
	user, err := cipher.Decrypt(userEnc)
	if err != nil {
		return nil, err
	}
	pass, err := cipher.Decrypt(passEnc)
	if err != nil {
		return nil, err
	}

	return mysqlRepo.New(host, portVal, user, pass, dbName, prefix)
}

func initTokenManager() (port.Token, error) {
	method := os.Getenv("JWT_METHOD")
	hmacKey := os.Getenv("JWT_HMAC_KEY")
	privateKeyPath := os.Getenv("JWT_RSA_PRIVATE_KEY_PATH")
	publicKeyPath := os.Getenv("JWT_RSA_PUBLIC_KEY_PATH")

	privateKey, err := utils.LoadRSAPrivKeyFromFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	publicKey, err := utils.LoadRSAPubKeyFromFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	return jwtToken.New(method, []byte(hmacKey), privateKey, publicKey), nil
}

func setupRouter(conf config.App, logger port.Logger, token port.Token, services domain.List) port.Router {
	apiKeys := conf.GetMiddlewareApiKeys()
	middleware := authMiddleware.New(apiKeys, token)
	handler := httpHandler.New(logger, token, services)
	apiConf := conf.GetApi()

	router := ginRouter.New()
	publicRoutes := router.NewGroup("")
	publicRoutes.UseBefore(middleware.ValidateApiKey())

	if apiConf.GetUserLoginEnabled() {
		method, route := apiConf.GetUserLoginProperties()
		publicRoutes.RegisterRoute(method, route, handler.UserLogin)
	}
	if apiConf.GetUserRegisterEnabled() {
		method, route := apiConf.GetUserRegisterProperties()
		publicRoutes.RegisterRoute(method, route, handler.UserRegister)
	}
	if apiConf.GetUserActivationEnabled() {
		method, route := apiConf.GetUserActivationProperties()
		publicRoutes.RegisterRoute(method, route, handler.UserActivation)
	}
	if apiConf.GetUserResendActivationEnabled() {
		method, route := apiConf.GetUserResendActivationProperties()
		publicRoutes.RegisterRoute(method, route, handler.UserResendActivation)
	}
	if apiConf.GetUserForgotPasswordEnabled() {
		method, route := apiConf.GetUserForgotPasswordProperties()
		publicRoutes.RegisterRoute(method, route, handler.UserForgotPassword)
	}
	if apiConf.GetUserPasswordResetEnabled() {
		method, route := apiConf.GetUserPasswordResetProperties()
		publicRoutes.RegisterRoute(method, route, handler.UserPasswordReset)
	}
	if apiConf.GetUserRefreshTokenValidateEnabled() {
		method, route := apiConf.GetUserRefreshTokenValidateProperties()
		publicRoutes.RegisterRoute(method, route, handler.UserRefreshTokenValidate)
	}

	protectedRoutes := router.NewGroup("")
	protectedRoutes.UseBefore(middleware.ValidateAccessToken())

	if apiConf.GetUserLogoutEnabled() {
		method, route := apiConf.GetUserLogoutProperties()
		protectedRoutes.RegisterRoute(method, route, handler.UserLogout)
	}

	return router
}
