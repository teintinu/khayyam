#!/usr/bin/env node

const readline = require('readline')

let argIdx = 2
let resp = ''
const wait = setTimeout(() => 0, 0xFFFF)

let afterRead = () => 0
const args = process.argv
const delay = parseInt(args[argIdx++], 10)

console.log('started')

readline.createInterface({
  input: process.stdin,
  output: process.stdout,
  terminal: false
}).on('line', function (line) {
  resp = line
  if (wait) {
    clearTimeout(wait)
    afterRead()
    nextArg()
  }
})

runArgs()

function runArgs () {
  let arg = args[argIdx++]
  if (arg === 'READ') {
    arg = args[argIdx++]
    console.log('INPUT: ' + arg)
    afterRead = () => {
      console.log(arg + ': ' + resp)
    }
  } else {
    console.log(arg)
    nextArg()
  }
}

function nextArg () {
  setTimeout(() => {
    if (argIdx < args.length) {
      runArgs()
    } else {
      console.log('finished')
      process.exit()
    }
  }, delay)
}
