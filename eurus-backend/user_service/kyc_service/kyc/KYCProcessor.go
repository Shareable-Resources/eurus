package kyc

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/crypto"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/log"
	"eurus-backend/secret"
	kyc_model "eurus-backend/user_service/kyc_service/kyc/model"

	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ethereum/go-ethereum/common"
)

const bucketName string = "eurus-kyc"

func RequestGetKYCCountryList(server *KYCServer, req *request.RequestBase) *response.ResponseBase {
	list, err := DbGetKYCCountryList(server.SlaveDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot request kyc country list")
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot request kyc country list")
	}
	var returnResult *[]kyc_model.KYCCountryCode = list
	return response.CreateSuccessResponse(req, returnResult)
}

// 1. Create KYC Status, first kyc_level would be 1
func RequestCreateKYCStatus(server *KYCServer, req *kyc_model.RequestCreateKYCStatus) *response.ResponseBase {
	var kycStatus = &kyc_model.UserKYCStatus{UserId: req.UserId, KYCCountryCode: req.KYCCountryCode, KYCDoc: req.KYCDoc}
	result, err := DbCreateKYCStatus(server.DefaultDatabase, kycStatus)
	if err != nil {

		log.GetLogger(log.Name.Root).Error("cannot create KYC Status")
		// Duplicate Key Error
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return response.CreateErrorResponse(req, foundation.UniqueViolationError, foundation.UniqueViolationError.String())
		}
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot create KYC Status")
	}
	var returnResult *kyc_model.ResponseCreateKYCStatus = &kyc_model.ResponseCreateKYCStatus{Id: result.Id}
	return response.CreateSuccessResponse(req, returnResult)
}

// 1. Get all user_kyc_status of a specific user
func RequestGetKYCStatusOfUser(server *KYCServer, req *kyc_model.RequestGetKYCStatusOfUser) *response.ResponseBase {
	kycStatusList, user, err := DbGetKYCStatusOfUser(server.SlaveDatabase, req.UserId)

	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot GetKYCStatusOfUser")
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot GetKYCStatusOfUser")
	}
	if kycStatusList == nil || user == nil {
		log.GetLogger(log.Name.Root).Error("No Record is found")
		return response.CreateErrorResponse(req, foundation.RecordNotFound, "No Record is found")
	}

	sess, err := getAwsSession(server.Config)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Cannot open AWS session: ", err)
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}

	svc := s3.New(sess)

	for _, userKycStatus := range *kycStatusList {
		for _, image := range userKycStatus.Images {

			req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(image.ImagePath),
			})

			urlStr, err := req.Presign(60 * time.Minute)
			if err == nil {
				image.ImagePath = urlStr
			}
		}
	}
	returnResult := &kyc_model.ResponseGetKYCStatusOfUser{Data: kycStatusList, KYCLevel: int(user.KycLevel)}
	return response.CreateSuccessResponse(req, returnResult)
}

//1. Get all kyc status list for approval
func RequestGetKYCStatusList(server *KYCServer, req *kyc_model.RequestGetKYCStatusList, page int, pageSize int) *response.ResponseBase {
	result, totalRows, err := DbGetKYCStatusList(server.SlaveDatabase, req.KYCStatus, req.Email, req.WalletAddress, page, pageSize)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot GetKYCStatusList")

		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot GetKYCStatusList")
	}

	returnResult := &kyc_model.ResponseGetKYCStatusList{Data: result, TotalRows: totalRows}
	return response.CreateSuccessResponse(req, returnResult)
}

