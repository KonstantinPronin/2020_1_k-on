server {

    listen 443 ssl http2; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/kino-on.ru-0005/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/kino-on.ru-0005/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot

	server_name kino-on.ru www.kino-on.ru;

	set $front_root /./frontend;

# секция сжатия
	gzip on;
    gzip_disable "msie6";
    gzip_types text/plain text/css application/json
    application/x-javascript image/jpeg application/octet-stream
    text/xml application/xml application/xml+rss text/javascript
    application/javascript;

#секция оптимизации отправки данных - отправка разом
    tcp_nopush on;

    location  /sw.js {
        root $front_root/dist;
	    try_files $uri $request_uri =404;    }

	location  /static/img {
	    root $front_root;
	    # разбиения файла на кусочки, получается
	    sendfile           on;
        sendfile_max_chunk 500k;
    	try_files $uri $request_uri =404;
    }

	location / {
	    root $front_root;
	        try_files $uri $request_uri /dist/index.html;
	}


	location /api/ {
	    proxy_redirect off;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Host $server_name;
        proxy_pass_header X-CSRF-TOKEN;
	    if ($request_uri ~* "/api/(.*)") {
           proxy_pass  http://127.0.0.1:8080/$1;
		}
    }
}

server {

    if ($host = www.kino-on.ru) {
        return 301 https://$host$request_uri;
    } # managed by Certbot


    if ($host = kino-on.ru) {
        return 301 https://$host$request_uri;
    } # managed by Certbot

    return 404; # managed by Certbot
}