CREATE TABLE `prod_log_upload_result` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `file_name` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '服务器存放的文件名，包含根路径，形如：/data/demo.log',
  `oss_file_name` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '阿里云对象存储服务器(OSS)存放的文件名，不含oss_bucket名称',
  `m_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '文件内容最后修改时间',
  `bytes` bigint unsigned NOT NULL DEFAULT '0' COMMENT '文件大小，字节数',
  `sha1` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '文件哈希sha1值',
  `origin_status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '原文件状态 0:未处理, 1:已上传到OSS, 2:已经删除原始文件, 99:原始文件为空文件，字节数为零，无法上传到oss，也不需要删除',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `source` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '来源，是从哪个系统上传的',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='生产环境服务器日志上传记录表';