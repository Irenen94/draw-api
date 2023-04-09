package identity_status

type IdentityStatus int64

const (
	SUCCESS IdentityStatus = 0
	FAIL    IdentityStatus = 1
	OTHER   IdentityStatus = 2
)
