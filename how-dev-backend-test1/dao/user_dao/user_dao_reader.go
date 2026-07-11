package user_dao

import (
	"context"
	"errors"
	hgorm "github.com/zhifenghao123/how-dev-go-framework/hdev-gorm"
	"gorm.io/gorm"
)

type UserReadDao struct {
	*hgorm.BaseDao `wired:"true"`
}

// Database database
func (d *UserReadDao) Database() string {
	return "hao_db3"
}

func (d *UserReadDao) GetUserById(ctx context.Context, id int64) (TableUser, error) {
	var user TableUser
	err := d.WithContext(ctx).
		Model(&TableUser{}).
		Where("id = ?", id).
		First(&user).
		Error
	// 如果是记录不存在的错误，返回空结构体和 nil
	// GORM 的 First() 方法，当查询记录不存在时，会返回 gorm.ErrRecordNotFound 错误。这里兼容处理一下gorm.ErrRecordNotFound和真正的DB错误。
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return TableUser{}, nil
	}
	return user, err
}

func (d *UserReadDao) GetUserByUserId(ctx context.Context, userId int64) (TableUser, error) {
	var user TableUser
	err := d.WithContext(ctx).
		Model(&TableUser{}).
		Where("user_id = ?", userId).
		First(&user).
		Error
	// 如果是记录不存在的错误，返回空结构体和 nil
	// GORM 的 First() 方法，当查询记录不存在时，会返回 gorm.ErrRecordNotFound 错误。这里兼容处理一下gorm.ErrRecordNotFound和真正的DB错误。
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return TableUser{}, nil
	}
	return user, err
}
