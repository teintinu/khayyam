
const fs = require("fs");
const stdinBuffer = fs.readFileSync(0, 'utf-8');
const versions = stdinBuffer
  .split('\n')
  .filter((l) => l.startsWith('v'))
  .map((l) => l.substr(1).split('.').map(n=>parseInt(n,10)))
  .sort((a, b) => {
    for (let i = 0; i <= 2; i++) {
      const d = b[ i ] - a[ i ]
      if (d !== 0) return d
    }
    return 0
  })
const last = versions[ 0 ]
const next = [ last[ 0 ], last[ 1 ], last[ 2 ] + 1 ]
console.log("v" + next.join('.'));