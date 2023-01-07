------------------------------
-- Table structure for setting
-- ---------------------------
DROP TABLE IF EXISTS `setting`;
CREATE TABLE `setting` -- 系统设置信息表
(
    `allow_reg_flag` TINYINT       DEFAULT 1, -- 允许用户注册标识，0-不允许，1-允许
    `build_level`    TINYINT       DEFAULT 2, -- 构建级别：1，当build_env空闲时，项目才进行构建（安全级别高）；2，随机选取一个build_env来构建（安全级别低）
    `build_envs`     VARCHAR(1024) DEFAULT '' -- 构建环境集
);
INSERT INTO `setting` (`allow_reg_flag`, `build_level`, `build_envs`) VALUES ('1', '2', 'default:');


-- ------------------------
-- Table structure for user
-- ------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` -- 用户信息表
(
    `id`       INTEGER PRIMARY KEY AUTOINCREMENT, -- 用户id
    `name`     VARCHAR(64)  NOT NULL,             -- 用户名
    `nickname` VARCHAR(64)  DEFAULT '',           -- 昵称
    `passwd`   VARCHAR(128) NOT NULL,             -- 密码
    `rem`      VARCHAR(256) DEFAULT '',           -- 备注
    `del_flag` TINYINT      DEFAULT 0,            -- 删除标识，0-正常，1-删除
    `add_time` INTEGER      DEFAULT 0,            -- 创建时间（时间戳，s）
    `upd_time` INTEGER      DEFAULT 0             -- 修改时间（时间戳，s）
);
INSERT INTO `user` (`name`, `nickname`, `passwd`, `add_time`) VALUES ('admin', 'admin', '75b17d369a5ce9b50e1a608bee111cac', 1671614960);
-- admin
-- prod, test, dev


-- ----------------------
-- Table structure for rx
-- ----------------------
DROP TABLE IF EXISTS `rx`;
CREATE TABLE `rx` -- rx信息表
(
    `id`        INTEGER PRIMARY KEY AUTOINCREMENT, -- 用户id
    `name`      VARCHAR(64) NOT NULL,              -- 名称
    `owner_id`  INTEGER     NOT NULL,              -- 拥有者id
    `sharer_id` INTEGER       DEFAULT 0,           -- 共享者id
    `item_ids`  VARCHAR(1024) DEFAULT '',          -- 共享item_id集合
    `rem`       VARCHAR(256)  DEFAULT '',          -- 备注
    `del_flag`  TINYINT       DEFAULT 0,           -- 删除标识，0-正常，1-删除
    `add_time`  INTEGER       DEFAULT 0,           -- 创建时间（时间戳，s）
    `upd_time`  INTEGER       DEFAULT 0            -- 修改时间（时间戳，s）
);


-- -----------------------
-- Table structure for git
-- -----------------------
DROP TABLE IF EXISTS `git`;
CREATE TABLE `git` -- Git信息表
(
    `id`       INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `user_id`  INTEGER      NOT NULL,             -- Git所属用户id
    `name`     VARCHAR(64)  NOT NULL,             -- 名称
    `user`     VARCHAR(64)  NOT NULL,             -- 用户
    `passwd`   VARCHAR(128) NOT NULL,             -- 密码
    `rem`      VARCHAR(256) DEFAULT '',           -- 备注
    `del_flag` TINYINT      DEFAULT 0,            -- 删除标识，0-正常，1-删除
    `add_time` INTEGER      DEFAULT 0,            -- 创建时间（时间戳，s）
    `upd_time` INTEGER      DEFAULT 0             -- 修改时间（时间戳，s）
);


-- --------------------------
-- Table structure for server
-- --------------------------
DROP TABLE IF EXISTS `server`;
CREATE TABLE `server` -- Server信息表
(
    `id`       INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `user_id`  INTEGER      NOT NULL,             -- Server所属用户id
    `name`     VARCHAR(64)  NOT NULL,             -- 名称
    `host`     VARCHAR(64)  NOT NULL,             -- HOST
    `port`     SMALLINT     DEFAULT 22,           -- 端口
    `user`     VARCHAR(64)  NOT NULL,             -- 用户
    `passwd`   VARCHAR(128) NOT NULL,             -- 密码
    `rem`      VARCHAR(256) DEFAULT '',           -- 备注
    `del_flag` TINYINT      DEFAULT 0,            -- 删除标识，0-正常，1-删除
    `add_time` INTEGER      DEFAULT 0,            -- 创建时间（时间戳，s）
    `upd_time` INTEGER      DEFAULT 0             -- 修改时间（时间戳，s）
);


-- ------------------------
-- Table structure for item
-- ------------------------
DROP TABLE IF EXISTS `item`;
CREATE TABLE `item` -- 项目信息表
(
    `id`        INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `user_id`   INTEGER      NOT NULL,             -- 项目所属用户id
    `name`      VARCHAR(64)  NOT NULL,             -- 名称
    `git_id`    INTEGER      DEFAULT 0,            -- 项目所属Git id
    `repo_url`  VARCHAR(512) NOT NULL,             -- Git仓库地址
    `branch`    VARCHAR(64)  NOT NULL,             -- 分支名
    `server_id` INTEGER      NOT NULL,             -- 项目所属Server id
    `script`    TEXT         DEFAULT NULL,         -- 脚本
    `rem`       VARCHAR(256) DEFAULT '',           -- 备注
    `del_flag`  TINYINT      DEFAULT 0,            -- 删除标识，0-正常，1-删除
    `add_time`  INTEGER      DEFAULT 0,            -- 创建时间（时间戳，s）
    `upd_time`  INTEGER      DEFAULT 0             -- 修改时间（时间戳，s）
);


-- --------------------------
-- Table structure for record
-- --------------------------
DROP TABLE IF EXISTS `record`;
CREATE TABLE `record` -- 项目部署记录信息表
(
    `id`            INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `item_id`       INTEGER NOT NULL,                  -- 所属项目id
    `pull_stime`    INTEGER      DEFAULT 0,            -- pull开始时间（时间戳，s）
    `pull_etime`    INTEGER      DEFAULT 0,            -- pull结束时间（时间戳，s）
    `pull_status`   TINYINT      DEFAULT 0,            -- pull状态，非0表示异常
    `pull_rem`      TEXT         DEFAULT NULL,         -- pull备注信息
    `commit_id`     VARCHAR(128) DEFAULT '',           -- 提交id
    `rev_msg`       TEXT         DEFAULT NULL,         -- 提交信息
    `build_stime`   INTEGER      DEFAULT 0,            -- build开始时间（时间戳，s）
    `build_etime`   INTEGER      DEFAULT 0,            -- build结束时间（时间戳，s）
    `build_status`  TINYINT      DEFAULT 0,            -- build状态，非0表示异常
    `build_rem`     TEXT         DEFAULT NULL,         -- build备注信息
    `pack_stime`    INTEGER      DEFAULT 0,            -- pack开始时间（时间戳，s）
    `pack_etime`    INTEGER      DEFAULT 0,            -- pack结束时间（时间戳，s）
    `pack_status`   TINYINT      DEFAULT 0,            -- pack状态，非0表示异常
    `pack_rem`      TEXT         DEFAULT NULL,         -- pack备注信息
    `ul_stime`      INTEGER      DEFAULT 0,            -- upload开始时间（时间戳，s）
    `ul_etime`      INTEGER      DEFAULT 0,            -- upload结束时间（时间戳，s）
    `ul_status`     TINYINT      DEFAULT 0,            -- upload状态，非0表示异常
    `ul_rem`        TEXT         DEFAULT NULL,         -- upload备注信息
    `unpack_stime`  INTEGER      DEFAULT 0,            -- unpack开始时间（时间戳，s）
    `unpack_etime`  INTEGER      DEFAULT 0,            -- unpack结束时间（时间戳，s）
    `unpack_status` TINYINT      DEFAULT 0,            -- unpack状态，非0表示异常
    `unpack_rem`    TEXT         DEFAULT NULL,         -- unpack备注信息
    `deploy_stime`  INTEGER      DEFAULT 0,            -- deploy开始时间（时间戳，s）
    `deploy_etime`  INTEGER      DEFAULT 0,            -- deploy结束时间（时间戳，s）
    `deploy_status` TINYINT      DEFAULT 0,            -- deploy状态，非0表示异常
    `deploy_rem`    TEXT         DEFAULT NULL,         -- deploy备注信息
    `status`        TINYINT NOT NULL,                  -- 状态，1-部署中，2-部署异常，3-部署成功
    `rem`           TEXT         DEFAULT NULL,         -- 备注信息
    `del_flag`      TINYINT      DEFAULT 0,            -- 删除标识，0-正常，1-删除
    `add_time`      INTEGER      DEFAULT 0,            -- 创建时间（时间戳，s）
    `upd_time`      INTEGER      DEFAULT 0             -- 修改时间（时间戳，s）
);
