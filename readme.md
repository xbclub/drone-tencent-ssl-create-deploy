# Drone 腾讯云 SSL 证书上传和部署插件

这是一个 Drone CI/CD 插件，用于自动化上传 SSL 证书到腾讯云并部署到指定的实例。 (支持配置EdgeOne）

## 功能特性

- 支持上传 SSL 证书到腾讯云
- 支持将证书部署到指定的实例
- 支持多个部署目标（支持配置EdgeOne）
- 自动化处理证书更新流程

## 系统要求

- 腾讯云账号和相关权限
- 有效的 SSL 证书文件（公钥和私钥）

## 安装

## 环境变量

插件通过以下环境变量进行配置：

| 环境变量 | 描述 | 是否必需 |
|---------|------|----------|
| `PLUGIN_SECRET_ID` | 腾讯云 API 密钥 ID | 是 |
| `PLUGIN_SECRET_KEY` | 腾讯云 API 密钥 Key | 是 |
| `PLUGIN_OLD_CERTIFICATE_ID` | 需要更新的旧证书 ID | 是 |
| `PLUGIN_RESOURCE_TYPE` | 资源类型（如 CDN、CLB 等） | 是 |
| `PLUGIN_PUBLIC_KEY` | 公钥文件路径 | 是 |
| `PLUGIN_PRIVATE_KEY` | 私钥文件路径 | 是 |

## 注意事项

- 确保腾讯云账号具有 SSL 证书相关的 API 权限
- 建议使用 Drone 的 `from_secret` 来管理敏感信息
- 插件会在证书文件读取失败时自动终止并输出错误信息

### 方式一：使用 Drone pipeline
* 以自动部署到EdgeOne 为例
```yaml
kind: pipeline
name: default

steps:
  - name: tencent-ssl
    image: xbclub/drone-tencent-ssl-create-deploy
    settings:
      secret_id:
        from_secret: tencent_secret_id
      secret_key:
        from_secret: tencent_secret_key
      resource_type: teo
      deploy_domain: "www.domain1.com,blog.domain1.com"
      public_key: /path/to/cert.pem
      private_key: /path/to/key.pem
```

### 方式二：使用 Docker
> 注意 不要使用 -- restart always 参数 否则将重复部署不会自动终止
```bash
bash docker run --rm \
-e PLUGIN_SECRET_ID="your-secret-id" \
-e PLUGIN_SECRET_KEY="your-secret-key" \
-e PLUGIN_OLD_CERTIFICATE_ID="your-old-certificate-id" \
-e PLUGIN_RESOURCE_TYPES="CDN" \
-e PLUGIN_PUBLIC_KEY="/certs/public.crt" \
-e PLUGIN_PRIVATE_KEY="/certs/private.key" \
-v /path/to/your/certs:/certs \
xbclub/drone-tencent-ssl-create-deploy
```

