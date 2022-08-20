// import { Bundler, ProcessParams, RunningProcess, System } from '@khayyam/domains'
// import { spawn, ChildProcessWithoutNullStreams } from 'child_process'
// import { Readable, Writable } from 'stream'
// import { cpus } from 'os'
// import { asap } from 'pjobs'

import { System } from '@khayyam/domains'

export function systemAdapterFactory (): System {
  return null as any
  // const fsListeners = new Set<Unscribe>()
  // let tmChanged: any
  // const sys: System = {
  //   concurrency: cpus().length,
  //   workspaceChanged: fsChanged,
  //   listenForWorkspaceChanges (callback) {
  //     fsListeners.add(callback)
  //     return () => {
  //       fsListeners.delete(callback)
  //     }
  //   },
  //   getRepository: fsGetRepository,
  //   loadWorkspace,
  //   notify: fsNotify,
  //   createProcess: fsCreateProcess
  // }
  // return sys
  // function fsChanged (): void {
  //   clearTimeout(tmChanged)
  //   tmChanged = setTimeout(() => {
  //     fsListeners.forEach((fn) => asap(fn))
  //   }, 1000)
  // }
  // function fsGetRepository (folder: string): string|undefined {
  //   return 'repo:' + folder
  // }
  // function loadWorkspace (folder:string) {
  //   const [sBundler, sPkg, sDeps, sLayer, sInvalid] = folder.split('/')
  //   let bundler: Bundler
  //   if (sBundler === 'npm') bundler = createFakeBundlerNPM(fakeLogger)
  //   else throw new Error('fake-fs-invalid-bundle')
  //   return createFakeWorkspace({
  //     logger: fakeLogger,
  //     bundler,
  //     packages: sPkg as FakePackages,
  //     deps: sDeps === 'deps',
  //     layers: sLayer === 'layers',
  //     invalid: sInvalid === 'invalid'
  //   })
  // }
  // function fsNotify (msg: string, pkgName: string) {
  //   fakeLogger.log('fsNotify:' + msg + ' on package: ' + pkgName)
  // }
  // function fsCreateProcess ({ cmd, args, handleOutput }: ProcessParams): RunningProcess {
  //   let childProcess:ChildProcessWithoutNullStreams|undefined
  //   let sendStream : Writable|undefined
  //   let receiveStream: Readable|undefined
  //   let exitCode = 0
  //   const promise = new Promise<void>((resolve, reject) => {
  //     const child = spawn(cmd, args, {
  //       cwd,
  //       shell
  //     })
  //     child.stdout.on('data', (chunk) => {
  //       const data = typeof chunk === 'string' ? chunk : chunk.toString('utf-8')
  //       handleOutput(data, false)
  //     })
  //     child.stderr.on('data', function (chunk) {
  //       const data = typeof chunk === 'string' ? chunk : chunk.toString('utf-8')
  //       handleOutput(data, false)
  //     })
  //     sendStream = child.stdin

  //     childProcess = child
  //     child.on('error', function (err) {
  //       reject(err)
  //     })
  //     child.on('close', async (code) => {
  //       await sleep(100)
  //       exitCode = code || 0
  //       asap(p.kill)
  //       if (exitCode) {
  //         const err = new Error('error level=' + exitCode)
  //         handleOutput(err.message, true).finally(() => {
  //           reject(err)
  //         })
  //       } else {
  //         resolve()
  //       }
  //     })
  //   })
  //   const p:RunningProcess = {
  //     promise,
  //     kill () {
  //       if (childProcess) {
  //         try {
  //           childProcess.kill()
  //         } catch (e) {
  //         //
  //         }
  //       }
  //       childProcess = undefined
  //       sendStream = undefined
  //       if (receiveStream) {
  //         receiveStream.destroy()
  //         receiveStream = undefined
  //       }
  //     },
  //     type (text) {
  //       if (sendStream) {
  //         sendStream.write(text)
  //       }
  //     }
  //   }
  //   return p
  // }
}
