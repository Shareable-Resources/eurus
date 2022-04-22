package auth

import (
	"encoding/json"
	"errors"
	"eurus-backend/foundation"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/crypto"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/ws/ws_message"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type AuthClient struct { //implements network.IAuth
	config         auth_base.IAuthBaseConfig
	IsStop         bool
	conn           *websocket.Conn
	sendMutex      *sync.Mutex
	requestMap     map[string]*waitableResponse
	sessionId      int64
	isLoggedIn     bool
	loginToken     string
	additionalInfo string
	loginHandler   func(auth_base.IAuth)
}

type waitableResponse struct {
	Mutex     *sync.Mutex
	Condition *sync.Cond
	Nonce     string
	Error     error
	Response  *ws_message.ResponseMessage
}

func NewAuthClient() *AuthClient {
	auth := new(AuthClient)
	auth.IsStop = false
	auth.requestMap = make(map[string]*waitableResponse)
	auth.sendMutex = new(sync.Mutex)

	return auth
}

func (me *AuthClient) IsLoggedIn() bool {
	return me.isLoggedIn
}

func (me *AuthClient) LoginAuthServer(config auth_base.IAuthBaseConfig) error {
	me.config = config
	go me.receiveProc()

	return nil
}

func (me *AuthClient) GetConfig() auth_base.IAuthBaseConfig {
	return me.config
}

func (me *AuthClient) Close() {
	if me.conn != nil {
		me.IsStop = true
		me.conn.Close()
	}
}

func (me *AuthClient) login() error {

	authReq := new(AuthenticateRequest)
	authReq.ServiceId = me.config.GetServiceId()

	req, err := ws_message.CreateMasterRequestMessage(authReq, "authenticate")
	if err != nil {
		return foundation.NewErrorWithMessage(foundation.InternalServerError, err.Error())
	}

	sign, err := crypto.GenerateRSASignFromBase64(me.config.GetPrivateKey(), req.Message)
	if err != nil {
		return foundation.NewErrorWithMessage(foundation.InvalidSignature, err.Error())
	}
	req.Sign = sign
	me.sendMutex.Lock()
	defer me.sendMutex.Unlock()

	reqRes := &waitableResponse{Mutex: nil, Condition: nil, Nonce: req.Nonce}
	me.requestMap[reqRes.Nonce] = reqRes

	err = me.sendRequest(req)
	if err != nil {
		return foundation.NewErrorWithMessage(foundation.InternalServerError, err.Error())
	}

	return nil
}

func (me *AuthClient) GenerateLoginToken(userId string) (auth_base.ILoginToken, error) {

	innerReq := &RequestLoginTokenRequest{}
	innerReq.UserId = userId

	req, err := ws_message.CreateMasterRequestMessage(innerReq, innerReq.MethodName())
	if err != nil {
		return nil, err
	}
	resp, err := me.sendAndWaitResponse(req)
	if err != nil {
		return nil, err
	}

	loginTokenResp := resp.Data.(auth_base.ILoginToken)
	loginTokenResp.SetUserId(userId)
	return loginTokenResp, nil

}

func (me *AuthClient) VerifyLoginToken(loginToken string, verifyMode auth_base.VerifyMode) (bool, auth_base.ILoginToken, *foundation.ServerError) {

	innerReq := &VerifyLoginTokenRequest{}
	innerReq.Token = loginToken
	innerReq.VerifyMode = int(verifyMode)

	req, err := ws_message.CreateMasterRequestMessage(innerReq, innerReq.MethodName())
	if err != nil {
		return false, nil, foundation.NewErrorWithMessage(foundation.InternalServerError, err.Error())
	}
	resp, err := me.sendAndWaitResponse(req)
	if err != nil {
		return false, nil, foundation.NewErrorWithMessage(foundation.NetworkError, err.Error())
	}
	if resp.ReturnCode != int64(foundation.Success) {
		return false, nil, foundation.NewErrorWithMessage(foundation.ServerReturnCode(resp.ReturnCode), resp.Message)
	}
	loginTokenResp := resp.Data.(auth_base.ILoginToken)
	loginTokenResp.SetToken(loginToken)
	return true, loginTokenResp, nil
}

func (me *AuthClient) RefreshLoginToken(loginToken string) (auth_base.ILoginToken, *foundation.ServerError) {
	innerReq := &RefreshLoginTokenRequest{}
	innerReq.Token = loginToken

	req, err := ws_message.CreateMasterRequestMessage(innerReq, innerReq.MethodName())
	if err != nil {
		return nil, foundation.NewErrorWithMessage(foundation.InternalServerError, err.Error())
	}
	resp, err := me.sendAndWaitResponse(req)
	if err != nil {
		return nil, foundation.NewErrorWithMessage(foundation.NetworkError, err.Error())
	}
	if resp.GetReturnCode() < int64(foundation.Success) {
		return nil, foundation.NewErrorWithMessage(foundation.ServerReturnCode(resp.GetReturnCode()), resp.GetMessage())
	}
	refreshTokenResp := resp.Data.(auth_base.ILoginToken)
	refreshTokenResp.SetToken(loginToken)

	return refreshTokenResp, nil
}

func (me *AuthClient) RevokeLoginToken(loginToken string) (auth_base.ILoginToken, *foundation.ServerError) {
	innerReq := &RevokeLoginTokenRequest{}
	innerReq.Token = loginToken
	req, err := ws_message.CreateMasterRequestMessage(innerReq, innerReq.MethodName())
	if err != nil {
		return nil, foundation.NewErrorWithMessage(foundation.InternalServerError, err.Error())
	}
	resp, err := me.sendAndWaitResponse(req)
	if err != nil {
		return nil, foundation.NewErrorWithMessage(foundation.NetworkError, err.Error())
	}
	if resp.GetReturnCode() < int64(foundation.Success) {
		return nil, foundation.NewErrorWithMessage(foundation.ServerReturnCode(resp.GetReturnCode()), resp.GetMessage())
	}

	revokeLoginTokenResp := resp.Data.(auth_base.ILoginToken)
	return revokeLoginTokenResp, nil
}

//TODO
func (me *AuthClient) RequestNonRefreshableLoginToken(userId string, duration int64, tokenType int16) (auth_base.ILoginToken, *foundation.ServerError) {
	innerReq := &NonRefreshableLoginTokenRequest{}
	innerReq.UserId = userId
	innerReq.Duration = duration
	innerReq.TokenType = TokenType(tokenType)
	req, err := ws_message.CreateMasterRequestMessage(innerReq, innerReq.MethodName())
	if err != nil {
		return nil, foundation.NewErrorWithMessage(foundation.InternalServerError, err.Error())
	}
	resp, err := me.sendAndWaitResponse(req)
	if err != nil {
		return nil, foundation.NewErrorWithMessage(foundation.NetworkError, err.Error())
	}
	if resp.GetReturnCode() < int64(foundation.Success) {
		return nil, foundation.NewErrorWithMessage(foundation.ServerReturnCode(resp.GetReturnCode()), resp.GetMessage())
	}

	//var a map[string]interface{}

	//d, _:=json.Marshal(resp.Data)
	//fmt.Println(d)

	paymentLoginTokenResp := resp.Data.(auth_base.ILoginToken)
	return paymentLoginTokenResp, nil
}

func (me *AuthClient) VerifySignature(data []byte, sign string, serviceId int64) (bool, *foundation.ServerError) {
	innerReq := new(VerifySignRequest)
	innerReq.Data = data
	innerReq.Sign = sign
	innerReq.ServiceId = serviceId

	req, err := ws_message.CreateMasterRequestMessage(innerReq, innerReq.MethodName())
	if err != nil {
		return false, foundation.NewErrorWithMessage(foundation.InternalServerError, err.Error())
	}
	resp, err := me.sendAndWaitResponse(req)
	if err != nil {
		return false, foundation.NewErrorWithMessage(foundation.NetworkError, err.Error())
	}
	if resp.GetReturnCode() < int64(foundation.Success) {
		return false, foundation.NewErrorWithMessage(foundation.ServerReturnCode(resp.GetReturnCode()), resp.GetMessage())
	}

	return resp.Data.(*VerifySignResponse).IsSuccess, nil
}

func (me *AuthClient) sendAndWaitResponse(req *ws_message.MasterRequestMessage) (*ws_message.ResponseMessage, error) {
	if me.conn == nil {
		return nil, errors.New("Connection not established")
	}
	waitMutex := new(sync.Mutex)
	waitSignal := sync.NewCond(waitMutex)

	reqRes := &waitableResponse{Mutex: waitMutex, Condition: waitSignal, Nonce: req.Nonce}

	me.sendMutex.Lock()
	me.requestMap[reqRes.Nonce] = reqRes
	me.sendMutex.Unlock()

	waitMutex.Lock()
	defer waitMutex.Unlock()

	err := me.sendRequest(req)
	if err != nil {
		return nil, err
	}
	waitSignal.Wait()

	if reqRes.Error != nil {
		return nil, reqRes.Error
	}

	return reqRes.Response, reqRes.Error
}

func (me *AuthClient) receiveProc() {

	var buff []byte = make([]byte, 4096)
	for !me.IsStop {

		for !me.IsStop && me.conn != nil {

			_, reader, err := me.conn.NextReader()
			if err != nil {
				me.clearLoginData()
				log.GetLogger(log.Name.Root).Errorln("Read Auth socket error: ", err.Error())
				break
			}
			//reader.Read will assign buff
			length, err := reader.Read(buff)
			if length == 0 && err != nil {
				me.clearLoginData()
				log.GetLogger(log.Name.Root).Errorln("Read socket error: ", err.Error())
				break
			}
			me.parseMessage(buff[:length])
		}

		authUrl := url.URL{Scheme: "ws",
			Host: me.config.GetAuthIp() + ":" + strconv.FormatUint(uint64(me.config.GetAuthPort()), 10),
			Path: me.config.GetAuthPath(),
		}
		log.GetLogger(log.Name.Root).Infoln("Going to connecton to ", authUrl.String())
		var res *http.Response
		var err error
		me.conn, res, err = websocket.DefaultDialer.Dial(authUrl.String(), nil)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Connect auth error: ", res, err)
			time.Sleep(5 * time.Second)
		} else {
			serverError := me.login()
			if serverError != nil {
				me.conn.Close()
				me.clearLoginData()
				time.Sleep(5 * time.Second)
			}
		}
	}
	log.GetLogger(log.Name.Root).Infoln("Receive proc exit")
}

