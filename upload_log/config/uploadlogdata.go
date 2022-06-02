package config

type UploadLogData struct {
	ServerName         string `mapstructure:"server-name" json:"server-name" yaml:"server-name"`
	BizRegExp          string `mapstructure:"biz-reg-exp" json:"biz-reg-exp" yaml:"biz-reg-exp"`
	BizSearchDir       string `mapstructure:"biz-search-dir" json:"biz-search-dir" yaml:"biz-search-dir"`
	BizSystemName      string `mapstructure:"biz-system-name" json:"biz-system-name" yaml:"biz-system-name"`
	OssAccessKeyID     string `mapstructure:"oss-access-key-id" json:"oss-access-key-id" yaml:"oss-access-key-id"`
	OssAccessKeySecret string `mapstructure:"oss-access-key-secret" json:"oss-access-key-secret" yaml:"oss-access-key-secret"`
	OssBucket          string `mapstructure:"oss-bucket" json:"oss-bucket" yaml:"oss-bucket"`
	OssEndpoint        string `mapstructure:"oss-endpoint" json:"oss-endpoint" yaml:"oss-endpoint"`
	OssObjectPrefix    string `mapstructure:"oss-object-prefix" json:"oss-object-prefix" yaml:"oss-object-prefix"`
}
