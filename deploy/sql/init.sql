CREATE DATABASE IF NOT EXISTS go_zero_ecommerce DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE go_zero_ecommerce;

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(64) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码',
  `nickname` varchar(64) DEFAULT NULL COMMENT '昵称',
  `mobile` varchar(20) DEFAULT NULL COMMENT '手机号',
  `email` varchar(128) DEFAULT NULL COMMENT '邮箱',
  `gender` tinyint NOT NULL DEFAULT '0' COMMENT '性别 0未知 1男 2女',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 0禁用 1正常',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_mobile` (`mobile`),
  UNIQUE KEY `uk_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

INSERT INTO `user` (`username`, `password`, `nickname`, `mobile`, `email`, `gender`, `status`) VALUES
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '管理员', '13800000000', 'admin@example.com', 1, 1),
('test', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '测试用户', '13800000001', 'test@example.com', 0, 1);

DROP TABLE IF EXISTS `category`;
CREATE TABLE `category` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `parent_id` bigint NOT NULL DEFAULT '0' COMMENT '父ID',
  `name` varchar(128) NOT NULL COMMENT '分类名称',
  `sort_order` int NOT NULL DEFAULT '0' COMMENT '排序',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 0禁用 1启用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品分类表';

INSERT INTO `category` (`parent_id`, `name`, `sort_order`, `status`) VALUES
(0, '数码电器', 1, 1),
(0, '服装鞋帽', 2, 1),
(0, '食品饮料', 3, 1),
(1, '手机数码', 1, 1),
(1, '电脑办公', 2, 1),
(2, '男装', 1, 1),
(2, '女装', 2, 1),
(3, '休闲食品', 1, 1),
(4, '智能手机', 1, 1),
(4, '平板电脑', 2, 1);

DROP TABLE IF EXISTS `product`;
CREATE TABLE `product` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `category_id` bigint NOT NULL COMMENT '分类ID',
  `name` varchar(255) NOT NULL COMMENT '商品名称',
  `subtitle` varchar(500) DEFAULT NULL COMMENT '副标题',
  `main_image` varchar(500) DEFAULT NULL COMMENT '主图',
  `sub_images` text COMMENT '副图JSON',
  `detail` text COMMENT '详情',
  `spec` varchar(500) DEFAULT NULL COMMENT '规格',
  `price` bigint NOT NULL COMMENT '销售价格 分',
  `original_price` bigint NOT NULL COMMENT '原价 分',
  `stock` int NOT NULL DEFAULT '0' COMMENT '库存',
  `sales` int NOT NULL DEFAULT '0' COMMENT '销量',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 0下架 1上架',
  `brand` varchar(128) DEFAULT NULL COMMENT '品牌',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_name` (`name`),
  KEY `idx_brand` (`brand`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品表';

INSERT INTO `product` (`category_id`, `name`, `subtitle`, `main_image`, `spec`, `price`, `original_price`, `stock`, `sales`, `status`, `brand`, `detail`) VALUES
(9, 'iPhone 15 Pro Max', 'A17 Pro芯片 钛金属 专业级摄影系统', 'https://example.com/iphone15.jpg', '256GB 原色钛金属', 999900, 1099900, 1000, 567, 1, 'Apple', 'Apple iPhone 15 Pro Max 详情...'),
(9, '华为 Mate 60 Pro', '麒麟芯片 卫星通话 超可靠玄武架构', 'https://example.com/mate60.jpg', '256GB 雅川青', 699900, 749900, 800, 1234, 1, 'HUAWEI', '华为 Mate 60 Pro 详情...'),
(9, '小米14 Ultra', '徕卡光学 Summilux镜头 骁龙8 Gen3', 'https://example.com/mi14.jpg', '256GB 黑色', 649900, 699900, 600, 456, 1, 'Xiaomi', '小米14 Ultra 详情...'),
(10, 'iPad Pro 12.9英寸', 'M2芯片 Liquid视网膜XDR显示屏', 'https://example.com/ipad.jpg', '256GB WiFi版 深空灰', 929900, 999900, 300, 234, 1, 'Apple', 'iPad Pro 详情...'),
(5, 'MacBook Pro 14英寸', 'M3 Pro芯片 18GB内存 512GB', 'https://example.com/macbook.jpg', 'M3 Pro 18GB 512GB', 1699900, 1799900, 200, 123, 1, 'Apple', 'MacBook Pro 详情...'),
(6, '男士商务休闲西装', '修身版型 羊毛混纺 四季可穿', 'https://example.com/suit.jpg', 'XL码 深蓝色', 129900, 199900, 500, 89, 1, 'HLA', '海澜之家男士西装详情...'),
(7, '女士夏季连衣裙', '优雅气质 真丝面料 修身显瘦', 'https://example.com/dress.jpg', 'M码 香槟色', 89900, 129900, 800, 345, 1, 'ONLY', 'ONLY女士连衣裙详情...'),
(8, '每日坚果礼盒', '30包混合坚果 健康零食 每日营养', 'https://example.com/nuts.jpg', '750g 30包', 9900, 14900, 5000, 2345, 1, '三只松鼠', '三只松鼠每日坚果详情...');

DROP TABLE IF EXISTS `stock`;
CREATE TABLE `stock` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_id` bigint NOT NULL COMMENT '商品ID',
  `total` int NOT NULL DEFAULT '0' COMMENT '总库存',
  `available` int NOT NULL DEFAULT '0' COMMENT '可用库存',
  `lock_stock` int NOT NULL DEFAULT '0' COMMENT '锁定库存',
  `sales` int NOT NULL DEFAULT '0' COMMENT '已售',
  `version` int NOT NULL DEFAULT '0' COMMENT '版本号 乐观锁',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存表';

INSERT INTO `stock` (`product_id`, `total`, `available`, `version`) VALUES
(1, 1000, 1000, 0),
(2, 800, 800, 0),
(3, 600, 600, 0),
(4, 300, 300, 0),
(5, 200, 200, 0),
(6, 500, 500, 0),
(7, 800, 800, 0),
(8, 5000, 5000, 0);

DROP TABLE IF EXISTS `order`;
CREATE TABLE `order` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `order_no` varchar(64) NOT NULL COMMENT '订单号',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `total_amount` bigint NOT NULL COMMENT '订单总金额 分',
  `pay_amount` bigint NOT NULL COMMENT '支付金额 分',
  `freight_amount` bigint NOT NULL DEFAULT '0' COMMENT '运费 分',
  `discount_amount` bigint NOT NULL DEFAULT '0' COMMENT '优惠金额 分',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '订单状态 0创建中 1待付款 2已付款 3已发货 4已取消 5已完成',
  `pay_type` tinyint NOT NULL DEFAULT '0' COMMENT '支付方式 0未支付 1支付宝 2微信 3银行卡',
  `pay_time` datetime DEFAULT NULL COMMENT '支付时间',
  `receiver_name` varchar(64) DEFAULT NULL COMMENT '收货人',
  `receiver_phone` varchar(20) DEFAULT NULL COMMENT '收货电话',
  `receiver_address` varchar(500) DEFAULT NULL COMMENT '收货地址',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';

DROP TABLE IF EXISTS `order_item`;
CREATE TABLE `order_item` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `order_id` bigint NOT NULL COMMENT '订单ID',
  `order_no` varchar(64) NOT NULL COMMENT '订单号',
  `product_id` bigint NOT NULL COMMENT '商品ID',
  `product_name` varchar(255) NOT NULL COMMENT '商品名称',
  `product_image` varchar(500) DEFAULT NULL COMMENT '商品图片',
  `price` bigint NOT NULL COMMENT '单价 分',
  `num` int NOT NULL COMMENT '数量',
  `total_price` bigint NOT NULL COMMENT '小计 分',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单明细表';
