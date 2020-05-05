package model

// DNS 解析
type DNS struct {
	// 当前需要解析的host名称
	Host string `json:"host"`
	// 是否为本机
	IsCurrent bool `json:"isCurrent"`
	// 如果不是本机则从邻居找mac地址相同的绑定
	Mac string `json:"mac"`
}

// Config 配置
type Config struct {
	// cloudflare 的邮箱
	AuthEmail string `json:"authEmail"`
	// 秘钥
	AuthKey string `json:"authKey"`
	// cloudflare 区域id
	ZoneID string `json:"zoneId"`
	// 顶级域名名称
	DomainName string `json:"domainName"`

	Ddns []DNS `json:"ddns"`
}

// GetConfigAllHost 获取全部host
func (c *Config) GetConfigAllHost() []string {
	tmp := make([]string, 0, len(c.Ddns))

	for _, dns := range c.Ddns {
		if dns.Host != "" {
			tmp = append(tmp, dns.Host)
		}
	}

	return tmp
}
