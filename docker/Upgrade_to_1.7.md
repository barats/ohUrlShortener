# ohUrlShortener v1.6 以下版本升级说明

ohUrlShortener 短链接系统 1.6 版本做了较大的结构性调整，关系到管理端统计数据项及统计策略。  
代码层面主要的变化可查阅：
1. https://github.com/barats/ohUrlShortener/commits/v1.6
1. https://gitee.com/barat/ohurlshortener/commits/v1.6 

升级到1.6版本的过程中：为了兼容已存在的旧数据，需要手动对数据库进行一定的结构性改造，其具体过程如下：

## 1. 升级数据库结构

### 1.1 登录至 PostgreSQL 容器环境中

```shell
sudo docker exec -it ohurlshortener_pg psql -U postgres
```

### 1.2 连接至 oh_url_shortener 数据库

```shell
\c oh_url_shortener
```

### 1.3 执行以下数据库脚本

```sql

-- Create table for top25 urls
CREATE TABLE public.stats_top25 (
	id serial4 NOT NULL,
	short_url varchar(200) NOT NULL,
	today_count int8 NOT NULL DEFAULT 0,
	d_today_count int8 NOT NULL DEFAULT 0,
	stats_time timestamp with time zone NOT NULL DEFAULT NOW(), 
	CONSTRAINT stats_tv_pk PRIMARY KEY (id)
);

-- Stored procedure for top25 urls 
CREATE FUNCTION p_stats_top25() RETURNS void AS $$
BEGIN
	RAISE NOTICE 'Procedure p_stats_top25() called';
	
	-- delete all records 
	DELETE FROM public.stats_top25 WHERE 1=1;

	-- insert fresh-new records
	INSERT INTO public.stats_top25(short_url,today_count,d_today_count,stats_time)  
		SELECT l.short_url AS short_url, COUNT(l.ip) AS today_count ,COUNT(DISTINCT(l.ip)) AS d_today_count, NOW() AS stats_time
		FROM public.access_logs l WHERE date(l.access_time) = date(NOW()) GROUP BY l.short_url ORDER BY today_count DESC LIMIT 25;
END; 
$$ LANGUAGE plpgsql;

-- Create table for sum view 
CREATE TABLE public.stats_sum (
	stats_key varchar(200) NOT NULL,
	stats_value int8 NOT NULL DEFAULT 0,
	CONSTRAINT stats_sum_key PRIMARY KEY (stats_key)
);

-- Insert pre-defined stats 
INSERT INTO public.stats_sum (stats_key,stats_value) VALUES  
	('today_count',0), ('d_today_count',0),
	('yesterday_count',0), ('d_yesterday_count',0),
	('last_7_days_count',0), ('d_last_7_days_count',0),
	('monthly_count',0), ('d_monthly_count',0);

-- Stored procedure for stats sum view 
CREATE FUNCTION p_stats_sum() RETURNS void AS $$ 
DECLARE 
	today_count int8;
	d_today_count int8;
	yesterday_count int8;
	d_yesterday_count int8;
	last_7_days_count int8;
	d_last_7_days_count int8;
	monthly_count int8;
	d_monthly_count int8;
BEGIN
	RAISE NOTICE 'Procedure p_stats_sum() called';
	
	SELECT COUNT(l.ip),COUNT(DISTINCT(l.ip)) INTO today_count,d_today_count 
	FROM public.access_logs l WHERE date(l.access_time) = date(NOW());

	SELECT COUNT(l.ip),COUNT(DISTINCT(l.ip)) INTO yesterday_count,d_yesterday_count 
	FROM public.access_logs l WHERE date(l.access_time) = (NOW() - INTERVAL '1 day')::date;

	SELECT COUNT(l.ip),COUNT(DISTINCT(l.ip)) INTO last_7_days_count,d_last_7_days_count 
	FROM public.access_logs l WHERE date(l.access_time) >= (NOW() - INTERVAL '7 day')::date;

	SELECT COUNT(l.ip),COUNT(DISTINCT(l.ip)) INTO monthly_count,d_monthly_count 
	FROM public.access_logs l WHERE DATE_PART('month', l.access_time) = DATE_PART('month',NOW());

	UPDATE public.stats_sum SET stats_value = 
	CASE
		WHEN stats_key = 'today_count' THEN today_count
		WHEN stats_key = 'd_today_count' THEN d_today_count
		WHEN stats_key = 'yesterday_count' THEN yesterday_count
		WHEN stats_key = 'd_yesterday_count' THEN d_yesterday_count
		WHEN stats_key = 'last_7_days_count' THEN last_7_days_count
		WHEN stats_key = 'd_last_7_days_count' THEN d_last_7_days_count
		WHEN stats_key = 'monthly_count' THEN monthly_count
		WHEN stats_key = 'd_monthly_count' THEN d_monthly_count
		ELSE 0
	END;	
END;
$$ LANGUAGE plpgsql;

-- Create table for ip url sum 
CREATE TABLE public.stats_ip_sum (
	short_url varchar(200) NOT NULL,
	today_count int8 NOT NULL DEFAULT 0,
	d_today_count int8 NOT NULL DEFAULT 0,
	yesterday_count int8 NOT NULL DEFAULT 0,
	d_yesterday_count int8 NOT NULL DEFAULT 0,
	last_7_days_count int8 NOT NULL DEFAULT 0,
	d_last_7_days_count int8 NOT NULL DEFAULT 0,
	monthly_count int8 NOT NULL DEFAULT 0,
	d_monthly_count int8 NOT NULL DEFAULT 0,
	total_count int8 NOT NULL DEFAULT 0,
	d_total_count int8 NOT NULL DEFAULT 0,
	CONSTRAINT stats_ip_sum_pk PRIMARY KEY (short_url)	
);

-- Stored procedure for ip url sum 
CREATE FUNCTION p_stats_ip_sum() RETURNS void AS $$ 
BEGIN 
	
	RAISE NOTICE 'Procedure p_stats_ip_sum() called';
	
	-- Delete all records
	DELETE FROM public.stats_ip_sum WHERE 1=1;

	-- Calculate new stats data 
	INSERT INTO public.stats_ip_sum(short_url,today_count,d_today_count,yesterday_count,d_yesterday_count,last_7_days_count,d_last_7_days_count,
		monthly_count,d_monthly_count,total_count,d_total_count)
		SELECT
			u.short_url,				
			(SELECT count(ip) FROM public.access_logs WHERE date(access_time) = date(NOW()) AND short_url = u.short_url),
			(SELECT count(DISTINCT(ip)) FROM public.access_logs WHERE date(access_time) = date(NOW()) AND short_url = u.short_url),
			
			(SELECT count(ip) FROM public.access_logs WHERE date(access_time) = (NOW() - INTERVAL '1 day')::date AND short_url = u.short_url),
			(SELECT count(DISTINCT(ip)) FROM public.access_logs WHERE date(access_time) = (NOW() - INTERVAL '1 day')::date AND short_url = u.short_url),
			
			(SELECT count(ip) FROM public.access_logs WHERE date(access_time) >= (NOW() - INTERVAL '7 day')::date AND short_url = u.short_url),	
			(SELECT count(DISTINCT(ip)) FROM public.access_logs WHERE date(access_time) >= (NOW() - INTERVAL '7 day')::date AND short_url = u.short_url),
			
			(SELECT count(ip) FROM public.access_logs WHERE DATE_PART('month',access_time) = DATE_PART('month',NOW()) AND short_url = u.short_url),
			(SELECT count(DISTINCT(ip)) FROM public.access_logs WHERE DATE_PART('month',access_time) = DATE_PART('month',NOW()) AND short_url = u.short_url),
			
			(SELECT count(ip) FROM public.access_logs WHERE short_url = u.short_url),
			(SELECT count(DISTINCT(ip)) FROM public.access_logs WHERE short_url = u.short_url)	
		FROM public.short_urls u 
			LEFT JOIN public.access_logs l ON u.short_url = l.short_url
		GROUP BY u.short_url;
END;
$$ LANGUAGE plpgsql;
```

### 1.4 测试数据库对象是否正确

执行以下 SQL 函数确认期是否能够正常执行

```sql
select * from p_stats_top25();
select * from p_stats_ip_sum();
select * from p_stats_sum();
```

### 1.5 退出 PostgreSQL 及 Docker 容器环境

```sql
\q
```

## 2. 重新启动 Docker 容器

### 2.1 停止当前运行的各容器

进入到 `ohurlshortener\docker` 目录中并执行

```shell
./stop_destory.sh
```

### 2.2 修改 `vars.env`

进入到 `ohurlshortener\docker` 目录中修改 `vars.env` 文件，将其中的 `OH_ADMIN_VERSION`、`OH_PORTAL_VERSION` 改到最新版

```shell
OH_PORTAL_VERSION=1.7
OH_ADMIN_VERSION=1.7
```

### 2.3 重启启动各容器

进入到 `ohurlshortener\docker` 目录中并执行

```shell
./one_step_start.sh
```

## 3. 检查容器及应用情况