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

type mongoConf struct {
	Uri         string `json:"uri"`
	DBName      string `json:"db-name"`
	CollectName string `json:"collect-name"`
}

type severConf struct {
	Log           Log       `json:"log" yaml:"log"`
	Consul        consul    `json:"consul" yaml:"consul"`
	CosBucketName string    `json:"cos-bucket-name"`
	Minio         minioConf `json:"minio"`
	Mongo         mongoConf `json:"mongo"`
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
		switch d.Key {
		case util.ConsulConfigBucketName:
			ServerConf.CosBucketName = string(decodeValue)
		case util.ConsulConfigMinio:
			err = json.Unmarshal(decodeValue, &ServerConf.Minio)
			if err != nil {
				return err
			}
		case util.ConsulConfigMongo:
			err = json.Unmarshal(decodeValue, &ServerConf.Mongo)
			if err != nil {
				return err
			}
		default:
			continue
		}
	}
	return err
}
