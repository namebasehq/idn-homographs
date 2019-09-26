package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"unicode/utf8"

	"golang.org/x/net/idna"
)

type homoglyphMap map[rune]map[rune]struct{}

func (hm *homoglyphMap) UnmarshalJSON(data []byte) error {
	var confusables map[string][]struct{ C string }
	if err := json.Unmarshal(data, &confusables); err != nil {
		return err
	}

	homoglyphs := make(map[rune]map[rune]struct{})
	for k, confusableArray := range confusables {
		key, _ := utf8.DecodeRuneInString(k)
		homoglyphs[key] = make(map[rune]struct{})
		homoglyphs[key][key] = struct{}{}
		for _, confusable := range confusableArray {
			r, _ := utf8.DecodeRuneInString(confusable.C)
			homoglyphs[key][r] = struct{}{}
		}
	}
	*hm = homoglyphs
	return nil
}

var (
	confusables homoglyphMap
	tlds        []string
)

func init() {
	f, err := os.Open("confusables.json")
	if err != nil {
		log.Fatal(err)
	}

	if err := json.NewDecoder(f).Decode(&confusables); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	f, err = os.Open("tlds_ascii.json")
	if err != nil {
		log.Fatal(err)
	}

	if err := json.NewDecoder(f).Decode(&tlds); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func recursiveHomographs(output io.Writer, canonical, homograph string) int {
	if len(canonical) == 0 {
		punycode, err := idna.ToASCII(homograph)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := fmt.Fprintln(output, homograph, punycode); err != nil {
			log.Fatal(err)
		}
		return 1
	}
	r, w := utf8.DecodeRuneInString(canonical)
	total := 0
	for homoglyph := range confusables[r] {
		total += recursiveHomographs(output, canonical[w:], homograph+string(homoglyph))
	}
	return total
}

type bufferedWriteCloser struct {
	io.Closer
	*bufio.Writer
}

func (b bufferedWriteCloser) Close() error {
	if err := b.Flush(); err != nil {
		return err
	}
	return b.Closer.Close()
}

func openOutput(tld string) io.WriteCloser {
	flags := os.O_CREATE | os.O_TRUNC | os.O_APPEND | os.O_WRONLY
	output, err := os.OpenFile(fmt.Sprintf("homographs/%s.txt", tld), flags, 0400)
	if err != nil {
		log.Fatal(err)
	}
	return bufferedWriteCloser{output, bufio.NewWriter(output)}
}

func main() {
	numHomographs := 0
	for _, tld := range tlds {
		canonical, err := idna.ToUnicode(tld)
		if err != nil {
			log.Fatal(err)
		}
		output := openOutput(tld)
		count := recursiveHomographs(output, canonical, "")
		fmt.Println(tld, count)
		numHomographs += count
		if err := output.Close(); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(numHomographs)
}
