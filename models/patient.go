package models

import (
	"time"

	"github.com/google/uuid"
)

// Patient - ข้อมูลผู้ป่วยที่ดึงมาจากระบบ HIS ของแต่ละโรงพยาบาล
// แต่ละ record ผูกกับโรงพยาบาลเดียว (multi-tenant by HospitalID)
type Patient struct {
	Model
	HospitalID   uuid.UUID  `json:"hospital_id" gorm:"type:uuid;not null;uniqueIndex:idx_hospital_patient_hn"`
	PatientHN    string     `json:"patient_hn" gorm:"size:50;uniqueIndex:idx_hospital_patient_hn;not null"`
	NationalID   string     `json:"national_id" gorm:"size:13;index"`
	PassportID   string     `json:"passport_id" gorm:"size:20;index"`
	FirstNameTH  string     `json:"first_name_th" gorm:"size:150"`
	MiddleNameTH string     `json:"middle_name_th" gorm:"size:150"`
	LastNameTH   string     `json:"last_name_th" gorm:"size:150"`
	FirstNameEN  string     `json:"first_name_en" gorm:"size:150"`
	MiddleNameEN string     `json:"middle_name_en" gorm:"size:150"`
	LastNameEN   string     `json:"last_name_en" gorm:"size:150"`
	DateOfBirth  *time.Time `json:"date_of_birth" gorm:"type:date"`
	PhoneNumber  string     `json:"phone_number" gorm:"size:20"`
	Email        string     `json:"email" gorm:"size:255"`
	Gender       string     `json:"gender" gorm:"size:1;check:gender IN ('M','F','')"`
	Hospital     Hospital   `json:"-" gorm:"foreignKey:HospitalID"`
}
