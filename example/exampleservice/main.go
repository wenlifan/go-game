package main

import (
	"github.com/mogud/snow/core/configuration/sources"
	"github.com/mogud/snow/core/host"
	builder "github.com/mogud/snow/core/host/builder"
	"github.com/mogud/snow/routines/node"
)

func main() {
	b := builder.NewDefaultBuilder()
	b.GetConfigurationManager().AddSource(&sources.JsonConfigurationSource{
		Path: "./example.json",
	})

	// 绑定 Node 配置到 node.Option，保证 LocalIP 等字段从配置中加载
	host.AddOption[*node.Option](b, "Node")

	node.AddNode(b, func() *node.RegisterOption {
		return &node.RegisterOption{
			ServiceRegisterInfos: []*node.ServiceRegisterInfo{
				node.CheckedServiceRegisterInfoName[AlphaService](2001, "Alpha"),
				node.CheckedServiceRegisterInfoName[BetaService](2002, "Beta"),
			},
		}
	})

	host.Run(b.Build())
}
