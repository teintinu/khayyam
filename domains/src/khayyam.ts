import { Bundler, Goal } from './bundler'
import { JobManager } from './job'
import { asap } from 'pjobs'
import { System, Unscribe } from './sys'
import { ByPackage, WalkFilter, Workspace } from './workspace'

export interface KhayyamDev {
  readonly workspace: Workspace;
  readonly bundlers: Bundler[]
  readonly jobManager: JobManager
  stop(): void
}

export async function khayyamCD (
  workspace: Workspace,
  jobManager:JobManager
) {
  await internalWorkspaceBuild(
    workspace,
    jobManager,
    'all',
    'production'
  )
}

export function khayyamDev ({
  sys, folder, jobManager
}:{
  sys: System,
  folder: string,
  jobManager:JobManager
}) {
  let watchers: Unscribe[] = []
  let workspace: Workspace
  let debounceBuild: any
  const changedPackages: ByPackage<boolean> = {}
  const khayyamFile = folder + '/khayyam.yaml'
  const stopMainWatch = sys.watch([khayyamFile], reloadWorkspace)
  return {
    stop () {
      stopPackagesWatchers()
      stopMainWatch()
    }
  }
  function reloadWorkspace () {
    workspace = sys.loadWorkspace(folder)
    workspace.packages.forEach((pkg) => {
      let paths: string[] = []
      workspace.bundlers.forEach((bundler) => {
        paths = paths.concat(bundler.getPathsToWatch(pkg))
      })
      sys.watch(paths, () => {
        changedPackages[pkg.name] = true
        rebuildWorkspace()
      })
    })
  }
  function rebuildWorkspace () {
    jobManager.killAll()
    if (debounceBuild) clearTimeout(debounceBuild)
    debounceBuild = setTimeout(() => {
      debounceBuild = undefined
      internalWorkspaceBuild(
        workspace,
        jobManager,
        changedPackages,
        'debug'
      )
    }, 2000)
  }
  function stopPackagesWatchers () {
    watchers.forEach(w => asap(w))
    watchers = []
  }

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
}

async function internalWorkspaceBuild (
  workspace: Workspace,
  jobManager:JobManager,
  filter: WalkFilter,
  goal: Goal
) {
  const build = workspace.walk(filter, true, (pkg, bundler) =>
    bundler.build(pkg, jobManager, goal)
  )
  const test = workspace.walk(filter, true, (pkg, bundler) =>
    bundler.test(pkg, jobManager)
  )
  const publish = workspace.walk(filter, false, (pkg, bundler) =>
    bundler.publish(pkg, jobManager, goal)
  )
  const lint = workspace.walk(filter, false, (pkg, bundler) =>
    bundler.lint(pkg, jobManager)
  )
  const measure = workspace.walk(filter, false, (pkg, bundler) =>
    bundler.measure(pkg, jobManager)
  )
  test.depends('each', [build])
  publish.depends('each', [build, test, lint, measure])
  await jobManager.execute()
}
