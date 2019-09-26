# idn-homographs

Tool to generate all IDN homographs of a set of names (TLDs) 

The core list [confusables.json](./confusables.json) is borrowed from https://github.com/vhf/confusable_homoglyphs and compiled into a more human-readable format.

[main.go](./main.go) enumerates all 3,096,305,646 homographs of ICANN TLDs.

[index.js](./index.js) is a sample implementation of how to detect a homograph in JavaScript.
