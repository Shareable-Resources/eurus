package user

import (
	"eurus-backend/foundation/database"
	"eurus-backend/marketing/banner"
	"eurus-backend/marketing/reward"
	"time"

	"gorm.io/gorm"
)

func DbQueryRewardList(db *database.ReadOnlyDatabase, userId uint64) ([]*reward.DistributedToken, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	var distributeTokenList []*reward.DistributedToken = make([]*reward.DistributedToken, 0)
	tx := dbConn.Model(reward.DistributedToken{}).Where("user_id = ? AND status = ?", userId, reward.DistributedStatusSuccess).Find(&distributeTokenList)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return distributeTokenList, nil

}

func DbQueryMarketingBanner(db *database.ReadOnlyDatabase, position uint64, langCode string) ([]*QueryMarketingBannerList, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()

	var marketingBannerList []*QueryMarketingBannerList = make([]*QueryMarketingBannerList, 0)
	sess := dbConn.Session(&gorm.Session{}).Model(banner.MarketingBanner{}).Select("*")
	sess = sess.Where(" position = ?", position)
	sess = sess.Where(" start_date <= ?", currentTime)
	sess = sess.Where(" end_date >= ?", currentTime)
	// sess = sess.Where(" is_default = ?", isDefault)

	if langCode != "" {
		sess = sess.Where(" lang_code = ?", langCode)
	} else {
		sess = sess.Where(" lang_code = ?", banner.English)
	}

	sess = sess.Joins("Inner JOIN marketing_banner_details as mbd ON marketing_banners.id = mbd.banner_id")
	sess = sess.Order("marketing_banners.seq")

	tx := sess.Find(&marketingBannerList)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return marketingBannerList, nil
}
