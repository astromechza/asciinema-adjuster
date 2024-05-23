package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const charInterval = time.Millisecond * 50
const lineInterval = time.Second

func jitter(i time.Duration, p float64) time.Duration {
	return i + time.Duration(rand.Float64()*p*float64(i))
}

func main() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), `Usage of %s: [options] <input file> <prompt pattern>
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
	if args := flag.Args(); len(args) != 2 {
		log.Fatalf("expected 2 positional args, got %d", len(args))
	} else if raw, err := os.ReadFile(args[0]); err != nil {
		log.Fatalf("failed to read input file '%s': %v", args[0], err)
	} else {
		rr := bytes.NewReader(raw)
		dec := json.NewDecoder(rr)
		out := new(bytes.Buffer)

		var header map[string]interface{}
		if err := dec.Decode(&header); err != nil {
			log.Fatalf("failed to read header line: %v", err)
		}
		enc := json.NewEncoder(out)
		_ = enc.Encode(&header)
		var lc int
		var total time.Duration
		var last float64
		var inline bool

		for {
			line := [3]interface{}{}
			if err := dec.Decode(&line); err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				log.Fatalf("failed to read char line %d: %v", lc, err)
			}
			tOffset := line[0].(float64)
			elapsed := time.Duration((tOffset - last) * float64(time.Second))
			last = tOffset

			if inline {
				elapsed = jitter(charInterval, 2)
				log.Print(elapsed)
			}

			segment := line[2].(string)
			if strings.HasSuffix(segment, "\r\n") {
				inline = false
			} else if strings.HasSuffix(segment, args[1]) {
				elapsed = lineInterval
				inline = true
				line[2] = "\u001b[2K\u001b[36m" + segment + "\u001b[0m"
			}

			total += elapsed
			line[0] = float64(total) / float64(time.Second)
			_ = enc.Encode(&line)
			lc += 1
		}
		log.Printf("read %d lines", lc)
		_, _ = fmt.Fprintf(os.Stdout, out.String())
	}
}
