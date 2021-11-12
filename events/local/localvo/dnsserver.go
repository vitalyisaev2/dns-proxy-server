package localvo

import (
    "fmt"
    "github.com/mageddo/dns-proxy-server/utils/iputils"
    "regexp"
    "strconv"
)

type DNSServer struct {
    Ip string
    Port int
}

func (s *DNSServer) GetIpArray() [4]byte {
    return *iputils.ToIpByteArray(&[4]byte{}, s.Ip)
}

func (s *DNSServer) GetAddress() string {
    return fmt.Sprintf("%s:%d", s.Ip, s.Port)
}

func ToIpsByteArray(servers []DNSServer) [][4]byte {
    byteServers := make([][4]byte, len(servers))
    for i, server := range servers {
        byteServers[i] = server.GetIpArray()
    }
    return byteServers
}

func ToIpsStringArray(servers []DNSServer) []string {
    byteServers := make([]string, len(servers))
    for i, server := range servers {
        byteServers[i] = server.GetAddress()
    }
    return byteServers
}

func ByteArrayToDnsServer(byteArrayServers [][4]byte) []DNSServer {
    servers := make([]DNSServer, len(byteArrayServers))
    for i, byteArrayServer := range byteArrayServers {
        servers[i] = toDnsServer(iputils.ToIpString(byteArrayServer))
    }
    return servers
}

func StringArrayToDnsServer(stringArrayServers []string) []DNSServer {
    servers := make([]DNSServer, len(stringArrayServers))
    for i, stringServer := range stringArrayServers {
        servers[i] = toDnsServer(stringServer)
    }
    return servers
}

func toDnsServer(dnsAddress string) DNSServer {
    regex := regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+):?(\d*)`)
    matches := regex.FindStringSubmatch(dnsAddress)
    return DNSServer{
        Ip:   matches[1],
        Port: parseIntDnsServer(matches[2]),
    }
}


func parseIntDnsServer(str string) int {
    if n, err := strconv.Atoi(str); err != nil{
        return 53
    } else {
        return n
    }
}
