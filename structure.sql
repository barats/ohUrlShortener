-- Database Structure For ohUrlShortener
CREATE DATABASE oh_url_shortener ENCODING 'UTF8';

-- Connect to database repostats
\c oh_url_shortener

CREATE TABLE public.short_urls (
  id serial4 NOT NULL,
	short_url text NOT NULL,
	dest_url varchar(200) NOT NULL,
	sha varchar(100) NOT NULL,
	created_at timestamp with time zone NOT NULL DEFAULT now(),
	is_valid bool NOT NULL DEFAULT true,	
	CONSTRAINT short_urls_pk PRIMARY KEY (id),
	CONSTRAINT short_urls_sha_un UNIQUE (sha)
);

CREATE TABLE public.access_logs (
	id serial4 NOT NULL,
	short_url varchar(200) NOT NULL,
	access_time timestamp with time zone NOT NULL DEFAULT NOW(),
	ip varchar(32) NULL,
	user_agent varchar(500) NULL,
	CONSTRAINT access_logs_pk PRIMARY KEY (id)
);
CREATE INDEX access_logs_short_url_idx ON public.access_logs (short_url);
CREATE INDEX access_logs_access_time_idx ON public.access_logs (access_time);
CREATE INDEX access_logs_ip_idx ON public.access_logs (ip);


CREATE VIEW public.url_puv AS
SELECT 
	l.SHORT_URL AS "Short_URL",
	count(l.ip) AS "IP_Count",
	count(DISTINCT(l.ip)) AS "Distinct_IP_Count"
FROM public.access_logs l
GROUP BY l.short_url;

CREATE VIEW public.puv_by_date AS 
SELECT
	date(l.access_time) AS "Access_Date",
	count(l.ip) AS "IP_Count",
	count(DISTINCT(l.ip)) AS "Distinct_IP_Count"
FROM public.access_logs l
GROUP BY date(l.access_time)
ORDER BY "Access_Date" DESC;