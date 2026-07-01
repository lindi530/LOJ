-- 1. 课程表
CREATE TABLE IF NOT EXISTS courses (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '课程ID',

    title VARCHAR(128) NOT NULL COMMENT '课程标题',
    description TEXT COMMENT '课程介绍',
    cover_url VARCHAR(512) COMMENT '课程封面',

    price DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '价格，0表示免费',
    is_free TINYINT NOT NULL DEFAULT 1 COMMENT '是否免费：0否 1是',

    lesson_count INT NOT NULL DEFAULT 0 COMMENT '课时数量',
    student_count INT NOT NULL DEFAULT 0 COMMENT '学习人数',

    status TINYINT NOT NULL DEFAULT 0 COMMENT '状态：0草稿 1上架 2下架',

    created_by BIGINT COMMENT '创建人/管理员ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX idx_status_sort (status, sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课程表';


-- 2. 课程章节表
CREATE TABLE IF NOT EXISTS course_chapters (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '章节ID',

    course_id BIGINT NOT NULL COMMENT '课程ID',
    title VARCHAR(128) NOT NULL COMMENT '章节标题',
    description TEXT COMMENT '章节介绍',

    sort INT NOT NULL DEFAULT 0 COMMENT '排序',

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课程章节表';


-- 3. 视频资源表
CREATE TABLE IF NOT EXISTS video_assets (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '视频资源ID',
    origin_path VARCHAR(500) NOT NULL COMMENT '原始视频路径',
    play_path VARCHAR(500) DEFAULT NULL COMMENT 'master.m3u8路径',
    cover_path VARCHAR(500) DEFAULT NULL COMMENT '封面路径',
    status TINYINT NOT NULL DEFAULT 0 COMMENT '
        0=上传中,
        1=转码中,
        2=可播放,
        3=转码失败
    ',
    duration INT NOT NULL DEFAULT 0 COMMENT '视频时长，单位秒',
    size_bytes BIGINT COMMENT '视频大小，单位字节'

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='视频资源表';


-- -- 4. 课程课时表
-- CREATE TABLE IF NOT EXISTS course_lessons (
--     id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '课时ID',
--
--     course_id BIGINT NOT NULL COMMENT '课程ID',
--     chapter_id BIGINT NOT NULL COMMENT '章节ID',
--     video_asset_id BIGINT COMMENT '视频资源ID',
--
--     title VARCHAR(128) NOT NULL COMMENT '课时标题',
--     duration INT NOT NULL DEFAULT 0 COMMENT '视频时长，单位秒',
--
--     is_free_preview TINYINT NOT NULL DEFAULT 0 COMMENT '是否试看：0否 1是',
--
--     sort INT NOT NULL DEFAULT 0 COMMENT '排序',
--     status TINYINT NOT NULL DEFAULT 0 COMMENT '状态：0草稿 1上架 2下架',
--
--     created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
--     updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
--
--     INDEX idx_course_chapter_sort (course_id, chapter_id, sort),
--     INDEX idx_video_asset_id (video_asset_id)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课程课时表';


-- 5. 课时题目关联表
CREATE TABLE IF NOT EXISTS chapter_problems (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '关联ID',

    course_id BIGINT NOT NULL COMMENT '课程ID',
    chapter_id BIGINT NOT NULL COMMENT '章节ID',
    problem_id BIGINT NOT NULL COMMENT 'OJ题目ID，对应 problem.id',
    sort INT NOT NULL DEFAULT 0 COMMENT '排序'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课时题目关联表';

-- 5. 视频关联表
CREATE TABLE IF NOT EXISTS chapter_videos (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '关联ID',

    course_id BIGINT NOT NULL COMMENT '课程ID',
    chapter_id BIGINT NOT NULL COMMENT '章节ID',
    video_id BIGINT NOT NULL COMMENT '视频ID，对应 video_assets.id'

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课时视频关联表';

-- -- 6. 课时学习进度表
-- CREATE TABLE IF NOT EXISTS lesson_progress (
--     id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '学习进度ID',
--
--     user_id BIGINT NOT NULL COMMENT '用户ID，对应 user.id',
--     course_id BIGINT NOT NULL COMMENT '课程ID',
--     lesson_id BIGINT NOT NULL COMMENT '课时ID',
--
--     last_position INT NOT NULL DEFAULT 0 COMMENT '上次观看到的位置，单位秒',
--     watched_seconds INT NOT NULL DEFAULT 0 COMMENT '已观看秒数',
--     progress_percent DECIMAL(5,2) NOT NULL DEFAULT 0.00 COMMENT '学习进度百分比',
--
--     is_finished TINYINT NOT NULL DEFAULT 0 COMMENT '是否完成：0否 1是',
--     finished_at DATETIME COMMENT '完成时间',
--
--     created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
--     updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
--
--     UNIQUE KEY uk_user_lesson (user_id, lesson_id),
--     INDEX idx_user_course (user_id, course_id),
--     INDEX idx_course_lesson (course_id, lesson_id)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课时学习进度表';

CREATE TABLE IF NOT EXISTS video_profile (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    video_id BIGINT NOT NULL,
    resolution VARCHAR(20) NOT NULL COMMENT '
    480P
    720P
    1080P
    ',
    width INT NOT NULL,
    height INT NOT NULL,
    bitrate INT NOT NULL COMMENT '码率(kbps)',
    playlist_path VARCHAR(500) NOT NULL COMMENT 'index.m3u8',
    file_size BIGINT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_video_resolution (
        video_id,
        resolution
    )
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='视频资源存放表';

CREATE TABLE IF NOT EXISTS course_orders (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_no VARCHAR(64) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    course_id BIGINT NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    status TINYINT NOT NULL DEFAULT 0, -- 0待支付 1已支付 2取消 3退款 4过期
    pay_channel VARCHAR(32),
    transaction_id VARCHAR(128),
    paid_at DATETIME NULL,
    expire_at DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_course (user_id, course_id),
    INDEX idx_status_expire (status, expire_at)
);

CREATE TABLE IF NOT EXISTS course_enrollments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    course_id BIGINT NOT NULL,
    order_id BIGINT NULL,
    status TINYINT NOT NULL DEFAULT 1, -- 1有效 2已失效/退款
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_course (user_id, course_id),
    INDEX idx_course (course_id)
);