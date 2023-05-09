/*
 Navicat Premium Data Transfer

 Source Server         : inner
 Source Server Type    : MySQL
 Source Server Version : 80011
 Source Host           : 192.168.180.122:3306
 Source Schema         : share

 Target Server Type    : MySQL
 Target Server Version : 80011
 File Encoding         : 65001

 Date: 07/02/2023 18:01:03
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for admin_menu
-- ----------------------------
DROP TABLE IF EXISTS `admin_menu`;
CREATE TABLE `admin_menu`  (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '菜单id',
  `pid` int(11) NOT NULL DEFAULT 0 COMMENT '父级菜单',
  `path` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '路由地址',
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '路由名称',
  `redirect` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '重定向地址',
  `component` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '视图文件路径',
  `icon` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '菜单图标',
  `title` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '菜单标题',
  `active_menu` tinyint(4) NOT NULL DEFAULT 1 COMMENT '当前路由为详情页时，需要高亮的菜单',
  `is_link` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '1' COMMENT '是否外链 不为空为外链地址',
  `is_hide` tinyint(4) NOT NULL DEFAULT 1 COMMENT '是否隐藏',
  `is_full` tinyint(4) NOT NULL DEFAULT 1 COMMENT '是否全屏',
  `is_affix` tinyint(4) NOT NULL DEFAULT 1 COMMENT '是否固定在 tabs nav',
  `is_keep_alive` tinyint(4) NOT NULL DEFAULT 1 COMMENT '是否外链',
  `auth_btn` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '菜单对应的页面的按钮组标识 json数组',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 20 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '菜单列表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_menu
-- ----------------------------
INSERT INTO `admin_menu` VALUES (1, 0, '/home/index', 'home', '', '/home/index', 'HomeFilled', '首页', 0, '', 0, 0, 1, 1, '[]', 0, '2023-01-11 15:31:37', '2023-01-11 15:31:37');
INSERT INTO `admin_menu` VALUES (11, 0, '/system', 'system', '/system/accountManage', '', 'Tools', '系统管理', 0, '', 0, 0, 0, 1, '[]', 0, '2023-01-11 17:11:36', '2023-01-11 17:11:36');
INSERT INTO `admin_menu` VALUES (12, 11, '/system/accountManage', 'accountManage', '', '/system/accountManage/index', 'Menu', '账号管理', 0, '', 0, 0, 0, 1, '[{\"name\":\"create\",\"desc\":\"新增账号\"},{\"name\":\"edit\",\"desc\":\"编辑账号\"},{\"name\":\"delete\",\"desc\":\"删除账号\"},{\"name\":\"status\",\"desc\":\"切换状态\"},{\"name\":\"role\",\"desc\":\"设置角色\"},{\"name\":\"reset_pwd\",\"desc\":\"重置密码\"}]', 0, '2023-01-11 17:12:38', '2023-01-11 17:12:38');
INSERT INTO `admin_menu` VALUES (13, 11, '/system/roleManage', 'roleManage', '', '/system/roleManage/index', 'Menu', '角色管理', 0, '', 0, 0, 0, 1, '[{\"name\":\"create\",\"desc\":\"新增角色\"},{\"name\":\"edit\",\"desc\":\"编辑角色\"},{\"name\":\"delete\",\"desc\":\"删除角色\"}]', 0, '2023-01-11 17:13:14', '2023-01-11 17:13:14');
INSERT INTO `admin_menu` VALUES (14, 11, '/system/menuMange', 'menuMange', '', '/system/menuMange/index', 'Menu', '菜单管理', 0, '', 0, 0, 0, 1, '[{\"name\":\"create\",\"desc\":\"新增根菜单\"},{\"name\":\"edit\",\"desc\":\"编辑菜单\"},{\"name\":\"delete\",\"desc\":\"删除菜单\"}]', 0, '2023-01-11 17:13:54', '2023-01-11 17:13:54');
INSERT INTO `admin_menu` VALUES (15, 11, '/system/permissionsMange', 'permissionsMange', '', '/system/permissionsMange/index', 'Menu', '权限管理', 0, '', 0, 0, 0, 1, '[]', 0, '2023-01-11 17:14:33', '2023-01-11 17:14:33');
INSERT INTO `admin_menu` VALUES (16, 11, '/system/departmentManage', 'departmentManage', '', '/system/departmentManage/index', 'Menu', '部门管理', 0, '', 0, 0, 0, 1, '[]', 0, '2023-01-11 17:15:17', '2023-01-11 17:15:17');

-- ----------------------------
-- Table structure for admin_permissions
-- ----------------------------
DROP TABLE IF EXISTS `admin_permissions`;
CREATE TABLE `admin_permissions`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `pid` int(11) NOT NULL COMMENT '父级id',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '权限名称',
  `route_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '接口路由url',
  `sort` int(10) NOT NULL DEFAULT 0 COMMENT '排序值',
  `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 33 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'Api权限表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_permissions
-- ----------------------------
INSERT INTO `admin_permissions` VALUES (1, 0, '系统管理', '/admin/sys', 0, '2023-01-10 10:19:51', '2023-01-10 10:19:51');
INSERT INTO `admin_permissions` VALUES (2, 1, '用户管理', '/admin/user', 0, '2023-01-10 10:34:03', '2023-01-10 10:34:03');
INSERT INTO `admin_permissions` VALUES (3, 1, '角色管理', '/admin/role', 0, '2023-01-10 10:34:13', '2023-01-10 10:34:13');
INSERT INTO `admin_permissions` VALUES (4, 1, '菜单管理', '/admin/menu', 0, '2023-01-10 10:34:22', '2023-01-10 10:34:22');
INSERT INTO `admin_permissions` VALUES (5, 1, '权限管理', '/admin/permissions', 0, '2023-01-10 10:34:28', '2023-01-10 10:34:28');
INSERT INTO `admin_permissions` VALUES (6, 2, '账号列表', '/admin/account/list', 0, '2023-01-10 10:36:40', '2023-01-10 10:36:40');
INSERT INTO `admin_permissions` VALUES (7, 2, '新增账号', '/admin/account/create', 0, '2023-01-10 10:37:13', '2023-01-10 10:37:13');
INSERT INTO `admin_permissions` VALUES (8, 2, '删除账号', '/admin/account/delete', 0, '2023-01-10 10:37:29', '2023-01-10 10:37:29');
INSERT INTO `admin_permissions` VALUES (9, 2, '编辑账号', '/admin/account/edit', 0, '2023-01-10 10:37:39', '2023-01-10 10:37:39');
INSERT INTO `admin_permissions` VALUES (10, 2, '重置密码', '/admin/account/reset_pwd', 0, '2023-01-10 10:37:51', '2023-01-10 10:37:51');
INSERT INTO `admin_permissions` VALUES (11, 2, '切换状态', '/admin/account/change_status', 0, '2023-01-10 10:39:01', '2023-01-10 10:39:01');
INSERT INTO `admin_permissions` VALUES (12, 3, '角色列表', '/admin/role/list', 0, '2023-01-10 10:39:39', '2023-01-10 10:39:39');
INSERT INTO `admin_permissions` VALUES (13, 3, '新增角色', '/admin/role/create', 0, '2023-01-10 10:41:10', '2023-01-10 10:41:10');
INSERT INTO `admin_permissions` VALUES (14, 3, '编辑角色', '/admin/role/update', 0, '2023-01-10 10:41:15', '2023-01-10 10:41:15');
INSERT INTO `admin_permissions` VALUES (15, 3, '删除角色', '/admin/role/delete', 0, '2023-01-10 10:41:42', '2023-01-10 10:41:42');
INSERT INTO `admin_permissions` VALUES (16, 3, '设置权限', '/admin/role/set_permissions', 0, '2023-01-10 10:43:30', '2023-01-10 10:43:30');
INSERT INTO `admin_permissions` VALUES (17, 4, '菜单列表', '/admin/menu/list', 0, '2023-01-10 10:53:28', '2023-01-10 10:53:28');
INSERT INTO `admin_permissions` VALUES (18, 4, '新增菜单', '/admin/menu/create', 0, '2023-01-10 10:53:41', '2023-01-10 10:53:41');
INSERT INTO `admin_permissions` VALUES (19, 4, '编辑菜单', '/admin/menu/update', 0, '2023-01-10 10:53:52', '2023-01-10 10:53:52');
INSERT INTO `admin_permissions` VALUES (20, 4, '删除菜单', '/admin/menu/delete', 0, '2023-01-10 10:54:04', '2023-01-10 10:54:04');
INSERT INTO `admin_permissions` VALUES (21, 5, '权限列表', '/admin/permissions/list', 0, '2023-01-10 10:54:48', '2023-01-10 10:54:48');
INSERT INTO `admin_permissions` VALUES (22, 5, '新增权限', '/admin/permissions/create', 0, '2023-01-10 10:55:03', '2023-01-10 10:55:03');
INSERT INTO `admin_permissions` VALUES (23, 5, '编辑权限', '/admin/permissions/update', 0, '2023-01-10 10:55:18', '2023-01-10 10:55:18');
INSERT INTO `admin_permissions` VALUES (24, 5, '删除权限', '/admin/permissions/delete', 0, '2023-01-10 10:55:33', '2023-01-10 10:55:33');

-- ----------------------------
-- Table structure for admin_role
-- ----------------------------
DROP TABLE IF EXISTS `admin_role`;
CREATE TABLE `admin_role`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '角色名称',
  `create_time` datetime(0) NULL DEFAULT NULL,
  `update_time` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_role
-- ----------------------------
INSERT INTO `admin_role` VALUES (1, '超级管理员', '2023-01-12 10:23:56', '2023-01-12 10:29:41');
INSERT INTO `admin_role` VALUES (2, '管理员', '2023-01-12 14:51:59', '2023-01-12 14:51:59');
INSERT INTO `admin_role` VALUES (3, '运营', '2023-01-13 18:34:48', '2023-02-01 16:58:34');

-- ----------------------------
-- Table structure for admin_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `admin_role_menu`;
CREATE TABLE `admin_role_menu`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `role_id` int(11) NULL DEFAULT NULL COMMENT '角色id',
  `menu_id` int(11) NULL DEFAULT NULL COMMENT '角色菜单id',
  `btn` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '角色菜单按钮json数组',
  `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 257 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '角色菜单表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_role_menu
-- ----------------------------
INSERT INTO `admin_role_menu` VALUES (201, 2, 1, NULL, '2023-01-12 15:13:55', '2023-01-12 15:13:55');
INSERT INTO `admin_role_menu` VALUES (211, 2, 11, NULL, '2023-01-12 15:13:55', '2023-01-12 15:13:55');
INSERT INTO `admin_role_menu` VALUES (212, 2, 12, NULL, '2023-01-12 15:13:55', '2023-01-12 15:13:55');
INSERT INTO `admin_role_menu` VALUES (213, 2, 13, NULL, '2023-01-12 15:13:55', '2023-01-12 15:13:55');
INSERT INTO `admin_role_menu` VALUES (214, 2, 14, NULL, '2023-01-12 15:13:55', '2023-01-12 15:13:55');
INSERT INTO `admin_role_menu` VALUES (215, 2, 15, NULL, '2023-01-12 15:13:55', '2023-01-12 15:13:55');
INSERT INTO `admin_role_menu` VALUES (217, 2, 16, NULL, '2023-01-12 15:14:28', '2023-01-12 15:14:28');
INSERT INTO `admin_role_menu` VALUES (237, 3, 1, '[]', '2023-02-02 11:14:32', '2023-02-02 11:14:32');
INSERT INTO `admin_role_menu` VALUES (246, 3, 11, '[]', '2023-02-02 14:41:14', '2023-02-02 14:41:14');
INSERT INTO `admin_role_menu` VALUES (252, 3, 12, '[\"create\",\"edit\",\"delete\",\"role\",\"reset_pwd\"]', '2023-02-02 16:13:09', '2023-02-02 16:13:09');
INSERT INTO `admin_role_menu` VALUES (256, 3, 13, '[\"create\",\"edit\",\"delete\"]', '2023-02-03 15:42:09', '2023-02-03 15:42:09');

-- ----------------------------
-- Table structure for admin_role_permissions
-- ----------------------------
DROP TABLE IF EXISTS `admin_role_permissions`;
CREATE TABLE `admin_role_permissions`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `role_id` int(11) NOT NULL COMMENT '角色id',
  `permissions_id` int(11) NOT NULL COMMENT '权限id',
  `create_time` datetime(0) NOT NULL COMMENT '创建时间',
  `update_time` datetime(0) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 416 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '角色权限表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_role_permissions
-- ----------------------------
INSERT INTO `admin_role_permissions` VALUES (361, 3, 2, '2023-02-02 11:11:52', '2023-02-02 11:11:52');
INSERT INTO `admin_role_permissions` VALUES (362, 3, 6, '2023-02-02 11:11:52', '2023-02-02 11:11:52');
INSERT INTO `admin_role_permissions` VALUES (363, 3, 7, '2023-02-02 11:11:52', '2023-02-02 11:11:52');
INSERT INTO `admin_role_permissions` VALUES (365, 3, 9, '2023-02-02 11:11:52', '2023-02-02 11:11:52');
INSERT INTO `admin_role_permissions` VALUES (366, 3, 10, '2023-02-02 11:11:52', '2023-02-02 11:11:52');
INSERT INTO `admin_role_permissions` VALUES (374, 2, 4, '2023-02-02 11:15:27', '2023-02-02 11:15:27');
INSERT INTO `admin_role_permissions` VALUES (375, 2, 17, '2023-02-02 11:15:27', '2023-02-02 11:15:27');
INSERT INTO `admin_role_permissions` VALUES (376, 2, 18, '2023-02-02 11:15:27', '2023-02-02 11:15:27');
INSERT INTO `admin_role_permissions` VALUES (377, 2, 19, '2023-02-02 11:15:27', '2023-02-02 11:15:27');
INSERT INTO `admin_role_permissions` VALUES (378, 2, 20, '2023-02-02 11:15:27', '2023-02-02 11:15:27');
INSERT INTO `admin_role_permissions` VALUES (379, 2, 1, '2023-02-02 11:15:57', '2023-02-02 11:15:57');
INSERT INTO `admin_role_permissions` VALUES (380, 2, 2, '2023-02-02 11:15:57', '2023-02-02 11:15:57');
INSERT INTO `admin_role_permissions` VALUES (381, 2, 6, '2023-02-02 11:15:57', '2023-02-02 11:15:57');
INSERT INTO `admin_role_permissions` VALUES (382, 2, 7, '2023-02-02 11:15:57', '2023-02-02 11:15:57');
INSERT INTO `admin_role_permissions` VALUES (383, 2, 8, '2023-02-02 11:15:57', '2023-02-02 11:15:57');
INSERT INTO `admin_role_permissions` VALUES (384, 2, 9, '2023-02-02 11:15:57', '2023-02-02 11:15:57');
INSERT INTO `admin_role_permissions` VALUES (385, 2, 10, '2023-02-02 11:15:57', '2023-02-02 11:15:57');
INSERT INTO `admin_role_permissions` VALUES (386, 2, 11, '2023-02-02 11:15:57', '2023-02-02 11:15:57');
INSERT INTO `admin_role_permissions` VALUES (387, 2, 3, '2023-02-02 11:15:57', '2023-02-02 11:15:57');
INSERT INTO `admin_role_permissions` VALUES (388, 2, 12, '2023-02-02 11:15:57', '2023-02-02 11:15:57');
INSERT INTO `admin_role_permissions` VALUES (389, 2, 13, '2023-02-02 11:15:57', '2023-02-02 11:15:57');
INSERT INTO `admin_role_permissions` VALUES (390, 2, 14, '2023-02-02 11:15:57', '2023-02-02 11:15:57');
INSERT INTO `admin_role_permissions` VALUES (391, 2, 15, '2023-02-02 11:15:57', '2023-02-02 11:15:57');
INSERT INTO `admin_role_permissions` VALUES (392, 2, 16, '2023-02-02 11:15:57', '2023-02-02 11:15:57');
INSERT INTO `admin_role_permissions` VALUES (393, 2, 5, '2023-02-02 11:15:57', '2023-02-02 11:15:57');
INSERT INTO `admin_role_permissions` VALUES (394, 2, 21, '2023-02-02 11:15:57', '2023-02-02 11:15:57');
INSERT INTO `admin_role_permissions` VALUES (395, 2, 22, '2023-02-02 11:15:57', '2023-02-02 11:15:57');
INSERT INTO `admin_role_permissions` VALUES (396, 2, 23, '2023-02-02 11:15:57', '2023-02-02 11:15:57');
INSERT INTO `admin_role_permissions` VALUES (397, 2, 24, '2023-02-02 11:15:57', '2023-02-02 11:15:57');
INSERT INTO `admin_role_permissions` VALUES (398, 3, 1, '2023-02-02 15:11:38', '2023-02-02 15:11:38');
INSERT INTO `admin_role_permissions` VALUES (410, 3, 11, '2023-02-02 16:36:28', '2023-02-02 16:36:28');
INSERT INTO `admin_role_permissions` VALUES (411, 3, 3, '2023-02-03 15:45:09', '2023-02-03 15:45:09');
INSERT INTO `admin_role_permissions` VALUES (412, 3, 12, '2023-02-03 15:45:09', '2023-02-03 15:45:09');
INSERT INTO `admin_role_permissions` VALUES (413, 3, 13, '2023-02-03 15:45:09', '2023-02-03 15:45:09');
INSERT INTO `admin_role_permissions` VALUES (414, 3, 14, '2023-02-03 15:45:09', '2023-02-03 15:45:09');
INSERT INTO `admin_role_permissions` VALUES (415, 3, 16, '2023-02-03 15:45:09', '2023-02-03 15:45:09');

-- ----------------------------
-- Table structure for admin_user
-- ----------------------------
DROP TABLE IF EXISTS `admin_user`;
CREATE TABLE `admin_user`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
  `password` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码',
  `salt` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码加盐key',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户头像',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '用户状态 0:禁用 1:启用',
  `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `username`(`username`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '管理员表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_user
-- ----------------------------
INSERT INTO `admin_user` VALUES (1, 'admin', '64066d2abd10e6192f7e2974c39bddd0', 'XCdS7R', 'http://192.168.180.121:5252/assets/upload/images/585bbada0910da59cfda3023f171cc1b.png', 1, '2023-01-12 16:13:56', '2023-01-12 16:13:56');
INSERT INTO `admin_user` VALUES (3, 'mac', '01d00d3781d2b9257aa6a91288fe9696', 'GINYUv', 'http://192.168.180.121:5252/assets/upload/images/bdf4202896dfd0ba371c9b4a8136f11c.png', 1, '2023-01-13 12:11:18', '2023-01-13 12:11:18');

-- ----------------------------
-- Table structure for admin_user_role
-- ----------------------------
DROP TABLE IF EXISTS `admin_user_role`;
CREATE TABLE `admin_user_role`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `admin_id` int(11) NULL DEFAULT NULL,
  `role_id` int(11) NULL DEFAULT NULL,
  `create_time` datetime(0) NULL DEFAULT NULL,
  `update_time` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `role_id`(`role_id`) USING BTREE,
  INDEX `admin_id`(`admin_id`) USING BTREE,
  INDEX `admin_role_id`(`admin_id`, `role_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 126 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_user_role
-- ----------------------------
INSERT INTO `admin_user_role` VALUES (1, 1, 1, '2023-01-12 16:13:56', '2023-01-12 16:13:56');
INSERT INTO `admin_user_role` VALUES (121, 3, 3, '2023-01-31 16:36:48', '2023-01-31 16:36:48');

-- ----------------------------
-- Table structure for app
-- ----------------------------
DROP TABLE IF EXISTS `app`;
CREATE TABLE `app`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '应用id',
  `secret` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '应用秘钥',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '应用名称',
  `package_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '应用报名',
  `salt` char(6) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '应用秘钥盐',
  `status` tinyint(4) NULL DEFAULT NULL COMMENT '状态 1:启用 0:禁用',
  `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10001 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '应用信息' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of app
-- ----------------------------
INSERT INTO `app` VALUES (10000, 'df2a7376fbd442b4d70f1b0f18059bea', '应用2', 'com.bac.app', 'JLZFxx', 1, '2023-02-06 15:03:12', '2023-02-06 15:03:12');

-- ----------------------------
-- Table structure for app_config
-- ----------------------------
DROP TABLE IF EXISTS `app_config`;
CREATE TABLE `app_config`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `app_id` int(11) NOT NULL DEFAULT 0 COMMENT '应用id',
  `user_day_payout_limit` int(11) NOT NULL DEFAULT 0 COMMENT '个人每日提现最高值',
  `pay_limit` bigint(20) NOT NULL DEFAULT 0 COMMENT '应用提现最高值',
  `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `appid`(`app_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10001 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'app配置' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of app_config
-- ----------------------------
INSERT INTO `app_config` VALUES (10000, 1, 22, 82, '2023-02-06 15:03:12', '2023-02-06 15:03:12');

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户id',
  `appid` int(11) NOT NULL COMMENT '应用id',
  `app_uid` bigint(20) NOT NULL COMMENT '第三方应用的用户uid',
  `phone` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',
  `invite_code` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户的邀请码',
  `invite_uid` bigint(20) NOT NULL DEFAULT 0 COMMENT '邀请人uid(上级uid)',
  `device_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '设备号',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '用户状态 1:启用 0:禁用',
  `online_status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '在线状态 1:在线 0:离线',
  `register_type` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '注册方式',
  `register_time` datetime(0) NULL DEFAULT NULL COMMENT '注册时间',
  `level` int(11) NULL DEFAULT 1 COMMENT '用户等级',
  `last_login_time` datetime(0) NULL DEFAULT NULL COMMENT '最后登录时间',
  `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `invite_code`(`invite_code`) USING BTREE,
  INDEX `app_uid`(`app_uid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (2, 10000, 27, '18191886055', 'EGE8S2', 0, '58', 1, 0, 'whatapps', '2023-02-06 17:19:44', 1, '2023-02-06 20:58:18', '2023-02-06 20:58:18', '2023-02-06 20:58:18');
INSERT INTO `user` VALUES (3, 10000, 27, '18191886055', '8G8S2D', 0, '58', 1, 0, 'whatapps', '2023-02-06 17:19:44', 1, '2023-02-07 16:47:28', '2023-02-07 16:47:28', '2023-02-07 16:47:28');

-- ----------------------------
-- Table structure for user_account
-- ----------------------------
DROP TABLE IF EXISTS `user_account`;
CREATE TABLE `user_account`  (
  `uid` bigint(20) NOT NULL COMMENT '用户id',
  `total_recharge` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '累计充值金额(分)',
  `total_commission` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '累计佣金金额(分)',
  `total_withdraw_commission` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '累计提现佣金(分)',
  `commission` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '当前可用佣金(分)',
  `freeze_commission` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '冻结中的佣金(分)',
  `child_total_recharge` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '下级累计充值金额(分)',
  PRIMARY KEY (`uid`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_account
-- ----------------------------
INSERT INTO `user_account` VALUES (2, 0, 0, 0, 0, 0, 0);
INSERT INTO `user_account` VALUES (3, 0, 0, 0, 0, 0, 0);

-- ----------------------------
-- Table structure for user_level_config
-- ----------------------------
DROP TABLE IF EXISTS `user_level_config`;
CREATE TABLE `user_level_config`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '等级id',
  `name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '等级名称',
  `percent_range` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '佣金比例范围10%-15%',
  `recharge` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '下线累计充值达到金额(分)',
  `gullak` int(11) NOT NULL DEFAULT 0 COMMENT '可获得罐子数上限',
  `withdraw_amount` bigint(20) NOT NULL DEFAULT 0 COMMENT '每日提现上限',
  `withdraw_times` int(11) NOT NULL DEFAULT 0 COMMENT '每日提现次数',
  `recharge_discount` int(11) NOT NULL DEFAULT 0 COMMENT '充值折扣20%',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_level_config
-- ----------------------------
INSERT INTO `user_level_config` VALUES (1, 'Bronze', '5%-10%', 0, 2, 100000, 2, 0);
INSERT INTO `user_level_config` VALUES (2, 'Silver', '10%-15%', 1000000, 4, 200000, 5, 0);
INSERT INTO `user_level_config` VALUES (3, 'Gold', '15%-20%', 100000000, 8, 500000, 10, 0);
INSERT INTO `user_level_config` VALUES (4, 'Platinum', '20%-25%', 1000000000, 15, 1000000, 20, 0);
INSERT INTO `user_level_config` VALUES (5, 'Diamond', '25%-30%', 10000000000, 20, 2000000, 40, 20);

SET FOREIGN_KEY_CHECKS = 1;
