DROP TABLE goods;
CREATE TABLE goods ( -- 商品表
                          id INT UNSIGNED AUTO_INCREMENT COMMENT '商品ID',
                          name VARCHAR(255) NOT NULL UNIQUE COMMENT '商品名称',
                          description TEXT COMMENT '商品描述',
                          price DECIMAL(10, 2) NOT NULL COMMENT '商品价格',
                          weight FLOAT NOT NULL COMMENT '商品重量（kg）',
                          quantity SMALLINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '商品库存数量',
                          is_active TINYINT(1) NOT NULL DEFAULT 1 COMMENT '商品是否激活（0 - 未激活，1 - 激活）',
                          rating DOUBLE COMMENT '商品评分',
                          created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '商品创建时间',
                          updated_at DATETIME ON UPDATE CURRENT_TIMESTAMP COMMENT '商品更新时间',
                          PRIMARY KEY (id),
                          INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';
