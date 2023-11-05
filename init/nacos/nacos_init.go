package nacos

import (
	"github.com/Monstergogo/beauty-share/init/logger"
	"github.com/Monstergogo/beauty-share/util"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
)

type Config struct {
	TimeoutMs           uint64 `json:"timeout_ms" yaml:"timeout_ms"`
	NotLoadCacheAtStart bool   `json:"not_load_cache_at_start" yaml:"not_load_cache_at_start"`
	LogDir              string `json:"log_dir" yaml:"log_dir"`
	CacheDir            string `json:"cache_dir" yaml:"cache_dir"`
	LogLever            string `json:"log_lever" yaml:"log_lever"`
	ServerAddr          string `json:"server_addr" yaml:"server_addr"`
	Port                uint64 `json:"port" yaml:"port"`
}

// GetConfFromYaml
// @desc 从yaml文件中读取配置
func GetConfFromYaml(path string, config interface{}) error {
	dataBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(dataBytes, config)
	return err
}

var (
	namingClient naming_client.INamingClient
	configClient config_client.IConfigClient
)

// InitNacos 初始化nacos
func InitNacos() {
	var (
		nacosConf Config
		err       error
	)
	if err = GetConfFromYaml(util.NacosConfPath, &nacosConf); err != nil {
		panic(err)
	}
	nacosConf.ServerAddr = os.Getenv(util.NacosServerAddrEnvKey)
	if nacosConf.ServerAddr == "" {
		nacosConf.ServerAddr = "127.0.0.1"
	}
	nacosServerPort := os.Getenv(util.NacosServerPortEnvKey)
	if nacosServerPort == "" {
		nacosConf.Port = 8848
	} else {
		port, err := strconv.ParseUint(nacosServerPort, 10, 64)
		if err != nil {
			panic(err)
		}
		nacosConf.Port = port
	}

	clientConfig := *constant.NewClientConfig(
		constant.WithTimeoutMs(nacosConf.TimeoutMs),
		constant.WithNotLoadCacheAtStart(nacosConf.NotLoadCacheAtStart),
		constant.WithLogDir(nacosConf.LogDir),
		constant.WithCacheDir(nacosConf.CacheDir),
		constant.WithLogLevel(nacosConf.LogLever),
	)

	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(
			nacosConf.ServerAddr,
			nacosConf.Port,
		),
	}
	// 创建服务发现客户端
	namingClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(err)
	}

	// 创建动态配置客户端
	configClient, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(err)
	}
	logger.GetLogger().Info("nacos init success")
}

func GetNacosNamingClient() naming_client.INamingClient {
	return namingClient
}

func GetNacosConfigClient() config_client.IConfigClient {
	return configClient
}
