/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package multi_sign

import (
	pbcommon "chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"

	"management_backend/src/ctrl/contract_management"
	"management_backend/src/db/chain_participant"
	loggers "management_backend/src/logger"
)

const (
	POLICY_ADMIN  = 0
	POLICY_CLIENT = 1
	POLICY_ALL    = 2
)

//const (
//	CONTRACT_PAYLOAD = 0
//	CHAIN_PAYLOAD    = 1
//)

var log = loggers.GetLogger(loggers.ModuleWeb)

func MultiSignInvoke(parameters string, multiSignType int, orgList []string, roleType, configStatus int) error {
	if multiSignType == contract_management.BLOCK_UPDATE {
		return ChainConfigModify(parameters, orgList, roleType)
	}

	if multiSignType == contract_management.PERMISSION_UPDATE {
		return ChainAuthModify(parameters, orgList, roleType)
	}

	if multiSignType == contract_management.INIT_CONTRACT || multiSignType == contract_management.UPGRADE_CONTRACT {
		return ContractInstallModify(parameters, orgList, roleType, multiSignType)
	}
	if multiSignType == contract_management.FREEZE_CONTRACT {
		return ContractFreezeModify(parameters, orgList, roleType)
	}
	if multiSignType == contract_management.UNFREEZE_CONTRACT {
		return ContractUnfreezeModify(parameters, orgList, roleType)
	}
	if multiSignType == contract_management.REVOKE_CONTRACT {
		return ContractRevokeModify(parameters, orgList, roleType)
	}
	return nil
}

func GetEndorsements(payload *pbcommon.Payload, orgList []string, roleType int) ([]*pbcommon.EndorsementEntry, error) {
	var endorsement *pbcommon.EndorsementEntry
	var endorsements []*pbcommon.EndorsementEntry

	for _, orgId := range orgList {
		if roleType == POLICY_CLIENT {
			roleType = chain_participant.CLIENT
		} else {
			roleType = chain_participant.ADMIN
		}
		cert, err := chain_participant.GetUserCertByOrgId(orgId, roleType)
		if err != nil {
			return nil, err
		}
		privateKeyBytes := []byte(cert.PrivateKey)
		crtBytes := []byte(cert.Cert)

		endorsement, err = sdk.SignPayload(privateKeyBytes, crtBytes, payload)
		if err != nil {
			return nil, err
		}
		endorsements = append(endorsements, endorsement)
	}

	return endorsements, nil
}
