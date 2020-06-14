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

git_remote_repository() {
  if [[ $(git_is_inside) != "true" ]]; then
    return
  fi

  local REMOTE_URL
  REMOTE_URL="$(git remote get-url --push "$(git remote show | head -1)")"

  if [[ ${REMOTE_URL} =~ ^.*@.*:(.*)\/(.*).git$ ]]; then
    printf "%s/%s" "${BASH_REMATCH[1]}" "${BASH_REMATCH[2]}"
  fi
}

git_remote_host() {
  if [[ $(git_is_inside) != "true" ]]; then
    return
  fi

  local REMOTE_URL
  REMOTE_URL="$(git remote get-url --push "$(git remote show | head -1)")"

  if [[ ${REMOTE_URL} =~ ^.*@(.*):.*\/.*.git$ ]]; then
    printf "%s" "${BASH_REMATCH[1]}"
  fi
}

algolia_index() {
  curl -X PUT \
     -H "X-Algolia-Application-Id: ${ALGOLIA_APPLICATION_ID}" \
     -H "X-Algolia-API-Key: ${ALGOLIA_API_KEY}" \
     --data-binary '{
      "searchableAttributes": [
        "repository",
        "hash",
        "type",
        "component",
        "content"
      ],
      "ranking": [
        "desc(date)",
        "typo",
        "geo",
        "words",
        "filters",
        "proximity",
        "attribute",
        "exact",
        "custom"
      ],
      "exactOnSingleWordQuery": "word"
     }' \
    "https://${ALGOLIA_APPLICATION_ID}.algolia.net/1/indexes/${ALGOLIA_INDEX}/settings" > /dev/null
}

algolia_latest() {
  local HTTP_OUTPUT="http_output.txt"
  local HTTP_STATUS

  HTTP_STATUS="$(curl -q -sSL \
    --max-time 10 \
    -o "${HTTP_OUTPUT}" \
    -w "%{http_code}" \
    -H "X-Algolia-Application-Id: ${ALGOLIA_APPLICATION_ID}" \
    -H "X-Algolia-API-Key: ${ALGOLIA_API_KEY}" \
    --get \
    --data-urlencode "query=${REPOSITORY}" \
    --data-urlencode "hitsPerPage=1" \
    "https://${ALGOLIA_APPLICATION_ID}-dsn.algolia.net/1/indexes/${ALGOLIA_INDEX}")"

  if [[ ${HTTP_STATUS} == "200" ]] && [[ $(python -c "import json; print(len(json.load(open('${HTTP_OUTPUT}'))['hits']))") -gt 0 ]]; then
    printf "HEAD...%s" "$(python -c "import json; print(json.load(open('${HTTP_OUTPUT}'))['hits'][0]['hash'])")"
    rm "${HTTP_OUTPUT}"
    return
  fi

  rm "${HTTP_OUTPUT}"
  git rev-parse --abbrev-ref HEAD
}

algolia_insert() {
  curl -X POST \
     -H "X-Algolia-Application-Id: ${ALGOLIA_APPLICATION_ID}" \
     -H "X-Algolia-API-Key: ${ALGOLIA_API_KEY}" \
     --data-binary "${1}" \
    "https://${ALGOLIA_APPLICATION_ID}.algolia.net/1/indexes/${ALGOLIA_INDEX}" > /dev/null
}

walk_log() {
  git_conventionnal_commits

  local REPOSITORY
  REPOSITORY="$(git_remote_repository)"

  local GIT_HOST
  GIT_HOST="$(git_remote_host)"

  local count=1
  IFS=$'\n'

  shopt -s nocasematch
  for hash in $(git log --no-merges --pretty=format:'%h' "$(algolia_latest)"); do
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

      count="$(( count + 1 ))"
      algolia_insert "$(printf '{"remote": "%s", "repository": "%s", "hash": "%s", "revert": %s, "date": %s, "type": "%s", "component": "%s", "content": "%s", "breaking": %s}\n' "${GIT_HOST}" "${REPOSITORY}" "${HASH}" "${REVERT}" "${DATE}" "${TYPE}" "${COMPONENT}" "${CONTENT}" "${BREAK}")"

      if [[ ${count} -gt 50 ]]; then
        printf "%bLimiting first insert to 50 commits%b\n" "${YELLOW}" "${RESET}"

        break
      fi
    fi
  done
}

main() {
  var_color

  if [[ $(git_is_inside) != "true" ]]; then
    printf "%bnot inside a git tree%b\n" "${YELLOW}" "${RESET}"
    return 2
  fi

  var_read ALGOLIA_APPLICATION_ID
  var_read ALGOLIA_API_KEY "" "secret"
  var_read ALGOLIA_INDEX "herodote"

  algolia_index

  walk_log
}

main "${@}"
