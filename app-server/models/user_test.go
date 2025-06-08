package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdateUserFields(t *testing.T) {
	createTime, _ := time.Parse("2006-01-02 15:04:05", "2024-01-01 10:00:00")
	origin := &User{
		Id:         1,
		Username:   "olduser",
		Password:   "oldpass",
		Name:       "Old Name",
		Gender:     "M",
		Addr:       "Old Addr",
		HeadImg:    "old.png",
		CreateTime: createTime,
		Lat:        10.0,
		Lon:        20.0,
	}

	// 只更新部分字段，ID为零值，不应被更新
	update := User{
		Username: "newuser",
		Name:     "New Name",
		Addr:     "New Addr",
		Lat:      30.0,
	}

	origin.UpdateUserFields(update)

	assert.Equal(t, int64(1), origin.Id) // ID不变
	assert.Equal(t, "newuser", origin.Username)
	assert.Equal(t, "oldpass", origin.Password) // 未更新
	assert.Equal(t, "New Name", origin.Name)
	assert.Equal(t, "M", origin.Gender) // 未更新
	assert.Equal(t, "New Addr", origin.Addr)
	assert.Equal(t, "old.png", origin.HeadImg)     // 未更新
	assert.Equal(t, createTime, origin.CreateTime) // 未更新
	assert.Equal(t, 30.0, origin.Lat)
	assert.Equal(t, 20.0, origin.Lon) // 未更新
}
