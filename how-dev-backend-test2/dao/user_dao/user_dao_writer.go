package user_dao

import (
	"context"
	hgorm "github.com/zhifenghao123/how-dev-go-framework/hdev-gorm"
)

type UserWriteDao struct {
	*hgorm.BaseDao `wired:"true"`
}

// Database database
func (d *UserWriteDao) Database() string {
	return "hao_db3"
}

func (d *UserWriteDao) Create(ctx context.Context, user *TableUser) (int64, error) {
	err := d.WithContext(ctx).
		Model(&TableUser{}).
		Create(user).Error
	return user.Id, err
}

func (d *UserWriteDao) DeleteByUserId(ctx context.Context, userId int64) (bool, error) {
	result := d.WithContext(ctx).
		Model(&TableUser{}).
		Where("user_id = ?", userId).
		Delete(&TableUser{})

	/**
	 * 在GORM中，Delete方法的行为是：
	 * 如果找到匹配记录：正常删除，返回nil错误
	 * 如果没有找到匹配记录：不会报错，影响行数为0，仍然返回nil错误
	 * 只有在数据库连接问题、SQL语法错误等真正异常时才会返回错误
	 */
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}
