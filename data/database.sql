-- ------------------------
-- Table structure for user
-- ------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` -- 用户信息表
(
    `id`          INTEGER PRIMARY KEY AUTOINCREMENT, -- 用户id
    `name`        VARCHAR(64)  NOT NULL,             -- 用户名
    `nickname`    VARCHAR(64)  DEFAULT '',           -- 昵称
    `passwd`      VARCHAR(256) NOT NULL,             -- 密码
    `rem`         VARCHAR(256) DEFAULT '',           -- 备注
    `del_flag`    TINYINT      DEFAULT 0,            -- 删除标识，0-正常，1-删除
    `create_time` INTEGER      DEFAULT 0,            -- 创建时间（时间戳，s）
    `update_time` INTEGER      DEFAULT 0             -- 修改时间（时间戳，s）
);
INSERT INTO `user` (`name`, `nickname`, `passwd`, `create_time`) VALUES ('admin', 'admin', 'admin', 1671614960);
INSERT INTO `user` (`name`, `nickname`, `passwd`, `create_time`) VALUES ('demo', 'demo', 'demo', 1671614960);
-- prod, test, dev


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
    `passwd`      VARCHAR(256) NOT NULL,             -- 密码
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
    `passwd`      VARCHAR(256) NOT NULL,             -- 密码
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
    `branch`      VARCHAR(64)  DEFAULT 'master',     -- 分支名
    `server_id`   INTEGER      NOT NULL,             -- 项目所属Server id
    `cmd`         VARCHAR(512) DEFAULT '',           -- 构建命令
    `script`      TEXT         DEFAULT NULL,         -- 脚本，目前支持 #!/dockerfile, #!/static 解析
    `rem`         VARCHAR(256) DEFAULT '',           -- 备注
    `del_flag`    TINYINT      DEFAULT 0,            -- 删除标识，0-正常，1-删除
    `create_time` INTEGER      DEFAULT 0,            -- 创建时间（时间戳，s）
    `update_time` INTEGER      DEFAULT 0             -- 修改时间（时间戳，s）
);
