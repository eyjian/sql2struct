DROP TABLE `t_commodity`;
CREATE TABLE `t_commodity` (
                          `f_id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '商品ID',
                          `f_name` VARCHAR(255) NOT NULL COMMENT '商品名称',
                          `f_description` TEXT,
                          `f_price` DECIMAL(10, 2) NOT NULL,
                          `f_weight` FLOAT NOT NULL COMMENT '商品重量（kg）',
                          `f_quantity` SMALLINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '商品库存数量',
                          `f_is_active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '商品是否激活（0 - 未激活，1 - 激活）',
                          `f_rating` DOUBLE COMMENT '商品评分',
                          `f_created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '商品创建时间',
                          `f_updated_at` DATETIME ON UPDATE CURRENT_TIMESTAMP COMMENT '商品更新时间',
                          UNIQUE INDEX idx_name_at (`f_name`),
                          INDEX idx_created_at (`f_created_at`),
                          KEY idx_updated_at (`f_updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';
