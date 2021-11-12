package service

import (
    "github.com/mageddo/go-logging"
    "testing"
    "io/ioutil"
    "github.com/stretchr/testify/assert"
    "fmt"
    "github.com/mageddo/dns-proxy-server/flags"
    "github.com/mageddo/dns-proxy-server/utils"
    "time"
)

func TestSetupFor_NormalModeInstallStartSuccess(t *testing.T) {

    if !flags.IsTestVersion() {
        logging.Infof("status=test-skiped")
        return
    }

    sc := NewService()
    cmd := "'sh -c \"echo hi && sleep 20 && echo bye\"'"
    err := sc.SetupFor(DNS_PROXY_SERVER_PATH, DNS_PROXY_SERVER_SERVICE, &Script{cmd})
    assert.Nil(t, err)

    bytes, err := ioutil.ReadFile(DNS_PROXY_SERVER_PATH)
    assert.Nil(t, err)
    assert.Equal(t, fmt.Sprintf(SERVICE_TEMPLATE, cmd), string(bytes))

    time.Sleep(time.Second)

    out, err := ioutil.ReadFile("/var/log/dns-proxy-server.log")
    assert.Nil(t, err)
    assert.Contains(t, string(out), "hi")

    out, err, code := utils.Exec("sh", "-c", "ps aux | grep \"echo hi\"")
    assert.Equal(t, 0, code)
    assert.Nil(t, err)
    assert.Contains(t, string(out), "echo hi && sleep 20")

}
