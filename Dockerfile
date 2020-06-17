FROM vibioh/viws

ENV VIWS_CSP "default-src 'self'; base-uri 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; connect-src 'self' *.algolia.net *.algolianet.com"
ENV VIWS_ENV ALGOLIA_APP,ALGOLIA_KEY,ALGOLIA_INDEX
ENV VIWS_HEADERS = X-UA-Compatible:ie=edge
ENV VIWS_PORT 1080
ENV VIWS_SPA true

ARG VERSION
ENV VERSION=${VERSION}

COPY build/ /www/
