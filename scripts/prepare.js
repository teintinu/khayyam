
const fs = require('fs')

const github_ref = process.env.GITHUB_REF || ''
const packageJsonFile = process.argv[2]

const matchesTag = /^refs\/tags\/v(.*)$/.exec(github_ref)

if (matchesTag && matchesTag[ 1 ]) {
  const version = matchesTag[ 1 ]
  const packageJson = JSON.parse(fs.readFileSync(packageJsonFile, 'utf8'))
  packageJson.version = version  
  fs.writeFileSync(packageJsonFile, JSON.stringify(packageJson, null, 2), 'utf8')
}
