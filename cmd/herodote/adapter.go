package main

import "github.com/ViBiOh/herodote/pkg/adapter"

type adapters struct {
	adapter adapter.App
}

func newAdapters(client clients) adapters {
	return adapters{
		adapter: adapter.New(client.redis, client.database),
	}
}
