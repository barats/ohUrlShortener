-- Database Structure For ohUrlShortener
CREATE DATABASE oh_url_shortener ENCODING 'UTF8';

-- Connect to database repostats
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

INSERT INTO public.short_urls(short_url, dest_url, created_at, is_valid, memo)
VALUES('AC7VgPE9', 'https://www.gitlink.org.cn/baladiwei/ohurlshortener', '2022-04-01 17:31:41.270', true, '短链接系统 gitlink 页面');

INSERT INTO public.short_urls(short_url, dest_url, created_at, is_valid, memo)
VALUES('AvTkHZP7', 'https://gitee.com/barat/ohurlshortener', '2022-04-01 17:31:55.899', true, '短链接系统 gitee 页面');

INSERT INTO public.short_urls(short_url, dest_url, created_at, is_valid, memo)
VALUES('gkT39tb5', 'https://github.com/barats/ohUrlShortener', '2022-04-01 17:32:13.209', true, '短链接系统 github 页面');


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


CREATE VIEW public.url_ip_count_stats AS
SELECT
	u.short_url AS short_url,	
	(SELECT count(ip) FROM public.access_logs WHERE date(ACCESS_TIME) = date(NOW()) AND short_url = u.short_url) AS today_count,
	(SELECT count(ip) FROM public.access_logs WHERE date(ACCESS_TIME) = (NOW() - INTERVAL '1 day')::date AND short_url = u.short_url) AS yesterday_count,
	(SELECT count(ip) FROM public.access_logs WHERE date(ACCESS_TIME) = (NOW() - INTERVAL '7 day')::date AND short_url = u.short_url) AS last_7_days_count,	
	(SELECT count(ip) FROM public.access_logs WHERE DATE_PART('month',access_time) = DATE_PART('month',NOW()) AND short_url = u.short_url) AS monthly_count,
	(SELECT count(ip) FROM public.access_logs WHERE short_url = u.short_url) AS total_count,
	(SELECT count(DISTINCT(ip)) FROM public.access_logs WHERE date(ACCESS_TIME) = date(NOW()) AND short_url = u.short_url) AS d_today_count,
	(SELECT count(DISTINCT(ip)) FROM public.access_logs WHERE date(ACCESS_TIME) = (NOW() - INTERVAL '1 day')::date AND short_url = u.short_url) AS d_yesterday_count,
	(SELECT count(DISTINCT(ip)) FROM public.access_logs WHERE date(ACCESS_TIME) = (NOW() - INTERVAL '7 day')::date AND short_url = u.short_url) AS d_last_7_days_count,		
	(SELECT count(DISTINCT(ip)) FROM public.access_logs WHERE DATE_PART('month',access_time) = DATE_PART('month',NOW()) AND short_url = u.short_url) AS d_monthly_count,
	(SELECT count(DISTINCT(ip)) FROM public.access_logs WHERE short_url = u.short_url) AS d_total_count	
FROM public.short_urls u 
	LEFT JOIN public.access_logs l ON u.short_url = l.short_url
GROUP BY u.short_url;

CREATE VIEW public.sum_url_ip_count_stats AS 
SELECT 
	COUNT(l.ip) AS today_count,	
	COUNT(DISTINCT(l.ip)) AS d_today_count
FROM public.access_logs l
WHERE date(l.access_time) = date(NOW());


CREATE VIEW public.total_count_top25 AS 
SELECT s.*, u.id,u.dest_url,u.created_at,u.is_valid,u.memo
FROM public.url_ip_count_stats s, public.short_urls u 
WHERE u.short_url = s.short_url
ORDER BY s.today_count DESC 
LIMIT 25;