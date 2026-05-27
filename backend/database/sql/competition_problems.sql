CREATE TABLE IF NOT EXISTS competition_problems (
     id BIGINT PRIMARY KEY AUTO_INCREMENT,
     competition_id BIGINT NOT NULL,
     problem_id BIGINT NOT NULL
);