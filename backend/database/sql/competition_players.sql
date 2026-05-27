CREATE TABLE IF NOT EXISTS competition_players (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '报名记录ID',
  competition_id BIGINT NOT NULL COMMENT '竞赛ID',
  user_id BIGINT NOT NULL COMMENT '用户ID',
  user_name VARCHAR(100) NOT NULL COMMENT '用户名',
  is_cancel TINYINT NOT NULL DEFAULT 0 COMMENT '是否取消：0=未取消，1=已取消',

  UNIQUE KEY uk_competition_user (competition_id, user_id)
) COMMENT='竞赛报名表';