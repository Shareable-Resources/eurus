package kyc_const

var RootPath = "/kyc"
var UserServerPath = "/user"
var AdminServerPath = "/admin"

var EndPoint = struct {
	LoginAdminUser      string
	CreateAdminUser     string
	GetKYCCountryList   string
	CreateKYCStatus     string
	SubmitKYCDocument   string
	SubmitKYCApproval   string
	GetKYCStatusOfUser  string
	GetKYCStatusList    string
	UpdateKYCStatus     string
	ResetKYCStatus      string
	ChangeAdminPassword string
	RefreshToken        string
}{
	LoginAdminUser:     "/login",
	CreateAdminUser:    "/createAdminUser",
	GetKYCCountryList:  "/getKYCCountryList",
	CreateKYCStatus:    "/createKYCStatus",
	SubmitKYCDocument:  "/submitKYCDocument",
	SubmitKYCApproval:  "/submitKYCApproval",
	GetKYCStatusOfUser: "/userKYCStatus/{userId}",
	GetKYCStatusList:   "/userKYCStatus",
	UpdateKYCStatus:    "/updateKYCStatus",
	//ApproveKYCImage:  "/approveKYCImage",
	//RejectKYCImage:   "/rejectKYCImage",
	//ApproveKYYStatus: "/approveKYCStatus",
	ResetKYCStatus:      "/resetKYCStatus",
	ChangeAdminPassword: "/changeAdminPassword",
	RefreshToken:        "/refreshToken",
}

var AdminUserKey = "astaxie12798akljzmknm.ahkjkljl;k"

var AdminCommonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
