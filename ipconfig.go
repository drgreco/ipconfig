package main

import (
    "fmt"
    "regexp"
    "strings"
    "net"
    "net/http"
    "github.com/zenazn/goji"
    "github.com/zenazn/goji/web"
    "github.com/zenazn/goji/bind"
)

var ipRE = regexp.MustCompile("^[0-9]+.[0-9]+.[0-9]+.[0-9]+")
var portRE = regexp.MustCompile("[0-9]+$")
var curlRE = regexp.MustCompile("curl")

func index (c web.C, w http.ResponseWriter, r *http.Request) {
    var context = strings.TrimSpace(c.URLParams["context"])
    if context == "" {
        if curlRE.MatchString(ua(r)) {
            fmt.Fprintf(w, "%s\n", ip(r))
        } else {
            fmt.Fprintf(w, "%s", fullPage(r))
        }

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

    } else if context == "keepalive" {
        fmt.Fprintf(w, "%s\n", keepalive(r))

    } else if context == "connection" {
        fmt.Fprintf(w, "%s\n", connection(r))

    } else if context == "encoding" {
        fmt.Fprintf(w, "%s\n", encoding(r))

    } else if context == "mime" {
        fmt.Fprintf(w, "%s\n", mime(r))

    } else if context == "charset" {
        fmt.Fprintf(w, "%s\n", charset(r))

    } else if context == "via" {
        fmt.Fprintf(w, "%s\n", via(r))

    } else if context == "forwarded" {
        fmt.Fprintf(w, "%s\n", forwarded(r))

    } else if context == "all" {
        fmt.Fprintf(w, "%s", allDump(r))
    } else if context == "all.xml" {

    } else if context == "all.json" {

    } else if context == "" {
        if curlRE.MatchString(ua(r)) == false {
            fmt.Fprintf(w, "%s", fullPage(r))
        } else if curlRE.MatchString(ua(r)) ==  true {
            fmt.Fprintf(w, "%s\n", ip(r))
        }
    } else {
        http.NotFound(w, r)
    }
}

func ip (r *http.Request) string {
    return fmt.Sprintf("%s", ipRE.FindString(r.RemoteAddr))
}

func host (r *http.Request) string {
    if ipRE.MatchString(r.Host) {
        var list []string
        hosts,err := net.LookupAddr(ipRE.FindString(r.Host))
        if err != nil {
            return ""
        }
        for _,host := range hosts {
            list = append(list, host)
        }
        return strings.Join(list, ", ")
    } else {
        return fmt.Sprintf("%s", ipRE.FindString(r.Host))
    }
}

func ua (r *http.Request) string {
    var header = r.Header
    return fmt.Sprintf("%s", header["User-Agent"])
}

func port (r *http.Request) string {
    return fmt.Sprintf("%s", portRE.FindString(r.RemoteAddr))
}

func proto (r *http.Request) string {
    return fmt.Sprintf("%s", r.Proto)
}

func lang (r *http.Request) string {
    return fmt.Sprintf("%s", getHeader(r, "Accept-Language"))
}

func keepalive (r *http.Request) string {
    //TODO figure out what this is from
    return ""
}

func method (r *http.Request) string {
    return fmt.Sprintf("%s", r.Method)
}

func connection (r *http.Request) string {
    return fmt.Sprintf("%s", getHeader(r, "Connection"))
}

func encoding (r *http.Request) string {
    return fmt.Sprintf("%s", getHeader(r, "Accept-Encoding"))
}

func mime (r *http.Request) string {
    //TODO parse multipartreader https://golang.org/pkg/net/http/#Request.MultipartReader
    return ""
}

func charset (r *http.Request) string {
    return fmt.Sprintf("%s", getHeader(r, "Accept-Charset"))
}

func via (r *http.Request) string {
    return fmt.Sprintf("%s", getHeader(r, "Via"))
}

func forwarded (r *http.Request) string {
    return fmt.Sprintf("%s", getHeader(r, "Forwarded"))
}

func all (r *http.Request) map[string]string {
    var allData map[string]string
    allData = make(map[string]string)
    allData["IP Address"] = ip(r)
    allData["Port"] = port(r)
    allData["User Agent"] = ua(r)
    allData["Hostname"] = host(r)
    allData["Language"] = lang(r)
    allData["Encoding"] = encoding(r)
    allData["Protocol"] = proto(r)
    allData["Keepalive"] = keepalive(r)
    allData["Method"] = method(r)
    allData["Connection"] = connection(r)
    allData["Mime"] = mime(r)
    allData["Charset"] = charset(r)
    allData["Via"] = via(r)
    allData["Forwarded"] = forwarded(r)

    return allData
}

func allDump (r *http.Request) string {
    var allData = all(r)
    var stringData []string

    //stringData = append(stringData, "\n--What We Know--\n")
    for k, v := range allData {
        stringData = append(stringData, strings.TrimSpace(k),":",strings.TrimSpace(v),"\n")
    }

    //stringData = append(stringData, "\n--What is still in Headers--\n")
    //for k, v := range r.Header {
    //    stringData = append(stringData, strings.TrimSpace(k),":",strings.TrimSpace(strings.Join(v,"")),"\n")
    //}

    //stringData = append(stringData, "\n--What is still in Trailers--\n")
    //for k, v := range r.Trailer {
    //    stringData = append(stringData, strings.TrimSpace(k),":",strings.TrimSpace(strings.Join(v,"")),"\n")
    //}

    return strings.Join(stringData, "")
}

func allXml (c web.C, w http.ResponseWriter, r *http.Request) string {
    return ""
}

func allJson (c web.C, w http.ResponseWriter, r *http.Request) string {
    return ""
}

func getHeader (r *http.Request, key string) string {
    return r.Header.Get(key)
}

func fullPage (r *http.Request) string {
    var allData = all(r)
    var stringData []string

    stringData = append(stringData, "<head><style>")
    stringData = append(stringData, "table{border-collapse: collapse;}")
    stringData = append(stringData, "table, td, th { border: 1px solid black;}")
    stringData = append(stringData, "</style></head>")
    stringData = append(stringData, "<body><table>")
    for k, v := range allData {
        stringData = append(stringData, fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", k, v))
    }
    stringData = append(stringData, "</table></body>")
    return strings.Join(stringData, "")
}

// Yay, we are gonna do the mains
func main() {
    goji.Get("/", index)
    goji.Get("/:context", index)
    goji.ServeListener(bind.Socket("0.0.0.0:8080"))
}
