SET
FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS image_dragon.algo_log;
CREATE TABLE image_dragon.algo_log
(
    `id`            int(8) AUTO_INCREMENT COMMENT '主键',
    `request_id`    varchar(64)  NOT NULL COMMENT '请求ID',
    `prompt`        varchar(512) NOT NULL COMMENT '文本描述',
    `task_id`       varchar(64)  NOT NULL COMMENT '任务活动id',
    `ddim_steps`    int(64) DEFAULT NULL COMMENT '通常steps越多，效果越好',
    `strength`      float                 DEFAULT NULL COMMENT '有初始图时，实际steps=ddim_steps*strength',
    `H`             int(64) DEFAULT NULL COMMENT '生成图片高度，有初始图时该参数无效',
    `W`             int(64) DEFAULT NULL COMMENT '生成图片宽度，有初始图时该参数无效',
    `n_samples`     int(64) DEFAULT NULL COMMENT '生成图片数',
    `seed`          int(64) DEFAULT NULL COMMENT '随机种子',
    `init_image`    varchar(512) NOT NULL COMMENT '原始图片url',
    `create_images` text         NOT NULL COMMENT '生成图片url',
    `snap_time`     datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '服务调用时间',
    `total_time`    int(32) DEFAULT NULL COMMENT '初始化总耗时（毫秒）',
    `code`          int(8) DEFAULT 0 COMMENT '状态码',
    `message`       varchar(64)           DEFAULT NULL COMMENT '状态信息',
    `details`       text                  DEFAULT NULL COMMENT '返回结果',
--     `create_time`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
--     `create_user`   varchar(32)           DEFAULT NULL COMMENT '创建人员',
    PRIMARY KEY (`id`) USING BTREE,
    KEY             `idx_request_id` (`request_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  ROW_FORMAT = COMPACT;

alter table image_dragon.algo_log
    add `artists` varchar(256) DEFAULT NULL COMMENT '艺术家';

alter table image_dragon.algo_log
    add `styles` varchar(256) DEFAULT NULL COMMENT '风格';

alter table image_dragon.algo_log
    add `image_type` varchar(32) DEFAULT NULL COMMENT '类型';

alter table image_dragon.algo_log
    add `smart_mode` varchar(32) DEFAULT NULL COMMENT '是否开启智能模式';

alter table image_dragon.algo_log
    add `negative` varchar(256) DEFAULT NULL COMMENT '反向词';

alter table image_dragon.algo_log
    add `size_type` int(8) DEFAULT NULL COMMENT '图片尺寸';