// CS admin login API
func LoginAdminUser(server *KYCServer, req *kyc_model.RequestLoginAdminUser) *response.ResponseBase {
	var token auth_base.ILoginToken
	var decryptedPassword string
	var err error
	var userInDb *kyc_model.AdminUser
	userInDb, err = DbGetAdminUserByUsername(server.SlaveDatabase, req.Username)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Admin user does not exist")
		return response.CreateErrorResponse(req, foundation.DatabaseError, "Admin user does not exist")
	}
	decryptedPassword, err = decrypt(userInDb.Password, server.Config.AdminAESKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Cannot decrypt")
		return response.CreateErrorResponse(req, foundation.DatabaseError, "Cannot decrypt")
	}
	if decryptedPassword != req.Password {
		log.GetLogger(log.Name.Root).Error("Password mismatch")
		return response.CreateErrorResponse(req, foundation.SignMatchError, "Password mismatch")
	}

	dataToBeSavedInDbWhenCreatingToken := &kyc_model.AdminUser{Username: req.AdminUser.Username}
	jsonByte, err := json.Marshal(dataToBeSavedInDbWhenCreatingToken)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot json stringify")
	}
	token, err = server.AuthClient.GenerateLoginToken(string(jsonByte))
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot GenerateLoginToken")
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot GenerateLoginToken")
	}
	var returnResult auth_base.ILoginToken = token
	return response.CreateSuccessResponse(req, returnResult)
}

// This function will update all kyc_image based on the request array objects
// If any of the kyc_image is rejected, then change kyc_status to rejected
// Admin user can enter reject reason
// The API will check whether the image updated property is true or not
// This function is for Admin User
// Only Valid for kycStatus.KYCStatus == kyc_model.KYCStatusApproved or kycStatus.KYCStatus == kyc_model.KYCStatusWaitingForResubmit
func UpdateKYCStatus(server *KYCServer, req *kyc_model.RequestUpdateKYCStatus) *response.ResponseBase {
	//Check Status argument
	if req.KYCStatus != kyc_model.KYCStatusApproved && req.KYCStatus != kyc_model.KYCStatusWaitingForResubmit {
		return response.CreateErrorResponse(req, foundation.InvalidArgument, "Invalid Argument for kyc status")
	}
	for _, image := range req.Images {
		if image.Status != kyc_model.KYCImageStatusApproved && image.Status != kyc_model.KYCImageStatusWaitingForResubmit {
			return response.CreateErrorResponse(req, foundation.InvalidArgument, "Invalid Argument for image status")
		}
	}
	levelPromoted, user, err := DbUpdateKYCStatusAndRelatedImages(server.DefaultDatabase, &req.UserKYCStatus)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot UpdateKYCStatus")
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	if levelPromoted > 0 {
		err = server.scProcessor.SetUserKYCLevel(common.HexToAddress(user.WalletAddress), levelPromoted)
		if err != nil {
			log.GetLogger(log.Name.Root).Error("Set user KYC level failed. User id: ", req.UserId)
			return response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
		}
	}

	res := new(kyc_model.ResponseUpdateKYCStatus)
	res.PromotedLevel = levelPromoted

	return response.CreateSuccessResponse(req, res)
}

// This function will reset user_kyc_status to pending, retry_count to 0, used by CS Admin
func ResetKYCStatus(server *KYCServer, req *kyc_model.RequestResetKYCStatus) *response.ResponseBase {
	var operatorId = database.NullString{sql.NullString{Valid: true, String: req.AdminUser.Username}}
	err := DbResetKYCStatus(server.DefaultDatabase, req.Id, operatorId)

	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot ResetKYCStatus")
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())

	}

	return response.CreateSuccessResponse(req, nil)
}

//This function will be used by user to submit KYC Request, used by normal user
func SubmitKYCApproval(server *KYCServer, req *kyc_model.RequestSubmitKYCApproval) *response.ResponseBase {

	var err error
	req.KYCStatus = kyc_model.KYCStatusWaitingForApproval
	err = DbUpdateKYCStatusToApproval(server.DefaultDatabase, req.Id, req.KYCStatus)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot SubmitKYCApproval")
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())

	}
	return response.CreateSuccessResponse(req, nil)
}

