/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package chain_management

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"chainmaker.org/chainmaker/pb-go/v2/consensus"

	"management_backend/src/config"
	"management_backend/src/ctrl/common"
	"management_backend/src/db/chain"
	"management_backend/src/db/chain_participant"
	dbcommon "management_backend/src/db/common"
	"management_backend/src/db/relation"
	"management_backend/src/entity"
	loggers "management_backend/src/logger"
	"management_backend/src/utils"

	"gopkg.in/yaml.v2"
)

const (
	SIGN = 0
	TLS  = 1
)

const (
	SIGN_USE = "sign"
	TLS_USE  = "tls"
)

const (
	CONF_LOCAL_PATH  = "configs"
	CONF_SERVER_PATH = "../configs"
)

const (
	DEPENDENCE_LOCAL_PATH  = "dependence"
	DEPENDENCE_SERVER_PATH = "../dependence"
)

const MONITOR_START = 1
const NO_TLS = 1
const NO_DOCKER_VM = 0

var log = loggers.GetLogger(loggers.ModuleWeb)

type DownloadChainConfigHandler struct{}

func (downloadChainConfigHandler *DownloadChainConfigHandler) LoginVerify() bool {
	return false
}

func (downloadChainConfigHandler *DownloadChainConfigHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindDownloadChainConfigHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	confYml := config.ConfEnvPath
	if confYml == "" {
		confYml = CONF_SERVER_PATH
	}
	if confYml == CONF_SERVER_PATH {
		confYml = DEPENDENCE_SERVER_PATH
	}
	if confYml == CONF_LOCAL_PATH {
		confYml = DEPENDENCE_LOCAL_PATH
	}
	chainId := params.ChainId

	//创建bc
	nodeIdMap, err := createBc(chainId, confYml)
	if err != nil {
		log.Error("CreateBc err : " + err.Error())
	}
	chainOrgNodes, err := relation.GetChainOrgByChainIdList(chainId)
	if err != nil {
		log.Error("GetChainOrgNode err : " + err.Error())
	}

	var chainName string
	var tls int
	var dockerVm int
	monitorStart := false
	chainInfo, err := chain.GetChainByChainId(chainId)
	if err != nil {
		log.Error("Get chainInfo by chainId err : " + err.Error())
		chainName = chainId
		tls = 0
		dockerVm = 0
	} else {
		chainName = chainInfo.ChainName
		tls = chainInfo.TLS
		dockerVm = chainInfo.DockerVm
		if chainInfo.Monitor == MONITOR_START {
			err = createLogAgent(confYml)
			if err != nil {
				log.Error("createLogAgent err : " + err.Error())
			}
		}
		monitorStart = true
	}

	var nodePaths string

	for _, chainOrgNode := range chainOrgNodes {
		//创建chainmaker
		err = createChainmaker(chainId, chainOrgNode.OrgId, chainOrgNode.NodeName, confYml, nodeIdMap, tls, dockerVm)
		if err != nil {
			log.Error("createChainmaker err : " + err.Error())
		}
		//创建bin lib
		err = createBinAndLib(chainOrgNode.OrgId, chainOrgNode.NodeName, confYml)
		if err != nil {
			log.Error("createBinAndLib err : " + err.Error())
		}
		//创建cert
		err = createCert(chainId, chainOrgNode.OrgId, chainOrgNode.NodeName)
		if err != nil {
			log.Error("createCert err : " + err.Error())
		}

		nodePaths = nodePaths + "./" + chainOrgNode.OrgId + "-" + chainOrgNode.NodeName + ","
	}

	nodePaths = strings.TrimRight(nodePaths, ",")

	err = createChainmakerAndScript(chainId, confYml, nodePaths, monitorStart)
	if err != nil {
		log.Error("createChainmakerAndScript err : " + err.Error())
	}

	defer func() {
		err = os.RemoveAll(chainId + ".zip")
		if err != nil {
			log.Error("remove zip err :", err.Error())
		}
	}()

	content, err := ioutil.ReadFile(chainId + ".zip")
	if err != nil {
		log.Error("ReadFile err : " + err.Error())
	}
	fileName := chainName + ".zip"

	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Header("Content-Disposition", "attachment; filename="+utils.Base64Encode([]byte(fileName)))
	ctx.Header("Content-Type", "application/zip")
	ctx.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
	ctx.Header("Access-Control-Expose-Headers", "Content-Disposition")
	_, err = ctx.Writer.Write(content)
	if err != nil {
		log.Error("ctx Write content err :", err.Error())
	}
}

