package models

type Hospital struct {
	Model
	Name     string    `json:"name" gorm:"size:255;not null;uniqueIndex"`
	Staffs   []Staff   `json:"-"`
	Patients []Patient `json:"-"`
}