func ProcessKYCImage(server *KYCServer, reader *multipart.Reader) (error, bool) {
	var err error
	for i := 0; i < server.Config.GetRetryCount(); i++ {
		var userKycStatus *kyc_model.UserKYCStatus

		part1, err := reader.NextPart()
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Unable to get multipart: ", err.Error())
			return err, false
		}
		request, err := io.ReadAll(part1)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Unable to read part 1: ", err.Error())
			return err, false
		}
		reqObj := new(kyc_model.KYCSubmitImage)
		err = json.Unmarshal(request, &reqObj)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Invalid request: ", err.Error())
			return err, false
		}

		if strings.HasPrefix(reqObj.FileExtension, ".") && len(reqObj.FileExtension) > 1 {
			reqObj.FileExtension = reqObj.FileExtension[1:]
		}

		part2, err := reader.NextPart()
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("No image found. Error: ", err, " request: ", string(request))
			return err, false
		}

		userKycStatus, err = DbGetKYCStatusById(&server.DefaultDatabase.ReadOnlyDatabase, reqObj.KYCStatusId)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln(err)
			time.Sleep(time.Second * time.Duration(server.Config.RetryInterval))
			continue
		}

		if userKycStatus.Id == 0 {
			log.GetLogger(log.Name.Root).Errorln("kyc status id not found, kyc document ignored. status id: ", reqObj.KYCStatusId, " input user id: ", reqObj.UserId)
			return nil, false
		} else if userKycStatus.UserId != reqObj.UserId {
			log.GetLogger(log.Name.Root).Errorln("Unmatch kyc status id and use id, kyc document ignored. status id: ", reqObj.KYCStatusId, " input user id: ", reqObj.UserId)
			return nil, false
		} else {
			kycImage, err := DbInsertKYCImage(server.DefaultDatabase, reqObj.KYCStatusId, reqObj.ImageType)
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln(err)
				time.Sleep(time.Second * time.Duration(server.Config.RetryInterval))
				continue
			}
			//TODO upload image to S3
			fileName := part2.FileName()
			index := strings.LastIndex(fileName, ".")
			finalFileName := fileName
			if index >= 0 && index+1 < len(fileName) {
				finalFileName = fileName[:index] + "_" + strconv.Itoa(kycImage.ImageSeq) + fileName[index:]
			} else {
				log.GetLogger(log.Name.Root).Errorln("Invalid file name: ", fileName, " User id: ", reqObj.UserId)
				return errors.New("Invalid file name"), false
			}
			file, err := os.Create(finalFileName)
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("Unable to write file: ", err, " User id: ", reqObj.UserId)
				return err, true
			}
			defer file.Close()
			var fileKeyName string
			var isUploaded bool
			for j := 0; i < server.Config.GetRetryCount(); j++ {
				if !isUploaded {
					fileKeyName, err = uploadImageToS3Bucket(server.Config, part2, finalFileName)
					if err != nil {
						log.GetLogger(log.Name.Root).Errorln("Unable to upload to S3: ", err, " User id: ", reqObj.UserId, " filename: ", finalFileName)
						time.Sleep(time.Second * time.Duration(server.Config.RetryInterval))
						continue
					}
					isUploaded = true
				}

				err = DbUpdateKYCImageStatusToUploaded(server.DefaultDatabase, fileKeyName, kycImage)
				if err != nil {
					log.GetLogger(log.Name.Root).Errorln("Unable update DB after upload to S3: ", err, " User id: ", reqObj.UserId, " filename: ", finalFileName)
					time.Sleep(time.Second * time.Duration(server.Config.RetryInterval))
					continue
				}
				break
			}

			if err != nil {
				return err, true
			}

			break
		}
	}
	return err, true
}

// This function will reset user_kyc_status to pending, retry_count to 0, used by CS Admin
func CreateAdminUser(server *KYCServer, req *kyc_model.RequestCreateAdminUser) *response.ResponseBase {
	//TODO - submit for reset KYC Status
	encryptedPassword, err := encrypt(req.Password, server.Config.AdminAESKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Cannot Encrypt")
		return response.CreateErrorResponse(req, foundation.DatabaseError, "Cannot Encrypt")
	}
	req.AdminUser.Password = encryptedPassword
	result, err := DbCreateAdminUser(server.DefaultDatabase, &req.AdminUser)
	if err != nil {

		log.GetLogger(log.Name.Root).Error("Cannot create Admin User")
		// Duplicate Key Error
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return response.CreateErrorResponse(req, foundation.UniqueViolationError, foundation.UniqueViolationError.String())
		}
		return response.CreateErrorResponse(req, foundation.DatabaseError, "Cannot create Admin User")
	}
	var returnResult kyc_model.ResponseCreateAdminUser = kyc_model.ResponseCreateAdminUser{Username: result.Username}
	return response.CreateSuccessResponse(req, returnResult)
}

