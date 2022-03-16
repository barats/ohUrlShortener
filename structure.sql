-- Database Structure For ohUrlShortener
CREATE DATABASE oh_url_shortener ENCODING 'UTF8';

-- Connect to database repostats
\c oh_url_shortener

CREATE TABLE public.short_urls (
  id serial4 NOT NULL,
	short_url text NOT NULL,
	dest_url varchar(200) NOT NULL,	
	created_at timestamp with time zone NOT NULL DEFAULT now(),
	is_valid bool NOT NULL DEFAULT true,	
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
	l.short_url AS short_url,
	(SELECT count(ip) FROM public.access_logs WHERE date(ACCESS_TIME) = date(NOW()) AND short_url = l.short_url) AS today_count,
	(SELECT count(ip) FROM public.access_logs WHERE date(ACCESS_TIME) = (NOW() - INTERVAL '1 day')::date AND short_url = l.short_url) AS yesterday_count,
	(SELECT count(ip) FROM public.access_logs WHERE DATE_PART('month',access_time) = DATE_PART('month',NOW()) AND short_url = l.short_url) AS monthly_count,
	(SELECT count(ip) FROM public.access_logs WHERE short_url = l.short_url) AS total_count,
	(SELECT count(DISTINCT(ip)) FROM public.access_logs WHERE date(ACCESS_TIME) = date(NOW()) AND short_url = l.short_url) AS d_today_count,
	(SELECT count(DISTINCT(ip)) FROM public.access_logs WHERE date(ACCESS_TIME) = (NOW() - INTERVAL '1 day')::date AND short_url = l.short_url) AS d_yesterday_count,
	(SELECT count(DISTINCT(ip)) FROM public.access_logs WHERE DATE_PART('month',access_time) = DATE_PART('month',NOW()) AND short_url = l.short_url) AS d_monthly_count,
	(SELECT count(DISTINCT(ip)) FROM public.access_logs WHERE short_url = l.short_url) AS d_total_count	
FROM public.access_logs l 
GROUP BY l.short_url;