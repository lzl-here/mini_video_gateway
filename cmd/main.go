package main

import (
	"context"
	"flag"
	"gateway/internal/middlewares"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/network/netpoll"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/gzip"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/swagger"
	config "github.com/lzl-here/mini-video-common/config"
	swaggerFiles "github.com/swaggo/files"
	// upb "github.com/lzl-here/mini-video-common/protobuf/grpc/user"
)

func main() {

	appConfig := loadConfig()

	h := server.New(
		server.WithHostPorts(":"+strconv.Itoa(appConfig.GatewayHttpPort)),
		server.WithTransport(netpoll.NewTransporter))
	// 网关不需要repo
	//initRepo(appConfig)
	initLog(appConfig)
	initMiddlewares(h)
	initRouter(h)
	// 自动注册到注册中心，并且使用signalWaiter监控服务异常
	h.Spin()
}

// 初始化配置
func loadConfig() *config.AppGatewayConfig {
	// 从环境变量中加载配置
	envPath := flag.String("env", "../.env.production", "环境变量文件路径")
	appConfig := config.LoadGateway(*envPath)
	return appConfig
}

// 初始化日志
func initLog(appConfig *config.AppGatewayConfig) {
	var l = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,            // 记录日志位置
		Level:       slog.LevelDebug, // 设置日志级别
		ReplaceAttr: nil,
	}))
	slog.SetDefault(l)
}

func initMiddlewares(h *server.Hertz) {

	h.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	h.Use(middlewares.JWTAuthMiddleware())
	h.Use(gzip.Gzip(gzip.DefaultCompression))

}

// 初始化路由
func initRouter(h *server.Hertz) {
	url := swagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	h.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, url))
	// ping
	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "pong")
	})

	// services
}
