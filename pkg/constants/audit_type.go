package constants

type AuditType struct {
	ID          int
	Description string
}

type AuditTypes struct {
	UserCreation           AuditType
	UserLogin              AuditType
	UserDeletion           AuditType
	SendPasswordResetEmail AuditType
	ResetPassword          AuditType
}

func GetAuditTypes() AuditTypes {
	return AuditTypes{
		UserCreation:           AuditType{0, "User created"},
		UserLogin:              AuditType{1, "User logged in"},
		UserDeletion:           AuditType{2, "User deleted"},
		SendPasswordResetEmail: AuditType{3, "User sent password reset email"},
		ResetPassword:          AuditType{4, "User reset password"},
	}
}
