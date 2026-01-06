package entities

import "gorm.io/gorm"

type ChatMessageReaction struct {
	gorm.Model
	MessageID uint   `gorm:"index;uniqueIndex:ux_msg_user_emoji"`
	UserID    uint   `gorm:"index;uniqueIndex:ux_msg_user_emoji"`
	Emoji     string `gorm:"size:32;uniqueIndex:ux_msg_user_emoji"`
}
