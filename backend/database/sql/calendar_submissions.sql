CREATE TABLE IF NOT EXISTS calendar_submissions (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  ac_date DATE NOT NULL,
  ac_count INT NOT NULL DEFAULT 0,
  UNIQUE KEY uk_user_date (user_id, ac_date)
);