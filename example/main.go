package main

import (
	"runtime"
	"server/data"
	"server/lib/ace_lib"
	"server/lib/app"
	"server/lib/conf"
	"server/lib/hope"
	"server/lib/log"
	"server/lib/red"
	"server/lib/sdk"
	"server/lib/sdk/pay"
	"server/lib/third_log/tgloglib"
	"server/lib/wmnet"
	"server/modules/garnet"
	"server/modules/ignore_input"
	"server/modules/l5"
	"server/modules/metrics"
	"server/sdk/txcos"
	"server/service/3rd/idip"
	"server/service/3rd/recharge_notify"
	"server/service/3rd/survey"
	"server/service/common/console"
	"server/service/common/db"
	"server/service/common/db/repo"
	"server/service/common/global"
	"server/service/common/http_server"
	"server/service/common/registry"
	"server/service/common/registry_cache"
	"server/service/common/ws_server"
	"server/service/cross/auction"
	"server/service/cross/camp"
	"server/service/cross/contact"
	"server/service/cross/contract"
	"server/service/cross/cross"
	"server/service/cross/dragon"
	"server/service/cross/rank"
	"server/service/game/activity"
	"server/service/game/admin"
	"server/service/game/agent"
	"server/service/game/agentmgr"
	"server/service/game/battlereport"
	"server/service/game/chat"
	"server/service/game/collocation"
	"server/service/game/complain"
	"server/service/game/game"
	"server/service/game/guild"
	"server/service/game/lantern"
	"server/service/game/mail"
	"server/service/game/rank"
	"server/service/game/recharge"
	"server/service/game/srv_role"
	"server/service/game/team"
	"server/service/game/yuewen"
	"server/service/gateway/gateway"
	"server/service/global/cache_loader"
	"server/service/global/dragon"
	"server/service/global/forward"
	"server/service/pull/client_env"
	"server/service/pull/idip_jx"
	"server/service/pull/playing_friend"
	"server/service/pull/share"
	"server/service/robot/robotmgr"
	"server/service/robot/robotmgr/robot"
	scene_service "server/service/scene/scene/service"
	"server/service/scene/scene_mgr"

	"github.com/mogud/snow/core/configuration/sources"
	"github.com/mogud/snow/core/host"
	"github.com/mogud/snow/core/host/builder"
	"github.com/mogud/snow/routines/node"
)

