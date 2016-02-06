package main

import (
    "fmt"
    "regexp"
    "strings"
    "net"
    "net/http"
    "io/ioutil"
    "github.com/zenazn/goji"
    "github.com/zenazn/goji/web"
    "github.com/zenazn/goji/bind"
)

// compile a few regex to keep things going smoothl
var ipRE = regexp.MustCompile("^[0-9]+.[0-9]+.[0-9]+.[0-9]+")
var portRE = regexp.MustCompile("[0-9]+$")
var curlRE = regexp.MustCompile("curl")

// This array is to set the order we want things output in
var dataOrder = []string{
    "ip_addr",
    "remote_host",
    "user_agent",
    "protocol",
    "port",
    "language",
    "referrer",
    "connection",
    "method",
    "encoding",
    "mime",
    "charset",
    "via",
    "forwarded",
}

/////////////////////////////////////////////////
//
// This is wher all the magic happens. 
// We are treating any subpage as a parameter 
// and parsing it here
//
/////////////////////////////////////////////////
func index (c web.C, w http.ResponseWriter, r *http.Request) {

    // grab param named context.
    // returns empty string if context param doesn't exist
    var context = strings.TrimSpace(c.URLParams["context"])

    // if there is no context, 
    // but the user-agent matches curl, print the ip
    if context == "" {
        if curlRE.MatchString(ua(r)) {
            fmt.Fprintf(w, "%s\n", ip(r))

        // if the user-agent is anything but curl, assume its a browser
        } else {
            var output = fullPage(r)
            if output == "err" {
                http.NotFound(w, r)
            } else {
                fmt.Fprintf(w, "%s", fullPage(r))
            }
        }

    // matches the context to the function
    // all of the following functions should
    // only return strings, which are then
    // printed in plain text
    } else if context == "ip" {
        fmt.Fprintf(w, "%s\n", ip(r))

    } else if context == "host" {
        fmt.Fprintf(w, "%s\n", host(r))

    } else if context == "ua" {
        fmt.Fprintf(w, "%s\n", ua(r))

    } else if context == "port" {
        fmt.Fprintf(w, "%s\n", port(r))

    } else if context == "proto" {
        fmt.Fprintf(w, "%s\n", proto(r))

    } else if context == "lang" {
        fmt.Fprintf(w, "%s\n", lang(r))

    } else if context == "ref" {
        fmt.Fprintf(w, "%s\n", referrer(r))

    } else if context == "connection" {
        fmt.Fprintf(w, "%s\n", connection(r))

    } else if context == "encoding" {
        fmt.Fprintf(w, "%s\n", encoding(r))

    } else if context == "charset" {
        fmt.Fprintf(w, "%s\n", charset(r))

    } else if context == "via" {
        fmt.Fprintf(w, "%s\n", via(r))

    } else if context == "forwarded" {
        fmt.Fprintf(w, "%s\n", forwarded(r))

    } else if context == "mime" {
        fmt.Fprintf(w, "%s\n", mime(r))

    } else if context == "all" {
        fmt.Fprintf(w, "%s", allTxt(r))

    } else if context == "all.xml" {
        fmt.Fprintf(w, "%s\n", allXml(r))

    } else if context == "all.json" {
        fmt.Fprintf(w, "%s\n", allJson(r))

    // we did our best, but no function has been defined
    } else {
        http.NotFound(w, r)
    }
}

//////////////////////////////////////
//
// Function Definitions for Context
//
//////////////////////////////////////


//
// returns ip (8.8.8.8)
//
func ip (r *http.Request) string {
    return fmt.Sprintf("%s", ipRE.FindString(r.RemoteAddr))
}

//
// returns result of reverse dns lookup
// if error, return whatever was
// in the Host request field (probably the IP)
// else returns a string of all returned hosts
// (google-public-dns-b.google.com.)
//
func host (r *http.Request) string {
    var list []string
    hosts,err := net.LookupAddr(ipRE.FindString(ip(r)))
    if err != nil {
        return ip(r)
    }
    for _,host := range hosts {
        list = append(list, host)
    }
    return strings.Join(list, ", ")
}

//
// returns whatever the client reported as its user-agent
// (Mozilla/3.01Gold (X11; I; SunOS 5.5.1 sun4m))
//
func ua (r *http.Request) string {
    var header = r.Header
    return fmt.Sprintf("%s", header["User-Agent"])
}

//
// returns client port. In most cases its a random negotiatd port
// and better be greater than 1024: (37906)
//
func port (r *http.Request) string {
    return fmt.Sprintf("%s", portRE.FindString(r.RemoteAddr))
}

//
// returns the negotiated protocol (HTTP/1.1)
func proto (r *http.Request) string {
    return fmt.Sprintf("%s", r.Proto)
}

