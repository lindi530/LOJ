CREATE TABLE IF NOT EXISTS competition_submissions (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    competition_id BIGINT UNSIGNED NOT NULL,
    problem_id INT UNSIGNED NOT NULL,
    submission_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    is_ac TINYINT(1) NOT NULL DEFAULT 0,
    is_first_ac TINYINT(1) NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    KEY idx_competition_user_problem (competition_id, user_id, problem_id),
    KEY idx_submission_id (submission_id)
) ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_unicode_ci;
