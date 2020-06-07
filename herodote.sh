#!/usr/bin/env bash

set -o nounset -o pipefail -o errexit

var_read() {
  local SECRET_ARG=""
  if [[ ${3:-} == "secret" ]]; then
    SECRET_ARG="-s"
  fi

  if [[ -z ${!1:-} ]]; then
    if [[ -z ${SCRIPTS_NO_INTERACTIVE:-} ]]; then
      read ${SECRET_ARG?} -r -p "${1}${2:+ [${2}]}=" READ_VALUE
    fi

    eval "${1}=${READ_VALUE:-${2}}"

    if [[ ${SECRET_ARG} == "-s" ]]; then
      printf "\n"
    fi
  elif [[ ${SECRET_ARG} != "-s" ]]; then
    printf "%s=%s\n" "${1}" "${!1}"
  else
    printf "Using secret %s\n" "${1}"
  fi
}

var_color() {
  export RED='\033[0;31m'
  export GREEN='\033[0;32m'
  export BLUE='\033[0;34m'
  export YELLOW='\033[33m'
  export RESET='\033[0m'
}

git_conventionnal_commits() {
  declare -gA CONVENTIONAL_COMMIT_SCOPES
  CONVENTIONAL_COMMIT_SCOPES['build']='Changes that affect the build system or external dependencies'
  CONVENTIONAL_COMMIT_SCOPES['chore']='Changes in the core of the repository'
  CONVENTIONAL_COMMIT_SCOPES['ci']='Changes in Continuous Integration configuration files and scripts'
  CONVENTIONAL_COMMIT_SCOPES['docs']='Documentation only changes'
  CONVENTIONAL_COMMIT_SCOPES['feat']=$(printf 'A new feature for user %b(production change)%b' "${RED}" "${RESET}")
  CONVENTIONAL_COMMIT_SCOPES['fix']=$(printf 'A bug fix for user %b(production change)%b' "${RED}" "${RESET}")
  CONVENTIONAL_COMMIT_SCOPES['perf']=$(printf 'A performance improvement for user %b(production change)%b' "${RED}" "${RESET}")
  CONVENTIONAL_COMMIT_SCOPES['refactor']=$(printf 'A change that is not a feature not a bug %b(production change)%b' "${RED}" "${RESET}")
  CONVENTIONAL_COMMIT_SCOPES['style']='A change that do not affect the meaning of the code'
  CONVENTIONAL_COMMIT_SCOPES['test']='A new test or correcting existing tests'
}

git_is_inside() {
  git rev-parse --is-inside-work-tree 2>&1
}

insert_algolia() {
  curl -X POST \
     -H "X-Algolia-API-Key: ${ALGOLIA_API_KEY}" \
     -H "X-Algolia-Application-Id: ${ALGOLIA_APPLICATION_ID}" \
     --data-binary "${1}" \
    "https://${ALGOLIA_APPLICATION_ID}.algolia.net/1/indexes/${ALGOLIA_INDEX}"
}

walk_log() {
  git_conventionnal_commits

  IFS=$'\n'

  shopt -s nocasematch
  for hash in $(git log --no-merges --pretty=format:'%h' "${1}...${2}"); do
    if [[ $(git show -s --format='%h %at %s' "${hash}") =~ ^([0-9a-f]{1,16})\ ([0-9]+)\ (revert )?($(IFS='|'; echo "${!CONVENTIONAL_COMMIT_SCOPES[*]}"))(\((.+)\))?(\!)?:\ (.*)$ ]]; then
      local HASH="${BASH_REMATCH[1]}"
      local DATE="${BASH_REMATCH[2]}"
      local REVERT="${BASH_REMATCH[3]}"
      local TYPE="${BASH_REMATCH[4]}"
      local COMPONENT="${BASH_REMATCH[6]}"
      local BREAK="${BASH_REMATCH[7]}"
      local CONTENT="${BASH_REMATCH[8]}"

      if [[ -n ${REVERT} ]]; then
        REVERT="true"
      else
        REVERT="false"
      fi

      if [[ -n ${BREAK} ]]; then
        BREAK="true"
      else
        BREAK="false"
      fi

      insert_algolia "$(printf '{"hash": "%s", "revert": %s, "date": %s, "type": "%s", "component": "%s", "content": "%s", "breaking": %s}\n' "${HASH}" "${REVERT}" "${DATE}" "${TYPE}" "${COMPONENT}" "${CONTENT}" "${BREAK}")"
    fi
  done
}

main() {
  var_color

  if [[ "${#}" -ne 2 ]]; then
    printf "%bUsage: ${0} [END_REF] [START_REF]%b\n" "${RED}" "${RESET}"
    return 1
  fi

  if [[ $(git_is_inside) != "true" ]]; then
    printf "%bnot inside a git tree%b\n" "${YELLOW}" "${RESET}"
    return 2
  fi

  var_read ALGOLIA_API_KEY
  var_read ALGOLIA_APPLICATION_ID
  var_read ALGOLIA_INDEX "herodote"

  walk_log "${1}" "${2}"
}

main "${@}"
