# running the service
The service is written in go, and is mostly a way for me to get more practice. There is a template.html too, in case you want that. I am using goji, which is amazingly simple and neat.
```
go get github.com/zenazn/goji
```

Then either run, or build, then run that executable:
```
go run ipconfig
```
or
```
go build ipconfig
./ipconfig
```

If you poke around the code, you will see (in main() at the end) that I decided to change the default port to 8080. This was a personal decision since I like running my socks proxies on 8000.

### example nginx config
Just in case you want it, its really nothing fancy. Seriously though, who remembers how to write this offhand?
```
server {
        listen       80;
        listen       [::]:80;
        server_name  HOSTNAME;
        root         /var/www/default/html;

        # Load configuration files for the default server block.
        include /etc/nginx/default.d/*.conf;

        location / {
                proxy_pass                 http://localhost:8080;
                proxy_set_header Host      $host;
                proxy_set_header X-Real-IP $remote_addr;
        }

}
```

###### Criticisms, suggestions, and questions are all welcome. 
