CREATE TABLE competitions (
      id BIGINT PRIMARY KEY AUTO_INCREMENT,
      name VARCHAR(100) NOT NULL COMMENT '竞赛名称',
      status TINYINT NOT NULL DEFAULT 1 COMMENT '1未开始 2进行中 3已结束',

      visibility TINYINT NOT NULL DEFAULT 1 COMMENT '1公开 2非公开',
      password_hash VARCHAR(255) DEFAULT NULL COMMENT '非公开竞赛密码哈希',

      player_count INT NOT NULL DEFAULT 0 COMMENT '总参赛人数/报名人数',

      start_time DATETIME,
      end_time DATETIME,

      creator_id BIGINT NOT NULL,
      creator_name VARCHAR(64) NOT NULL,

      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

      INDEX idx_status_time (status, start_time),
      INDEX idx_creator_id (creator_id),
      INDEX idx_visibility (visibility)
);