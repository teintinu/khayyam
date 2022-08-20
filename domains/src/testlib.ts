import { createJobManager, JobCallbackArgs, JobManager, Job, Node } from './job'
import { asap, defer } from 'pjobs'
import { Bundler } from './bundler'
import { ByPackage, createWorkspace, Package, WalkedJobs } from './workspace'
import { createProgress } from './progress'
import { System } from './sys'

export const cwd = process.cwd()
const fakeProcess = './scripts/fake-process.js'
const shell = true

export type TreeOfNode = [
  job: string,
  dependencies: TreeOfNode[],
]

export type Logger=ReturnType<typeof createFakeLog>

export function createFakeLog (verbose = false) {
  const logged: string[] = []
  const aOk = defer<void>()
  const bOk = defer<void>()
  const cOk = defer<void>()
  const handleOutput = jest.fn<void, JobCallbackArgs>()
    .mockImplementation(async (job, chunk, error) => {
      const lines = chunk.trim().split('\n')
      lines.forEach((s) => {
        if (s === 'a' || s === 'INPUT: a') setTimeout(aOk.resolve, 100)
        if (s === 'b' || s === 'INPUT: b') setTimeout(bOk.resolve, 100)
        if (s === 'c' || s === 'INPUT: c') setTimeout(cOk.resolve, 100)
        job.progress.update({
          state: 'working',
          message: (error ? ' (error): ' : '') + s
        })
      })
    })
  const jobManager = createJobManager()
  const fsListeners = new Set<()=>void>()
  let tmChanged: any
  const fakeSys: System & { simulateChange():void} = {
    workspaceChanged: fsChanged,
    listenForWorkspaceChanges (callback) {
      fsListeners.add(callback)
      return () => {
        fsListeners.delete(callback)
      }
    },
    getRepository: fsGetRepository,
    simulateChange: fsSimulateChange,
    loadWorkspace,
    notify: fsNotify
  }
  const fakeLogger = {
    jobManager,
    handleOutput,
    aOk,
    bOk,
    cOk,
    fakeSys,
    verbose,
    log (...args: string[]) {
      logged.push(args.join(' '))
      if (verbose) { console.info(args.join(' ')) }
    },
    logged,
    tree () {
      const done: {[id: number]:boolean} = {}
      const tree = jobManager.getTree()
      return tree.map(treeOfNode)
      function treeOfNode (node: Node): TreeOfNode {
        const id = node.job.id
        if (done[id]) {
          return [
            node.job.title + '(rec)',
            []
          ]
        }
        done[id] = true
        const n: TreeOfNode = [
          node.job.title + '=>[' +
          node.dependents.map(d => d.job.title).join() + ']',
          node.dependencies.map(treeOfNode)
        ]
        done[id] = false
        return n
      }
    }
  }
  return fakeLogger
  function fsChanged (): void {
    clearTimeout(tmChanged)
    tmChanged = setTimeout(() => {
      fsListeners.forEach((fn) => asap(fn))
    }, 1000)
  }
  function fsGetRepository (folder: string): string|undefined {
    return 'repo:' + folder
  }
  function fsSimulateChange ():void {
    fsChanged()
  }
  function loadWorkspace (folder:string) {
    const [sBundler, sPkg, sDeps, sLayer, sInvalid] = folder.split('/')
    let bundler:Bundler
    if (sBundler === 'npm') bundler = createFakeBundlerNPM(fakeLogger)
    else throw new Error('fake-fs-invalid-bundle')
    return createFakeWorkspace({
      logger: fakeLogger,
      bundler,
      packages: sPkg as FakePackages,
      deps: sDeps === 'deps',
      layers: sLayer === 'layers',
      invalid: sInvalid === 'invalid'
    })
  }
  function fsNotify (msg: string, pkgName: string) {
    fakeLogger.log('fsNotify:' + msg + ' on package: ' + pkgName)
  }
}

export function createFakeProgress (logger: Logger, title: string, manual: boolean) {
  const pg = createProgress(manual)
  pg.on((status) => {
    logger.log(title + ' ' + status.state + ' ' + status.message)
  })
  return pg
}

export function createFakeJob (logger: Logger, delay: number, title: string, args: string[], manual: boolean = false) {
  logger.log('creating job [' + title + ']')
  const job = logger.jobManager.createJob({
    title,
    cwd,
    command: 'node',
    args: [fakeProcess, String(delay), ...args],
    shell,
    queue: 'sequential',
    manual
  })
  job.progress.on((status) => {
    logger.log('pg ' + title + ': ' + status.state + ' ' + status.message)
  })
  job.listen(logger.handleOutput)
  return job
}

