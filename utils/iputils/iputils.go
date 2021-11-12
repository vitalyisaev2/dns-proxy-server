package iputils

import (
    "fmt"
    "strconv"
    "strings"
)

func ToIpsByteArray(ips []string) [][4]byte {
    byteArrayIps := make([][4]byte, len(ips))
    for i, ip := range ips {
        ToIpByteArray(&byteArrayIps[i], ip)
    }
    return byteArrayIps
}

func ToIpByteArray(byteArray *[4]byte, ip string) *[4]byte {
    if len(strings.TrimSpace(ip)) == 0 {
        return byteArray
    }
    index := strings.Index(ip, ":")
    if index >= 0 {
        ip = ip[:index]
    }
    ipStringArray := strings.Split(ip, ".")

    for j, ipPiece := range ipStringArray {
        if intIpPiece, err := strconv.Atoi(ipPiece); err != nil {
            panic(err)
        } else {
            byteArray[j] = byte(intIpPiece)
        }
    }
    return byteArray
}

func ToIpStringArray(ipsArray [][4]byte) []string {
    ips := make([]string, len(ipsArray))
    for i := range ips {
        ips[i] = ToIpString(ipsArray[i])
    }
    return ips
}

func ToIpString(ipArray [4]byte) string {
    return fmt.Sprintf("%d.%d.%d.%d", ipArray[0], ipArray[1], ipArray[2], ipArray[3])
}
