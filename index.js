const punycode = require('punycode');
const homoglyphs = require('./homoglyphs.json');
const tlds = require('./tlds_ascii.json');

const matchStrs = []
for (const tld of tlds) {
    const unicode = punycode.toUnicode(tld);
    let matchStr = ''
    for (const char of unicode) {
        const matches = homoglyphs[char] ? homoglyphs[char].join("|") : char
        matchStr += `(${matches})`
    }
    matchStrs.push(`(${matchStr})`)
}

const homographRegexp = new RegExp(`^${matchStrs.join('|')}$`, 'u')

function isTldHomograph(name) {
    return homographRegexp.test(punycode.toUnicode(name))
}

const testHomograph = name => {
    console.log(name, isTldHomograph(name))
}

testHomograph('')
testHomograph('asdfkasldflaksjdflaksd')
testHomograph('xn--kpu38a')
testHomograph('con')
testHomograph('corn')
testHomograph('xn-o-eka')