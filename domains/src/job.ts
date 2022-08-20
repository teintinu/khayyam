import { createProgress, Progress, State } from './progress'
import p, { asap, QueuePromises, queuePromises } from 'pjobs'
import { RunningProcess, System } from './sys'

export type JobQueue = 'persistent'|'sequential'|'parallel'

export interface JobManager {
  progress: Progress
  getTree(): Node[]
  createJob({
    title, cwd, cmd, args, shell, queue, manual
  } :JobCreate): Job
  execute(): Promise<void>
  waitFor(): Promise<void>
}

export interface Job {
  readonly id: number
  readonly title: string,
  readonly progress: Progress,
  readonly dependencies: ReadonlyArray<Job>,
  readonly defer: p.Defer<void>,
  readonly queue: JobQueue
  readonly startedAt: number
  readonly stoppedAt: number
  kill(): void,
  type(text: string): void,
  listen(callback: JobCallback): ()=>void,
  depends(...dependencies: Job[]): void
  execute(): Promise<void>,
}

export interface JobCreate {
  title: string,
  cwd: string,
  cmd: string,
  args: string[],
  shell?: boolean | string | undefined,
  queue: JobQueue
  manual: boolean
}

export type JobCallbackArgs = [
  job: Job,
  chunk: string,
  error?: boolean
]
export type JobCallback= (...args:JobCallbackArgs) => void

