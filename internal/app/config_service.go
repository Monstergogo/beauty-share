package app

import (
	"context"
	"fmt"
	"github.com/Monstergogo/beauty-share/init/conf"
	"github.com/Monstergogo/beauty-share/init/logger"
	"github.com/Monstergogo/beauty-share/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"net/http"
)

type ConfigServiceImpl struct {
}

func (c *ConfigServiceImpl) readKeyFromConsul(ctx context.Context, key string) ([]byte, error) {
	reqUrl := fmt.Sprintf(fmt.Sprintf("%s/v1/kv/%s/%s?raw", conf.Consul.Endpoint, util.ConsulConfigPrefix, key))
	resp, err := http.Get(reqUrl)
	if err != nil {
		logger.LogWithTraceId(ctx, zapcore.ErrorLevel, "read key from consul err", zap.Any("err_msg", err))
		return nil, err
	}

	body, _ := io.ReadAll(resp.Body)
	return body, err
}

func (c *ConfigServiceImpl) GetStringValueByKeyName(ctx context.Context, key string) (string, error) {
	if value, ok := conf.Cache.Load(key); ok {
		return string(value.([]byte)), nil
	}
	resp, err := c.readKeyFromConsul(ctx, key)
	if err != nil {
		return "", err
	}
	return string(resp), err
}
