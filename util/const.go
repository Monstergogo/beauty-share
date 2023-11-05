package util

const (
	BucketName       = "photos"
	MinioEndpoint    = "localhost:9000"
	MinioNetProtocol = "http"
	GinServerPort    = 5008
	GrpcServerPort   = 5018
	GrpcServiceName  = "ShareService"
)

type HttpMethod string

const (
	HttpMethodGet    = HttpMethod("GET")
	HttpMethodPost   = HttpMethod("POST")
	HttpMethodPut    = HttpMethod("PUT")
	HttpMethodDelete = HttpMethod("DELETE")
	HttpMethodPatch  = HttpMethod("PATCH")
)

const (
	CtxTraceID  = "traceId"
	MinioID     = "dEYjKPxNoRqeBbsmz8ui"
	MinioSecret = "R15wFBw328BqRBWPSVj5m0UdLcc1VZtz0Wpe63Gb"
)

const FileUploadRequestName = "files"

const (
	MongoURI              = "mongodb://root:root123@localhost:27017/share"
	MongoShareDBName      = "share"
	MongoShareCollectName = "shares"
)

const (
	LogPath = "logs/access.log"
	ErrPath = "logs/err.log"
)

const (
	NacosConfPath = "conf/nacos.yaml"
)

const (
	NacosServerAddrEnvKey = "NACOS_SERVER_ADDR"
	NacosServerPortEnvKey = "NACOS_SERVER_PORT"
)