export function create2FakeJob (logger: Logger, delay: number, mainTitle: string, mainArgs: string[], depTitle: string, depArgs: string[]) {
  logger.log('creating main job [' + mainTitle + ']')
  const main = logger.jobManager.createJob({
    title: mainTitle,
    cwd,
    command: 'node',
    args: [fakeProcess, String(delay), ...mainArgs],
    shell,
    queue: 'persistent',
    manual: false
  })
  logger.log('creating dep job [' + depTitle + ']')
  const dep = logger.jobManager.createJob({
    title: depTitle,
    cwd,
    command: 'node',
    args: [fakeProcess, String(delay), ...depArgs],
    shell,
    queue: 'persistent',
    manual: false
  })
  main.depends(dep)
  main.progress.on((status) => {
    logger.log('pg-main ' + status.state + ' ' + status.message)
  })
  dep.progress.on((status) => {
    logger.log('pg-dep ' + status.state + ' ' + status.message)
  })
  main.listen(logger.handleOutput)
  dep.listen(logger.handleOutput)
  return {
    main, dep
  }
}

function createFakeBundlerNPM (logger:Logger) {
  const fakeBundler: Bundler = {
    name: 'fakeBundlerNPM',
    watch (pkg) {
      return ['fakeBundlerNPM']
    },
    build (pkg, jobManager, goal) {
      return createFakeBundlerCommand(pkg, jobManager, 'npm-build-' + goal)
    },
    test (pkg, jobManager) {
      return createFakeBundlerCommand(pkg, jobManager, 'npm-test')
    },
    serve (pkg, jobManager) {
      return createFakeBundlerCommand(pkg, jobManager, 'npm-serve')
    },
    lint (pkg, jobManager) {
      return createFakeBundlerCommand(pkg, jobManager, 'npm-lint')
    },
    publish (pkg, jobManager) {
      return createFakeBundlerCommand(pkg, jobManager, 'npm-publish')
    }
  }
  return fakeBundler
  function createFakeBundlerCommand (pkg: Package, cmdJobManager: JobManager, cmd: string) {
    expect(cmdJobManager).toBe(logger.jobManager)
    return createFakeJob(
      logger,
      100,
      pkg.name + '(' + cmd + ')',
      [
        cmd
      ]
    )
  }
}

type FakePackages='x'|'xy'|'abcde'

function createFakeWorkspace (
  { logger, bundler, packages, deps, layers, invalid }:
  { logger: Logger, bundler: Bundler, packages:FakePackages, deps: boolean; layers: boolean; invalid: boolean }) {
  const pkgs:{
      [n:string]:Package
    } = {}
  for (const n of packages) {
    const pkg: Readonly<Package > = {
      name: n,
      layer: layers ? n : '',
      dependencies: [],
      bundlers: [bundler.name],
      folder: n
    }
    pkgs[n] = pkg
  }
  if (deps) {
    if (packages === 'xy') {
      pkgs.x = {
        ...pkgs.x,
        dependencies: invalid ? ['z'] : ['y']
      }
    }
    if (packages === 'abcde') {
      pkgs.a = {
        ...pkgs.a,
        dependencies: invalid ? ['b', 'c', 'z'] : ['b', 'c']
      }
      pkgs.b = {
        ...pkgs.b,
        dependencies: ['c', 'd']
      }
      pkgs.c = {
        ...pkgs.c,
        dependencies: ['d']
      }
      pkgs.d = {
        ...pkgs.d,
        dependencies: ['e']
      }
    }
  }
  const ws = createWorkspace({
    sys: logger.fakeSys,
    layers: [],
    packages: Object.keys(pkgs).map((n) => pkgs[n]),
    bundlers: [bundler]
  })
  return ws
}

export function jobDebug (job: Job) {
  return job.title + '(' +
    [...job.dependencies]
      .sort((a, b) => a.title.localeCompare(b.title))
      .map(dep => dep.title) +
    ')'
}

export function walkedDebug (walkJobs: WalkedJobs) {
  const ret: ByPackage<unknown> = {}
  const pkgs = Object.keys(walkJobs.jobs).sort()
  pkgs.forEach((pkg) => {
    ret[pkg] = ([...walkJobs.jobs[pkg]] || [])
      .sort((a, b) => a.title.localeCompare(b.title))
      .map(jobDebug)
  })
  return ret
}
