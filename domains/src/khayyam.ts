import { Bundler } from './bundler'
import { JobManager } from './job'
import { System } from './sys'
import { Workspace } from './workspace'

export interface KhayyamDev {
  readonly workspace: Workspace;
  readonly bundlers: Bundler[]
  readonly jobManager: JobManager
  stop(): void
}

export async function khayyamCI (
  sys: System,
  folder: string,
  jobManager:JobManager
) {
  const workspace = sys.loadWorkspace(folder)
  const build = workspace.walk('all', true, (pkg, bundler) =>
    bundler.build(pkg, jobManager, 'production')
  )
  const test = workspace.walk('all', true, (pkg, bundler) =>
    bundler.test(pkg, jobManager)
  )
  const publish = workspace.walk('all', false, (pkg, bundler) =>
    bundler.publish(pkg, jobManager)
  )
  const lint = workspace.walk('all', false, (pkg, bundler) =>
    bundler.lint(pkg, jobManager)
  )
  const measure = workspace.walk('all', false, (pkg, bundler) =>
    bundler.measure(pkg, jobManager)
  )
  test.depends('each', [build])
  publish.depends('each', [build, test, lint, measure])
  await jobManager.execute()
}

// export function khayyamDev ({
//   jobManager, workspace, watchGlobs
// }:{
//   jobManager:JobManager,
//   path: string
//   bundlers: Bundler[]
//   watchGlobs(globs: string[]): ()=>void
// }) {
//   let watching: Array<()=>void> = []
//   return {
//     stop
//   }
//   function stop () {
//     const close = watching
//     watching = []
//     close.forEach((w) => w())
//   }
//   workspace.bundlers.forEach(wsBundler => {
//     workspace.walk((pkg, bundler) => {
//       const paths = bundler.watch(pkg)
//       const w = watch(paths)
//       watchers.push(w)
//     })
//   })
// //   const queue = queuePromises({
// //     concurrency: cpus().length
// //   })
// //   const km = bootKhayyam(bundlers, queue)
// //   return {
// //     km,
// //     updateWorkspace (workspace: Workspace) {
// //       killWatchers()
// //       km.workspace.packages.forEach((pkg) => {
// //         pkg.bundlers.forEach((bundlerName) => {
// //           const bundler = findBundler(bundlerName)
// //           if (bundler) {
// //             watchers.push(bundler.watch(pkg, (fn) => fn(), queue.promise))
// //           } else {
// //             pkg.progress[bundlerName].update({
// //               state: 'failed',
// //               message: 'Invalid bundle'
// //             })
// //           }
// //         })
// //       })
// //     }
// //   }
// }

// // export function bootKhayyam (
// //   bundlers: Bundler[],
// //   jobManager: JobManager
// // ) : Khayyam {
// //   let workspace: Workspace = {
// //     packages: [],
// //     findBundler,
// //     findPackage
// //   }
// //   let watchers: Job[] = []
// //   const km = {
// //     get workspace () {
// //       return workspace
// //     },
// //     bundlers,
// //     jobs: [],
// //     updateWorkspace (ws: Workspace): void {
// //       workspace = ws
// //     },
// //     updateBundlers (b: Bundler[]): void {
// //       bundlers = b
// //     },
// //     stop () {
// //       killWatchers()
// //     }
// //   }
// //   return km

// //   function killWatchers () {
// //     const old = watchers
// //     watchers = []
// //     old.forEach((w) => w.kill())
// //   }

// // }
