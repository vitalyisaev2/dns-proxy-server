package local

import (
    "github.com/mageddo/dns-proxy-server/cache/store"
    "github.com/mageddo/dns-proxy-server/events/local/localvo"
    "github.com/stretchr/testify/assert"
    "io/ioutil"
    "testing"
)

func TestSaveConfiguration_ClearCacheAfterChangeConfiguration(t *testing.T) {

    // arrange

    ResetConf()

    expectedHostname := "github.io"

    conf, err := LoadConfiguration()
    assert.Nil(t, err, "could not load conf")

    cache := store.GetInstance()
    env, _ := conf.GetActiveEnv()

    assert.False(t, cache.ContainsKey(expectedHostname))
    foundHostname, _ := env.GetHostname(expectedHostname)
    assert.Nil(t, foundHostname)

    // setting the host
    cache.Put(expectedHostname, foundHostname)
    assert.True(t, cache.ContainsKey(expectedHostname))
    assert.Nil(t, cache.Get(expectedHostname))

    // changing value for the hostname at configuration database
    hostname := localvo.Hostname{Ip: "192.168.0.2", Ttl:30, Hostname: expectedHostname, Type:"A"}
    assert.Nil(t, AddHostname( "", hostname))

    // cache must be clear after add a hostname in conf
    assert.False(t, cache.ContainsKey(expectedHostname))

    conf, err = LoadConfiguration()
    env, _  = conf.GetActiveEnv()

    foundHostname, _ = env.GetHostname(expectedHostname)
    assert.Equal(t, "192.168.0.2", foundHostname.Ip)

}

func TestShouldSaveARecord(t *testing.T) {

    // arrange
    ResetConf()

    expectedHostname := "github.io"

    conf, err := LoadConfiguration()
    assert.Nil(t, err, "could not load conf")

    // act
    assert.Nil(t, conf.AddHostname( "", localvo.Hostname{Ip: "192.168.0.2", Ttl:30, Hostname: expectedHostname, Type:localvo.A}))

    // assert

    env, _ := conf.GetActiveEnv()
    hostnameVo, _ := env.GetHostname("github.io")
    assert.Equal(t, expectedHostname, hostnameVo.Hostname)
    assert.Equal(t, localvo.A, hostnameVo.Type)

}


func TestShouldSaveCnameRecord(t *testing.T) {
    // arrange
    ResetConf()
    expectedHostname := "github.io"

    // act
    assert.Nil(t, AddHostname( "", localvo.Hostname{Ip: "192.168.0.2", Ttl:30, Hostname: expectedHostname, Type:localvo.CNAME}))

    // assert
    conf, err := LoadConfiguration()
    assert.Nil(t, err, "could not load conf")
    env, _ := conf.GetActiveEnv()
    hostnameVo, _ := env.GetHostname("github.io")
    assert.Equal(t, expectedHostname, hostnameVo.Hostname)
    assert.Equal(t, localvo.CNAME, hostnameVo.Type)
}

func TestShouldLoadV2ConfigurationVO(t *testing.T){
    // arrange
    ResetConf()
    const configJson = `
    {
        "version": 2,
        "remoteDnsServers": ["7.6.5.3", "7.6.5.4:54"],
        "envs": [
            {
                "name": "",
                "hostnames": [
                    {
                        "id": 1556725127318816137,
                        "hostname": "github.com",
                        "ip": "192.168.0.1",
                        "ttl": 50,
                        "type": "A"
                    }
                ]
            }
        ]
    }
    `
    assert.Nil(t, ioutil.WriteFile(confPath, []byte(configJson), 0766))
    expectedHostname := localvo.Hostname{Ip: "192.168.0.2", Ttl: 30, Hostname: "github.io", Type: localvo.CNAME}

    // act
    assert.Nil(t, AddHostname( "", expectedHostname))

    // assert
    conf, err := LoadConfiguration()
    assert.Nil(t, err, "could not load conf")
    env, _ := conf.GetActiveEnv()
    hostnameVo, _ := env.GetHostname("github.io")
    assert.Equal(t, expectedHostname.Hostname, hostnameVo.Hostname)
    assert.Equal(t, expectedHostname.Ip, hostnameVo.Ip)
    assert.Equal(t, expectedHostname.Target, hostnameVo.Target)
    assert.Equal(t, expectedHostname.Type, hostnameVo.Type)
    assert.Equal(t, 2, len(conf.RemoteDnsServers))

    assert.Equal(t, "7.6.5.3", conf.RemoteDnsServers[0].Ip)
    assert.Equal(t, 53, conf.RemoteDnsServers[0].Port)

    assert.Equal(t, "7.6.5.4", conf.RemoteDnsServers[1].Ip)
    assert.Equal(t, 54, conf.RemoteDnsServers[1].Port)
}

func TestShouldLoadV1ConfigurationVO(t *testing.T){
    // arrange
    ResetConf()
    const configJson = `
    {
        "remoteDnsServers": [[7,6,5,3], [7,6,5,4]],
        "envs": [
            {
                "name": "",
                "hostnames": [
                    {
                        "id": 1556725127318816137,
                        "hostname": "github.com",
                        "ip": [192,168,0,1],
                        "ttl": 50,
                        "type": "A"
                    }
                ]
            }
        ]
    }
    `
    assert.Nil(t, ioutil.WriteFile(confPath, []byte(configJson), 0766))
    expectedHostname := localvo.Hostname{Ip: "192.168.0.2", Ttl: 30, Hostname: "github.io", Type: localvo.CNAME}

    // act
    assert.Nil(t, AddHostname( "", expectedHostname))

    // assert
    conf, err := LoadConfiguration()
    assert.Nil(t, err, "could not load conf")
    env, _ := conf.GetActiveEnv()
    hostnameVo, _ := env.GetHostname("github.io")
    assert.Equal(t, expectedHostname.Hostname, hostnameVo.Hostname)
    assert.Equal(t, expectedHostname.Ip, hostnameVo.Ip)
    assert.Equal(t, expectedHostname.Target, hostnameVo.Target)
    assert.Equal(t, expectedHostname.Type, hostnameVo.Type)
    assert.Equal(t, 2, len(conf.RemoteDnsServers))

    assert.Equal(t, "7.6.5.3", conf.RemoteDnsServers[0].Ip)
    assert.Equal(t, 53, conf.RemoteDnsServers[0].Port)

    assert.Equal(t, "7.6.5.4", conf.RemoteDnsServers[1].Ip)
    assert.Equal(t, 53, conf.RemoteDnsServers[1].Port)
}

func TestDefaultStorageApiVersion(t *testing.T){
    // arrange
    ResetConf()

    expectedHostname := localvo.Hostname{Ip: "192.168.0.2", Ttl: 30, Hostname: "github.io", Type: localvo.CNAME}

    // act
    assert.Nil(t, AddHostname( "", expectedHostname))

    // assert
    conf, err := LoadConfiguration()
    assert.Nil(t, err, "could not load conf")
    assert.Equal(t, 2, conf.Version)
}
