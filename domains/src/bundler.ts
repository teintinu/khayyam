import { Job, JobManager } from './job'
import { Package } from './workspace'

export interface Bundler {
    name: string
    watch(pkg: Package): string[]
    build(pkg: Package, jobManager: JobManager, goal: 'debug'|'production'): Job
    test(pkg: Package, jobManager: JobManager): Job
    serve(pkg: Package, jobManager: JobManager): Job
    lint(pkg: Package, jobManager: JobManager): Job
    publish(pkg: Package, jobManager: JobManager): Job
}
