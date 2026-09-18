package handler

import "github.com/google/uuid"

type RegisterRequest struct {
	Username   string    `json:"username" binding:"required"`
	Password   string    `json:"password" binding:"required,min=6"`
	HospitalID uuid.UUID `json:"hospital_id" binding:"required"`
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
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginStaffResponse struct {
    Token string `json:"token"`
}