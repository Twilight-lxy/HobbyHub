package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateUserFields(t *testing.T) {
	origin := &User{
		ID:         1,
		Username:   "olduser",
		Password:   "oldpass",
		Name:       "Old Name",
		Gender:     "M",
		Addr:       "Old Addr",
		HeadImg:    "old.png",
		CreateTime: "2024-01-01",
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

	assert.Equal(t, int64(1), origin.ID) // ID不变
	assert.Equal(t, "newuser", origin.Username)
	assert.Equal(t, "oldpass", origin.Password) // 未更新
	assert.Equal(t, "New Name", origin.Name)
	assert.Equal(t, "M", origin.Gender) // 未更新
	assert.Equal(t, "New Addr", origin.Addr)
	assert.Equal(t, "old.png", origin.HeadImg)       // 未更新
	assert.Equal(t, "2024-01-01", origin.CreateTime) // 未更新
	assert.Equal(t, 30.0, origin.Lat)
	assert.Equal(t, 20.0, origin.Lon) // 未更新
}
