/* eslint-disable no-redeclare */
import { khayyamCI } from './khayyam'
import { createFakeLog } from './testlib'

describe.skip('Khayyam domain', () => {
  describe('CI', () => {
    it('x package', async () => {
      const logger = createFakeLog(true)
      const p = khayyamCI(
        logger.fakeSys,
        'npm/x/deps',
        logger.jobManager
      )
      expect(logger.logged).toMatchSnapshot('log')
      expect(logger.tree()).toMatchSnapshot('tree')
      await p
    })
    it('x package', async () => {
      const logger = createFakeLog()
      const p = khayyamCI(
        logger.fakeSys,
        'npm/xy/deps',
        logger.jobManager
      )
      expect(logger.logged).toMatchSnapshot('log')
      expect(logger.tree()).toMatchSnapshot('tree')
      await p
    })
    it('a,b,c,d,e packages', async () => {
      const logger = createFakeLog()
      const p = khayyamCI(
        logger.fakeSys,
        'npm/abcde/deps',
        logger.jobManager
      )
      expect(logger.logged).toMatchSnapshot('log')
      expect(logger.tree()).toMatchSnapshot('tree')
      await p
    }, 10000)
  })
  // describe('dev', () => {
  //   it('x package', async () => {
  //     const logger = createFakeLog()
  //     const bundler = createFakeBundlerNPM(logger)
  //     const ws = createFakeWorkspaceX(logger, bundler)
  //     await khayyamDev(
  //       logger.jobManager,
  //       ws
  //     )
  //     expect(logger.logged).toMatchSnapshot('log')
  //     expect(logger.tree()).toMatchSnapshot('tree')
  //   })
  //   it('a,b,c,d,e packages', async () => {
  //     const logger = createFakeLog()
  //     const bundler = createFakeBundlerNPM(logger)
  //     const ws = createFakeWorkspaceABCDE(logger, bundler)
  //     await khayyamDev(
  //       logger.jobManager,
  //       ws
  //     )
  //     expect(logger.logged).toMatchSnapshot('log')
  //     expect(logger.tree()).toMatchSnapshot('tree')
  //   }, 10000)
  // })
})
