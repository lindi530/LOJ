CREATE TABLE IF NOT EXISTS user_wallets (
    user_id BIGINT NOT NULL COMMENT '用户ID，对应用户表主键；每个用户一条钱包记录',
    balance DECIMAL(18,2) NOT NULL DEFAULT 0.00 COMMENT '当前虚拟币余额，最小单位为0.01',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '钱包记录创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '钱包记录最后更新时间',
    PRIMARY KEY (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户虚拟币钱包表';

CREATE TABLE IF NOT EXISTS wallet_transactions (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '钱包流水自增ID',
    tx_no VARCHAR(64) NOT NULL COMMENT '钱包流水号，全局唯一，用于幂等和排查',
    user_id BIGINT NOT NULL COMMENT '用户ID，对应用户表主键',
    change_amount DECIMAL(18,2) NOT NULL COMMENT '本次虚拟币变动金额，正数表示增加，负数表示扣减，最小单位为0.01',
    balance_after DECIMAL(18,2) NOT NULL COMMENT '本次变动后的虚拟币余额，最小单位为0.01',
    biz_type VARCHAR(32) NOT NULL COMMENT '业务类型，例如course_order、signin、admin_grant',
    biz_id VARCHAR(64) NOT NULL COMMENT '业务ID，例如课程订单号或签到日期',
    remark VARCHAR(255) COMMENT '流水备注，用于展示或后台排查',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '流水创建时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_tx_no (tx_no) COMMENT '保证流水号唯一',
    UNIQUE KEY uk_biz_user (biz_type, biz_id, user_id) COMMENT '保证同一业务对同一用户只记账一次',
    INDEX idx_user_created_at (user_id, created_at) COMMENT '支持按用户和时间查询钱包流水'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户虚拟币钱包流水表';
