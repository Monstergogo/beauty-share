package util

const (
	MinioNetProtocol = "http"
	GinServerPort    = 5008
	GrpcServerPort   = 5018
	GrpcServiceName  = "share.ShareService"
)

type HttpMethod string

const (
	HttpMethodGet    = HttpMethod("GET")
	HttpMethodPost   = HttpMethod("POST")
	HttpMethodPut    = HttpMethod("PUT")
	HttpMethodDelete = HttpMethod("DELETE")
	HttpMethodPatch  = HttpMethod("PATCH")
)

const CtxTraceID = "traceId"

const FileUploadRequestName = "files"

const (
	MongoShareDBName       = "share"
	MongoShareCollectName  = "shares"
	MongoUriDataID         = "mongo_uri"
	MinioEndpointDataID    = "minio_endpoint"
	MinioIDDataID          = "minio_id"
	MinioSecretDataID      = "minio_secret"
	MinioShareBucketDataID = "share_bucket_name"
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
