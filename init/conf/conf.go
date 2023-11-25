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

type log struct {
	LogPath string `json:"log_path" yaml:"log_filepath"`
	ErrPath string `json:"err_path" yaml:"error_log_filepath"`
}

type consul struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

type mongo struct {
	Uri string `json:"uri" yaml:"mongo_uri"`
}

type minio struct {
	ID       string `json:"id"`
	Secret   string `json:"secret"`
	Endpoint string `json:"endpoint"`
	Bucket   string `json:"bucket"`
}

type consulKVStoreReadResp struct {
	CreateIndex int64  `json:"CreateIndex"`
	ModifyIndex int64  `json:"ModifyIndex"`
	LockIndex   int64  `json:"LockIndex"`
	Flags       int    `json:"Flags"`
	Key         string `json:"Key"`
	Value       string `json:"Value"`
}

type yamlConf struct {
	Log    *log    `json:"log" yaml:"log"`
	Consul *consul `json:"consul" yaml:"consul"`
}

// Cache 缓存配置
var (
	Cache  sync.Map
	Minio  *minio
	Mongo  *mongo
	Log    *log
	Consul *consul
)

func init() {
	conf := &yamlConf{}
	if err := util.GetConfFromYaml(util.ConfPath, &conf); err != nil {
		panic(err)
	}

	Log = conf.Log
	Consul = conf.Consul
	consulServerAddr := os.Getenv(util.ConsulServerAddrEnvKey)
	if consulServerAddr != "" {
		conf.Consul.Endpoint = consulServerAddr
	}

	err := readKeyFromConsulKVStore(fmt.Sprintf("%s/v1/kv/%s?recurse=true", conf.Consul.Endpoint, util.ConsulConfigPrefix))
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
		case util.ConsulConfigMinio:
			Minio = &minio{}
			err = json.Unmarshal(decodeValue, &Minio)
			if err != nil {
				return err
			}
		case util.ConsulConfigMongo:
			Mongo = &mongo{}
			Mongo.Uri = string(decodeValue)
		default:
			continue
		}
	}
	return err
}
