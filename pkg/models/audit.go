package models

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"workout-tracker-go-app/pkg/constants"
	"workout-tracker-go-app/pkg/initializers"
)

type Audit struct {
	gorm.Model
	AuditType     string `gorm:"not null"`
	UserId        uint   `gorm:"not null"`
	CustomMessage string
}

func CreateAudit(auditType constants.AuditType, userId uint, customMessage string) {
	audit := Audit{
		AuditType:     auditType.Description,
		UserId:        userId,
		CustomMessage: customMessage,
	}

	result := initializers.DB.Create(&audit)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Failed to create audit: %s for user with ID: %s: %s", auditType.Description, userId, result.Error.Error()))
	}
}