func main() {
	b := builder.NewDefaultBuilder()
	b.GetConfigurationManager().AddSource(&sources.JsonConfigurationSource{
		Path:           "conf_tmpl/app/app_tmpl.json5",
		Optional:       true,
		ReloadOnChange: false,
	})
	b.GetConfigurationManager().AddSource(&sources.JsonConfigurationSource{
		Path:           "conf/app/app.json5",
		Optional:       true,
		ReloadOnChange: false,
	})

	//enableCSClient := configuration.Get[bool](b.GetConfigurationManager(), "ConfigurationSystem:EnableClient")
	//if enableCSClient {
	//	source := configuration.Get[*cs_source.ConfigurationSystemSource](b.GetConfigurationManager(), "ConfigurationSystem:Source")
	//	b.GetConfigurationManager().AddSource(source)
	//}
	//
	//enableCSServer := configuration.Get[bool](b.GetConfigurationManager(), "ConfigurationSystem:EnableServer")
	//if enableCSServer {
	//	host.AddOption[*cs_server.Option](b, "ConfigurationSystem:Repository")
	//	host.AddHostedRoutine[*cs_server.CSServer](b)
	//}

	host.AddHostedRoutine[*ignore_input.IgnoreInput](b)

	// Log
	host.AddLogFormatter(b, "GameConsoleColor", log.GameColorLogFormatter)

	// Metrics
	host.AddOption[*metrics.Option](b, "Metrics")
	host.AddHostedRoutine[*metrics.Meter](b)

	// TGLog
	host.AddOption[*tgloglib.TGLog](b, "ThirdLog:TGLog")
	host.AddSingleton[*tgloglib.TGLogLib](b)

	// Hope 防沉迷
	host.AddOption[*hope.HopeLibOption](b, "Hope")
	host.AddSingleton[*hope.HopeLib](b)

	// 游戏服用的服务器下单
	host.AddOption[*pay.RechargeOption](b, "Recharge")
	// 全局服务器，接收发货回调通知
	host.AddOption[*pay.RechargeNotifyOption](b, "RechargeNotify")
	// ace
	host.AddOption[*ace_lib.AceLibOption](b, "Ace")
	host.AddHostedRoutine[*ace_lib.AceLib](b)

	//l5
	host.AddOption[*l5.L5Option](b, "L5")
	host.AddHostedRoutine[*l5.L5Lib](b)
	host.AddOption[*yuewen.Option](b, "YWReport")

	host.AddOption[*http_server.Option](b, "HttpServer")
	host.AddOption[*ws_server.Option](b, "WsServer")
	host.AddOption[*registry.Option](b, "Registry")
	host.AddOption[*scene_mgr.Option](b, "SceneManager")
	host.AddOption[*gateway.Option](b, "Gateway")
	host.AddOption[*sdk.Option](b, "SDK")
	host.AddOption[*agentmgr.Option](b, "AgentManager")
	host.AddOption[*data.MapleOptions](b, "Maple")
	host.AddOption[*repo.Option](b, "Db")
	host.AddOption[*cache_loader.Option](b, "CacheLoader")

	host.AddOption[*client_env.Option](b, "ClientEnv")

	host.AddOption[*idip.Option](b, "IDIP")
	host.AddOption[*idip_jx.Option](b, "IDIP_JX")
	host.AddOption[*share.Option](b, "Share")
	host.AddOption[*txcos.Option](b, "TXCos")
	host.AddSingleton[*txcos.TXCos](b)
	host.AddOption[*complain.ComplainOption](b, "Complain")
	host.AddOption[*forward.Option](b, "Forward")

	// robot
	host.AddOption[*robotmgr.Option](b, "RobotMgr")

	host.AddOption[*red.Option](b, "RedisClient")
	host.AddSingleton[*red.Redis](b)

	if runtime.GOOS == "windows" {
		host.AddOption[*garnet.Option](b, "Garnet")
		host.AddHostedLifecycleRoutine[*garnet.Garnet](b)
	}

	// Game
	host.AddOption[*global.Option](b, "Global")
	host.AddOption[*cross.Option](b, "Cross")
	host.AddOption[*game.Option](b, "Game")

	// App
	host.AddOption[*app.Option](b, "App")
	app.Add(b)

	// Node
	host.AddOption[*node.Option](b, "Node")
	node.AddNode(b, func() *node.RegisterOption {
		return &node.RegisterOption{
			ServiceRegisterInfos: []*node.ServiceRegisterInfo{
				// 通用服务
				node.CheckedServiceRegisterInfo[db.DB](100),
				node.CheckedServiceRegisterInfo[console.Console](101),
				node.CheckedServiceRegisterInfo[registry.Registry](103),
				node.CheckedServiceRegisterInfo[http_server.HttpServer](104),
				node.CheckedServiceRegisterInfoName[dragon.GlobalDragonMatch](105, "DragonMatch"),
				node.CheckedServiceRegisterInfo[ws_server.WsServer](106),

				// 唯一服务
				node.CheckedServiceRegisterInfo[registry_cache.RegistryCache](200),
				node.CheckedServiceRegisterInfo[cache_loader.CacheLoader](201),
				node.CheckedServiceRegisterInfo[idip.IDIP](203),
				node.CheckedServiceRegisterInfo[survey.Survey](205),
				node.CheckedServiceRegisterInfo[playing_friend.PlayingFriend](206),
				node.CheckedServiceRegisterInfo[recharge_notify.RechargeNotify](207),
				node.CheckedServiceRegisterInfoName[idip_jx.IDIPJX](208, "IDIP_JX"),
				node.CheckedServiceRegisterInfo[forward.Forward](209),

				// 节点服务
				node.CheckedServiceRegisterInfo[global.Global](300),
				node.CheckedServiceRegisterInfo[cross.Cross](301),
				node.CheckedServiceRegisterInfo[game.Game](302),

				// 客户端 http 服务
				node.CheckedServiceRegisterInfo[client_env.ClientEnv](350),

				// 跨服服务
				node.CheckedServiceRegisterInfo[auction.Auction](402),
				node.CheckedServiceRegisterInfo[team.Team](404),
				node.CheckedServiceRegisterInfo[battlereport.BattleReport](405),
				node.CheckedServiceRegisterInfo[contract.Contract](406),
				node.CheckedServiceRegisterInfo[cross_rank.CrossRank](407),
				node.CheckedServiceRegisterInfoName[cross_dragon.CrossDragon](408, "CrossDragonMatch"),
				node.CheckedServiceRegisterInfo[cross_area.CrossArea](409),

				// 跨服和单服都有
				node.CheckedServiceRegisterInfo[scene_mgr.SceneManager](450),
				node.CheckedServiceRegisterInfo[scene_service.Scene](451),

				// 单服服务
				node.CheckedServiceRegisterInfo[admin.Admin](502),
				node.CheckedServiceRegisterInfo[agentmgr.AgentManager](504),
				node.CheckedServiceRegisterInfo[agent.Agent](505),
				node.CheckedServiceRegisterInfo[recharge.Recharge](506),
				node.CheckedServiceRegisterInfoName[mail.Email](507, "Mail"),
				node.CheckedServiceRegisterInfo[chat.Chat](508),
				node.CheckedServiceRegisterInfo[contact.Contact](509),
				node.CheckedServiceRegisterInfo[collocation.Collocation](512),
				node.CheckedServiceRegisterInfo[camp.Camp](513),
				node.CheckedServiceRegisterInfo[guild.Guild](514),
				node.CheckedServiceRegisterInfo[lantern.Lantern](515),
				node.CheckedServiceRegisterInfo[complain.Complain](517),
				node.CheckedServiceRegisterInfo[rank.Rank](518),
				node.CheckedServiceRegisterInfo[activity.Activity](520),
				node.CheckedServiceRegisterInfo[share.Share](521),
				node.CheckedServiceRegisterInfo[yuewen.YWReport](522),
				node.CheckedServiceRegisterInfo[srv_role.ServiceRole](526),
				node.CheckedServiceRegisterInfo[gateway.Gateway](550),

				// robot
				node.CheckedServiceRegisterInfo[robotmgr.RobotMgr](600),
				node.CheckedServiceRegisterInfo[robot.Robot](601),
			},
			ClientHandlePreprocessor: wmnet.NodeClientPreprocessor,
			ServerHandlePreprocessor: wmnet.NodeServerPreprocessor,
			//MetricCollector:          host.GetRoutine[*metrics.Meter](b.GetRoutineProvider()),
			PostInitializer: func() {
				if len(app.GameConfigPath) > 0 {
					_ = conf.GetConfig()
				}
			},
		}
	})

	h := b.Build()

	host.Run(h)
}
