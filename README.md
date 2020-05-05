### 1.根据一个根据路由器的架构，构建一个可执行文件

GOOS 构建的系统
GOARCH 架构
GOARM arm 版本

```bash
GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "-s -w"
```

### 2.配置文件config.json放在运行的目录下即可

authEmail authKey 字段在 [https://dash.cloudflare.com/profile/api-tokens](cloudflare获取)

domainName 字段为在 cloudflare 下的二级域名

ddns 下的 isCurrent 为 true，指用当前主机 DDNS; 为 false 则在网络邻居中查找 mac 地址对应的 ipv6 DDNS

```json
{
  "authEmail": "",
  "authKey": "",
  "domainName": "pjunjie.cc",
  "ddns": [
    {
      "host": "router.pjunjie.cc",
      "isCurrent": true,
      "mac": null
    },
    {
      "host": "test.pjunjie.cc",
      "isCurrent": false,
      "mac": "dd:ee:ff:aa:bb:cc"
    }
  ]
}
```
