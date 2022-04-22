package foundation

import (
	"strconv"
)

type ServerReturnCode int64

const (
	Success                      ServerReturnCode = 0
	UnauthorizedAccess           ServerReturnCode = -1
	MethodNotFound               ServerReturnCode = -2
	InvalidSignature             ServerReturnCode = -3
	RequestMalformat             ServerReturnCode = -4
	InternalServerError          ServerReturnCode = -5
	InvalidSession               ServerReturnCode = -6
	LoginTokenExired             ServerReturnCode = -7
	LoginTokenInvalid            ServerReturnCode = -8
	DatabaseError                ServerReturnCode = -9
	SignMatchError               ServerReturnCode = -10
	RefreshTokenError            ServerReturnCode = -11
	RevokeTokenError             ServerReturnCode = -12
	NetworkError                 ServerReturnCode = -13
	RequestParamsValidationError ServerReturnCode = -14
	EthereumError                ServerReturnCode = -15
	TimestampError               ServerReturnCode = -16
	ServerIdUnmatch              ServerReturnCode = -17
	RequestTooFrequenct          ServerReturnCode = -18
	BadRequest                   ServerReturnCode = -20
	UserNotFound                 ServerReturnCode = -21
	InvalidArgument              ServerReturnCode = -22
	RecordNotFound               ServerReturnCode = -23
	UserStatusForbidden          ServerReturnCode = -24
	ServerTokenInvalid           ServerReturnCode = -25
	UniqueViolationError         ServerReturnCode = -26
	NotImplementError            ServerReturnCode = -27
	InvalidEncryption            ServerReturnCode = -28
	InsufficientAccountBalance   ServerReturnCode = -29
	TargetGasAmountReach         ServerReturnCode = -30
	LoginTokenTypeInvalid        ServerReturnCode = -31
	UserAlreadyRegistered        ServerReturnCode = -32
	IncorrectPassword            ServerReturnCode = -33
	HttpStatusCodeBegin          ServerReturnCode = -1000
	HttpStatusCodeEnd            ServerReturnCode = -1999
)

var returnCodeMap = map[ServerReturnCode]string{
	Success:                      "Success",
	UnauthorizedAccess:           "Unauthorized access",
	MethodNotFound:               "Method not found",
	InvalidSignature:             "Invalid signature",
	RequestMalformat:             "Request malformat",
	InternalServerError:          "Internal server error",
	InvalidSession:               "Invalid Session",
	LoginTokenExired:             "Login token expired",
	LoginTokenInvalid:            "Invalid login token",
	DatabaseError:                "Database error",
	SignMatchError:               "Sign match error",
	HttpStatusCodeBegin:          "Http Status Code ",
	HttpStatusCodeEnd:            "Http Status Code end",
	RefreshTokenError:            "Refresh Token Error",
	RevokeTokenError:             "Revoke Token error",
	NetworkError:                 "Network error",
	RequestParamsValidationError: "Request Parameters Validation Error",
	RequestTooFrequenct:          "Request too frequenct",
	BadRequest:                   "Bad request",
	UserNotFound:                 "User not found",
	InvalidArgument:              "Invalid argument",
	RecordNotFound:               "Record not found",
	UserStatusForbidden:          "Forbidden user status",
	UniqueViolationError:         "Record already exists",
	NotImplementError:            "Not implemented",
	InvalidEncryption:            "Invalid encryption",
	InsufficientAccountBalance:   "Insufficient account balance",
	TargetGasAmountReach:         "Target gas amount already reach",
	LoginTokenTypeInvalid:        "Invalid login token type",
	UserAlreadyRegistered:        "User already registered",
	IncorrectPassword:            "Incorrect password",
}

func (me ServerReturnCode) String() string {

	str, ok := returnCodeMap[me]
	if !ok {
		if int64(me) < int64(HttpStatusCodeBegin) && int64(me) > int64(HttpStatusCodeEnd) {
			return returnCodeMap[HttpStatusCodeBegin] + strconv.FormatInt(-(int64(me)+int64(HttpStatusCodeBegin)), 10)
		}

		return strconv.Itoa(int(me))
	}
	return str
}
