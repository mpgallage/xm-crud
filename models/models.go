package models

import (
	"github.com/gofrs/uuid"
)

type User struct {
	Username string `gorm:"type:varchar(100);not null;unique"`
	Password string `gorm:"type:varchar(100)"`
}

type JwtToken struct {
	Token string `json:"token"`
}

type Exception struct {
	Message string `json:"message"`
}

type CompanyType string

const (
	Corporations       CompanyType = "Corporations"
	NonProfit          CompanyType = "NonProfit"
	Cooperative        CompanyType = "Cooperative"
	SoleProprietorship CompanyType = "SoleProprietorship"
)

type Company struct {
	ID            uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name          string      `gorm:"type:varchar(15);not null;unique"`
	Description   string      `gorm:"type:varchar(3000)"`
	EmployeeCount int         `gorm:"type:int;not null"`
	Registered    bool        `gorm:"type:boolean;not null"`
	Type          CompanyType `gorm:"type:varchar(50);not null"` //Corporations | NonProfit | Cooperative | Sole Proprietorship
}

var companyTypes = map[CompanyType]bool{
	Corporations:       true,
	NonProfit:          true,
	Cooperative:        true,
	SoleProprietorship: true,
}

func IsValidCompanyType(t CompanyType) bool {
	return companyTypes[t]
}
