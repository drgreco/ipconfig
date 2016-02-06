# ipconfig
web-service for command line to return certain IP information. Clone of ifconfig.me

### supported endpoints
**/** 

if your user agent contains the string 'curl', returns your public ip address in plain text
if your user agent does not contain 'curl', gives you a webpage built off of template.html will all known info

**/ip**
returns your public ip in plain text

**/host**
attempts to do a reverse dns lookup on your public ip, returns all responses. if that fails, you get your ip back again.

**/ua**
returns the user-agent string you reported

**/proto**
returns the HTTP protocol you have negotiated

**/port**
returns the client side port that has been negotiated. I really hope thats its somewhere between 1025 and 65535

**/lang **
returns the langauge that your client reports using

**/ref**
returns any referres (if you happened to be redirected to the service)

**/connection**
returns status of HTTP connection

**/method**
returns the method which you queried the service (GET,POST, etc)

**/encoding**
returns the HTTP encoding schemes your client wants to use

**/mime**
returns any mime-types reported

**/charset**
returns requested character sets

**/via**
returns forwarding information

**/forwarded**
returns proxying information

**/all**
returns all of the above information formatted for console output

**/all.xml**
returns all of the above information formatted in json

**/all.json**
returns all of the above information formatted in json


## running the service
The service is written in go, and is mostly a way for me to get more practice. There is a template.html too, in case you want that. I am using goji, which is amazingly simple and neat.
```
go get github.com/zenazn/goji
```

then either run, or build, then run that executable
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

###### Criticisims, suggestions, and questions are all welcome. 