func createLogAgent(confYml string) error {
	logAgentFile, err := os.Create("release/cmlogagentd")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer func() {
		err = logAgentFile.Close()
		if err != nil {
			log.Error("close file err :", err.Error())
		}
	}()
	err = logAgentFile.Chmod(0777)
	if err != nil {
		log.Error(err.Error())
	}
	_, err = utils.CopyFile("release/cmlogagentd", confYml+"/bin/cmlogagentd")
	if err != nil {
		log.Error("CopyFile bin/cmlogagentd err : " + err.Error())
	}

	return nil
}

func createChainmakerAndScript(chainId, confYml, nodePaths string, monitorStart bool) error {
	chainmakerFile, err := os.Create("release/chainmaker")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer func() {
		err = chainmakerFile.Close()
		if err != nil {
			log.Error("close file err :", err.Error())
		}
	}()
	err = chainmakerFile.Chmod(0777)
	if err != nil {
		log.Error(err.Error())
	}
	_, err = utils.CopyFile("release/chainmaker", confYml+"/bin/chainmaker")
	if err != nil {
		log.Error("CopyFile bin/chainmaker err : " + err.Error())
	}

	//创建启动脚本
	startFile, err := os.Create("release/start.sh")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer func() {
		err = startFile.Close()
		if err != nil {
			log.Error("close file err :", err.Error())
		}
	}()
	err = startFile.Chmod(0777)
	if err != nil {
		log.Error(err.Error())
	}
	if monitorStart {
		_, err = utils.CopyFile("release/start", confYml+"/bin/logagentd_start.sh")
		if err != nil {
			log.Error("CopyFile bin/start.sh err : " + err.Error())
		}

		err = utils.RePlace("release/start", "{node_paths}", nodePaths)
		if err != nil {
			log.Error("rePlace release/start.sh err : " + err.Error())
		}
	} else {
		_, err = utils.CopyFile("release/start", confYml+"/bin/start.sh")
		if err != nil {
			log.Error("CopyFile bin/start.sh err : " + err.Error())
		}
	}

	quickStopFile, err := os.Create("release/quick_stop.sh")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer func() {
		err = quickStopFile.Close()
		if err != nil {
			log.Error("close file err :", err.Error())
		}
	}()
	err = quickStopFile.Chmod(0777)
	if err != nil {
		log.Error(err.Error())
	}
	_, err = utils.CopyFile("release/quick_stop.sh", confYml+"/bin/quick_stop.sh")
	if err != nil {
		log.Error("CopyFile bin/quick_stop.sh err : " + err.Error())
	}

	err = utils.Zip("release", chainId+".zip")
	if err != nil {
		log.Error("zip file err :", err.Error())
	}

	return nil
}

func createCert(chainId, orgId, nodeName string) error {

	err := createNodeCert(orgId, nodeName)
	if err != nil {
		log.Error("createNodeCert err : " + err.Error())
		return err
	}

	err = createUserCert(orgId, nodeName)
	if err != nil {
		log.Error("createUserCert err : " + err.Error())
		return err
	}

	err = createOrgCert(chainId, orgId, nodeName)
	if err != nil {
		log.Error("createOrgCert err : " + err.Error())
		return err
	}

	return nil
}

