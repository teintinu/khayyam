import { TodoEntity } from "./todoEntity";

export interface TodoRepository {
    add(description: string): Promise<void>
    changeDescription(id: string, description: string):Promise<void>
    changeFlag(id: string, flag: boolean):Promise<void>
    delete(id: string):Promise<void>
}
