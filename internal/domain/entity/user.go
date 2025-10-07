package entity

type User struct {
	Email    string `gorm:"primaryKey"`
	Password string `gorm:"not null"`
}