func createNodeCert(orgId, nodeName string) error {
	nodeCertList, err := chain_participant.GetNodeCert(nodeName)
	if err != nil {
		log.Error("GetNodeCert erBlockr : " + err.Error())
		return err
	}

	nodeId, err := os.Create(nodeName + ".nodeid")
	defer func() {
		err = nodeId.Close()
		if err != nil {
			log.Error("close file err :", err.Error())
		}
	}()
	if err != nil {
		log.Error(err.Error())
	}
	nodeInfo, err := chain_participant.GetNodeByNodeName(nodeName)
	if err != nil {
		log.Error("GetNodeByNodeName err : " + err.Error())
		return err
	}
	_, err = nodeId.Write([]byte(nodeInfo.NodeId))
	if err != nil {
		log.Error("nodeId Write err : " + err.Error())
	}

	err = os.MkdirAll("release/"+orgId+"-"+nodeName+"/config/"+orgId+"/certs/node/"+nodeName, os.ModePerm)
	if err != nil {
		log.Error("Mkdir org certs/node path err : " + err.Error())
	}

	err = os.Rename(nodeName+".nodeid", "release/"+
		orgId+"-"+nodeName+"/config/"+orgId+"/certs/node/"+nodeName+"/"+nodeName+".nodeid")
	if err != nil {
		log.Error("Rename nodeid err : " + err.Error())
	}

	for _, nodeCert := range nodeCertList {
		if nodeCert.CertUse == SIGN {
			err := mkdirNodeCert(orgId, nodeName, SIGN_USE, nodeCert)
			if err != nil {
				log.Error("Mkdir certs/node path err : " + err.Error())
			}

		} else if nodeCert.CertUse == TLS {
			err := mkdirNodeCert(orgId, nodeName, TLS_USE, nodeCert)
			if err != nil {
				log.Error("Mkdir certs/node path err : " + err.Error())
			}
		}
	}
	return nil
}

func createUserCert(orgId, nodeName string) error {
	userCertList, _, err := chain_participant.GetUserCertList(orgId)
	if err != nil {
		log.Error("GetUserCertList err : " + err.Error())
		return err
	}

	for _, userCert := range userCertList {
		userName := userCert.CertUserName
		err = os.MkdirAll("release/"+orgId+"-"+nodeName+"/config/"+orgId+"/certs/user/"+userName, os.ModePerm)
		if err != nil {
			log.Error("Mkdir org certs/user path err : " + err.Error())
		}
		if userCert.CertUse == SIGN {
			err = mkdirUserCert(userName, orgId, nodeName, SIGN_USE, userCert)
			if err != nil {
				log.Error("mkdirUserCerterr : " + err.Error())
			}
		} else {
			err = mkdirUserCert(userName, orgId, nodeName, TLS_USE, userCert)
			if err != nil {
				log.Error("mkdirUserCerterr : " + err.Error())
			}
		}
	}

	return nil
}

func mkdirNodeCert(orgId, nodeName, certUse string, nodeCert *dbcommon.Cert) error {
	nodeCertName := nodeName + "." + certUse + ".crt"
	nodeKeyName := nodeName + "." + certUse + ".key"

	nodeTlsCrt, err := os.Create(nodeCertName)
	if err != nil {
		log.Error(err.Error())
	}
	defer func() {
		err = nodeTlsCrt.Close()
		if err != nil {
			log.Error("close file err :", err.Error())
		}
	}()

	nodeTlsKey, err := os.Create(nodeKeyName)
	if err != nil {
		log.Error(err.Error())
	}
	defer func() {
		err = nodeTlsKey.Close()
		if err != nil {
			log.Error("close file err :", err.Error())
		}
	}()

	_, err = nodeTlsCrt.Write([]byte(nodeCert.Cert))
	if err != nil {
		log.Error("file Write err :", err.Error())
	}
	_, err = nodeTlsKey.Write([]byte(nodeCert.PrivateKey))
	if err != nil {
		log.Error("file Write err :", err.Error())
	}

	err = os.Rename(nodeCertName, "release/"+
		orgId+"-"+nodeName+"/config/"+orgId+"/certs/node/"+nodeName+"/"+nodeCertName)
	if err != nil {
		log.Error("Rename node sign.crt err : " + err.Error())
	}

	err = os.Rename(nodeKeyName, "release/"+
		orgId+"-"+nodeName+"/config/"+orgId+"/certs/node/"+nodeName+"/"+nodeKeyName)
	if err != nil {
		log.Error("Rename node sign.key err : " + err.Error())
	}
	return nil
}

