#!/bin/bash

set -e

modify_smf_config() {
    local CONFIG_FILE="$1"
    cp "$CONFIG_FILE" "${CONFIG_FILE}.bak"
    
    awk '
    BEGIN { in_interfaces = 0; in_n3_block = 0; in_endpoints = 0; }
    /^        interfaces:/ { in_interfaces = 1; }
    /^          - interfaceType: N3/ { if (in_interfaces) in_n3_block = 1; }
    /^          - interfaceType:/ && !/N3/ { in_n3_block = 0; in_endpoints = 0; }
    /^            endpoints:/ { if (in_n3_block) in_endpoints = 1; }
    /^            [a-z]/ && !/^            endpoints:/ && !/^              -/ { in_endpoints = 0; }
    /^      [a-z]/ && !/^        interfaces:/ { in_interfaces = 0; in_n3_block = 0; in_endpoints = 0; }
    {
        if (in_n3_block && in_endpoints && /^              - 127\.0\.0\.8/) {
            gsub(/127\.0\.0\.8/, "127.0.0.1");
        }
        print $0;
    }
    ' "$CONFIG_FILE" > "${CONFIG_FILE}.tmp"
    
    mv "${CONFIG_FILE}.tmp" "$CONFIG_FILE"
}

modify_upf_config() {
    local CONFIG_FILE="$1"
    cp "$CONFIG_FILE" "${CONFIG_FILE}.bak"
    
    awk '
    BEGIN { in_iflist = 0; }
    /^  ifList:/ { in_iflist = 1; }
    /^    - addr:/ { 
        if (in_iflist) current_addr = $0;
    }
    /^      type: N3/ {
        if (in_iflist && current_addr) {
            gsub(/127\.0\.0\.8/, "127.0.0.1", current_addr);
        }
    }
    /^[^ ]/ && !/^  ifList:/ { in_iflist = 0; }
    {
        if ($0 ~ /^    - addr:/ && in_iflist) {
            current_addr_line = NR;
            next;
        }
        if (current_addr_line == NR - 1) {
            print current_addr;
            current_addr_line = 0;
        }
        print $0;
    }
    END {
        if (current_addr_line > 0) print current_addr;
    }
    ' "$CONFIG_FILE" > "${CONFIG_FILE}.tmp"
    
    mv "${CONFIG_FILE}.tmp" "$CONFIG_FILE"
}

modify_config() {
    [ ! -f "$1" ] && return 1
    
    if grep -q "smfName:\|nsmf-pdusession" "$1"; then
        modify_smf_config "$1"
    elif grep -q "gtpu:\|forwarder:" "$1"; then
        modify_upf_config "$1"
    else
        return 1
    fi
}

main() {
    if [ $# -eq 0 ]; then
        echo "Usage: $0 <free5gc/config/smfcfg.yaml | free5gc/config/upfcfg.yaml>"
        exit 1
    fi

    if ! modify_config "$1"; then
        echo "Failed to modify config file: $1"
        exit 1
    fi
}

main "$@"