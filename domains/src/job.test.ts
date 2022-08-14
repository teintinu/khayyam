import { create2FakeJob, createFakeJob, createFakeLog } from './testlib'

describe('job', () => {
  const delay = 150
  it('without dependencies', async () => {
    const logger = createFakeLog()
    const job = createFakeJob(logger, delay, 'job', ['a'])
    expect(job.title).toBe('job')
    expect(job.progress.actual.state).toBe('created')
    expect(logger.handleOutput).not.toBeCalled()
    expect(logger.logged?.length).toBe(0)
    const p = job.execute()
    await logger.aOk.promise
    expect(job.progress.actual.state).toBe('working')
    expect(logger.logged).toMatchSnapshot('01:log-working')
    expect(logger.tree()).toMatchSnapshot('02:tree-working')
    expect(job.progress.actual.state).toBe('working')
    expect(job.progress.actual.message).toBe('a')
    expect(await p).toBeUndefined()
    expect(job.progress.actual.state).toBe('finished')
    expect(logger.logged).toMatchSnapshot('03:log-finished')
    expect(logger.tree()).toMatchSnapshot('04:tree-finished')
  })
  it('multiline', async () => {
    const logger = createFakeLog()
    const job = createFakeJob(logger, delay, 'multi', ['a', 'b'])
    expect(job.title).toBe('multi')
    expect(job.progress.actual.state).toBe('created')
    expect(logger.handleOutput).not.toBeCalled()
    expect(logger.logged?.length).toBe(0)
    const p = job.execute()
    await logger.aOk.promise
    expect(logger.logged).toMatchSnapshot('01:log-aOk')
    expect(logger.tree()).toMatchSnapshot('02:tree-aOk')
    expect(job.progress.actual.state).toBe('working')
    expect(job.progress.actual.message).toBe('a')
    await logger.bOk.promise
    expect(logger.logged).toMatchSnapshot('03:log-bOk')
    expect(logger.tree()).toMatchSnapshot('04:tree-bOk')
    expect(job.progress.actual.state).toBe('working')
    expect(job.progress.actual.message).toBe('b')
    expect(await p).toBeUndefined()
    expect(job.progress.actual.state).toBe('finished')
    expect(logger.logged).toMatchSnapshot('05:log-finished')
    expect(logger.tree()).toMatchSnapshot('06:tree-finished')
  })
  it('typing', async () => {
    const logger = createFakeLog()
    const job = createFakeJob(
      logger,
      delay,
      'typing',
      [
        'READ',
        'a',
        'READ',
        'b',
        'READ',
        'c'
      ]
    )
    expect(job.title).toBe('typing')
    expect(job.progress.actual.state).toBe('created')
    expect(logger.handleOutput).not.toBeCalled()
    expect(logger.logged?.length).toBe(0)
    const p = job.execute()
    await logger.aOk.promise
    expect(logger.logged).toMatchSnapshot('01:log-aOk')
    expect(logger.tree()).toMatchSnapshot('02:tree-aOk')
    job.type('123\n')
    await logger.bOk.promise
    expect(logger.logged).toMatchSnapshot('03:log-bOk')
    expect(logger.tree()).toMatchSnapshot('04:tree-bOk')
    job.type('456\n')
    await logger.cOk.promise
    expect(logger.logged).toMatchSnapshot('05:log-cOk')
    expect(logger.tree()).toMatchSnapshot('06:tree-cOk')
    job.type('789\n')
    expect(await p).toBeUndefined()
    expect(logger.logged).toMatchSnapshot('07:log-finished')
    expect(logger.tree()).toMatchSnapshot('08:tree-finished')
    expect(job.progress.actual.state).toBe('finished')
  })
  it('with dependencies', async () => {
    const logger = createFakeLog()
    const { main, dep } = create2FakeJob(
      logger,
      delay,
      'main', ['READ', 'a'],
      'dep', ['READ', 'b']
    )
    expect(main.title).toBe('main')
    expect(dep.title).toBe('dep')
    expect(main.progress.actual.state).toBe('created')
    expect(dep.progress.actual.state).toBe('created')
    expect(logger.handleOutput).not.toBeCalled()
    expect(logger.logged?.length).toBe(0)
    const pMain = main.execute()
    const pDep = dep.execute()
    await logger.bOk.promise
    expect(dep.progress.actual.state).toBe('working')
    expect(main.progress.actual.state).toBe('waiting')
    expect(logger.logged).toMatchSnapshot('01:log-dep working')
    expect(logger.tree()).toMatchSnapshot('02:tree-dep working')
    dep.type('dep-ok\n')
    await logger.aOk.promise
    expect(dep.progress.actual.state).toBe('finished')
    expect(main.progress.actual.state).toBe('working')
    expect(logger.logged).toMatchSnapshot('03:log-dep finished')
    expect(logger.tree()).toMatchSnapshot('04:tree-dep finished')
    main.type('main-ok\n')
    expect(await pDep).toBeUndefined()
    expect(await pMain).toBeUndefined()
    expect(dep.progress.actual.state).toBe('finished')
    expect(main.progress.actual.state).toBe('finished')
    expect(logger.logged).toMatchSnapshot('05:log-dep working')
    expect(logger.tree()).toMatchSnapshot('06:tree-dep working')
  })

  it('auto-start dependencies', async () => {
    const logger = createFakeLog()
    const { main, dep } = create2FakeJob(
      logger,
      delay,
      'main', ['READ', 'a'],
      'dep', ['READ', 'b']
    )
    expect(main.title).toBe('main')
    expect(dep.title).toBe('dep')
    expect(main.progress.actual.state).toBe('created')
    expect(dep.progress.actual.state).toBe('created')
    expect(logger.handleOutput).not.toBeCalled()
    expect(logger.logged?.length).toBe(0)
    const p1 = main.execute()
    const p2 = dep.execute()
    await logger.bOk.promise
    expect(dep.progress.actual.state).toBe('working')
    expect(main.progress.actual.state).toBe('waiting')
    expect(logger.logged).toMatchSnapshot('01:log-dep working')
    expect(logger.tree()).toMatchSnapshot('02:tree-dep working')
    dep.type('dep-ok\n')
    await logger.aOk.promise
    expect(dep.progress.actual.state).toBe('finished')
    expect(main.progress.actual.state).toBe('working')
    expect(logger.logged).toMatchSnapshot('03:log-dep finished')
    expect(logger.tree()).toMatchSnapshot('04:tree-dep finished')
    main.type('main-ok\n')
    expect(await p1).toBeUndefined()
    expect(await p2).toBeUndefined()
    expect(dep.progress.actual.state).toBe('finished')
    expect(main.progress.actual.state).toBe('finished')
    expect(logger.logged).toMatchSnapshot('05:log-finished')
    expect(logger.tree()).toMatchSnapshot('06:tree-finished')
  })
})
