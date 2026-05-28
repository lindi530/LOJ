CREATE TABLE IF NOT EXISTS competition_problems (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    competition_id BIGINT NOT NULL,
    problem_number CHAR(1) NOT NULL,
    problem_id INT NOT NULL,
    ac_count INT NOT NULL DEFAULT 0,
    submit_count INT NOT NULL DEFAULT 0
);