//
// returns the accepted-language header
// (en-US,en;q=0.5)
//
func lang (r *http.Request) string {
    return fmt.Sprintf("%s", getHeader(r, "Accept-Language"))
}

//
// returns the referrer if any
// (https://www.google.com)
// yes, you have to check for both
// because someone mispelled it in 
// the original spec
//
func referrer (r *http.Request) string {
    var iDontSpellGood = r.Referer()
    if iDontSpellGood != "" {
        return fmt.Sprintf("%s", r.Referer())
    } else {
        return fmt.Sprintf("%s", getHeader(r, "Referrer"))
    }
}

//
// returns method verb
// (GET)
//
func method (r *http.Request) string {
    return fmt.Sprintf("%s", r.Method)
}

//
// returns the connection header
// (keep-alive)
//
func connection (r *http.Request) string {
    return fmt.Sprintf("%s", getHeader(r, "Connection"))
}

//
// returns string of acceptable encodings
// (gzip, deflate)
//
func encoding (r *http.Request) string {
    return fmt.Sprintf("%s", getHeader(r, "Accept-Encoding"))
}

//
// retruns string of acceptable character sets
// (utf-8)
//
func charset (r *http.Request) string {
    return fmt.Sprintf("%s", getHeader(r, "Accept-Charset"))
}

//
// returns any known proxies through which the request was sent
// (1.0 fred, 1.1 example.com (Apache/1.1))
//
func via (r *http.Request) string {
    return fmt.Sprintf("%s", getHeader(r, "Via"))
}

//
// original information of client connecting through
// an HTTP prox
// (for=192.0.2.60;proto=http;by=203.0.113.43)
func forwarded (r *http.Request) string {
    return fmt.Sprintf("%s", getHeader(r, "Forwarded"))
}

//
// returns mime-type header
// (text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8)
//
func mime (r *http.Request) string {
    return fmt.Sprintf("%s", getHeader(r, "Content-Type"))
}

// 
// runs all functions, and returns them in a map
// this is used when rendering html, xml, json, and text
//
func all (r *http.Request) map[string]string {
    var allData map[string]string
    allData = make(map[string]string)
    allData["ip_addr"] = ip(r)
    allData["port"] = port(r)
    allData["user_agent"] = ua(r)
    allData["remote_host"] = host(r)
    allData["language"] = lang(r)
    allData["referrer"] = referrer(r)
    allData["encoding"] = encoding(r)
    allData["mime"] = mime(r)
    allData["protocol"] = proto(r)
    allData["method"] = method(r)
    allData["connection"] = connection(r)
    allData["charset"] = charset(r)
    allData["via"] = via(r)
    allData["forwarded"] = forwarded(r)

    return allData
}

//
// formats the map from the function all for text output
//
func allTxt (r *http.Request) string {
    var allData = all(r)
    var stringData []string
    for _, key := range dataOrder {
        stringData = append(
            stringData,
            strings.TrimSpace(key),
            ": ",
            strings.TrimSpace(allData[key]),"\n",
        )
    }

    return strings.Join(stringData, "")
}

//
// formats the map from the function all for text output
// overloaded so we don't have to call all()
//
func allTxtnc (allData map[string]string) string {
    var stringData []string
    for _, key := range dataOrder {
        stringData = append(
            stringData,
            strings.TrimSpace(key),
            ": ",
            strings.TrimSpace(allData[key]),"\n",
        )
    }

    return strings.Join(stringData, "")
}

//
// formats the map from the function all for XML output
//
func allXml (r *http.Request) string {
    var allData = all(r)
    var stringData []string
    stringData = append(
        stringData,
        "<info>\n",
    )
    for _, key := range dataOrder {
        stringData = append(
            stringData,
            "  <",
            strings.TrimSpace(key),
            ">",
            strings.TrimSpace(allData[key]),
            "</",
            strings.TrimSpace(key),
            ">\n",
        )
    }

    stringData = append(
        stringData,
        "</info>",
    )

    return strings.Join(stringData, "")
}

//
// formats the map from the function all for XML output
// overloaded to not call back to all()
//
func allXmlnc (allData map[string]string) string {
    var stringData []string
    stringData = append(
        stringData,
        "<info>\n",
    )
    for _, key := range dataOrder {
        stringData = append(
            stringData,
            "  <",
            strings.TrimSpace(key),
            ">",
            strings.TrimSpace(allData[key]),
            "</",
            strings.TrimSpace(key),
            ">\n",
        )
    }

    stringData = append(
        stringData,
        "</info>",
    )

    return strings.Join(stringData, "")
}

