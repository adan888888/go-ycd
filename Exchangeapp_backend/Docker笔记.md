# Docker 学习笔记

## 1. Dockerfile vs docker-compose.yaml

### Dockerfile
- **作用**：定义如何构建 Docker 镜像
- **内容**：构建步骤、依赖安装、代码编译等
- **类比**：就像"菜谱"，告诉 Docker 如何"烹饪"你的应用

### docker-compose.yaml
- **作用**：定义如何运行和管理容器
- **内容**：服务配置、端口映射、卷挂载、网络等
- **类比**：就像"餐厅管理"，告诉 Docker 如何"经营"你的应用

### 关系
- `docker-compose.yaml` 依赖 `Dockerfile` 来构建镜像
- `Dockerfile` 定义了"做什么"
- `docker-compose.yaml` 定义了"怎么做"

## 2. 工作流程

### 构建阶段
```bash
docker-compose up -d --build
```
- Docker Compose 读取 `docker-compose.yaml`
- 发现 `build: .` 配置
- 在当前目录查找 `Dockerfile`
- 使用 `Dockerfile` 构建镜像 `ycd:v1.0`

### 运行阶段
```bash
docker-compose up -d
```
- 使用已构建的镜像 `ycd:v1.0`
- 根据 `docker-compose.yaml` 配置运行容器
- 挂载卷、映射端口等

## 3. 常用命令

### Docker Compose 命令
```bash
# 启动服务（构建+运行）
docker-compose up -d --build

# 启动服务（仅运行）
docker-compose up -d

# 停止并删除容器
docker-compose down

# 停止但不删除容器
docker-compose stop

# 重启服务
docker-compose restart

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### Docker 命令
```bash
# 查看所有容器
docker ps -a

# 查看所有镜像
docker images -a

# 删除悬空镜像
docker image prune -f

# 删除所有未使用的镜像
docker image prune -a -f

# 查看系统占用
docker system df

# 清理系统
docker system prune -f
```

## 4. 配置文件解析

### docker-compose.yaml 配置
```yaml
name: MyYcd # docker-compose 容器组名称
services:
  app:
    build: . # 使用当前目录的 Dockerfile 构建镜像
    image: ycd:v1.0 # 指定镜像名称和标签
    container_name: ycd # 指定容器名称
    volumes:
      - /etc/localtime:/etc/localtime:ro # 挂载主机时间配置
      - ./config.yml:/app/config.yml # 挂载配置文件
    ports:
      - "3000:3000" # 端口映射
```

### Dockerfile 关键点
```dockerfile
# 多阶段构建
FROM golang:alpine AS builder
# 构建阶段...

FROM alpine
# 运行阶段...
```

## 5. 常见问题解决

### 问题1：docker-compose down 不生效
**原因**：容器是用 `docker run` 启动的，不是 `docker-compose` 启动的
**解决**：
```bash
# 停止并删除容器
docker stop ycd && docker rm ycd

# 使用 docker-compose 启动
docker-compose up -d --build
```

### 问题2：出现 <none> 镜像
**原因**：Docker 构建过程中的中间镜像
**解决**：
```bash
# 清理悬空镜像
docker image prune -f
```

### 问题3：代码修改后不生效
**原因**：Go 是编译语言，需要重新构建
**解决**：
```bash
# 重新构建并启动
docker-compose up -d --build
```

## 6. 最佳实践

### 开发环境
- 使用 `docker-compose up -d --build` 启动
- 配置文件通过卷挂载，修改后重启即可
- 定期清理未使用的镜像

### 生产环境
- 使用 `docker-compose up -d` 启动
- 确保配置文件正确挂载
- 监控容器状态和日志

## 7. 文件结构
```
项目根目录/
├── Dockerfile              # 镜像构建文件
├── docker-compose.yaml     # 容器编排文件
├── config.yml              # 应用配置文件
├── main.go                 # Go 应用入口
└── Docker笔记.md           # 本文档
```

## 8. 端口映射详解

### 端口映射格式
```yaml
ports:
  - "宿主机端口:容器端口"
```

### 实际例子
```yaml
ports:
  - "3000:3000"  # 宿主机3000 → 容器3000
  - "8080:3000"  # 宿主机8080 → 容器3000
  - "9000:80"    # 宿主机9000 → 容器80
```

### Go 应用端口配置
**重要：Go 应用应该监听容器端口（第二个端口）**

```go
// main.go
port := global.AppConfig.App.Port  // 从配置文件读取
srv := &http.Server{
    Addr:    port,  // 监听容器端口
    Handler: r,
}
```

```yaml
# config.yml
app:
  port: :3000  # Go 应用监听容器3000端口
```

### 端口映射原理
```
浏览器 → localhost:宿主机端口 → 宿主机端口 → 容器端口 → Go应用
```

### 为什么 Go 应用写容器端口？
- ✅ Go 应用运行在容器内部
- ✅ 容器有自己的网络空间
- ✅ Docker 负责端口转发
- ✅ 应用不需要知道宿主机端口
- ✅ 更解耦、更灵活

### 访问方式
- 宿主机访问：`http://localhost:宿主机端口`
- 容器内访问：`http://localhost:容器端口`

### 记忆口诀
> Go 应用写容器端口，Docker 负责端口转发！

## 9. 总结

- **Dockerfile** 和 **docker-compose.yaml** 都是必需的
- **Dockerfile** 负责构建镜像
- **docker-compose.yaml** 负责运行容器
- **端口映射**：Go 应用监听容器端口，Docker 负责转发
- 两者配合使用，实现完整的容器化部署
- 开发时使用 `docker-compose` 统一管理，避免命令混乱

---
*最后更新：2024年12月*
