package main

import (
	"fmt"
	"time"

	"github.com/mogud/snow/routines/node"
)

type AlphaService struct {
	node.Service
	betaProxy node.IProxy
}

type PingRequest struct {
	Message string `json:"Message"`
}

func (s *AlphaService) Start(_ any) {
	s.Infof("Alpha service ready at %s", node.Config.CurNodeAddr)
	s.EnableRpc()
	s.EnableHttpRpc()
	s.betaProxy = s.CreateProxy("Beta")
	if s.betaProxy == nil {
		s.Warnf("Beta proxy unavailable")
	}
}

func (s *AlphaService) RpcPing(ctx node.IRpcContext, req *PingRequest) {
	if req == nil || len(req.Message) == 0 {
		ctx.Error(fmt.Errorf("message is required"))
		return
	}

	ctx.Return(map[string]any{
		"from":    "Alpha",
		"echo":    req.Message,
		"tsMilli": time.Now().UnixMilli(),
	})
}

func (s *AlphaService) HttpRpcPing(ctx node.IRpcContext, req *PingRequest) {
	s.RpcPing(ctx, req)
}

func (s *AlphaService) RpcBetaSum(ctx node.IRpcContext, req *SumRequest) {
	if req == nil || len(req.Values) == 0 {
		ctx.Error(fmt.Errorf("values is required"))
		return
	}

	if s.betaProxy == nil || !s.betaProxy.Avail() {
		ctx.Error(fmt.Errorf("beta service unavailable"))
		return
	}

	s.betaProxy.Call("Sum", req).
		Then(func(result map[string]any) {
			ctx.Return(result)
		}).
		Catch(func(err error) {
			ctx.Error(err)
		}).Done()
}

func (s *AlphaService) HttpRpcBetaSum(ctx node.IRpcContext, req *SumRequest) {
	s.RpcBetaSum(ctx, req)
}

type BetaService struct {
	node.Service
}

type SumRequest struct {
	Values []int64 `json:"Values"`
}

func (s *BetaService) Start(_ any) {
	s.Infof("Beta service ready at %s", node.Config.CurNodeAddr)
	s.EnableRpc()
	s.EnableHttpRpc()
}

func (s *BetaService) RpcSum(ctx node.IRpcContext, req *SumRequest) {
	if req == nil || len(req.Values) == 0 {
		ctx.Error(fmt.Errorf("values is required"))
		return
	}

	var total int64
	for _, v := range req.Values {
		total += v
	}

	ctx.Return(map[string]any{
		"from":  "Beta",
		"count": len(req.Values),
		"sum":   total,
	})
}

func (s *BetaService) HttpRpcSum(ctx node.IRpcContext, req *SumRequest) {
	s.RpcSum(ctx, req)
}

/*
测试示例：调用Alpha的Ping方法
curl -X POST http://127.0.0.1:19090/node/rpc/Alpha     -H "Content-Type: application/json"     -d '{"Func":"Ping","Post":false,"Args":[{"Message":"hi"}]}'
{"Result":[{"from":"Alpha","echo":"hi","tsMilli":1764231153878}]}

测试示例：调用Alpha的BetaSum方法 to Beta Service
curl -X POST http://127.0.0.1:19090/node/rpc/Alpha   -H "Content-Type: application/json"   -d '{"Func":"BetaSum","Post":false,"Args":[{"Values":[4,5,6]}]}'
{"Result":[{"from":"Beta","count":3,"sum":15}]}
*/
