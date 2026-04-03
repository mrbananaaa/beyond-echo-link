package apperror

type Type string

const (
	TypeValidation     Type = "validation_error"
	TypeBusiness       Type = "business_error"
	TypeInfrastructure Type = "infrastructure_error"
	TypeExternal       Type = "externall_error"
	TypeDB             Type = "db_error"
)
