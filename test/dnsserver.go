// Copyright 2011 Miek Gieben. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Reflect is a small name server which sends back the IP address of its client, the
// recursive resolver.
// When queried for type A (resp. AAAA), it sends back the IPv4 (resp. v6) address.
// In the additional section the port number and transport are shown.
//
// Basic use pattern:
//
//  dig @localhost -p 8053 whoami.miek.nl A
//
//  ;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 2157
//  ;; flags: qr rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1
//  ;; QUESTION SECTION:
//  ;whoami.miek.nl.            IN  A
//
//  ;; ANSWER SECTION:
//  whoami.miek.nl.     0   IN  A   127.0.0.1
//
//  ;; ADDITIONAL SECTION:
//  whoami.miek.nl.     0   IN  TXT "Port: 56195 (udp)"
//
// Similar services: whoami.ultradns.net, whoami.akamai.net. Also (but it
// is not their normal goal): rs.dns-oarc.net, porttest.dns-oarc.net,
// amiopen.openresolvers.org.
//
// Original version is from: Stephane Bortzmeyer <stephane+grong@bortzmeyer.org>.
//
// Adapted to Go (i.e. completely rewritten) by Miek Gieben <miek@miek.nl>.
package main

import (
    "flag"
    "fmt"
    "os"
    "os/signal"
    "runtime/pprof"
    "strings"
    "syscall"
    "github.com/miekg/dns"
    "net"
    "errors"
    "github.com/mageddo/go-logging"
)

var (
    cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
    compress   = flag.Bool("compress", false, "compress replies")
    tsig       = flag.String("tsig", "", "use MD5 hmac tsig: keyname:base64")
)

func handleReflect(respWriter dns.ResponseWriter, reqMsg *dns.Msg) {

    defer func() {
        err := recover()
        if err != nil {
            logging.Errorf("status=error, error=%v", err)
        }
    }()

    var questionName string
    if len(reqMsg.Question) != 0 {
        questionName = reqMsg.Question[0].Name
    }   else {
        questionName = "null"
    }

    logging.Infof("questions=%d, 1stQuestion=%s", len(reqMsg.Question), questionName)

    resp := SolveName(questionName)
    resp.SetReply(reqMsg)
    resp.Compress = *compress


    var firstAnswer dns.RR
    if len(resp.Answer) != 0 {
        firstAnswer = resp.Answer[0]
    }

    logging.Infof("resp=%v", firstAnswer)
    respWriter.WriteMsg(resp)

}

func serve(net, name, secret string) {
    switch name {
    case "":
        server := &dns.Server{Addr: ":53", Net: net, TsigSecret: nil}
        if err := server.ListenAndServe(); err != nil {
            fmt.Printf("Failed to setup the "+net+" server: %s\n", err.Error())
        }
    default:
        server := &dns.Server{Addr: ":53", Net: net, TsigSecret: map[string]string{name: secret}}
        if err := server.ListenAndServe(); err != nil {
            fmt.Printf("Failed to setup the "+net+" server: %s\n", err.Error())
        }
    }
}

func main2() {
    var name, secret string
    flag.Usage = func() {
        flag.PrintDefaults()
    }
    flag.Parse()
    if *tsig != "" {
        a := strings.SplitN(*tsig, ":", 2)
        name, secret = dns.Fqdn(a[0]), a[1] // fqdn the name, which everybody forgets...
    }
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            logging.Error(err)
            os.Exit(-1)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }

    dns.HandleFunc(".", handleReflect)
    go serve("tcp", name, secret)
    go serve("udp", name, secret)
    sig := make(chan os.Signal)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    s := <-sig
    fmt.Printf("Signal (%s) received, stopping\n", s)
}


// reference https://miek.nl/2014/August/16/go-dns-package/
func SolveName(hostname string) *dns.Msg {


    //config, _ := dns.ClientConfigFromFile("/etc/resolv.conf")     
    c := new(dns.Client)

    m := new(dns.Msg)
    m.SetQuestion(dns.Fqdn(hostname), dns.TypeA) // CAN BE A, AAA, MX, etc.     
    m.RecursionDesired = true

    //r, _, err := c.Exchange(m, net.JoinHostPort(config.Servers[0], config.Port)) // server and port to ask        
    r, _, err := c.Exchange(m, net.JoinHostPort("8.8.8.8", "53")) // server and port to ask     

    // if the answer not be returned        
    if r == nil {
        panic(err)
    }

    // what the code of the return message ?        
    if r.Rcode != dns.RcodeSuccess {
        panic(errors.New(fmt.Sprintf(" *** invalid answer name %s after MX query for %s", hostname, hostname)))
    }

    // looping through the anwsers      
    return r

}
