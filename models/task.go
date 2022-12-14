package models

import (
	"time"

	"gorm.io/datatypes"
)

type Task struct {
	Id    int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Date  time.Time      `json:"date"`
	Tasks datatypes.JSON `json:"tasks"`
	// Tasks     string         `json:"tasks"`
	// StartTime datatypes.Time `json:"start_time"`
	// EndTime   datatypes.Time `json:"end_time"`
}
