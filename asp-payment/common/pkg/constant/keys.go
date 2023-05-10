package constant

import (
	"time"

	"github.com/go-redis/redis/v7"
)

// redis key

// redis key注意使用命名空间的方式,方便查询和拆分

const (
	Prefix = "asp:" // 项目key前缀

	SnPrefix        = "snRound:" // 自动生成 sn 时间戳前缀
	NoExpiredPrefix = "app:"     // 系统一直固定在缓存中 没有过期时间

	KEY_APPID_STRING                                           = "appid_%v"                                             // 商户的key %idnum%
	KEY_MERCHANT_PROJECT_CONFIG_INFO_STRING                    = "merchant_project_config_info_%v"                      // cp项目配置的key %mch_project_id%
	KEY_MERCHANT_PROJECT_CURRENCY_INFO_STRING                  = "merchant_project_currency_info_%v"                    // cp项目币种金额相关配置的key %mch_project_id%
	KEY_DEPART_LIST_STRING                                     = "depart_list"                                          // 内部商户列表的缓存
	KEY_CHANNEL_LIST_STRING                                    = "channel_list"                                         // 渠道的缓存
	KEY_DEPART_INFO_STRING                                     = "depart_info_%v"                                       // 商户详情的key %id%
	KEY_MERCHANT_PROJECT_INFO_STRING                           = "merchant_project_info_%v"                             // cp项目详情的key %idnum%
	KEY_CHANNEL_DEPART_PROVIDER_TRADETYPE_LIST                 = "channel_depart_tradetype_list_%v_%v_%v_%v"            // 商户渠道列表缓存的key %merchantProjectId% %currency% %provider% %trade_type%
	KEY_CHANNEL_DEPART_INFO_STRING                             = "channel_depart_info_%v_%v"                            // 商户渠道详情的key %channelids%  %departid%
	KEY_DEPART_MERCHANT_PROJECT_LINK_LIST_STRING               = "depart_merchant_project_link_list_%v"                 // 内部商户&cp项目关联表 %merchantProjectId%
	Nil                                                        = redis.Nil                                              // redis 判断是否为空的 Nil
	EXPIRED_TIME                                               = time.Second * 1800                                     // 过期时间 60 秒
	NO_EXPIRED_TIME                                            = time.Second * -1                                       // 过期时间 60 秒
	KEY_ORDER_NOTIFY_QUEUE                                     = "notify:order_notify_queue"                            // 加入缓存队列 发送下游回调通知
	KEY_ORDER_NOTIFY_MANUAL_QUEUE                              = "notify_manual:order_notify_queue"                     // 手动的加入缓存队列 发送下游回调通知
	KEY_MERCHANT_PROJECT_AMOUNT_QUEUE                          = "merchant_project:amount_queue"                        // 加入缓存队列 账户金额变动队列
	KEY_MERCHANT_PROJECT_PRE_FLOW_QUEUE                        = "merchant_project:pre_flow_queue"                      // 加入缓存队列 账户金额预扣款变动队列
	KEY_SYSTEM_CONFIG_PAYOUT_STATUS                            = "system_config:system_payout_status"                   // 系统代付开关
	KEY_MERCHANT_PROJECT_CHANNEL_DEPART_TRADE_TYPE_INFO_STRING = "merchant_project_channel_depart_trade_type_info%v_%v" // cp渠道内部商户支付方式信息 %merchantProjectId%trade_type
	KEY_CHANNEL_DEPART_TRADE_TYPE_INFO_STRING                  = "channel_depart_trade_type_info_%v_%v_%v"              // 渠道内部商户支付方式信息 %channel_id%depart_id%trade_type
	KEY_PAYOUT_AUDIT_LIST                                      = "payout_checkout_list"                                 // 代付审核队列key
)

// GetRedisKey 给redis key加上前缀
func GetRedisKey(key string) string {
	return Prefix + key
}

// GetNoExpiredRedisKey 给redis key加上前缀
func GetNoExpiredRedisKey(key string) string {
	return NoExpiredPrefix + key
}