func mkdirUserCert(userName, orgId, nodeName, certUse string, userCert *dbcommon.Cert) error {
	userCertName := userName + "." + certUse + ".crt"
	userKeyName := userName + "." + certUse + ".key"

	userSignCrt, err := os.Create(userCertName)
	if err != nil {
		log.Error(err.Error())
	}
	defer func() {
		err = userSignCrt.Close()
		if err != nil {
			log.Error("close file err :", err.Error())
		}
	}()
	userSignKey, err := os.Create(userKeyName)
	if err != nil {
		log.Error(err.Error())
	}
	defer func() {
		err = userSignKey.Close()
		if err != nil {
			log.Error("close file err :", err.Error())
		}
	}()
	_, err = userSignCrt.Write([]byte(userCert.Cert))
	if err != nil {
		log.Error("write file err :", err.Error())
	}
	_, err = userSignKey.Write([]byte(userCert.PrivateKey))
	if err != nil {
		log.Error("write file err :", err.Error())
	}

	err = os.Rename(userCertName, "release/"+
		orgId+"-"+nodeName+"/config/"+orgId+"/certs/user/"+userName+"/"+userCertName)
	if err != nil {
		log.Error("Rename user tls.crt err : " + err.Error())
	}

	err = os.Rename(userKeyName, "release/"+
		orgId+"-"+nodeName+"/config/"+orgId+"/certs/user/"+userName+"/"+userKeyName)
	if err != nil {
		log.Error("Rename user tls.key err : " + err.Error())
	}
	return nil
}

func createOrgCert(chainId, orgId, nodeName string) error {
	chainOrgs, err := relation.GetChainOrgList(chainId)
	if err != nil {
		log.Error("GetChainOrgByChainIdList err : " + err.Error())
		return err
	}

	for _, orgInfo := range chainOrgs {
		err = os.MkdirAll("release/"+orgId+"-"+nodeName+"/config/"+orgId+"/certs/ca/"+orgInfo.OrgId, os.ModePerm)
		if err != nil {
			log.Error("Mkdir org certs/ca path err : " + err.Error())
		}
		orgCert, err := chain_participant.GetOrgCaCert(orgInfo.OrgId)
		if err != nil {
			log.Error("GetOrgCaCert err : " + err.Error())
			return err
		}
		f, err := os.Create("ca.crt")
		if err != nil {
			log.Error(err.Error())
			return err
		}
		_, err = f.Write([]byte(orgCert.Cert))
		if err != nil {
			log.Error("write file err :", err.Error())
		}

		err = os.Rename("ca.crt", "release/"+
			orgId+"-"+nodeName+"/config/"+orgId+"/certs/ca/"+orgInfo.OrgId+"/ca.crt")
		if err != nil {
			log.Error("Rename chainmaker.yml err : " + err.Error())
		}
		defer func() {
			err = f.Close()
			if err != nil {
				log.Error("close file err :", err.Error())
			}
		}()
	}

	return nil
}

