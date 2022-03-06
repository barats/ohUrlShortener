# ohUrlShortener

一个适合中小型社区网站使用的短域名服务系统，支持短域名成产、查询及302转向，并顺带简单的点击量统计

## 所有短域名存储在 Redis 中

本系统适用于「中小型社区网站」所需要的短域名服务，所有短域名以 `short_url -> original_url` 的形式存储在 Redis 中  
但是，考虑到可扩展性，多封装了一层 `service`，以便需要的时候在业务逻辑层进行自定义扩展  

## 短域名生产规则

```
func GenerateShortLink(initialLink string) (string, error) {
	if utils.EemptyString(initialLink) {
		return "", fmt.Errorf("empty string")
	}
	urlHash, err := sha256Of(initialLink)
	if err != nil {
		return "", err
	}
	number := new(big.Int).SetBytes(urlHash).Uint64()
	str := encode([]byte(fmt.Sprintf("%d", number)))
	return str[:8], nil
}

func sha256Of(input string) ([]byte, error) {
	algorithm := sha256.New()
	_, err := algorithm.Write([]byte(strings.TrimSpace(input)))
	if err != nil {
		return nil, err
	}
	return algorithm.Sum(nil), nil
}

func encode(data []byte) string {
	return base58.Encode(data)
}
```

## 定时器每个1分钟清理一次访问日志

`main` 函数启动时，启动了 `Ticker` 每个65秒清理 Redis 中的访问日志

```
func setupTicker() {
	//sleep for 30s to make sure main process is gon
	time.Sleep(35 * time.Second)

	//Clear redis cache every 65 seconds
	ticker := time.NewTicker(65 * time.Second)
	for range ticker.C {
		log.Println("[StoreAccessLog] Start.")
		if err := service.StoreAccessLog(); err != nil {
			log.Printf("Error while trying to store access_log %s", err)
		}
		log.Println("[StoreAccessLog] Finish.")
	}
}
```

## 管理功能通过 API 请求完成

查看所有短域名  
`Get /admin/shorturl`

查看短域名访问统计  
`Get /admin/shorturl/:url/stats`

生成新的短域名
`POST shorturl`

## 关于短域名的一些引用材料

详情请参见 [ohUrlShortener Linked References](references.md)

