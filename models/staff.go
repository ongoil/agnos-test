package models

import (
	"github.com/google/uuid"
)

// Staff - เจ้าหน้าที่โรงพยาบาล ใช้ login เข้าระบบเพื่อค้นหาข้อมูลผู้ป่วย
// สิทธิ์การเข้าถึงถูกจำกัดเฉพาะผู้ป่วยในโรงพยาบาลเดียวกับตัวเอง (ผ่าน HospitalID)
type Staff struct {
	Model
	HospitalID   uuid.UUID `json:"hospital_id" gorm:"type:uuid;not null;index"`
	Username     string    `json:"username" gorm:"size:50;uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"size:255;not null"`
	Hospital     Hospital  `json:"-" gorm:"foreignKey:HospitalID"`
}
