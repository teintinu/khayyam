
import { doSomething } from './index'

it('test something', () => {
  expect('something').toEqual(doSomething())
})
