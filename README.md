# ipconfig
web-service for command line to return certain IP information. Clone of ifconfig.me

### supported endpoints
**/** 

If your user agent contains the string 'curl', returns your public IP address in plain text.
If your user agent does not contain 'curl', gives you a webpage built off of template.html with all known information.

**/ip**
Returns your public IP in plain text.

**/host**
Attempts to do a reverse DNS lookup on your public IP, returns all responses. If that fails, you get your IP back again.

**/ua**
Returns the user-agent string you reported.

**/proto**
Returns the HTTP protocol you have negotiated.

**/port**
Returns the client-side port that has been negotiated. I really hope thats its somewhere between 1025 and 65535.

**/lang **
Returns the langauge that your client reports using.

**/ref**
Returns any referrers (if you happened to be redirected to the service).

**/connection**
Returns the status of HTTP connection.

**/method**
Returns the method which you queried the service (GET, POST, etc).

**/encoding**
Returns the HTTP encoding schemes your client wants to use.

**/mime**
Returns any mime-types reported.

**/charset**
Returns requested character sets.

**/via**
Returns forwarding information.

**/forwarded**
Returns proxying information.

**/all**
Returns all of the above information formatted for console output.

**/all.xml**
Returns all of the above information formatted in XML.

**/all.json**
Returns all of the above information formatted in JSON.


## running the service
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
