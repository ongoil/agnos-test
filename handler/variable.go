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
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginStaffResponse struct {
	Token string `json:"token"`
}

type CreatePatientRequest struct {
	NationalID   string `json:"national_id"`
	PassportID   string `json:"passport_id"`
	FirstNameTH  string `json:"first_name_th"`
	MiddleNameTH string `json:"middle_name_th"`
	LastNameTH   string `json:"last_name_th"`
	FirstNameEN  string `json:"first_name_en"`
	MiddleNameEN string `json:"middle_name_en"`
	LastNameEN   string `json:"last_name_en"`
	DateOfBirth  string `json:"date_of_birth"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	Gender       string `json:"gender"`
}
