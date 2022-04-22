package user

import kyc_model "eurus-backend/user_service/kyc_service/kyc/model"

type UserStore struct {
	KYCCountryCodeList []*kyc_model.KYCCountryCode
}
