import { Job, JobManager } from './job'
import { Package } from './workspace'

export interface Bundler {
    name: string
    getPathsToWatch(pkg: Package): string[]
    build(pkg: Package, jobManager: JobManager, goal: Goal): Job
    test(pkg: Package, jobManager: JobManager): Job
    publish(pkg: Package, jobManager: JobManager, goal: Goal): Job
    lint(pkg: Package, jobManager: JobManager): Job
    measure(pkg: Package, jobManager: JobManager): Job
}

export type Goal = 'debug'|'production'
