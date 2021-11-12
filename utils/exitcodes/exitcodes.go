package exitcodes

import (
    "os"
    "github.com/mageddo/dns-proxy-server/utils"
    "syscall"
    "github.com/mageddo/go-logging"
)

const (
    SUCCESS = iota
    FAIL_SET_DNS_AS_DEFAULT = iota
    FAIL_START_WEB_SERVER = iota
    FAIL_START_DNS_SERVER = iota
)

func Exit(code int){
    logging.Errorf("status=exiting, code=%d", code)
    if code != SUCCESS {
        utils.Sig <- syscall.Signal(code)
        logging.Info("status=msg-posted")
    } else {
        os.Exit(code)
    }
}
