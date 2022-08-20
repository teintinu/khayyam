import { Workspace } from './workspace'

export type Unscribe=()=>void
export interface System {
    readonly concurrency: number
    getRepository(folder: string): string|undefined
    loadWorkspace(folder: string): Workspace
    notify(msg: string, packageName: string):void
    createProcess(process: ProcessParams): RunningProcess
    watch(globs: string[], callback:()=>void): Unscribe
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
