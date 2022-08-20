import { Workspace } from './workspace'

export interface System {
    readonly concurrency: number
    workspaceChanged(): void
    listenForWorkspaceChanges(callback:()=>void):()=>void
    getRepository(folder: string): string|undefined
    loadWorkspace(folder: string): Workspace
    notify(msg: string, packageName: string):void
    createProcess(process: ProcessParams): RunningProcess
}

export interface ProcessParams {
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