func createBinAndLib(orgId, nodeName, confYml string) error {
	err := os.MkdirAll("release/"+orgId+"-"+nodeName+"/bin", os.ModePerm)
	if err != nil {
		log.Error("Mkdir bin path err : " + err.Error())
	}
	restartFile, err := os.Create("release/" + orgId + "-" + nodeName + "/bin/restart.sh")
	if err != nil {
		log.Error(err.Error())
	}
	defer func() {
		err = restartFile.Close()
		if err != nil {
			log.Error("close file err :", err.Error())
		}
	}()
	err = restartFile.Chmod(0777)
	if err != nil {
		log.Error(err.Error())
	}
	stopFile, err := os.Create("release/" + orgId + "-" + nodeName + "/bin/stop.sh")
	if err != nil {
		log.Error(err.Error())
	}
	defer func() {
		err = stopFile.Close()
		if err != nil {
			log.Error("close file err :", err.Error())
		}
	}()
	err = stopFile.Chmod(0777)
	if err != nil {
		log.Error(err.Error())
	}

	_, err = utils.CopyFile("release/"+
		orgId+"-"+nodeName+"/bin/restart", confYml+"/bin/restart.sh")
	if err != nil {
		log.Error("CopyFile bin/restart.sh err : " + err.Error())
	}

	err = utils.RePlace("release/"+orgId+"-"+nodeName+"/bin/restart", "{org_id}", orgId)
	if err != nil {
		log.Error("rePlace bin/restart.sh err : " + err.Error())
	}

	_, err = utils.CopyFile("release/"+orgId+"-"+nodeName+"/bin/stop", confYml+"/bin/stop.sh")
	if err != nil {
		log.Error("CopyFile bin/stop.sh err : " + err.Error())
	}
	err = utils.RePlace("release/"+orgId+"-"+nodeName+"/bin/stop", "{org_id}", orgId)
	if err != nil {
		log.Error("rePlace bin/stop.sh err : " + err.Error())
	}

	err = createLib(orgId, nodeName, confYml)
	if err != nil {
		log.Error("createLib err : " + err.Error())
	}

	return nil
}

func createLib(orgId, nodeName, confYml string) error {
	err := os.MkdirAll("release/"+orgId+"-"+nodeName+"/lib", os.ModePerm)
	if err != nil {
		log.Error("Mkdir lib/chainmaker err : " + err.Error())
	}

	dylibFile, err := os.Create("release/" + orgId + "-" + nodeName + "/lib/libwasmer.so")
	if err != nil {
		log.Error(err.Error())
	}
	defer func() {
		err = dylibFile.Close()
		if err != nil {
			log.Error("close file err :", err.Error())
		}
	}()
	err = dylibFile.Chmod(0777)
	if err != nil {
		log.Error(err.Error())
	}

	wxdecFile, err := os.Create("release/" + orgId + "-" + nodeName + "/lib/wxdec")
	if err != nil {
		log.Error(err.Error())
	}
	defer func() {
		err = wxdecFile.Close()
		if err != nil {
			log.Error("close file err :", err.Error())
		}
	}()
	err = wxdecFile.Chmod(0777)
	if err != nil {
		log.Error(err.Error())
	}

	_, err = utils.CopyFile("release/"+
		orgId+"-"+nodeName+"/lib/libwasmer.so", confYml+"/lib/libwasmer.so")
	if err != nil {
		log.Error("CopyFile lib/chainmaker err : " + err.Error())
	}

	_, err = utils.CopyFile("release/"+orgId+"-"+nodeName+"/lib/wxdec", confYml+"/lib/wxdec")
	if err != nil {
		log.Error("CopyFile lib/chainmaker.service err : " + err.Error())
	}
	return nil
}

