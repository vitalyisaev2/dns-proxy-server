package main

import (
    "fmt"
    "os/exec"
    "syscall"
)

func main3(){

    //config, _ := dns.ClientConfigFromFile("/etc/resolv.conf")

    //err := syscall.Exec("/sbin/resolvconf", []string{}, []string{})

    cmd := exec.Command("/sbin/resolvconf", "-u")
    //bytes := make([]byte, 10000)
    out, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Errorf("err=%v", err)
    }
    fmt.Printf("pid=%d, exited=%t, exitCode=%d, out=%+v", cmd.ProcessState.Pid(), cmd.ProcessState.Exited(),
        cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), string(out))

    //for _, s := range config.Servers {
    //  fmt.Println(s)
    //}

}
