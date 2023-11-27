package util

const (
	MinioNetProtocol      = "http"
	GinServerPort         = 5008
	GrpcServerPort        = 5018
	GrpcServiceName       = "share.ShareService"
	HttpServiceName       = "share-service-http-metric"
	TracingServiceName    = "share-service"
	ServiceVersion        = "v0.1.2"
	OTLPCollectorGrpcAddr = "localhost:4317"
)

type HttpMethod string

const (
	HttpMethodGet    = HttpMethod("GET")
	HttpMethodPost   = HttpMethod("POST")
	HttpMethodPut    = HttpMethod("PUT")
	HttpMethodDelete = HttpMethod("DELETE")
	HttpMethodPatch  = HttpMethod("PATCH")
)

const CtxTraceID = "trace_id"

const FileUploadRequestName = "files"

const (
	MongoShareDBName       = "share"
	MongoShareCollectName  = "shares"
	MinioEndpointDataID    = "minio_endpoint"
	MinioShareBucketDataID = "share_bucket_name"
)

const (
	NacosConfPath = "conf/nacos.yaml"
	ConfPath      = "conf/conf.yaml"
)

const (
	NacosServerAddrEnvKey = "NACOS_SERVER_ADDR"
	NacosServerPortEnvKey = "NACOS_SERVER_PORT"
)

const (
	ConsulServerAddrEnvKey = "CONSUL_SERVER_ADDR"
	ConsulConfigPrefix     = "beauty-share"
	ConsulConfigMinio      = "beauty-share/minio"
	ConsulConfigMongo      = "beauty-share/mongo_uri"
)

var RouterFilter = map[string]bool{
	"/metrics": true,
	"/v1/ping": true,
}
