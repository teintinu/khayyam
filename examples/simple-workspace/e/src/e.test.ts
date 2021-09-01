import { functionE } from './'

it('module e', async () => {
  await new Promise((resolve) => setTimeout(resolve, 300))
  expect(functionE()).toEqual({ e: 'e' })
})
