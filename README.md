# asciinema-adjuster
A little Go lib for optimising .cast files produced by asciinema with new timing and prompt highlighting.

See <https://docs.asciinema.org/getting-started/> for instructions for recording and playing a `.cast` file.

Then:

```
$ go install github.com/astromechza/asciinema-adjuster
$ asciinema-adjuster example.cast '$ '
```

Adjust `$ ` to be the suffix of your prompt statement.
