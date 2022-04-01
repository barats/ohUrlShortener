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
		SUM(s.today_count) AS today_count,
		SUM(s.yesterday_count) AS yesterday_count,
		SUM(s.last_7_days_count) AS last_7_days_count,
		SUM(s.monthly_count) AS monthly_count,
		SUM(s.total_count) AS total_count,
		SUM(s.d_today_count) AS d_today_count,
		SUM(s.d_yesterday_count) AS d_yesterday_count,
		SUM(s.d_last_7_days_count) AS d_last_7_days_count,
		SUM(s.d_monthly_count) AS d_monthly_count,
		SUM(s.d_total_count) AS d_total_count
	FROM public.url_ip_count_stats s;


CREATE VIEW public.total_count_top25 AS 
SELECT s.*, u.id,u.dest_url,u.created_at,u.is_valid,u.memo
FROM public.url_ip_count_stats s, public.short_urls u 
WHERE u.short_url = s.short_url
ORDER BY s.today_count DESC 
LIMIT 25;