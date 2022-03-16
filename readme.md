# ohUrlShortener

一个适合中小型社区网站使用的短域名服务系统，支持短域名成产、查询及302转向，并顺带简单的点击量统计

## 所有短域名存储在 Redis 中

所有短域名以 `Key(short_url) -> Value(original_url)` 的形式存储在 Redis 中，理论上说：如果 Redis 服务器的内存较大的话，存储 10w 个 Key 应该也是可以的。但是，考虑到可扩展性，多封装了一层 `service`，以便需要的时候在业务逻辑层进行自定义扩展，eg：将 key 查询改成数据库查询等。  

## 短域名生产过程

```
func GenerateShortLink(initialLink string) (string, error) {
	if utils.EemptyString(initialLink) {
		return "", fmt.Errorf("empty string")
	}
	urlHash, err := utils.Sha256Of(initialLink)
	if err != nil {
		return "", err
	}
	number := new(big.Int).SetBytes(urlHash).Uint64()
	str := utils.Base58Encode([]byte(fmt.Sprintf("%d", number)))
	return str[:8], nil
}
```

## 定时器3分钟清理一次访问日志

```
	ticker := time.NewTicker(3 * time.Minute)
	for range ticker.C {
		log.Println("[StoreAccessLog] Start.")
		if err := service.StoreAccessLog(); err != nil {
			log.Printf("Error while trying to store access_log %s", err)
		}
		log.Println("[StoreAccessLog] Finish.")
	}
```

## 关于短域名的一些引用材料

详情请参见 [ohUrlShortener Linked References](references.md)

