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

configure_algolia() {
  curl -q -sSL --max-time 10 \
    -X PUT \
    -H "X-Algolia-Application-Id: ${ALGOLIA_APP}" \
    -H "X-Algolia-API-Key: ${ALGOLIA_KEY}" \
    --data-binary '{
      "searchableAttributes": [
        "repository",
        "hash",
        "type",
        "component",
        "content"
      ],
      "attributesForFaceting": [
        "searchable(repository)",
        "searchable(type)",
        "searchable(component)"
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
    "https://${ALGOLIA_APP}.algolia.net/1/indexes/${ALGOLIA_INDEX}/settings" >/dev/null
}

latest_commit() {
  local LATEST_HASH

  if [[ -n ${HERODOTE_API} ]]; then
    HTTP_STATUS="$(curl -q -sSL --max-time 10 \
      -o "${HTTP_OUTPUT}" \
      -w "%{http_code}" \
      --get \
      --data-urlencode "repository=${GIT_REPOSITORY}" \
      --data-urlencode "pageSize=1" \
      "${HERODOTE_API}/commits")"

    if [[ ${HTTP_STATUS} -eq 200 ]] && [[ $(jq --raw-output '.total' "${HTTP_OUTPUT}") -gt 0 ]]; then
      LATEST_HASH="$(jq --raw-output '.results[0].hash' "${HTTP_OUTPUT}")"
    fi
  else
    HTTP_STATUS="$(curl -q -sSL --max-time 10 \
      -o "${HTTP_OUTPUT}" \
      -w "%{http_code}" \
      -H "X-Algolia-Application-Id: ${ALGOLIA_APP}" \
      -H "X-Algolia-API-Key: ${ALGOLIA_KEY}" \
      --get \
      --data-urlencode "filters=repository:${GIT_REPOSITORY}" \
      --data-urlencode "hitsPerPage=1" \
      "https://${ALGOLIA_APP}-dsn.algolia.net/1/indexes/${ALGOLIA_INDEX}")"

    if [[ ${HTTP_STATUS} -eq 200 ]] && [[ $(jq --raw-output '.hits[] | length' "${HTTP_OUTPUT}") -gt 0 ]]; then
      LATEST_HASH="$(jq --raw-output '.hits[0].hash' "${HTTP_OUTPUT}")"
    fi
  fi

  if [[ ${HTTP_STATUS} -ge 400 ]]; then
    printf "%bunable to get latest commit from backend%b\n\t%bHTTP_STATUS:%b %d%b\n\t%bHTTP_OUTPUT:%b %s%b\n" "${RED}" "${RESET}" "${BLUE}" "${YELLOW}" "${HTTP_STATUS}" "${RESET}" "${BLUE}" "${YELLOW}" "$(cat "${HTTP_OUTPUT}")" "${RESET}" 1>&2
  fi

  rm "${HTTP_OUTPUT}"

  if [[ -n ${LATEST_HASH:-} ]]; then
    printf "HEAD...%s" "${LATEST_HASH}"
  else
    git rev-parse --abbrev-ref HEAD
  fi
}

insert_commit() {
  local PAYLOAD="${1:-}"

  if [[ -n ${HERODOTE_API} ]]; then
    HTTP_STATUS="$(curl -q -sSL --max-time 10 \
      -X POST \
      -o "${HTTP_OUTPUT}" \
      -w "%{http_code}" \
      -H "Authorization: ${HERODOTE_SECRET}" \
      -H "Content-Type: application/json" \
      "${HERODOTE_API}/commits" \
      -d "${PAYLOAD}")"
  else
    HTTP_STATUS="$(curl -q -sSL --max-time 10 \
      -X POST \
      -o "${HTTP_OUTPUT}" \
      -w "%{http_code}" \
      -H "X-Algolia-Application-Id: ${ALGOLIA_APP}" \
      -H "X-Algolia-API-Key: ${ALGOLIA_KEY}" \
      --data-binary "${PAYLOAD}" \
      "https://${ALGOLIA_APP}.algolia.net/1/indexes/${ALGOLIA_INDEX}")"
  fi

  if [[ ${HTTP_STATUS} -gt 299 ]]; then
    printf "%bunable to insert commit%b\n\t%bHTTP_STATUS:%b %d%b\n\t%bHTTP_OUTPUT:%b %s%b\n" "${RED}" "${RESET}" "${BLUE}" "${YELLOW}" "${HTTP_STATUS}" "${RESET}" "${BLUE}" "${YELLOW}" "$(cat "${HTTP_OUTPUT}")" "${RESET}" 1>&2
    rm "${HTTP_OUTPUT}"
    return 1
  fi

  rm "${HTTP_OUTPUT}"
}

walk_log() {
  git_conventionnal_commits

  local count=1
  IFS=$'\n'

  local SCOPES
  SCOPES="$(printf "%s|" "${!CONVENTIONAL_COMMIT_SCOPES[@]}")"
  SCOPES="${SCOPES%|}"

  shopt -s nocasematch
  for hash in $(git log --pretty=format:'%h' "$(latest_commit)"); do
    if [[ $(git show -s --format='%h %at %s' "${hash}") =~ ^([0-9a-f]{1,16})\ ([0-9]+)\ (.*)$ ]]; then
      local HASH="${BASH_REMATCH[1]}"
      local DATE="${BASH_REMATCH[2]}"
      local DESCRIPTION="${BASH_REMATCH[3]}"

      local CONTENT=""
      local TYPE=""
      local REVERT=""
      local COMPONENT=""
      local BREAK=""

      if [[ ${DESCRIPTION} =~ ^(revert )?(${SCOPES})(\((.+)\))?(\!)?:\ (.*)$ ]]; then
        REVERT="${BASH_REMATCH[1]}"
        TYPE="${BASH_REMATCH[2]}"
        COMPONENT="${BASH_REMATCH[4]}"
        BREAK="${BASH_REMATCH[5]}"
        CONTENT="${BASH_REMATCH[6]}"
      elif [[ ${DESCRIPTION} =~ Merge\ (pull\ request|branch) ]]; then
        TYPE="merge"
        CONTENT="${DESCRIPTION}"
      fi

      if [[ -z ${CONTENT:-} ]]; then
        continue
      fi

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

      count="$((count + 1))"

      local PAYLOAD
      PAYLOAD="$(
        jq -c -n \
          --arg hash "${HASH}" \
          --arg type "${TYPE}" \
          --arg component "${COMPONENT}" \
          --argjson revert "${REVERT}" \
          --argjson breaking "${BREAK}" \
          --arg content "${CONTENT}" \
          --argjson date "${DATE}" \
          --arg remote "${GIT_HOST}" \
          --arg repository "${GIT_REPOSITORY}" \
          '{
          "hash": $hash,
          "type": $type,
          "component": $component,
          "revert": $revert,
          "breaking": $breaking,
          "content": $content,
          "date": $date,
          "remote": $remote,
          "repository": $repository
        }'
      )"

      insert_commit "${PAYLOAD}"
      printf "%b%s inserted!%b\n" "${BLUE}" "${HASH}" "${RESET}"

      if [[ ${count} -gt 500 ]]; then
        printf "%bLimiting first insert to 500 commits%b\n" "${YELLOW}" "${RESET}"

        if [[ -n ${HERODOTE_API} ]]; then
          printf "%bRefreshing materialized view%b\n" "${BLUE}" "${RESET}"
          curl -q -sSL --max-time 30 -X POST -H "Authorization: ${HERODOTE_SECRET}" "${HERODOTE_API}/refresh"
        fi

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

  local HTTP_OUTPUT="http_output.txt"
  local HTTP_STATUS

  var_read GIT_HOST "$(git_remote_host)"
  var_read GIT_REPOSITORY "$(git_remote_repository)"

  var_read HERODOTE_API ""
  var_read HERODOTE_SECRET "" "secret"

  if [[ -z ${HERODOTE_API:-} ]]; then
    var_read ALGOLIA_APP ""
    var_read ALGOLIA_KEY "" "secret"
    var_read ALGOLIA_INDEX "herodote"

    configure_algolia
  fi

  walk_log
}

main "${@}"
