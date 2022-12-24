-- ------------------------
-- Table structure for user
-- ------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` -- 用户信息表
(
    `id`          INTEGER PRIMARY KEY AUTOINCREMENT, -- 用户id
    `name`        VARCHAR(64)  NOT NULL,             -- 用户名
    `nickname`    VARCHAR(64)  DEFAULT '',           -- 昵称
    `passwd`      VARCHAR(128) NOT NULL,             -- 密码
    `rem`         VARCHAR(256) DEFAULT '',           -- 备注
    `del_flag`    TINYINT      DEFAULT 0,            -- 删除标识，0-正常，1-删除
    `create_time` INTEGER      DEFAULT 0,            -- 创建时间（时间戳，s）
    `update_time` INTEGER      DEFAULT 0             -- 修改时间（时间戳，s）
);
INSERT INTO `user` (`name`, `nickname`, `passwd`, `create_time`) VALUES ('admin', 'admin', 'admin', 1671614960);
INSERT INTO `user` (`name`, `nickname`, `passwd`, `create_time`) VALUES ('demo', 'demo', 'demo', 1671614960);
-- prod, test, dev


-- ----------------------
-- Table structure for rx
-- ----------------------
DROP TABLE IF EXISTS `rx`;
CREATE TABLE `rx` -- rx信息表
(
    `id`          INTEGER PRIMARY KEY AUTOINCREMENT, -- 用户id
    `name`        VARCHAR(64) NOT NULL,              -- 名称
    `owner_id`    INTEGER     NOT NULL,              -- 拥有者id
    `sharer_id`   INTEGER      DEFAULT 0,            -- 共享者id
    `rem`         VARCHAR(256) DEFAULT '',           -- 备注
    `del_flag`    TINYINT      DEFAULT 0,            -- 删除标识，0-正常，1-删除
    `create_time` INTEGER      DEFAULT 0,            -- 创建时间（时间戳，s）
    `update_time` INTEGER      DEFAULT 0             -- 修改时间（时间戳，s）
);


-- -----------------------
-- Table structure for git
-- -----------------------
DROP TABLE IF EXISTS `git`;
CREATE TABLE `git` -- Git信息表
(
    `id`          INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `user_id`     INTEGER      NOT NULL,             -- Git所属用户id
    `name`        VARCHAR(64)  NOT NULL,             -- 名称
    `user`        VARCHAR(64)  NOT NULL,             -- 用户
    `passwd`      VARCHAR(128) NOT NULL,             -- 密码
    `rem`         VARCHAR(256) DEFAULT '',           -- 备注
    `del_flag`    TINYINT      DEFAULT 0,            -- 删除标识，0-正常，1-删除
    `create_time` INTEGER      DEFAULT 0,            -- 创建时间（时间戳，s）
    `update_time` INTEGER      DEFAULT 0             -- 修改时间（时间戳，s）
);


-- --------------------------
-- Table structure for server
-- --------------------------
DROP TABLE IF EXISTS `server`;
CREATE TABLE `server` -- Server信息表
(
    `id`          INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `user_id`     INTEGER      NOT NULL,             -- Server所属用户id
    `name`        VARCHAR(64)  NOT NULL,             -- 名称
    `host`        VARCHAR(64)  NOT NULL,             -- HOST
    `port`        SMALLINT     DEFAULT 22,           -- 端口
    `user`        VARCHAR(64)  NOT NULL,             -- 用户
    `passwd`      VARCHAR(128) NOT NULL,             -- 密码
    `rem`         VARCHAR(256) DEFAULT '',           -- 备注
    `del_flag`    TINYINT      DEFAULT 0,            -- 删除标识，0-正常，1-删除
    `create_time` INTEGER      DEFAULT 0,            -- 创建时间（时间戳，s）
    `update_time` INTEGER      DEFAULT 0             -- 修改时间（时间戳，s）
);


-- ------------------------
-- Table structure for item
-- ------------------------
DROP TABLE IF EXISTS `item`;
CREATE TABLE `item` -- 项目信息表
(
    `id`          INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `user_id`     INTEGER      NOT NULL,             -- 项目所属用户id
    `name`        VARCHAR(64)  NOT NULL,             -- 名称
    `git_id`      INTEGER      DEFAULT 0,            -- 项目所属Git id
    `repo_url`    VARCHAR(512) NOT NULL,             -- Git仓库地址
    `branch`      VARCHAR(64)  NOT NULL,             -- 分支名
    `server_id`   INTEGER      NOT NULL,             -- 项目所属Server id
    `ini`         TEXT         DEFAULT NULL,         -- 配置
    `rem`         VARCHAR(256) DEFAULT '',           -- 备注
    `del_flag`    TINYINT      DEFAULT 0,            -- 删除标识，0-正常，1-删除
    `create_time` INTEGER      DEFAULT 0,            -- 创建时间（时间戳，s）
    `update_time` INTEGER      DEFAULT 0             -- 修改时间（时间戳，s）
);


-- --------------------------
-- Table structure for record
-- --------------------------
DROP TABLE IF EXISTS `record`;
CREATE TABLE `record` -- 项目部署记录信息表
(
    `id`           INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `item_id`      INTEGER      NOT NULL,             -- 所属项目id
    `pull_stime`   INTEGER DEFAULT 0,                 -- pull开始时间（时间戳，s）
    `pull_etime`   INTEGER DEFAULT 0,                 -- pull结束时间（时间戳，s）
    `pull_rem`     TEXT    DEFAULT NULL,              -- pull备注信息
    `commit_id`    VARCHAR(128) NOT NULL,             -- 提交id
    `rev_msg`      TEXT    DEFAULT NULL,              -- 提交信息
    `build_stime`  INTEGER DEFAULT 0,                 -- build开始时间（时间戳，s）
    `build_etime`  INTEGER DEFAULT 0,                 -- build结束时间（时间戳，s）
    `build_rem`    TEXT    DEFAULT NULL,              -- build备注信息
    `pack_stime`   INTEGER DEFAULT 0,                 -- pack开始时间（时间戳，s）
    `pack_etime`   INTEGER DEFAULT 0,                 -- pack结束时间（时间戳，s）
    `pack_rem`     TEXT    DEFAULT NULL,              -- pack备注信息
    `ul_stime`     INTEGER DEFAULT 0,                 -- upload开始时间（时间戳，s）
    `ul_etime`     INTEGER DEFAULT 0,                 -- upload结束时间（时间戳，s）
    `ul_rem`       TEXT    DEFAULT NULL,              -- upload备注信息
    `deploy_stime` INTEGER DEFAULT 0,                 -- deploy开始时间（时间戳，s）
    `deploy_etime` INTEGER DEFAULT 0,                 -- deploy结束时间（时间戳，s）
    `deploy_rem`   TEXT    DEFAULT NULL,              -- deploy备注信息
    `status`       TINYINT      NOT NULL,             -- 状态，1-部署中，2-部署异常，3-部署成功
    `rem`          TEXT    DEFAULT NULL,              -- 备注信息
    `del_flag`     TINYINT DEFAULT 0,                 -- 删除标识，0-正常，1-删除
    `create_time`  INTEGER DEFAULT 0,                 -- 创建时间（时间戳，s）
    `update_time`  INTEGER DEFAULT 0                  -- 修改时间（时间戳，s）
);