func createChainmaker(chainId, orgId, nodeName, confYml string, nodeIdMap map[string]int, tls int, dockerVm int) error {
	nodeInfo, err := chain_participant.GetNodeByNodeName(nodeName)
	if err != nil {
		log.Error("GetNodeByNodeName err : " + err.Error())
		return err
	}
	conf := new(config.Chainmaker)
	yamlFile, _ := ioutil.ReadFile(confYml + "/config_tpl/chainmaker.yml")
	_ = yaml.Unmarshal(yamlFile, conf)

	err = setChainmaker(chainId, orgId, nodeIdMap, conf, nodeInfo, tls, dockerVm)
	if err != nil {
		log.Error("setChainmaker err : " + err.Error())
	}

	chainmakerBytes, _ := yaml.Marshal(conf)
	chainmaker, err := os.Create("chainmaker.yml")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer func() {
		err = chainmaker.Close()
		if err != nil {
			log.Error("close file err :", err.Error())
		}
	}()

	_, err = chainmaker.Write(chainmakerBytes)
	if err != nil {
		log.Error(err.Error())
	}
	err = os.MkdirAll("release/"+orgId+"-"+nodeName+"/config/"+orgId, os.ModePerm)
	if err != nil {
		log.Error("Mkdir org chainmaker path err : " + err.Error())
	}
	err = os.Rename("chainmaker.yml", "release/"+
		orgId+"-"+nodeName+"/config/"+orgId+"/chainmaker.yml")
	if err != nil {
		log.Error("Rename chainmaker.yml err : " + err.Error())
		return err
	}

	logFile, err := os.Create("release/" + orgId + "-" + nodeName + "/config/" + orgId + "/log.yml")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer func() {
		err = logFile.Close()
		if err != nil {
			log.Error("close file err :", err.Error())
		}
	}()
	_, err = utils.CopyFile("release/"+
		orgId+"-"+nodeName+"/config/"+orgId+"/log.yml", confYml+"/config_tpl/log.yml")
	if err != nil {
		log.Error("copy log.yml err : " + err.Error())
		return err
	}

	return nil
}

