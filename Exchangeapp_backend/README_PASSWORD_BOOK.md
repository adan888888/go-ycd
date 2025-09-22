# 数字密码本 API 文档

## 概述
数字密码本功能允许用户安全地存储和管理密码信息，支持增删改查、搜索、批量操作等功能。

## 功能特性
- ✅ 创建密码项
- ✅ 获取密码列表（支持分页和搜索）
- ✅ 获取单个密码项
- ✅ 更新密码项
- ✅ 删除密码项
- ✅ 批量删除密码项
- ✅ 用户权限控制（只能访问自己的密码）
- ✅ 搜索功能（支持标题、用户名、网站、备注搜索）

## API 端点

### 1. 创建密码项
```
POST /api/password-book
```

**请求体：**
```json
{
  "title": "小幺鸡",
  "username": "shitlaoge@gmail.com",
  "password": "password123",
  "website": "http://example.com",
  "notes": "备注信息"
}
```

**响应：**
```json
{
  "code": 200,
  "msg": "创建成功",
  "data": {
    "id": 1,
    "title": "小幺鸡",
    "username": "shitlaoge@gmail.com",
    "password": "password123",
    "website": "http://example.com",
    "notes": "备注信息",
    "created_at": "2025-01-22T10:30:00Z",
    "updated_at": "2025-01-22T10:30:00Z"
  }
}
```

### 2. 获取密码列表
```
GET /api/password-book?page=1&size=10&keyword=小幺鸡
```

**查询参数：**
- `page`: 页码（默认：1）
- `size`: 每页数量（默认：10）
- `keyword`: 搜索关键词（可选）

**响应：**
```json
{
  "code": 200,
  "msg": "查询成功",
  "data": {
    "total": 1,
    "items": [
      {
        "id": 1,
        "title": "小幺鸡",
        "username": "shitlaoge@gmail.com",
        "password": "password123",
        "website": "http://example.com",
        "notes": "备注信息",
        "created_at": "2025-01-22T10:30:00Z",
        "updated_at": "2025-01-22T10:30:00Z"
      }
    ]
  }
}
```

### 3. 获取单个密码项
```
GET /api/password-book/{id}
```

**响应：**
```json
{
  "code": 200,
  "msg": "查询成功",
  "data": {
    "id": 1,
    "title": "小幺鸡",
    "username": "shitlaoge@gmail.com",
    "password": "password123",
    "website": "http://example.com",
    "notes": "备注信息",
    "created_at": "2025-01-22T10:30:00Z",
    "updated_at": "2025-01-22T10:30:00Z"
  }
}
```

### 4. 更新密码项
```
PUT /api/password-book/{id}
```

**请求体：**
```json
{
  "title": "小幺鸡（更新）",
  "username": "shitlaoge@gmail.com",
  "password": "newpassword123",
  "website": "http://example.com",
  "notes": "更新后的备注"
}
```

### 5. 删除密码项
```
DELETE /api/password-book/{id}
```

**响应：**
```json
{
  "code": 200,
  "msg": "删除成功",
  "data": null
}
```

### 6. 批量删除密码项
```
POST /api/password-book/batch-delete
```

**请求体：**
```json
[1, 2, 3]
```

**响应：**
```json
{
  "code": 200,
  "msg": "批量删除成功",
  "data": {
    "deleted_count": 3
  }
}
```

## 认证
所有密码本 API 都需要用户认证。请在请求头中包含有效的认证令牌：

```
Authorization: Bearer <your-token>
```

## 错误响应
```json
{
  "code": 400,
  "msg": "请求参数错误",
  "data": null
}
```

## 数据库表结构
```sql
CREATE TABLE password_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(500) NOT NULL,
    website VARCHAR(500) DEFAULT '',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_title (title),
    INDEX idx_username (username),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

## 安全考虑
1. **用户隔离**：每个用户只能访问自己的密码项
2. **数据加密**：建议在存储前对密码进行加密
3. **HTTPS**：生产环境必须使用 HTTPS
4. **访问控制**：所有 API 都需要认证

## 使用示例

### 创建密码项
```bash
curl -X POST http://localhost:8080/api/password-book \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{
    "title": "小幺鸡",
    "username": "shitlaoge@gmail.com",
    "password": "password123",
    "website": "http://example.com",
    "notes": "备注信息"
  }'
```

### 获取密码列表
```bash
curl -X GET "http://localhost:8080/api/password-book?page=1&size=10&keyword=小幺鸡" \
  -H "Authorization: Bearer your-token"
```

### 更新密码项
```bash
curl -X PUT http://localhost:8080/api/password-book/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{
    "title": "小幺鸡（更新）",
    "username": "shitlaoge@gmail.com",
    "password": "newpassword123",
    "website": "http://example.com",
    "notes": "更新后的备注"
  }'
```

### 删除密码项
```bash
curl -X DELETE http://localhost:8080/api/password-book/1 \
  -H "Authorization: Bearer your-token"
```

## 集成到 Flutter 应用
在 Flutter 应用中，您需要：

1. **更新 API 服务**：添加密码本相关的 API 调用
2. **修改数据模型**：使用服务器数据而不是本地存储
3. **添加网络错误处理**：处理网络请求失败的情况
4. **实现数据同步**：支持在线/离线模式切换

## 部署说明
1. 确保数据库已创建 `password_items` 表
2. 启动服务器：`./main`
3. 访问 Swagger 文档：`http://localhost:8080/swagger/index.html`
4. 测试 API 端点