export interface Node {
  job: Job
  dependencies: Node[]
  dependents: Node[]
}
export function createJobManager (sys: System) {
  const totalProgress = createProgress(false)
  let nextId = 1
  let allJobs: Job[] = []
  const queues: {
    // eslint-disable-next-line no-unused-vars
    [queue in JobQueue]: QueuePromises
  } = {
    persistent: queuePromises({ concurrency: 0xFFFFFFF }),
    sequential: queuePromises({ concurrency: 1 }),
    parallel: queuePromises({ concurrency: sys.concurrency })
  }

  const manager: JobManager = {
    createJob (opts) {
      return internalCreateJob(opts)
    },
    execute,
    async waitFor (): Promise<void> {
      await execute()
      await queues.persistent.waitFor()
      await queues.sequential.waitFor()
      await queues.parallel.waitFor()
    },
    get progress () {
      updateTotalProgress()
      return totalProgress
    },
    getTree
  }
  return manager

  async function execute () {
    await Promise.all(
      allJobs.map((j) => j.execute())
    )
  }

  function internalCreateJob ({
    title, cwd, cmd, args, shell, queue, manual
  }:JobCreate): Job {
    nextId = (nextId + 1) % 0xFFFFFFF
    const callbacks = new Set<JobCallback>()
    const progress = createProgress(manual)
    let _dependencies: Array<Job> = []
    const defer = p.defer<void>()
    const job: Job & {setDeps(deps: Job[]):void} = {
      id: nextId,
      title,
      progress,
      get dependencies () {
        return _dependencies
      },
      defer,
      execute,
      kill,
      type,
      listen,
      queue,
      startedAt: Date.now(),
      stoppedAt: 0,
      depends (...deps) {
        deps.forEach(dep => setDeps([..._dependencies, dep]))
      },
      setDeps
    }
    allJobs.push(job)
    let process: RunningProcess|undefined
    return job

    async function execute () {
      if (!['created', 'manual'].includes(progress.actual.state)) {
        return defer.promise
      }
      await progress.update({
        state: 'waiting'
      })
      try {
        if (_dependencies) {
          await Promise.all(_dependencies.map((dep) => {
            if (dep.progress.actual.state === 'created') {
              return dep.execute()
            }
            return dep.defer.promise
          }))
        }
        await progress.update({
          state: 'working'
        })
        const q = queues[queue]
        await q.promise(createProcess)
        await waitDependencies()
        await progress.update({
          state: 'finished'
        }).finally(() => {
          defer.resolve()
        })
      } catch (e) {
        const err = e as Error
        progress.update({
          state: 'failed',
          message: err.message || err.toString(),
          detail: err.stack?.toString()
        }).finally(() => {
          defer.reject(err)
        })
      }
      return defer.promise
    }
    function setDeps (deps: Job[]) {
      deps.forEach(checkCircular)
      _dependencies = deps
      function checkCircular (dep: Job) {
        if (dep === job) throw new Error('Circular job dependency')
        dep.dependencies.forEach(checkCircular)
      }
    }
    async function createProcess () {
      let lastOutputHandling : Promise<void>
      process = sys.createProcess({
        cwd,
        cmd,
        args,
        shell,
        handleOutput (data, error) {
          const fn = () => handleOutput(job, data, error)
          if (lastOutputHandling) lastOutputHandling = lastOutputHandling.finally(fn)
          else lastOutputHandling = fn()
          return lastOutputHandling
        }
      })
      return process.promise
    }
    function kill () {
      if (process) {
        process.kill()
        process = undefined
      }
    }
    function type (text: string) {
      if (process) {
        process.type(text)
      }
    }
    async function waitDependencies () {
      let oldDeps: Job[] = []
      return new Promise<void>((resolve) => {
        loop()
        function loop () {
          if (_dependencies.length) {
            oldDeps = _dependencies
            Promise.all(_dependencies.map((dep) => dep.defer.promise))
              .then(() => {
                if (oldDeps === _dependencies) resolve()
                else { asap(loop) }
              })
          } else resolve()
        }
      })
    }
    function listen (callback: JobCallback): ()=>void {
      callbacks.add(callback)
      return () => {
        callbacks.delete(callback)
      }
    }
    async function handleOutput (...args:JobCallbackArgs) {
      const promises: Array<Promise<void>> = []
      callbacks.forEach((fn) => {
        promises.push(new Promise((resolve) => {
          asap(() => {
            try {
              fn(...args)
            } catch (e) {
            // ignore
            }
            resolve()
          })
        }))
      })
      await Promise.all(promises)
    }
  }

  function updateTotalProgress () {
    let wsTotal = 0
    let wsDone = 0
    let wsState: State = 'idle'
    const msg: string[] = []
    allJobs.forEach((job) => {
      const actual = job.progress.actual
      const bundleState = actual.state || 'failed'
      switch (bundleState) {
        case 'manual':
          break
        case 'created':
        case 'waiting':
          wsTotal++
          if (wsState === 'idle') {
            wsState = 'working'
          }
          break

        case 'working' :
          wsTotal++
          if (actual?.percent) {
            wsDone += (actual?.percent / 100.0)
          }
          if (wsState === 'idle') {
            wsState = 'working'
          }
          msg.push(job.title)
          break

        case 'finished':
          wsTotal++
          wsDone++
          if (wsState !== 'failed') {
            wsState = 'working'
          }
          break

        case 'failed':
        case 'killed':
          wsTotal++
          wsDone++
          wsState = 'failed'
          break
      }
    })
    totalProgress.update({
      state: wsState,
      percent: Math.round(wsDone / wsTotal * 100),
      message: msg.join()
    })
  }
  function getTree (): Node[] {
    const nodes: {
      [jobId: number]: Node
    } = {}
    const tmDispose = Date.now() - 5000
    const disposables: number[] = []
    allJobs = allJobs.filter((job) => {
      const disposable = job.stoppedAt && job.stoppedAt < tmDispose
      if (disposable) {
        disposables.push(job.id)
        return false
      } else {
        const node = getNode(job)
        job.dependencies.forEach(dep => {
          const depNode = getNode(dep)
          node.dependencies.push(depNode)
          depNode.dependents.push(node)
        })
        return true
      }
    })
    return allJobs
      .map((job) => {
        (job as unknown as {dependencies: Job[]}).dependencies =
          job.dependencies.filter((dep) => !disposables.includes(dep.id))
        return nodes[job.id]
      })
      .filter((node) => node.dependents.length === 0)
    function getNode (job: Job):Node {
      return nodes[job.id] || (
        nodes[job.id] = {
          job,
          dependencies: [],
          dependents: []
        }
      )
    }
  }
}
