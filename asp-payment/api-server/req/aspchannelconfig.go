package req

import (
	"asp-payment/common/model"
)

type ChannelConfigList struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Appid         string `json:"appid"`
	Secret        string `json:"secret"`
	Key           string `json:"key"`
	Cert          string `json:"cert"`
	CertKey       string `json:"cert_key"`
	Currency      string `json:"currency"`
	ChannelConfig string `json:"channel_config"`
	SerialNo      string `json:"serial_no"`
	Status        int    `json:"status"`
	Sort          int    `json:"sort"`
}

func (s *ChannelConfigList) Generate(model *model.AspChannelConfig) {
	model.Id = s.Id
	model.Name = s.Name
	model.Appid = s.Appid
	model.Secret = s.Secret
	model.Key = s.Key
	model.Cert = s.Cert
	model.CertKey = s.CertKey
	model.Currency = s.Currency
	model.ChannelConfig = s.ChannelConfig
	model.SerialNo = s.SerialNo
	model.Status = s.Status
	model.Sort = s.Sort
}
