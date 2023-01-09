#!/usr/bin/env bash
printf '\n[GNUPG:] SIG_CREATED ' >&${1#--status-fd=}
signore sign --file /dev/stdin --signer $3 2>/dev/null
