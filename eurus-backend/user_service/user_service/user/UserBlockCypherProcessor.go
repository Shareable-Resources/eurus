package user

import (
	"eurus-backend/foundation"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/log"
)

func processGetBlockCypherAccessToken(server *UserServer, req *BlockCypherAccessTokenRequest) *response.ResponseBase {

	blockCypherToken, err := DbGetAndUpdateBlockCypherToken(req.Coin, server.Config.BlockCypherChain, server.blockCypherDb)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("DbGetAndUpdateBlockCypherToken error: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}
	if blockCypherToken.Token == "" {
		return response.CreateErrorResponse(req, foundation.RecordNotFound, "Coin not found")
	}
	res := new(BlockCypherAccessTokenResponse)
	res.AccessToken = blockCypherToken.Token
	return response.CreateSuccessResponse(req, res)

}
