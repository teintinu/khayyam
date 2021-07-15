import { TodoRepository } from "./todoRepository";
import { TodoEntity } from './todoEntity'
import { TodoUseCase } from "./todoUseCase";

export class AddTodoUseCase extends TodoUseCase {
    eventId() {
        return "AddTodoUseCase"
    }
    run(description: string) {
        const newTodo: TodoEntity = {
            id: this.todoRepository.uuid(),
            description,
            checked: false
        }
        this.todoRepository.save(newTodo)
    }
}
