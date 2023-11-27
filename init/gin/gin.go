package gin

import (
	"github.com/Monstergogo/beauty-share/init/logger"
	"github.com/Monstergogo/beauty-share/init/tracing"
	"github.com/Monstergogo/beauty-share/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type HttpServerRouter struct {
	Path    string
	Method  util.HttpMethod
	Handler func(context *gin.Context)
}

var (
	ginServer  *gin.Engine
	apiCounter metric.Int64Counter
	histogram  metric.Int64Histogram
)

func Init() *gin.Engine {
	ginServer = gin.Default()
	ginServer.Use(metricMiddleware())
	ginServer.Use(reqLoggerMiddleware())
	ginServer.Use(otelgin.Middleware(util.TracingServiceName))

	initMetricKind()
	// prometheus metrics
	ginServer.GET("/metrics", gin.WrapH(promhttp.Handler()))
	return ginServer
}

func initMetricKind() {
	var err error
	apiCounter, err = tracing.Meter.Int64Counter(
		"gin_api_counter_total",
		metric.WithDescription("Number of API calls."),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		panic(err)
	}
	histogram, err = tracing.Meter.Int64Histogram(
		"gin_api_duration",
		metric.WithDescription("The duration of api execution."),
		metric.WithUnit("ms"),
	)
	if err != nil {
		panic(err)
	}
}

func reqLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.GetHeader(util.CtxTraceID)
		if traceId == "" {
			uuID, _ := uuid.NewRandom()
			traceId = uuID.String()
		}
		c.Set(util.CtxTraceID, traceId)
		reqBody, _ := c.GetRawData()

		if !util.RouterFilter[c.Request.RequestURI] {
			logger.LogWithTraceId(c, zapcore.InfoLevel, "req info", zap.Any("method", c.Request.Method), zap.Any("req_uri", c.Request.RequestURI),
				zap.Any("req_body", reqBody))
		}
		c.Next()
	}
}

func metricMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiCounter.Add(c, 1)
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		histogram.Record(c, duration.Milliseconds())
	}
}

func toGinHandler(handler func(*gin.Context) (interface{}, error)) func(*gin.Context) {
	return func(ctx *gin.Context) {
		res, err := handler(ctx)
		if err != nil {
			util.GinErrResponse(ctx, err)
			return
		}
		util.GinSuccessResponse(ctx, res)
	}
}

func InitRouter() {
	registerGinRouter(HttpServerRouter{
		Path:    "/v1/oss/upload",
		Method:  util.HttpMethodPost,
		Handler: toGinHandler(OssUpload),
	})
	registerGinRouter(HttpServerRouter{
		Path:    "v1/ping",
		Method:  util.HttpMethodGet,
		Handler: toGinHandler(Ping),
	})
}

func registerGinRouter(router HttpServerRouter) {
	switch router.Method {
	case util.HttpMethodGet:
		ginServer.GET(router.Path, router.Handler)
	case util.HttpMethodPost:
		ginServer.POST(router.Path, router.Handler)
	case util.HttpMethodPut:
		ginServer.PUT(router.Path, router.Handler)
	case util.HttpMethodDelete:
		ginServer.DELETE(router.Path, router.Handler)
	case util.HttpMethodPatch:
		ginServer.PATCH(router.Path, router.Handler)
	}
}
