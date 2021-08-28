
const fs = require('fs')

const github_ref = process.env.GITHUB_REF || ''
const packageJsonFile = process.argv[2]

const matchesTag = /^refs\/tags\/v(.*)$/.exec(github_ref)

if (matchesTag && matchesTag[ 1 ]) {
  const version = matchesTag[ 1 ]
  const packageJson = {
    "name": "khayyam",
    "version": version,
    "description": "A clean mono workspace automation",
    "main": "index.js",
    "scripts": {
      "postinstall": "go-npm install",
      "preuninstall": "go-npm uninstall"
    },
    "goBinary": {
      "name": "khayyam",
      "path": "./bin",
      "url": "https://github.com/teintinu/khayyam/releases/download/v{{version}}/khayyam_{{version}}_{{platform}}_{{arch}}.tar.gz"
    }    
  }  
  fs.writeFileSync(packageJsonFile, JSON.stringify(packageJson, null, 2), 'utf8')
}
