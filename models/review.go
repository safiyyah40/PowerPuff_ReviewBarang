package models

import "time"

type Review struct {
    ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
    UsernameID    int       `json:"username_id" binding:"required"`
    Username      string    `json:"username" binding:"required"`
    IsAnonymous   bool      `json:"is_anonymous"`
    ProductID     int       `json:"product_id" binding:"required"`
    ProductName   string    `json:"product_name" binding:"required"`
    Category      string    `json:"category" binding:"required"`
    Rating        int       `json:"rating" binding:"gte=1,lte=5"`
    TextReview    string    `json:"text_review" binding:"required"`
    Likes         int       `json:"likes" binding:"required"`
    CreatedAt     time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP" binding:"required"`
}
