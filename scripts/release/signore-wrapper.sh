#!/usr/bin/env bash
# Copyright IBM Corp. 2020, 2026
# SPDX-License-Identifier: MPL-2.0

printf '\n[GNUPG:] SIG_CREATED ' >&${1#--status-fd=}
signore sign --file /dev/stdin --signer $3 2>/dev/null