func (me *AuthClient) parseMessage(data []byte) {

	res := new(ws_message.ResponseMessage)
	err := json.Unmarshal(data, res)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to parse message: ", err, string(data))
		return
	}

	me.sendMutex.Lock()

	waitRes, ok := me.requestMap[res.Nonce]

	if !ok {
		me.sendMutex.Unlock()
		return
	}

	waitRes.Response = res
	delete(me.requestMap, res.Nonce)
	me.sendMutex.Unlock()
	if res.MethodName != ApiAuth {
		waitRes.Condition.Signal()
	} else {
		me.handleLoginResponse(res)
	}
}

func (me *AuthClient) handleLoginResponse(res *ws_message.ResponseMessage) {
	if res.ReturnCode == int64(foundation.Success) {
		me.isLoggedIn = true
		actualRes := res.Data.(*AuthenticateResponse)
		me.sessionId = actualRes.SessionId
		me.loginToken = actualRes.Token
		me.additionalInfo = actualRes.AdditionalInfo
		log.GetLogger(log.Name.Root).Infoln("Login auth server successful")
		if me.loginHandler != nil {
			go me.loginHandler(me)
		}
	} else {
		log.GetLogger(log.Name.Root).Errorln("Login auth server failed. Code: ", res.ReturnCode, " message: ", res.Message)
		me.conn.Close()
		me.clearLoginData()
		time.Sleep(5 * time.Second)
	}
}

