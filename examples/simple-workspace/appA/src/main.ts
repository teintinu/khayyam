import {functionA} from '@simple-workspace/a'

export async function main() {
  console.log('a:', functionA())
}

main().catch(console.log)
