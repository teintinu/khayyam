import { Workspace } from './workspace'

export interface System {
    workspaceChanged(): void
    listenForWorkspaceChanges(callback:()=>void):()=>void
    getRepository(folder: string): string|undefined
    loadWorkspace(folder: string): Workspace
    notify(msg: string, packageName: string):void
}
