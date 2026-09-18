package models

import (
	"time"

	"github.com/google/uuid"
)

type Hospital struct {
	Model
	Name     string    `json:"name" gorm:"size:255;not null"` // ชื่อโรงพยาบาล
	Staffs   []Staff   `json:"staffs,omitempty"`              // เจ้าหน้าที่ทั้งหมดในโรงพยาบาลนี้
	Patients []Patient `json:"patients,omitempty"`            // ผู้ป่วยทั้งหมดในโรงพยาบาลนี้
}

// Patient - ข้อมูลผู้ป่วยที่ดึงมาจากระบบ HIS ของแต่ละโรงพยาบาล
// แต่ละ record ผูกกับโรงพยาบาลเดียว (multi-tenant by HospitalID)
type Patient struct {
	Model
	HospitalID   uuid.UUID  `json:"hospital_id" gorm:"type:uuid;not null;uniqueIndex:idx_hospital_patient_hn"` // FK อ้างถึง Hospital, ใช้ filter สิทธิ์การเข้าถึงของ staff
	PatientHN    string     `json:"patient_hn" gorm:"size:50;uniqueIndex:idx_hospital_patient_hn;not null"`    // เลข HN จากระบบ HIS ต้นทาง unique ต่อโรงพยาบาล (คนละ HN กันได้ระหว่างโรงพยาบาล)
	NationalID   string     `json:"national_id" gorm:"size:13;index"`                                          // เลขบัตรประชาชนไทย 13 หลัก ใช้เป็นเงื่อนไขค้นหา (Hospital A API: param "id")
	PassportID   string     `json:"passport_id" gorm:"size:20;index"`                                          // เลขพาสปอร์ต ใช้แทน national_id กรณีเป็นชาวต่างชาติ
	FirstNameTH  string     `json:"first_name_th" gorm:"size:150"`
	MiddleNameTH string     `json:"middle_name_th" gorm:"size:150"`
	LastNameTH   string     `json:"last_name_th" gorm:"size:150"`
	FirstNameEN  string     `json:"first_name_en" gorm:"size:150"`
	MiddleNameEN string     `json:"middle_name_en" gorm:"size:150"`
	LastNameEN   string     `json:"last_name_en" gorm:"size:150"`
	DateOfBirth  *time.Time `json:"date_of_birth" gorm:"type:date"`                  // วันเกิด เก็บเป็น date จริง รองรับ query ช่วงวันที่/คำนวณอายุ ใช้ pointer เพราะบางระบบ HIS อาจไม่ส่งข้อมูลนี้มา (null ได้)
	PhoneNumber  string     `json:"phone_number" gorm:"size:20"`                     // เบอร์โทรศัพท์ (รองรับรูปแบบมี +66 หรือ - คั่น)
	Email        string     `json:"email" gorm:"size:255"`                           // อีเมล
	Gender       string     `json:"gender" gorm:"size:10"`                           // เพศ ค่าที่รับได้คือ "M" หรือ "F" ตามสเปค
	Hospital     Hospital   `json:"hospital,omitempty" gorm:"foreignKey:HospitalID"` // relation กลับไปที่ Hospital (preload เฉพาะตอนจำเป็น)
}

// Staff - เจ้าหน้าที่โรงพยาบาล ใช้ login เข้าระบบเพื่อค้นหาข้อมูลผู้ป่วย
// สิทธิ์การเข้าถึงถูกจำกัดเฉพาะผู้ป่วยในโรงพยาบาลเดียวกับตัวเอง (ผ่าน HospitalID)
type Staff struct {
	Model
	HospitalID   uuid.UUID `json:"hospital_id" gorm:"type:uuid;not null;index"`     // FK อ้างถึง Hospital ที่ staff คนนี้สังกัด
	Username     string    `json:"username" gorm:"size:50;uniqueIndex;not null"`    // ใช้สำหรับ login ต้อง unique ทั้งระบบ
	PasswordHash string    `json:"-" gorm:"size:255;not null"`                      // เก็บ hash เท่านั้น (bcrypt/argon2) ห้ามส่งออกใน JSON response
	Hospital     Hospital  `json:"hospital,omitempty" gorm:"foreignKey:HospitalID"` // relation กลับไปที่ Hospital
}