func setChainmaker(chainId, orgId string, nodeIdMap map[string]int, conf *config.Chainmaker,
	nodeInfo *dbcommon.Node, tls int, dockerVm int) error {
	confingFile := conf.ChainLogConf.ConfigFile
	conf.ChainLogConf.ConfigFile = strings.Replace(confingFile, "{org_path}", orgId, -1)

	blockchainConf := config.BlockchainConf{}
	blockchainConf.ChainId = chainId
	blockchainConf.Genesis = "../config/" + orgId + "/chainconfig/bc1.yml"
	chainList := []*config.BlockchainConf{}
	chainList = append(chainList, &blockchainConf)
	conf.BlockchainConf = chainList

	conf.NodeConf.OrgId = orgId
	certFile := conf.NodeConf.CertFile
	certFile = strings.Replace(certFile, "{org_path}", orgId, -1)
	certFile = strings.Replace(certFile, "{node_cert_path}",
		"node/"+nodeInfo.NodeName+"/"+nodeInfo.NodeName+".sign", -1)
	conf.NodeConf.CertFile = certFile

	privKeyFile := conf.NodeConf.PrivKeyFile
	privKeyFile = strings.Replace(privKeyFile, "{org_path}", orgId, -1)
	privKeyFile = strings.Replace(privKeyFile, "{node_cert_path}",
		"node/"+nodeInfo.NodeName+"/"+nodeInfo.NodeName+".sign", -1)
	conf.NodeConf.PrivKeyFile = privKeyFile

	seedList := []string{}
	for nodeId, port := range nodeIdMap {
		seedList = append(seedList, "/ip4/0.0.0.0/tcp/"+strconv.Itoa(port)+"/p2p/"+nodeId)
	}
	conf.NetConf.Seeds = seedList

	tlsCertFile := conf.NetConf.Tls.CertFile
	tlsCertFile = strings.Replace(tlsCertFile, "{org_path}", orgId, -1)
	tlsCertFile = strings.Replace(tlsCertFile, "{net_cert_path}",
		"node/"+nodeInfo.NodeName+"/"+nodeInfo.NodeName+".tls", -1)
	conf.NetConf.Tls.CertFile = tlsCertFile

	tlsKeyFile := conf.NetConf.Tls.PrivKeyFile
	tlsKeyFile = strings.Replace(tlsKeyFile, "{org_path}", orgId, -1)
	tlsKeyFile = strings.Replace(tlsKeyFile, "{net_cert_path}",
		"node/"+nodeInfo.NodeName+"/"+nodeInfo.NodeName+".tls", -1)
	conf.NetConf.Tls.PrivKeyFile = tlsKeyFile

	port, _ := strconv.Atoi(nodeInfo.NodePort)
	conf.RpcConf.Port = port
	conf.RpcConf.Tls.CertFile = tlsCertFile
	conf.RpcConf.Tls.PrivKeyFile = tlsKeyFile
	if tls == NO_TLS {
		conf.RpcConf.Tls.Mode = "disable"
	}

	conf.NetConf.ListenAddr = "/ip4/0.0.0.0/tcp/" + strconv.Itoa(nodeIdMap[nodeInfo.NodeId])

	storePath := conf.StorageConf.StorePath
	conf.StorageConf.StorePath = strings.Replace(storePath, "{org_id}", orgId+"-"+nodeInfo.NodeName, -1)

	blockPath := conf.StorageConf.BlockdbConfig.LeveldbConfig.StorePath
	conf.StorageConf.BlockdbConfig.LeveldbConfig.StorePath = strings.Replace(blockPath,
		"{org_id}", orgId+"-"+nodeInfo.NodeName, -1)

	statePath := conf.StorageConf.StatedbConfig.LeveldbConfig.StorePath
	conf.StorageConf.StatedbConfig.LeveldbConfig.StorePath = strings.Replace(statePath,
		"{org_id}", orgId+"-"+nodeInfo.NodeName, -1)

	historyPath := conf.StorageConf.HistorydbConfig.LeveldbConfig.StorePath
	conf.StorageConf.HistorydbConfig.LeveldbConfig.StorePath = strings.Replace(historyPath,
		"{org_id}", orgId+"-"+nodeInfo.NodeName, -1)

	resultPath := conf.StorageConf.ResultdbConfig.LeveldbConfig.StorePath
	conf.StorageConf.ResultdbConfig.LeveldbConfig.StorePath = strings.Replace(resultPath,
		"{org_id}", orgId+"-"+nodeInfo.NodeName, -1)

	if dockerVm == NO_DOCKER_VM {
		conf.VmConf.EnableDockervm = false
	}

	logPath := conf.VmConf.DockervmLogPath
	conf.VmConf.DockervmLogPath = strings.Replace(logPath, "{org_id}", orgId+"-"+nodeInfo.NodeName, -1)

	mountPath := conf.VmConf.DockervmMountPath
	conf.VmConf.DockervmMountPath = strings.Replace(mountPath, "{org_id}", orgId+"-"+nodeInfo.NodeName, -1)

	containerName := conf.VmConf.DockervmContainerName
	conf.VmConf.DockervmContainerName = strings.Replace(containerName, "dockervm_container_name", "chainmaker-vm-docker-go-container"+orgId+"-"+nodeInfo.NodeName, -1)

	return nil
}

