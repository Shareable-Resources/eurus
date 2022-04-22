package api

import (
	"bytes"
	"encoding/json"

	"eurus-backend/foundation"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/log"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const UserAgent = "Eurus/1.0"

//API package singleton
func SetApiLogger(logger *logrus.Logger) {
	apiLogger = logger
}

var apiLogger *logrus.Logger

type RequestResponse struct {
	Req          request.IRequest
	Res          response.IResponse
	RetrySetting foundation.IRetrySetting
}

func NewRequestResponse(req request.IRequest, res response.IResponse) *RequestResponse {
	reqRes := new(RequestResponse)
	reqRes.Req = req
	reqRes.Res = res
	return reqRes
}
func (me *RequestResponse) RequestToJson() ([]byte, error) {
	str, err := json.Marshal(me.Req)
	if err != nil {
		return nil, err
	}
	return str, nil
}

func (me *RequestResponse) RequestJsonToModel(jsonData []byte) error {
	err := json.Unmarshal(jsonData, me.Req)
	return err
}

func (me *RequestResponse) ResponseToJson() ([]byte, error) {
	return json.Marshal(me.Res)
}

//unmarshal json to the me.Res
func (me *RequestResponse) ResponseJsonToModel(jsonStr []byte) error {
	err := json.Unmarshal(jsonStr, me.Res)
	return err
}

func SendApiRequest(targetUrl url.URL, reqRes *RequestResponse, authClient auth_base.IAuth) (*RequestResponse, error) {
	var err error
	var retryCount int = 1
	var retryInterval time.Duration = 1

	if reqRes.RetrySetting != nil {
		retryCount = reqRes.RetrySetting.GetRetryCount()
		retryInterval = reqRes.RetrySetting.GetRetryInterval()
	}

	if reqRes.Req.GetNonce() == "" {
		reqRes.Req.SetNonce(uuid.New().String())
	}
	for i := 0; i < retryCount; i++ {
		reqRes, err = sendApiRequestInternal(targetUrl, reqRes, authClient)
		if err != nil {
			if i < retryCount-1 {
				time.Sleep(time.Second * retryInterval)
			}
			continue
		} else {
			break
		}
	}
	return reqRes, err
}

func sendApiRequestInternal(targetUrl url.URL, reqRes *RequestResponse, authClient auth_base.IAuth) (*RequestResponse, error) {
	var err error

	var content []byte
	content, err = reqRes.RequestToJson()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Requst data error ", err.Error())
		return nil, err
	}

	var req *http.Request
	var resp *http.Response
	switch reqRes.Req.GetMethod() {

	case http.MethodPost:

		targetUrl.Path = reqRes.Req.GetRequestPath()

		req, err = http.NewRequest(reqRes.Req.GetMethod(), targetUrl.String(), strings.NewReader(string(content)))
		if err != nil {
			return reqRes, err
		}
		addAuthorizationHeader(req, authClient)
		addUserAgentHeader(req)
		req.Header.Add("Content-Type", "application/json")
		resp, err = http.DefaultClient.Do(req)
	case http.MethodGet:

		targetUrl.Path = reqRes.Req.GetRequestPath()

		dict := make(map[string]interface{})

		err = json.Unmarshal(content, &dict)
		if err != nil {
			return reqRes, err
		}
		var queryStr string
		for key, val := range dict {
			valStr := fmt.Sprintf("%v", val)
			queryStr += fmt.Sprintf("%v=%v&", url.QueryEscape(key), url.QueryEscape(valStr))
		}
		targetUrl.RawQuery = queryStr
		req, err = http.NewRequest(reqRes.Req.GetMethod(), targetUrl.String(), strings.NewReader(string(content)))
		if err != nil {
			return reqRes, err
		}

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		addAuthorizationHeader(req, authClient)
		addUserAgentHeader(req)
		resp, err = http.DefaultClient.Do(req)

	case http.MethodDelete:
		targetUrl.Path = reqRes.Req.GetRequestPath()
		// Create client
		client := &http.Client{}

		// Create request
		req, err := http.NewRequest("DELETE", targetUrl.String(), nil)
		if err != nil {
			return reqRes, err
		}
		addAuthorizationHeader(req, authClient)
		addUserAgentHeader(req)
		// Fetch Request
		resp, err = client.Do(req)
		if err != nil {
			return reqRes, err
		}
	default:
		return reqRes, errors.New("Method not supported")
	}

	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to send request to", targetUrl.String(), ".", err.Error())
		return reqRes, err
	}
	_, err = HttpResponseToModel(resp, reqRes.Res)

	return reqRes, err
}

