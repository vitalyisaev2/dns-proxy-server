package main

import (
    "os"
    "os/signal"
    "syscall"
    "fmt"
)

func main5() {


    sig := make(chan os.Signal)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    go func (){
        sig <- syscall.Signal(3)
    }()
    s := <-sig
    fmt.Println(s)

}
