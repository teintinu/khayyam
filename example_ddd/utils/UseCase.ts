export abstract class UseCase {
    abstract eventId(): string
    abstract run(): Promise<void>    
}