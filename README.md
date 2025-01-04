# GVA (Go Vue Admin) 服务端

一个基于 Go + Gin + GORM + Redis 的现代化后台管理系统服务端框架。采用 DDD (Domain-Driven Design) 设计模式，实现了清晰的代码结构和业务逻辑分层。

## 🌟 特性

- 基于 DDD 架构，代码结构清晰，易于维护
- JWT 认证，安全可靠的用户鉴权
- RBAC 权限管理，灵活的角色权限控制
- 集成 Redis 缓存，提升系统性能
- 完整的单元测试，保证代码质量
- 标准的 RESTful API 设计

## 🏗️ 技术栈

- **Web 框架**: [Gin](https://gin-gonic.com/)
- **ORM**: [GORM](https://gorm.io/)
- **缓存**: [Redis](https://redis.io/)
- **认证**: JWT
- **配置**: YAML
- **测试**: Go testing + testify

## 📁 项目结构
```bash
gva-server/
├── cmd/                # 主程序入口
├── configs/            # 配置文件
├── internal/           # 内部代码
│   ├── domain/         # 领域层
│   │   ├── entity/     # 实体定义
│   │   ├── repository/ # 仓储接口
│   │   └── service/    # 领域服务
│   ├── infrastructure/ # 基础设施层
│   │   ├── config/     # 配置管理
│   │   ├── database/   # 数据库操作
│   │   └── redis/      # Redis操作
│   ├── interfaces/     # 接口层
│   │   ├── handler/    # 请求处理器
│   │   └── middleware/ # 中间件
│   └── application/    # 应用层
│       └── dto/        # 数据传输对象
└── pkg/                # 公共包
    ├── jwt/           # JWT工具
    └── utils/         # 通用工具
```

## 🚀 快速开始

### 环境要求
- Go 1.20+
- MySQL 5.7+
- Redis 6.0+

### 安装步骤

1. 克隆项目
\```bash
git clone https://github.com/yourusername/gva-server.git
cd gva-server
\```

2. 安装依赖
```bash
go mod tidy
```

3. 配置环境
```bash
cp configs/config.example.yaml configs/config.yaml
# 修改配置文件中的数据库和Redis连接信息
```

4. 运行项目
```bash
go run cmd/server/main.go
```

## 📚 API 文档

### 用户模块

#### 注册
```http
POST /api/v1/register
Content-Type: application/json

{
    "username": "testuser",
    "password": "password123"
}
```

#### 登录
```http
POST /api/v1/login
Content-Type: application/json

{
    "username": "testuser",
    "password": "password123"
}
```

## ✨ 当前功能

### 用户管理
- [x] 用户注册
- [x] 用户登录
- [x] JWT 认证
- [x] 角色绑定

### 权限管理
- [x] RBAC 基础结构
- [x] 角色定义
- [x] 权限定义

## 🚀 开发计划

### 近期计划
1. 用户管理模块完善
   - [ ] 用户信息修改
   - [ ] 密码重置
   - [ ] 头像上传

2. 权限管理增强
   - [ ] 动态权限分配
   - [ ] 菜单权限管理
   - [ ] 数据权限控制

3. 系统功能扩展
   - [ ] 操作日志记录
   - [ ] 系统监控
   - [ ] 定时任务

## 🧪 测试

运行所有测试：
```bash
go test ./...
```

运行特定测试：
```bash
go test ./internal/interfaces/handler -v
```

## 📝 开发规范

1. 代码规范
   - 遵循 Go 官方代码规范
   - 使用 gofmt 格式化代码
   - 添加必要的注释

2. Git 提交规范
   ```
   feat: 添加新功能
   fix: 修复问题
   docs: 修改文档
   style: 修改代码格式
   refactor: 代码重构
   test: 添加测试
   chore: 修改构建过程或辅助工具
   ```

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'feat: Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

## 📄 许可证

[MIT License](LICENSE)

## 🙏 鸣谢

感谢所有为项目做出贡献的开发者！

---

> 注：本项目仍在积极开发中，欢迎加入一起完善！ 