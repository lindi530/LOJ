CREATE TABLE IF NOT EXISTS competition_results (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    competition_id BIGINT NOT NULL,
    competition_name varchar(100) NOT NULL,
    user_id BIGINT NOT NULL,
    competition_rank INT NOT NULL,
    solved_count INT NOT NULL DEFAULT 0,
    penalty INT NOT NULL DEFAULT 0,
    rating_before INT NOT NULL,
    rating_after INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_competition_user (competition_id, user_id),
    INDEX idx_user_created_at (user_id, created_at)
);
