-- 1. 用户参与竞赛信息表
CREATE TABLE IF NOT EXISTS competition_summary (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    rank_score INT NOT NULL DEFAULT 1000 COMMENT '当前rank分',
    competition_count INT NOT NULL DEFAULT 0 COMMENT '参与比赛场数',

    UNIQUE KEY uk_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户参与竞赛信息表';