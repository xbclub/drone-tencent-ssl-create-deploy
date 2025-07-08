package main

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"strings"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

type Plugin struct {
	SecretId       string
	SecretKey      string
	PublicKey      string
	PrivateKey     string
	InstanceIdList []string
	ResourceType   string
}

func (p *Plugin) Execute() error {
	credential := common.NewCredential(p.SecretId, p.SecretKey)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ssl.tencentcloudapi.com"

	client, _ := ssl.NewClient(credential, "", cpf)
	request := ssl.NewUploadCertificateRequest()

	request.CertificatePublicKey = common.StringPtr(p.PublicKey)
	request.CertificatePrivateKey = common.StringPtr(p.PrivateKey)

	response, err := client.UploadCertificate(request)
	if err != nil {
		if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
			return fmt.Errorf("API error: %v", sdkErr)
		}
		return err
	}

	logx.Infof("Certificate updated successfully: %s", response.ToJsonString())
	if response.Response.CertificateId == nil || len(*response.Response.CertificateId) == 0 {
		return fmt.Errorf("certificate ID NOT FOUND")
	}
	deployRequest := ssl.NewDeployCertificateInstanceRequest()
	deployRequest.CertificateId = response.Response.CertificateId
	deployRequest.InstanceIdList = common.StringPtrs(p.InstanceIdList)
	deployRequest.ResourceType = common.StringPtr(p.ResourceType)

	deployResponse, err := client.DeployCertificateInstance(deployRequest)
	if err != nil {
		if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
			return fmt.Errorf("Certificate Deploy API error: %v", sdkErr)
		}
		return err
	}
	logx.Infof("Certificate Deploy successfully: %s", deployResponse.ToJsonString())
	return nil
}

// readCertificateFile 读取证书文件内容
func readCertificateFile(filePath string) (string, error) {
	if filePath == "" {
		return "", fmt.Errorf("certificate file path is empty")
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read certificate file '%s': %v", filePath, err)
	}

	return string(content), nil
}

func main() {
	logconfig := logx.LogConf{
		Level:    "info",
		Mode:     "console",
		Encoding: "plain",
		Stat:     false,
	}
	logx.MustSetup(logconfig)
	// 从环境变量获取证书文件路径
	publicKeyPath := os.Getenv("PLUGIN_PUBLIC_KEY")
	privateKeyPath := os.Getenv("PLUGIN_PRIVATE_KEY")

	// 读取公钥文件
	publicKeyContent, err := readCertificateFile(publicKeyPath)
	if err != nil {
		logx.Errorf("Error reading public key: %s", err)
		os.Exit(1)
	}

	// 读取私钥文件
	privateKeyContent, err := readCertificateFile(privateKeyPath)
	if err != nil {
		logx.Errorf("Error reading private key: %s", err)
		os.Exit(1)
	}
	instanceIdList := strings.Split(os.Getenv("PLUGIN_DEPLOY_DOMAIN"), ",")

	plugin := Plugin{
		SecretId:       os.Getenv("PLUGIN_SECRET_ID"),
		SecretKey:      os.Getenv("PLUGIN_SECRET_KEY"),
		InstanceIdList: instanceIdList,
		ResourceType:   os.Getenv("PLUGIN_RESOURCE_TYPE"),
		PublicKey:      publicKeyContent,
		PrivateKey:     privateKeyContent,
	}
	if err := plugin.Execute(); err != nil {
		logx.Errorf("Error: %s", err)
		os.Exit(1)
	}
}
