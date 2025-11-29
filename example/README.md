# Example 目录说明

本目录包含各种使用示例，每个示例都有独立的子目录，方便管理和扩展。

## 目录结构

```
example/
├── README.md          # 本文件
├── main.go            # 主项目示例（游戏服务器）
└── snow/              # Snow 框架使用示例
    ├── main.go        # 主程序入口
    ├── services.go    # 服务定义
    └── config.json    # 配置文件
```

## 使用说明

### Snow 框架示例

Snow 框架使用示例位于 `snow/` 目录中。

#### 运行示例

```bash
# 进入 snow 示例目录
cd example/snow

# 运行示例
go run main.go services.go
```

#### 测试 HTTP RPC

```bash
# 测试 Echo 方法
curl -X POST http://127.0.0.1:20090/node/rpc/Echo \
  -H "Content-Type: application/json" \
  -d '{"Func":"Echo","Post":false,"Args":[{"Message":"hello"}]}'

# 测试 GetTime 方法
curl -X POST http://127.0.0.1:20090/node/rpc/Echo \
  -H "Content-Type: application/json" \
  -d '{"Func":"GetTime","Post":false,"Args":[{}]}'
```

## 添加新示例

当需要添加新的示例时，请遵循以下规范：

1. 在 `example/` 目录下创建新的子目录，目录名应该清晰描述示例的用途
2. 每个示例目录应该包含：
   - `README.md`: 示例说明文档（可选）
   - 必要的源代码文件
   - 配置文件（如果需要）
3. 更新本 `README.md` 文件，添加新示例的说明

## 示例列表

- **snow**: Snow 框架使用示例，展示如何创建服务节点、实现 RPC 方法等

