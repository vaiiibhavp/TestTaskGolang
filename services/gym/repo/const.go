package repo

const (
	// session bucket
	TOKEN_PREFIX        = `SESSION:`
	OTP_USED_FORCEFULLY = `FORCEFULLY_YES`
	OTP_USED_YES        = `YES`
	OTP_USED_NO         = `NO`
)

const ( //role ids
	SUPER_ADMIN_ROLE_ID    = 1
	SELLER_ROLE_ID         = 2
	CUSTOMER_ROLE_ID       = 3
	SALESMAN_ROLE_ID       = 4
	ADMIN_ROLE_ID          = 5
	SELLER_MANAGER_ROLE_ID = 6
)

const ( //role name
	CUSTOMER_ROLE_NAME = "Customer"
	SELLER_ROLE_NAME   = "Seller"
	RECORDS_PER_PAGE   = 5
)

const (
	OTP_STATUS_SENT          = "SENT"
	OTP_STATUS_NOT_SENT      = "NOT_SENT"
	OTP_STATUS_DELIVERED     = "DELIVERED"
	OTP_STATUS_NOT_DELIVERED = "NOT_DELIVERED"
)
