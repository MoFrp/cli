package sub

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/pkg/config"
	"github.com/fatedier/frp/pkg/config/source"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/policy/security"
	"github.com/fatedier/frp/pkg/util/banner"
	"github.com/fatedier/frp/pkg/util/log"
)

var (
	quickstartToken  string
	quickstartMaster string
)

func init() {
	quickstartCmd.Flags().StringVarP(&quickstartToken, "token", "t", "", "quickstart token (required)")
	quickstartCmd.Flags().StringVarP(&quickstartMaster, "master", "m", "https://api.mofrp.moiu.cn", "master backend address")
	_ = quickstartCmd.MarkFlagRequired("token")
	rootCmd.AddCommand(quickstartCmd)
}

var quickstartCmd = &cobra.Command{
	Use:   "quickstart",
	Short: "Quick start frpc with token from master backend",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runQuickstart(quickstartToken, quickstartMaster)
	},
}

type QuickstartResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		NodeName      string  `json:"node_name"`
		TunnelName    string  `json:"tunnel_name"`
		Type          string  `json:"type"`
		Config        string  `json:"config"`
		ConnectPublic *string `json:"connect_public,omitempty"`
	} `json:"data"`
}

func runQuickstart(token, masterAddr string) error {
	log.InitLogger("console", "info", 3, false)

	banner.DisplayCLIBanner()

	parts := strings.SplitN(token, ":", 2)
	var tunnelID, actualToken string
	if len(parts) == 2 {
		tunnelID = parts[0]
		actualToken = parts[1]
	} else {
		actualToken = token
	}

	url := fmt.Sprintf("%s/api/v1/frpc/config?token=%s&tunnel_id=%s", masterAddr, actualToken, tunnelID)

	log.Infof("正在从主控获取配置: %s", masterAddr)

	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Get(url)
	if err != nil {
		log.Errorf("从主控获取配置失败: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("读取响应失败: %v", err)
		os.Exit(1)
	}

	var result QuickstartResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Errorf("解析响应失败: %v", err)
		os.Exit(1)
	}

	if result.Code != 0 {
		log.Errorf("主控返回错误: %s", result.Msg)
		os.Exit(1)
	}

	log.Infof("已获取到配置文件: 隧道=%s, 类型=%s, 节点=%s",
		result.Data.TunnelName, result.Data.Type, result.Data.NodeName)

	var allCfg v1.ClientConfig
	if err := config.LoadConfigure([]byte(result.Data.Config), &allCfg, true); err != nil {
		log.Errorf("解析配置失败: %v", err)
		os.Exit(1)
	}

	common := allCfg.ClientCommonConfig
	if err := common.Complete(); err != nil {
		log.Errorf("配置补全失败: %v", err)
		os.Exit(1)
	}

	log.Infof("正在启动隧道: %s", result.Data.TunnelName)

	var proxyCfgs []v1.ProxyConfigurer
	for _, c := range allCfg.Proxies {
		proxyCfgs = append(proxyCfgs, c.ProxyConfigurer)
	}

	visitorCfgs := make([]v1.VisitorConfigurer, 0, len(allCfg.Visitors))
	for _, c := range allCfg.Visitors {
		visitorCfgs = append(visitorCfgs, c.VisitorConfigurer)
	}

	proxyCfgs = config.CompleteProxyConfigurers(proxyCfgs)

	configSource := source.NewConfigSource()
	if err := configSource.ReplaceAll(proxyCfgs, visitorCfgs); err != nil {
		log.Errorf("设置配置源失败: %v", err)
		os.Exit(1)
	}

	aggregator := source.NewAggregator(configSource)
	unsafeFeatures := security.NewUnsafeFeatures(nil)

	var connectPublic string
	if result.Data.ConnectPublic != nil {
		connectPublic = *result.Data.ConnectPublic
	}

	svr, err := client.NewService(client.ServiceOptions{
		Common:                 &common,
		ConfigSourceAggregator: aggregator,
		UnsafeFeatures:         unsafeFeatures,
		ConfigFilePath:         "",
		ConnectPublic:          connectPublic,
	})
	if err != nil {
		log.Errorf("创建客户端服务失败: %v", err)
		os.Exit(1)
	}

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		log.Infof("正在关闭...")
		svr.GracefulClose(500 * time.Millisecond)
	}()

	if err := svr.Run(context.Background()); err != nil {
		log.Errorf("%v", err)
		os.Exit(1)
	}
	return nil
}
