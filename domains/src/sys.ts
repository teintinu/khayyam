import { PackageState, Workspace } from './workspace'

export interface System {
    readonly concurrency: number
    getRepository(folder: string): string|undefined
    loadWorkspace(folder: string): Promise<Workspace>
    notify(msg: string, packageName?: string):void
    createProcess(process: ProcessParams): RunningProcess
    watch(id:string, globs: string[], callback:(args:{id: string, path: string})=>void): void
    killWatch(id:string):void
    schedule(id: string, callback:(args:{id: string})=>void, ms: number): void
    killSchedule(id:string):void
    packageState(pkgName:string): Promise<PackageState>
    updatePackageState<P extends keyof PackageState>(
        pkgName:string,
        prop: P,
        value: PackageState[P]
    ): Promise<void>
}

export interface ProcessParams {
    title: string
    cwd: string
    cmd: string
    args: string[]
    shell?: string|boolean,
    handleOutput(data: string, error: boolean): Promise<void>
}

export interface RunningProcess {
    readonly promise: Promise<void>
    type(text: string):void
    kill(): void
}
