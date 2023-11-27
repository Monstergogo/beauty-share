package app

import (
	"encoding/json"
	"fmt"
	"github.com/Monstergogo/beauty-share/init/conf"
	"net/http"
	"strings"
)

type RegisterPayload struct {
	ID                string                 `json:"ID"`
	Name              string                 `json:"Name"`
	Tags              []string               `json:"Tags"`
	Address           string                 `json:"Address"`
	Port              int                    `json:"Port"`
	Meta              map[string]string      `json:"Meta"`
	EnableTagOverride bool                   `json:"EnableTagOverride"`
	Check             map[string]interface{} `json:"Check"`
	Weights           map[string]int         `json:"Weights"`
}

type ConsulServiceImpl struct {
}

func (c *ConsulServiceImpl) RegisterService(payload RegisterPayload) error {
	registerUrl := fmt.Sprintf("%s/v1/agent/service/register?replace-existing-checks=true", conf.Consul.Endpoint)
	payloadMarshal, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("PUT", registerUrl, strings.NewReader(string(payloadMarshal)))
	req.Header.Add("Content-Type", "application/json")
	_, err = http.DefaultClient.Do(req)
	return err
}

func (c *ConsulServiceImpl) DeregisterService(serviceName string) error {
	registerUrl := fmt.Sprintf("%s/v1/agent/service/deregister/%s", conf.Consul.Endpoint, serviceName)
	req, _ := http.NewRequest("PUT", registerUrl, nil)
	req.Header.Add("Content-Type", "application/json")
	_, err := http.DefaultClient.Do(req)
	return err
}
