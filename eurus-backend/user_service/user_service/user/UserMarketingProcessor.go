package user

import (
	"eurus-backend/foundation"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/log"
)

func ProcessQueryRewardList(server *UserServer, req *request.RequestBase) *response.ResponseBase {

	userID, err := UnmarshalUserIdFromLoginToken(req)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.RequestParamsValidationError, "Incorrect wallet address")
	}

	list, err := DbQueryRewardList(server.SlaveDatabase, userID)
	if err != nil {

		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	var rewardList []*UserReward = make([]*UserReward, 0)
	for _, item := range list {
		userReward := new(UserReward)
		userReward.RewardType = item.DistributedType
		userReward.Amount = item.Amount.BigInt()
		userReward.AssetName = item.AssetName
		userReward.CreateDate = item.CreatedDate
		userReward.TxHash = item.TxHash
		rewardList = append(rewardList, userReward)
	}

	return response.CreateSuccessResponse(req, rewardList)
}

func ProcessQueryRewardScheme(server *UserServer, req *request.RequestBase) *response.ResponseBase {
	res := new(RewardSchemeResponse)
	res.Data = server.Config.MarketingRewardSchemeJson

	return response.CreateSuccessResponse(req, res)
}

func ProcessQueryMarketingBanner(server *UserServer, req *MarketingBannerRequest) *response.ResponseBase {
	marketingBannerList, err := DbQueryMarketingBanner(server.SlaveDatabase, req.Position, req.LangCode)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("QueryMarketingBanner error: ", err.Error(), " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.DatabaseError, "QueryMarketingBanner error: "+err.Error())
	}
	if marketingBannerList == nil || len(marketingBannerList) <= 0 {
		log.GetLogger(log.Name.Root).Errorln("Cannot query default MarketingBanner ", " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.DatabaseError, "Cannot query default MarketingBanner ")
	}

	var modifiedMarketingBannerList []*QueryMarketingBannerList = make([]*QueryMarketingBannerList, 0)
	for i := 0; i < len(marketingBannerList); i++ {
		nextCount := i + 1
		if nextCount < len(marketingBannerList) {
			marketingBanner := marketingBannerList[i]
			nextMarketingBanner := marketingBannerList[nextCount]

			if marketingBanner.Seq == nextMarketingBanner.Seq {
				if marketingBanner.IsDefault {
					modifiedMarketingBannerList = append(modifiedMarketingBannerList, nextMarketingBanner)
				} else {
					modifiedMarketingBannerList = append(modifiedMarketingBannerList, marketingBanner)
				}

				i = nextCount
			} else {
				modifiedMarketingBannerList = append(modifiedMarketingBannerList, marketingBannerList[i])
			}
		} else {
			marketingBanner := marketingBannerList[i]
			lastModifiedMarketingBannerList := modifiedMarketingBannerList[len(modifiedMarketingBannerList)-1]

			if marketingBanner.Seq != lastModifiedMarketingBannerList.Seq {
				modifiedMarketingBannerList = append(modifiedMarketingBannerList, marketingBanner)
			}
		}
	}

	var queryMarketingBannerList []*QueryMarketingBanner = make([]*QueryMarketingBanner, 0)
	for _, s := range modifiedMarketingBannerList {
		queryMarketingBanner := &QueryMarketingBanner{
			Seq:        s.Seq,
			BannerType: s.BannerType,
			Icon: Icon{
				Mobile: s.IconUrlMobile,
			},
			Content: s.Content,
			BannerImage: BannerImage{
				Mobile: s.BannerUrlMobile,
			},
			Link: Link{
				Mobile: s.LinkMobile,
			},
		}

		queryMarketingBannerList = append(queryMarketingBannerList, queryMarketingBanner)
	}

	return response.CreateSuccessResponse(req, queryMarketingBannerList)
}