//
// formats the map from the function all for JSON output
//
func allJson (r *http.Request) string {
    var allData = all(r)
    var stringData []string

    stringData = append(stringData, "{")
    for _, key := range dataOrder {
        stringData = append(
            stringData,
            "\"",
            strings.TrimSpace(key),
            "\":\"",
            strings.TrimSpace(allData[key]),
            "\",",
        )
    }
    // remove the last character of teh last slice
    stringData[len(stringData)-1] = strings.TrimSuffix(stringData[len(stringData)-1], ",")

    stringData = append(stringData, "}")

    return strings.Join(stringData, "")
}

//
//
// formats the map from the function all for JSON output
// overloaded to not call back to all()
//
func allJsonnc (allData map[string]string) string {
    var stringData []string

    stringData = append(stringData, "{")
    for _, key := range dataOrder {
        stringData = append(
            stringData,
            "\"",
            strings.TrimSpace(key),
            "\":\"",
            strings.TrimSpace(allData[key]),
            "\",",
        )
    }
    // remove the last character of teh last slice
    stringData[len(stringData)-1] = strings.TrimSuffix(stringData[len(stringData)-1], ",")

    stringData = append(stringData, "}")

    return strings.Join(stringData, "")
}

// facilitates grabbing known headers
// returns all known values for header as a single csv string
//
func getHeader (r *http.Request, key string) string {
    return r.Header.Get(key)
}

//
// format the map from the fuction all for HTML output
//
func fullPage (r *http.Request) string {
    var allData = all(r)

    templateBytes,err := ioutil.ReadFile("template.html")

    if err != nil {
        return "err"
    }

    // compile some more regex for templating the html
    var urlRE = regexp.MustCompile("##URL##")
    var ip_addrRE = regexp.MustCompile("##ip_addr##")
    var remote_hostRE = regexp.MustCompile("##remote_host##")
    var user_agentRE = regexp.MustCompile("##user_agent##")
    var protocolRE = regexp.MustCompile("##protocol##")
    var portRE = regexp.MustCompile("##port##")
    var languageRE = regexp.MustCompile("##language##")
    var referrerRE = regexp.MustCompile("##referrer##")
    var connectionRE = regexp.MustCompile("##connection##")
    var methodRE = regexp.MustCompile("##method##")
    var encodingRE = regexp.MustCompile("##encoding##")
    var mimeRE = regexp.MustCompile("##mime##")
    var charsetRE = regexp.MustCompile("##charset##")
    var viaRE = regexp.MustCompile("##via##")
    var forwardedRE = regexp.MustCompile("##forwarded##")
    var allRE = regexp.MustCompile("##all##")
    var allxmlRE = regexp.MustCompile("##allxml##")
    var alljsonRE = regexp.MustCompile("##alljson##")

/*    // create a few regex we know we will need
    var urlRE = regexp.MustCompile("##URL##")
    var allRE = regexp.MustCompile("##all##")
    var allxmlRE = regexp.MustCompile("##allxml##")
    var alljsonRE = regexp.MustCompile("##alljson##")
*/
    // convert the byte array into a string
    template := string(templateBytes)

    // apply our few known regex
    template = urlRE.ReplaceAllString(template, r.Host)
    template = allRE.ReplaceAllString(template, allTxtnc(allData))
    template = allxmlRE.ReplaceAllString(template, allXmlnc(allData))
    template = alljsonRE.ReplaceAllString(template, allJsonnc(allData))

    template = ip_addrRE.ReplaceAllString(template, allData["ip_addr"])
    template = remote_hostRE.ReplaceAllString(template, allData["remote_host"])
    template = user_agentRE.ReplaceAllString(template, allData["user_agent"])
    template = protocolRE.ReplaceAllString(template, allData["protocol"])
    template = portRE.ReplaceAllString(template, allData["port"])
    template = languageRE.ReplaceAllString(template, allData["language"])
    template = referrerRE.ReplaceAllString(template, allData["referrer"])
    template = connectionRE.ReplaceAllString(template, allData["connection"])
    template = methodRE.ReplaceAllString(template, allData["method"])
    template = encodingRE.ReplaceAllString(template, allData["encoding"])
    template = mimeRE.ReplaceAllString(template, allData["mime"])
    template = charsetRE.ReplaceAllString(template, allData["charset"])
    template = viaRE.ReplaceAllString(template, allData["via"])
    template = forwardedRE.ReplaceAllString(template, allData["forwarded"])

/*    // loop through our special list
    // and apply the rest of the regex
    for _, key := range dataOrder {
        var stringData []string
        stringData = append(stringData, "##", key, "##")
        var re = regexp.MustCompile(strings.Join(stringData, ""))
        re.ReplaceAllString(template, allData[key])
        fmt.Println(fmt.Sprintf("##%s: %s##", key, allData[key]))
    }
*/
    return template
}

//
// Yay, we are gonna do the mains
//
func main() {
    goji.Get("/", index)
    goji.Get("/:context", index)
    goji.ServeListener(bind.Socket("0.0.0.0:8080"))
}
