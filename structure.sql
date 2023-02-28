-- Database Structure For ohUrlShortener
CREATE DATABASE oh_url_shortener ENCODING 'UTF8';

-- Connect to database oh_url_shortener
\c oh_url_shortener

CREATE TABLE public.short_urls (
  id serial4 NOT NULL,
	short_url varchar(200) NOT NULL,
	dest_url text NOT NULL,	
	created_at timestamp with time zone NOT NULL DEFAULT now(),
	is_valid bool NOT NULL DEFAULT true,	
	memo text,
	CONSTRAINT short_urls_pk PRIMARY KEY (id),
	CONSTRAINT short_urls_un UNIQUE (short_url)
);

-- Insert new data
INSERT INTO public.short_urls(short_url, dest_url, created_at, is_valid, memo) VALUES
	('AC7VgPE9', 'https://www.gitlink.org.cn/baladiwei/ohurlshortener', NOW(), true, '短链接系统 gitlink 页面'),
	('AvTkHZP7', 'https://gitee.com/barat/ohurlshortener', NOW(), true, '短链接系统 gitee 页面'),
	('gkT39tb5', 'https://github.com/barats/ohUrlShortener', NOW(), true, '短链接系统 github 页面'),
	('9HtCr7YN', 'https://www.ohurls.cn', NOW(), true, 'ohUrlShortener 短链接系统首页');


CREATE TABLE public.access_logs (
	id serial4 NOT NULL,
	short_url varchar(200) NOT NULL,
	access_time timestamp with time zone NOT NULL DEFAULT NOW(),
	ip varchar(64) NULL,
	user_agent varchar(1000) NULL,
	CONSTRAINT access_logs_pk PRIMARY KEY (id)
);
CREATE INDEX access_logs_short_url_idx ON public.access_logs (short_url);
CREATE INDEX access_logs_access_time_idx ON public.access_logs (access_time);
CREATE INDEX access_logs_ip_idx ON public.access_logs (ip);
CREATE INDEX access_logs_ua_idx ON public.access_logs (user_agent);

CREATE TABLE public.users (
  id serial4 NOT NULL,
	account varchar(200) NOT NULL,
	password text NOT NULL,			
	CONSTRAINT users_pk PRIMARY KEY (id),
	CONSTRAINT users_account_un UNIQUE (account)
);

-- account: ohUrlShortener password: -2aDzm=0(ln_9^1
INSERT INTO public.users (account, "password") VALUES('ohUrlShortener', 'EZ2zQjC3fqbkvtggy9p2YaJiLwx1kKPTJxvqVzowtx6t');


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