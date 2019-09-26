const punycode = require('punycode');
const fs = require('fs');
const confusables = require('./confusables.json');
const tlds = require('./tlds_ascii.json');

let totalCount = 0;
const tldHomographs = {};
for (const tld of tlds {
    const unicodeDomain = punycode.toUnicode(tld); // For IDN TLDs
    let homographs = [''] // Array of Unicode strings
    for (const char of unicodeDomain) {
        const homoglyphSet = new Set([char]);
        confusables[char].forEach(({ c }) => homoglyphSet.add(c));

        const newHomographs = new Array(homographs.length * homoglyphSet.size);
        Array.from(homoglyphSet).forEach((homoglyph, i) => {
            homographs.forEach((partial, j) => {
                newHomographs[j * homoglyphSet.size + i] = partial + homoglyph;
            })
        })
        homographs = Array.from(new Set(newHomographs));
    }
    console.log(tld, homographs.length)
    totalCount += homographs.length
    fs.writeFile(
        `homographs/${tld}.json`,
        JSON.stringify(homographs, null, 2),
        'utf8',
        err => { if (err) console.log(err) });
}
console.log(totalCount)