func createBc(chainId, confYml string) (map[string]int, error) {
	err := os.RemoveAll("release/")
	if err != nil {
		log.Error("Remove org path err : " + err.Error())
	}
	chainInfo, err := chain.GetChainByChainId(chainId)
	if err != nil {
		log.Error("GetChainByChainId err : " + err.Error())
		return nil, err
	}

	chainOrgs, err := relation.GetChainOrgList(chainId)
	if err != nil {
		log.Error("GetChainOrgList err : " + err.Error())
		return nil, err
	}

	trustList := []*config.TrustRootsConf{}
	nodeList := []*config.NodesConf{}
	nodeIdMap := map[string]int{}
	conf := new(config.Chainmaker)
	yamlFile, _ := ioutil.ReadFile(confYml + "/config_tpl/chainmaker.yml")
	_ = yaml.Unmarshal(yamlFile, conf)
	address := conf.NetConf.ListenAddr
	portList := strings.Split(address, "/")
	portStr := portList[len(portList)-1]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 12301
	}

	for _, chainOrg := range chainOrgs {
		var chainOrgNodes []*dbcommon.ChainOrgNode
		chainOrgNodes, err = relation.GetChainOrg(chainOrg.OrgId, chainId)
		if err != nil {
			log.Error("GetChainOrg err : " + err.Error())
			return nil, err
		}
		nodes := config.NodesConf{}
		nodes.OrgId = chainOrg.OrgId

		for _, orgNode := range chainOrgNodes {
			nodeInfo, err := chain_participant.GetConsensusNodeByNodeName(orgNode.NodeName)
			if err != nil {
				log.Error("GetConsensusNodeByNodeName err : " + err.Error())
				continue
			}
			if nodeInfo.Type == chain_participant.NODE_CONSENSUS {
				nodes.NodeId = []string{nodeInfo.NodeId}
				port++
				nodeIdMap[nodeInfo.NodeId] = port
			}
		}
		nodeList = append(nodeList, &nodes)
	}

	bcConf := new(config.Bc)
	bcFile, err := ioutil.ReadFile(confYml + "/config_tpl/chainconfig/bc1.yml")
	if err != nil {
		log.Error(err.Error())
	}
	_ = yaml.Unmarshal(bcFile, bcConf)

	bcConf.ChainId = chainId
	bcConf.Block.TxTimeout = chainInfo.TxTimeout
	bcConf.Block.BlockTxCapacity = chainInfo.BlockTxCapacity
	bcConf.Consensus.Nodes = nodeList
	bcConf.TrustRoots = trustList
	bcConf.Consensus.Type = consensus.ConsensusType_value[chainInfo.Consensus]

	chainOrgNodes, err := relation.GetChainOrgByChainIdList(chainId)
	if err != nil {
		log.Error("GetChainOrgNode err : " + err.Error())
		return nil, err
	}
	err = removeBc(bcConf, chainOrgNodes, trustList, chainOrgs)
	if err != nil {
		log.Error("removeBc err : " + err.Error())
		return nil, err
	}
	return nodeIdMap, nil
}

func removeBc(bcConf *config.Bc, chainOrgNodes []*dbcommon.ChainOrgNode, trustList []*config.TrustRootsConf,
	chainOrgs []*dbcommon.ChainOrg) error {
	for _, chainOrgNode := range chainOrgNodes {
		trustList = trustList[0:0]
		for _, orgInfo := range chainOrgs {
			trust := config.TrustRootsConf{}
			trust.OrgId = orgInfo.OrgId
			trust.Root = []string{"../config/" + chainOrgNode.OrgId + "/certs/ca/" + orgInfo.OrgId + "/ca.crt"}
			trustList = append(trustList, &trust)
		}
		bcConf.TrustRoots = trustList
		bcBytes, _ := yaml.Marshal(bcConf)
		bc1, err := os.Create("bc1.yml")
		if err != nil {
			log.Error(err.Error())
			return err
		}
		defer func() {
			err = bc1.Close()
			if err != nil {
				log.Error("close file err :", err.Error())
			}
		}()

		_, err = bc1.Write(bcBytes)
		if err != nil {
			log.Error(err.Error())
		}
		err = os.MkdirAll("release/"+chainOrgNode.OrgId+"-"+chainOrgNode.NodeName+
			"/config/"+chainOrgNode.OrgId+"/chainconfig", os.ModePerm)
		if err != nil {
			log.Error("Mkdir bc1 path err : " + err.Error())
		}
		err = os.Rename("bc1.yml", "release/"+chainOrgNode.OrgId+"-"+chainOrgNode.NodeName+
			"/config/"+chainOrgNode.OrgId+"/chainconfig"+"/bc1.yml")
		if err != nil {
			log.Error("Rename bc1.yml err : " + err.Error())
		}
	}
	return nil
}
