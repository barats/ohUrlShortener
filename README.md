# ohUrlShortener

适合中小型社区网站使用的短链接服务系统，支持短链接成产、查询及302转向，并顺带简单的点击量统计


![Screenshot](screenshot.jpg)

## 配置文件说明
根目录下 `config.ini` 中存放着关于 ohUrlShortener 短链接系统的一些必要配置，请在启动应用之前确保这些配置的正确性

```
[app]

应用是否以 debug 模式启动，主要作用会在go-gin 框架上体现（eg：日志输出等）
debug = false   

短链接系统本地启动端口
port = 9091

短链接系统管理后台本地启动端口
admin_port = 9092

短链接系统的完整 url 前缀，eg：https://t.cn/ 是前缀(不要忘记最后一个/符号)
url_prefix = http://localhost:9091/
```

## 短链接在应用启动时会存入 Redis 中

所有短链接再系统启动时会以 `Key(short_url) -> Value(original_url)` 的形式存储在 Redis 中。

### 1. 为什么要这么做？  

当短链接的查询请求进入应用时，为了能够更快、更准确的将用户请求转向到目标链接，与传统的方式从数据库中查询相比，直接从 Redis 中获取目标链接就会显得更有价值。

### 2. 这种处理方式有什么缺点？

理论上来说，如果 Redis 所在的服务器的内存较大的话，存储10w个Key也是可以的。但是，硬件条件不允许的情况下，就需要控制 Redis 中的 Key 数量（主要是怕机器扛不住，Redis 本身的性能不会有问题）。这部分的功能扩展，考虑在将来的某个版本中实现并允许配置管理。

### 3. 万一 

考虑到可扩展性，多封装了一层 `service`，以便需要的时候在业务逻辑层进行自定义扩展，eg：将 key 查询改成数据库查询等。  

## 短链接生产过程相关代码

所在文件 `core/short_url.go` 

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

## 定时器1分钟清理一次访问日志

所在文件 `main.go` 

```
const ACCESS_LOG_CLEAN_INTERVAL = 1 * time.Minute 

func startTicker() error {
	ticker := time.NewTicker(ACCESS_LOG_CLEAN_INTERVAL)
	for range ticker.C {
		log.Println("[StoreAccessLog] Start.")
		if err := service.StoreAccessLogs(); err != nil {
			log.Printf("Error while trying to store access_log %s", err)
		}
		log.Println("[StoreAccessLog] Finish.")
	}
	return nil
}
```

## Give Thanks To

由衷感谢以下开源软件、框架等（包括但不限于）

1. [gin-gonic/gin](https://github.com/gin-gonic/gin) 
2. [FomanticUI](https://fomantic-ui.com/)
3. [dchest/captcha](https://github.com/dchest/captcha) 
4. [Masterminds/sprig](https://github.com/Masterminds/sprig)
5. [go-redis/redis](https://github.com/go-redis/redis/) 
6. [jmoiron/sqlx](https://github.com/jmoiron/sqlx)
7. [go-ini/ini](https://github.com/go-ini/ini)
