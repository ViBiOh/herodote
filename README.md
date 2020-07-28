# herodote

Git [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) historian accross repositories.

- Full-text search with the help of Algolia as a backend
- Light frontend (<60kb gzipped) with desktop and responsive UI
- Filters on repository, type or component.
- Light shell script for loading data into backend
- Github Actions provided for integration

[![Build Status](https://travis-ci.com/ViBiOh/herodote.svg?branch=master)](https://travis-ci.com/ViBiOh/herodote)
[![codecov](https://codecov.io/gh/ViBiOh/herodote/branch/master/graph/badge.svg)](https://codecov.io/gh/ViBiOh/herodote)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ViBiOh_herodote&metric=alert_status)](https://sonarcloud.io/dashboard?id=ViBiOh_herodote)

## Concepts

Herodote aims to provide a quick view of activitiy of all your repositories. It's the changelog of your organizations.

Herodote only understands conventionnal commits. Commits that don't match expectations are ignored. To ensure you have conventionnal commits, you can use [`commit-msg hooks`](https://github.com/ViBiOh/scripts/blob/master/hooks/commit-msg) and/or a [simple Github Action that checks it](.github/workflows/branch_clean.yml). You can also have a look at [`MeilleursAgents/JudCoCo`](https://github.com/MeilleursAgents/JudCoCo).

Herodote loads data with its own script, which is idempotent. On cold start, with an empty index, it only loads last 50 commits.

## Getting started

### Algolia

Herodote use an [Algolia](https://www.algolia.com) index as a backend storage. You need to [create an application](https://www.algolia.com/account/applications) in your account.

For a personnal use, the free-tier is enough with 10k search and index by month.

### CI Integration

Herodote is fed by its own script: [herodote.sh](herodote.sh). The script automatically configure the Algolia index.

It automatically detects last commit's SHA in index and add only new ones of repository.

The script needs the following variables to be set (or will prompt you for):

- `ALGOLIA_APP`: Application ID of Algolia, can be found from the _API Keys_ section on your app's dashboard
- `ALGOLIA_KEY`: Admin API Key of Algolia, can be found from the _API Keys_ section on your app's dashboard
- `ALGOLIA_INDEX`: Index name when commits will be written (default to `herodote`)
- `GIT_HOST`: Name of your git provider (e.g. `github.com`). It's guessed from `git remote get-url --push origin` if you are in a git folder
- `GIT_REPOSITORY`: Name of your repository (e.g. `ViBiOh/herodote`). It's guessed from `git remote get-url --push origin` if you are in a git folder

If you execute your script in a non-interactive environment, set the `SCRIPTS_NO_INTERACTIVE=1` for disabling prompt.

#### Github Actions

You can use the following Github Actions for pushing your commits to Algolia index on merge to `master`.

```yaml
---
name: Herodote
on:
  push:
    branches:
      - master
jobs:
  build:
    name: Feed
    runs-on: ubuntu-latest
    steps:
      - name: Fetch history
        uses: actions/checkout@v2
        with:
          ref: ${{ github.event.pull_request.head.sha }}
          fetch-depth: 0

      - name: Push history
        run: |
          curl -q -sSL --max-time 10 "https://raw.githubusercontent.com/ViBiOh/herodote/master/herodote.sh" | bash
        env:
          ALGOLIA_APP: ${{ secrets.HERODOTE_ALGOLIA_APP }}
          ALGOLIA_KEY: ${{ secrets.HERODOTE_ALGOLIA_KEY }}
          GIT_HOST: github.com
          GIT_REPOSITORY: ${{ github.repository }}
          SCRIPTS_NO_INTERACTIVE: "1"
```

You **have to** add secrets in your repository in the repository's settings: https://github.com/YOUR_NAME/YOUR_REPOSITORY/settings/secrets

- `HERODOTE_ALGOLIA_APP`: `ALGOLIA_APP` from [#ci-integration](#ci-integration)
- `HERODOTE_ALGOLIA_KEY`: `ALGOLIA_KEY` from [#ci-integration](#ci-integration)

### Frontend

You can deploy Herodote's frontend by using the given Docker container: [vibioh/herodote](https://hub.docker.com/r/vibioh/herodote/tags?page=1&name=latest)

Your **have to** provide environment variable in order to make it work:

- `ALGOLIA_APP`: Same value as in [#ci-integration](#ci-integration)
- `ALGOLIA_INDEX`: Same value as in [#ci-integration](#ci-integration) (there is no default here, you have to provide value)
- `ALGOLIA_KEY`: Search-Only API Key of Algolia, can be found from the _API Keys_ section on your app's dashboard. **⚠️ don't provide the admin key, the variable is sent to the client, it's public! ⚠️**

For more detailed configuration of container, you can have a look at the [`ViBiOh/viws`](https://github.com/ViBiOh/viws) project.
