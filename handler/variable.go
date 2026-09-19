package handler

import "github.com/google/uuid"

type RegisterRequest struct {
	Username   string    `json:"username" binding:"required"`
	Password   string    `json:"password" binding:"required,min=6"`
	Hospital   string    `json:"hospital"`
	HospitalID uuid.UUID `json:"hospital_id"`
}

type CreateStaffResponse struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	HospitalID uuid.UUID `json:"hospital_id"`
}

type bodyHospitalRequest struct {
	HospitalName string `json:"hospital_name" binding:"required"`
}

type LoginStaffRequest struct {
	Username   string    `json:"username" binding:"required"`
	Password   string    `json:"password" binding:"required"`
	Hospital   string    `json:"hospital"`
	HospitalID uuid.UUID `json:"hospital_id"`
}

type LoginStaffResponse struct {
	Token string `json:"token"`
}

type PatientSearchRequest struct {
	NationalID  string `form:"national_id"`
	PassportID  string `form:"passport_id"`
	FirstName   string `form:"first_name"`
	MiddleName  string `form:"middle_name"`
	LastName    string `form:"last_name"`
	DateOfBirth string `form:"date_of_birth"`
	PhoneNumber string `form:"phone_number"`
	Email       string `form:"email"`
}
