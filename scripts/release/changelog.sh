#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

set -uo pipefail

CHANGIE_VERSION="${CHANGIE_VERSION:-1.24.0}"
SEMVER_VERSION="${SEMVER_VERSION:-7.7.3}"

function usage {
  cat <<-'EOF'
Usage: ./changelog.sh <command> [<options>]

Description:
  This script will update CHANGELOG.md with the given version and date and add the changelog entries.
  It will also set the internal/version/version.go file to the correct version.

Commands:
  generate <release-type>
    generate will create a new section in the CHANGELOG.md file for the given release
    type.
    `dev`: will update the changelog with the latest unreleased changes.
    `beta`: will generate a new beta release.
    `minor`: will make the initial minor release for this branch.
    `patch`: will generate a new patch release

  getVersion:
    Returns the current version (assumes 'generate' already ran).
EOF
}

function pleaseUseGNUsed {
    echo "Please install GNU sed to your PATH as 'sed'."
    exit 1
}

function generate {
    RELEASE_TYPE="${1:-}"

    if [[ -z "$RELEASE_TYPE" ]]; then
        echo "missing <release-type> argument"
        usage
        exit 1
    fi

    sed --version > /dev/null || pleaseUseGNUsed

    case "$RELEASE_TYPE" in

        dev)
        LATEST_VERSION=$(npx -y changie@$CHANGIE_VERSION latest -r --skip-prereleases)

        # Check if we released this version already
        if git tag -l "v$LATEST_VERSION" | grep -q "v$LATEST_VERSION"; then
            LATEST_VERSION=$(npx -y semver@$SEMVER_VERSION -i patch $LATEST_VERSION)
        fi

        COMPLETE_VERSION="$LATEST_VERSION-dev"

        npx -y changie@$CHANGIE_VERSION merge -u "# $LATEST_VERSION (Unreleased)"
        ;;

        beta)
        NEXT_VERSION=$(npx -y changie@$CHANGIE_VERSION next minor)
        BETA_NUMBER=$(git tag -l "v$NEXT_VERSION-beta*" | wc -l)
        BETA_NUMBER=$((BETA_NUMBER + 1))
        HUMAN_DATE=$(date +"%B %d, %Y") # Date in Janurary 1st, 2022 format
        COMPLETE_VERSION="$NEXT_VERSION-beta$BETA_NUMBER"

        npx -y changie@$CHANGIE_VERSION merge -u "## $COMPLETE_VERSION ($HUMAN_DATE)"
        ;;

        patch)
        COMPLETE_VERSION=$(npx -y changie@$CHANGIE_VERSION next patch)
        npx -y changie@$CHANGIE_VERSION batch patch
        npx -y changie@$CHANGIE_VERSION merge
        ;;

        minor)
        COMPLETE_VERSION=$(npx -y changie@$CHANGIE_VERSION next minor)
        npx -y changie@$CHANGIE_VERSION batch $COMPLETE_VERSION
        npx -y changie@$CHANGIE_VERSION merge
        ;;

        *)
        echo "invalid <release-type> argument"
        usage
        exit 1

        ;;
    esac

    setVersion $COMPLETE_VERSION
}

function main {
  case "$1" in
    generate)
    generate "${@:2}"
      ;;

    getVersion)
    getVersion
      ;;

    *)
      usage
      exit 1

      ;;
  esac
}

function getVersion {
    cat ./internal/version/version.go | sed -r -n 's/const version = "([^"]+)"/\1/p'
}

function setVersion {
    sed -i "s/const version =.*/const version = \"$1\"/" internal/version/version.go
}

main "$@"
exit $?
