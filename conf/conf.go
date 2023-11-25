package conf

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Monstergogo/beauty-share/init/logger"
	"github.com/Monstergogo/beauty-share/util"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Log struct {
	LogPath string `json:"log_path" yaml:"log_filepath"`
	ErrPath string `json:"err_path" yaml:"error_log_filepath"`
}

type consul struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

type minioConf struct {
	ID       string `json:"id"`
	Secret   string `json:"secret"`
	Endpoint string `json:"endpoint"`
}

type severConf struct {
	Log           Log       `json:"log" yaml:"log"`
	Consul        consul    `json:"consul" yaml:"consul"`
	CosBucketName string    `json:"cos-bucket-name"`
	Minio         minioConf `json:"minio"`
	MongoUri      string    `json:"mongo-uri" yaml:"mongo-uri"`
}

type consulKVStoreReadResp struct {
	CreateIndex int64  `json:"CreateIndex"`
	ModifyIndex int64  `json:"ModifyIndex"`
	LockIndex   int64  `json:"LockIndex"`
	Flags       int    `json:"Flags"`
	Key         string `json:"Key"`
	Value       string `json:"Value"`
}

var ServerConf *severConf

// Cache 缓存配置
var Cache sync.Map

func init() {
	ServerConf = &severConf{}
	if err := util.GetConfFromYaml(util.ConfPath, &ServerConf); err != nil {
		panic(err)
	}
	consulServerAddr := os.Getenv(util.ConsulServerAddrEnvKey)
	if consulServerAddr != "" {
		ServerConf.Consul.Endpoint = consulServerAddr
	}

	err := readKeyFromConsulKVStore(fmt.Sprintf("%s/v1/kv/%s?recurse=true", ServerConf.Consul.Endpoint, util.ConsulConfigPrefix))
	if err != nil {
		panic(err)
	}
}

func readKeyFromConsulKVStore(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		logger.GetLogger().Error("read key from consul err", zap.Any("err_msg", err))
		return err
	}

	body, _ := io.ReadAll(resp.Body)
	var respData []consulKVStoreReadResp
	json.Unmarshal(body, &respData)

	for _, d := range respData {
		decodeValue, _ := base64.StdEncoding.DecodeString(d.Value)
		Cache.Store(strings.Split(d.Key, "/")[1], decodeValue)
		switch d.Key {
		case util.ConsulConfigBucketName:
			ServerConf.CosBucketName = string(decodeValue)
		case util.ConsulConfigMinio:
			err = json.Unmarshal(decodeValue, &ServerConf.Minio)
			if err != nil {
				return err
			}
		case util.ConsulConfigMongo:
			ServerConf.MongoUri = string(decodeValue)
		default:
			continue
		}
	}
	return err
}