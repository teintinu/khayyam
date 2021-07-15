import { getValues } from "./internal"

it('internal', async() => {
    expect(await getValues()).toEqual([123, 456, 789])
})
