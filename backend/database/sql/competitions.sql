CREATE TABLE IF NOT EXISTS competitions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '竞赛名称',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '1未开始 2进行中 3已结束',

    visibility TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否公开：1公开，0非公开',
    password_hash VARCHAR(255) DEFAULT NULL COMMENT '非公开竞赛密码哈希',

    player_count INT NOT NULL DEFAULT 0 COMMENT '总参赛人数/报名人数',

    start_time DATETIME DEFAULT NULL,
    end_time DATETIME DEFAULT NULL,

    creator_id BIGINT NOT NULL,
    creator_name VARCHAR(64) NOT NULL,

    announcement TEXT COMMENT '公告',

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_status_time (status, start_time),
    INDEX idx_creator_id (creator_id),

    CHECK (status IN (1, 2, 3)),
    CHECK (visibility IN (0, 1))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='竞赛表';