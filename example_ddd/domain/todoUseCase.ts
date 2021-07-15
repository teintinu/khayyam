import { UseCase } from "@hexagonal-todo/utils";
import { TodoRepository } from "./todoRepository";

export abstract class TodoUseCase extends UseCase {
    constructor(protected todoRepository: TodoRepository) {
        super()
    }
}