func addAuthorizationHeader(req *http.Request, authClient auth_base.IAuth) {
	if authClient != nil && authClient.IsLoggedIn() {
		req.Header.Add("Authorization", "Bearer "+authClient.GetLoginToken())
	}
}

func addUserAgentHeader(req *http.Request) {
	req.Header.Add("User-Agent", UserAgent)
}

func HttpResponseToModel(resp *http.Response, res response.IResponse) (response.IResponse, error) {
	if resp.StatusCode != http.StatusOK {
		var returnCode int64 = int64(foundation.HttpStatusCodeBegin) - int64(resp.StatusCode)
		log.GetLogger(log.Name.Root).Errorln("Config service response error status: ", resp.StatusCode, resp.Status)
		return nil, foundation.NewErrorWithMessage(foundation.ServerReturnCode(returnCode), resp.Status)
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, res)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func HttpRequestToModel(httpReq *http.Request, req request.IRequest, withLoginToken bool) error {
	var contentLength int64 = -1
	contentLengthStr := httpReq.Header.Get("Content-Length")
	if contentLengthStr != "" {
		length, err := strconv.ParseInt(contentLengthStr, 10, 64)
		if err == nil {
			contentLength = length
		}
	}

	contentType := httpReq.Header.Get("Content-Type")
	if contentType != "" && !strings.HasPrefix(contentType, "application/json") && !strings.HasPrefix(contentType, "text/plain") && !strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		return errors.New("Content type not suported")
	}
	if contentLength > 1024*8 {
		return errors.New("Content to large")
	}
	var data []byte
	var err error
	if httpReq.Method == http.MethodGet {
		query := httpReq.URL.Query()
		var queryMap map[string]interface{} = make(map[string]interface{})
		for key, value := range query {
			if len(value) == 0 {
				queryMap[key] = ""
			} else if len(value) == 1 {
				queryMap[key] = value[0]
			} else {
				queryMap[key] = value
			}
		}
		data, err = json.Marshal(queryMap)
		if err != nil {
			return errors.Wrap(err, "Unable to marshal query string")
		}
	} else {
		data, err = ioutil.ReadAll(httpReq.Body)
		if err != nil {
			return errors.Wrap(err, "Reading request body failed")
		}
	}

	if apiLogger != nil {
		apiLogger.Infof("Request received. Remote address IP: %s, Method: %s, URL: %s, Body: %s, Content length: %d", httpReq.RemoteAddr, httpReq.Method, httpReq.RequestURI, string(data), contentLength)
	}

	if contentLength > 0 && len(data) < int(contentLength) {
		return errors.New("Invalid content length received")
	}
	if len(data) == 0 {
		if withLoginToken {
			err = RetrieveLoginToken(httpReq, req)
			if err != nil {
				return errors.Wrap(err, "RetrieveLoginToken failed")
			}
			validate := validator.New()
			err = validate.Struct(req)

			return err
		} else {
			validate := validator.New()
			err = validate.Struct(req)
			return err
		}
	}

	err = json.Unmarshal(data, req)

	if err == nil && withLoginToken {
		err = RetrieveLoginToken(httpReq, req)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return errors.Wrap(err, "Unmarshal request body failed: "+string(data))
	}

	validate := validator.New()
	err = validate.Struct(req)

	return err
}

func RetrieveLoginToken(httpReq *http.Request, req request.IRequest) error {
	loginTokenObj, ok := httpReq.Context().Value("loginToken").(auth_base.ILoginToken)
	if ok {
		req.SetLoginToken(loginTokenObj)
	} else {
		return errors.New("Invalid Login token format")
	}
	return nil
}

