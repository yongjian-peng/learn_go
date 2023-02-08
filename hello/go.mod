module example.com/hello

go 1.17

replace example.com/greetings => ../greetings

require (
	example.com/morestrings v0.0.0-00010101000000-000000000000
	github.com/spf13/cast v1.5.0
)

replace example.com/morestrings => ./morestrings
