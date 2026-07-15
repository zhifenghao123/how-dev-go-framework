package user_dao

import "time"

/**
 * 用户数据访问层
 `
CREATE TABLE `user` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` varchar(64) NOT NULL COMMENT '用户唯一标识',
  `username` varchar(64) NOT NULL COMMENT '用户名',
  `password` varchar(256) NOT NULL COMMENT '密码（加密存储）',
  `email` varchar(64) DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `status` int NOT NULL DEFAULT '1' COMMENT '状态（1-正常，0-禁用）',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id` (`user_id`),
  UNIQUE KEY `idx_username` (`username`),
  UNIQUE KEY `idx_email` (`email`),
  UNIQUE KEY `idx_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';
*/

type TableUser struct {
	Id         int64     `gorm:"column:id;type:bigint;not null;autoIncrement;comment:主键ID"`
	UserId     int64     `gorm:"column:user_id;type:bigint(20);not null;comment:用户唯一标识"`
	Username   string    `gorm:"column:username;type:varchar(64);not null;comment:用户名"`
	Password   string    `gorm:"column:password;type:varchar(256);not null;comment:密码（加密存储）"`
	Email      string    `gorm:"column:email;type:varchar(64);default:null;comment:邮箱"`
	Phone      string    `gorm:"column:phone;type:varchar(20);default:null;comment:手机号"`
	Status     int       `gorm:"column:status;type:int;not null;default:1;comment:状态（1-正常，0-禁用）"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间"`
}

func (TableUser) TableName() string {
	return "user"
}