func RequestToModel(httpReq *http.Request, req request.IRequest) response.IResponse {
	err := HttpRequestToModel(httpReq, req, true)

	if err != nil {
		msg := formatValidationError(err)
		return response.CreateErrorResponse(req, foundation.InternalServerError, msg)

	}
	return nil
}

func RequestToModelNoLoginToken(httpReq *http.Request, req request.IRequest) response.IResponse {
	err := HttpRequestToModel(httpReq, req, false)

	if err != nil {
		errStr := formatValidationError(err)
		return response.CreateErrorResponse(req, foundation.InternalServerError, errStr)
	}
	return nil
}

func formatValidationError(err error) string {
	var msg string
	if valErr, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range valErr {
			switch fieldErr.Tag() {
			case "required":
				msg += fmt.Sprintf("%s is a mandatory field. ", fieldErr.Field())
			case "min":
				switch fieldErr.Type().Kind() {
				case reflect.String:
					msg += fmt.Sprintf("%s length should be greater or equal to %s. ", fieldErr.Field(), fieldErr.Param())
				default:
					msg += fmt.Sprintf("%s value should be greater or equal to %s. ", fieldErr.Field(), fieldErr.Param())
				}
			case "max":
				switch fieldErr.Type().Kind() {
				case reflect.String:
					msg += fmt.Sprintf("%s length should be less or equal to %s. ", fieldErr.Field(), fieldErr.Param())
				default:
					msg += fmt.Sprintf("%s value should be less or equal to %s. ", fieldErr.Field(), fieldErr.Param())
				}
			case "len":
				msg += fmt.Sprintf("%s should be in %s characters", fieldErr.Field(), fieldErr.Param())
			case "excludesall":
				msg += fmt.Sprintf("%s should not contain any character of '%s'. ", fieldErr.Field(), fieldErr.Param())
			case "excludesrune":
				msg += fmt.Sprintf("%s should not contain any character of '%s'. ", fieldErr.Field(), fieldErr.Param())
			case "printascii":
				msg += fmt.Sprintf("%s contain invalid character(s). ", fieldErr.Field())
			default:
				msg += fmt.Sprintf("%s does not fulfiled criteria of %s with param %s. ", fieldErr.Field(), fieldErr.Tag(), fieldErr.Param())
			}
		}
	} else {
		msg = fmt.Sprintf("%v", err)
	}
	return msg
}

func HttpWriteResponse(writer http.ResponseWriter, req request.IRequest, res response.IResponse) error {
	return HttpWriteResponseWithStatusCode(writer, req, res, http.StatusOK)
}

func HttpWriteResponseWithStatusCode(writer http.ResponseWriter, req request.IRequest, res response.IResponse, statusCode int) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	var resData []byte
	err := enc.Encode(res)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to marshal response model. Nonce: ", req.GetNonce(), " Error: ", err)
		errStr := fmt.Sprintf("%v", err)
		res := response.CreateErrorResponse(req, foundation.InternalServerError, errStr)
		resData, _ = json.Marshal(&res)
	} else {
		resData = buf.Bytes()
	}

	length := len(resData)
	var totalWritten int = 0
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(statusCode)

	for {
		lenWritten, err := writer.Write(resData)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("HTTP write response error: ", err.Error(), " Nonce: ", req.GetNonce())
			return err
		}

		totalWritten += lenWritten
		if lenWritten < length {

			resData = resData[totalWritten:]
			continue
		}
		break
	}

	return nil
}

func HttpWriteRawResponse(writer http.ResponseWriter, contentType string, statusCode int, content []byte) error {

	length := len(content)
	var totalWritten int = 0
	writer.Header().Set("Content-Type", contentType)
	writer.WriteHeader(statusCode)
	var writeContent []byte = content
	for {
		lenWritten, err := writer.Write(content)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("HTTP write response error: ", err.Error())
			return err
		}

		totalWritten += lenWritten
		if lenWritten < length {

			writeContent = writeContent[totalWritten:]
			continue
		}
		break
	}
	return nil
}
