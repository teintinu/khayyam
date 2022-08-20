import { Bundler, Goal } from './bundler'
import { JobManager } from './job'
import { asap, defer, Defer, sleep } from 'pjobs'
import { System } from './sys'
import { ByPackage, Package, WalkFilter, Workspace } from './workspace'

export interface KhayyamDev {
  readonly workspace: Workspace;
  readonly bundlers: Bundler[]
  readonly jobManager: JobManager
  stop(): void
}

export async function khayyamCD (
  sys: System,
  folder: string,
  jobManager:JobManager
) {
  await internalWorkspaceBuild(
    await sys.loadWorkspace(folder),
    jobManager,
    'all',
    'production'
  )
}

export async function khayyamDev (
  sys: System,
  folder: string,
  jobManager:JobManager
) {
  const MAIN_DEBOUNCE = 'loadAndBuildWorkspace'
  const PKGS_DEBOUNCE = 'rebuildPackages'
  const REBUILD_SCHEDULE = 3000
  const pkgWatchers = new Set<string>()
  let changedPackages: 'all' |ByPackage<boolean> = 'all'
  let buildDefer:Defer<void> = defer<void>()
  loadAndBuildWorkspace()
  return {
    async waitBuild () {
      await buildDefer.promise
    },
    stop
  }
  function loadAndBuildWorkspace () {
    stopListeners()
    sys.schedule(MAIN_DEBOUNCE, loadAndBuildWorkspaceDebounced, REBUILD_SCHEDULE)
    function loadAndBuildWorkspaceDebounced () {
      sys.loadWorkspace(folder)
        .then(watchWsAndPackages)
      function watchWsAndPackages (workspace: Workspace) {
        sys.watch(MAIN_DEBOUNCE, [workspace.khayyamFile], () => {
          loadAndBuildWorkspace()
        })
        pkgWatchers.clear()
        workspace.packages.forEach(watchPkg)
        changedPackages = 'all'
        rebuildPackages()
        function watchPkg (pkg:Package) {
          let paths: string[] = []
          workspace.bundlers.forEach((bundler) => {
            paths = paths.concat(bundler.getPathsToWatch(pkg))
          })
          pkgWatchers.add(pkg.name)
          sys.watch(pkg.name, paths, () => {
            if (changedPackages !== 'all') {
              changedPackages[pkg.name] = true
            }
            rebuildPackages()
          })
        }
        function rebuildPackages () {
          jobManager.killAll()
          sys.schedule(PKGS_DEBOUNCE, () => {
            const changes = changedPackages
            changedPackages = {}
            internalWorkspaceBuild(
              workspace,
              jobManager,
              changes,
              'debug'
            ).then(() => {
              sys.notify('Workspace built', '')
              buildDefer.resolve()
              buildDefer = defer<void>()
            }, (reason) => {
              sys.notify('Workspace build failed', '')
              buildDefer.reject(reason)
              buildDefer = defer<void>()
            })
          }, REBUILD_SCHEDULE)
        }
      }
    }
  }
  async function stop () {
    buildDefer.reject(new Error('stopped'))
    stopListeners()
    buildDefer.promise.catch(() => undefined)
    await sleep(800)
  }
  function stopListeners () {
    jobManager.killAll()
    sys.killWatch(MAIN_DEBOUNCE)
    sys.killSchedule(MAIN_DEBOUNCE)
    sys.killSchedule(PKGS_DEBOUNCE)
    pkgWatchers.forEach(w => asap(() => sys.killWatch(w)))
    pkgWatchers.clear()
  }
}

async function internalWorkspaceBuild (
  workspace: Workspace,
  jobManager:JobManager,
  filter: WalkFilter,
  goal: Goal
) {
  const build = workspace.walk(filter, true, true, (pkg, bundler) =>
    bundler.build(pkg, jobManager, goal)
  )
  const test = workspace.walk(filter, true, true, (pkg, bundler) =>
    bundler.test(pkg, jobManager)
  )
  const publish = workspace.walk(filter, true, false, (pkg, bundler) =>
    bundler.publish(pkg, jobManager, goal)
  )
  const lint = workspace.walk(filter, true, false, (pkg, bundler) =>
    bundler.lint(pkg, jobManager)
  )
  const measure = workspace.walk(filter, true, false, (pkg, bundler) =>
    bundler.measure(pkg, jobManager)
  )
  test.depends('each', [build])
  publish.depends('each', [build, test, lint, measure])
  await jobManager.execute()
}