func uploadImageToS3Bucket(config *KYCConfig, src io.Reader, fileName string) (string, error) {

	sess, err := getAwsSession(config)
	if err != nil {
		return "", err
	}
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)
	fileKey := "images/" + secret.Tag + "/" + fileName
	// Upload the file to S3.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:               aws.String(bucketName),
		Key:                  aws.String(fileKey),
		ServerSideEncryption: aws.String("AES256"),
		Body:                 src,
	})
	if err != nil {
		return "", err
	}

	return fileKey, nil
}

func ChangeAdminPassword(server *KYCServer, req *kyc_model.RequestChangeAdminPassword) *response.ResponseBase {
	//TODO - submit for reset KYC Status
	var decryptedPassword string
	var err error
	var userInDb *kyc_model.AdminUser
	userInDb, err = DbGetAdminUserByUsername(&server.DefaultDatabase.ReadOnlyDatabase, req.Username)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Admin user does not exist")
		return response.CreateErrorResponse(req, foundation.DatabaseError, "Admin user does not exist")
	}
	decryptedPassword, err = decrypt(userInDb.Password, server.Config.AdminAESKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Cannot decrypt")
		return response.CreateErrorResponse(req, foundation.DatabaseError, "Cannot decrypt")
	}
	if decryptedPassword != req.Password {
		log.GetLogger(log.Name.Root).Error("Password mismatch")
		return response.CreateErrorResponse(req, foundation.SignMatchError, "Password mismatch")
	}
	req.AdminUser.Password, err = encrypt(req.NewPassword, server.Config.AdminAESKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Cannot Encrypt")
		return response.CreateErrorResponse(req, foundation.DatabaseError, "Cannot Encrypt")
	}
	result, err := DbChangeAdminPassword(server.DefaultDatabase, &req.AdminUser)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Cannot change admin user password")
		return response.CreateErrorResponse(req, foundation.DatabaseError, "Cannot change admin user password")
	}
	var returnResult kyc_model.ResponseChangeAdminPassword = kyc_model.ResponseChangeAdminPassword{Username: result.Username}
	return response.CreateSuccessResponse(req, returnResult)
}

func encrypt(password string, adminAesKey string) (string, error) {

	aesKey, err := base64.StdEncoding.DecodeString(adminAesKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to decode aes key: ", err)
		return "", err
	}
	cipher, err := crypto.EncryptAES([]byte(password), aesKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to encrypt by AES: ", err)
		return "", err
	}

	return cipher, err

}

func decrypt(encryptedPassword string, adminAesKey string) (string, error) {

	aesKey, err := base64.StdEncoding.DecodeString(adminAesKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to decode aes key: ", err)
		return "", err
	}
	data, _ := base64.StdEncoding.DecodeString(encryptedPassword)
	base64Encoded, _ := crypto.DecryptAES(data, aesKey)
	result, _ := base64.StdEncoding.DecodeString(base64Encoded)
	return string(result), err
}

func getAwsSession(config *KYCConfig) (*session.Session, error) {
	var verbose bool = true

	sess, err := session.NewSession(&aws.Config{
		Region:                        aws.String(config.S3BucketZone),
		Credentials:                   credentials.NewStaticCredentials(config.AwsAccessKeyId, config.AwsAccessSecretAccessKey, ""),
		CredentialsChainVerboseErrors: &verbose,
	})

	return sess, err

}

func RefreshToken(authClient auth_base.IAuth, req *kyc_model.RefreshTokenRequest) *response.ResponseBase {
	loginToken, err := authClient.RefreshLoginToken(req.LoginToken.GetToken())
	if err != nil {
		return response.CreateErrorResponse(req, err.GetReturnCode(), err.Error())
	}

	res := kyc_model.RefreshTokenResponse{}
	res.ExpireTime = loginToken.GetExpiredTime()
	res.CreatedDate = loginToken.GetCreatedDate()
	res.LastModifiedDate = loginToken.GetLastModifiedDate()
	res.Token = loginToken.GetToken()

	return response.CreateSuccessResponse(req, res)
}