func (me *AuthClient) clearLoginData() {
	me.isLoggedIn = false
	me.conn = nil
	me.loginToken = ""
	me.sessionId = 0
}

func (me *AuthClient) sendRequest(req *ws_message.MasterRequestMessage) error {
	writer, err := me.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return nil
	}
	defer writer.Close()
	reqData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	for {
		totalLen := len(reqData)
		length, err := writer.Write(reqData)
		if err != nil {
			return err
		}
		if length >= totalLen {
			break
		}
		reqData = reqData[length:]
	}
	return nil

}

func (me *AuthClient) GetLoginToken() string {
	return me.loginToken
}

func (me *AuthClient) GetSessionId() int64 {
	return me.sessionId
}

func (me *AuthClient) GetAdditionalInfo() string {
	return me.additionalInfo
}

func (me *AuthClient) SetLoginHandler(loginHandler func(auth_base.IAuth)) {
	me.loginHandler = loginHandler
}

func (me *AuthClient) SetLoginToken(loginToken string) {
	me.loginToken = loginToken
	me.isLoggedIn = true
}

///Domain specific function, does not exists in network.IAuth interface
func (me *AuthClient) GetServiceIdFromServerLoginToken(loginToken auth_base.ILoginToken) (int64, error) {
	userIdStr := loginToken.GetUserId()
	var tokenMap map[string]interface{} = make(map[string]interface{})
	err := json.Unmarshal([]byte(userIdStr), &tokenMap)
	if err != nil {
		return 0, err
	}

	serviceId, ok := tokenMap["serviceId"].(float64)
	if !ok {
		return 0, errors.New("Invalid login token")
	}
	return int64(serviceId), nil

}
