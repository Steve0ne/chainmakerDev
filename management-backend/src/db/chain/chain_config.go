package chain

import (
	"management_backend/src/db/common"
	"management_backend/src/db/connection"
)

func CreateChainConfigRecord(config *common.ChainConfig) error {
	if err := connection.DB.Create(&config).Error; err != nil {
		log.Error("save chain config record failed: " + err.Error())
		return err
	}
	return nil
}

func GetLastChainConfigRecord(chainId string, beforeTime int64) (*common.ChainConfig, error) {
	var (
		configList []*common.ChainConfig
		err        error
	)

	err = connection.DB.Model(&common.ChainConfig{}).Where("chain_id = ?", chainId).
		Where("block_time < ?", beforeTime).Order("block_time desc").Limit(1).Find(&configList).Error
	if err != nil {
		log.Error("GetLastChainConfigRecord Failed: " + err.Error())
		return nil, err
	}

	return configList[0], nil
}
