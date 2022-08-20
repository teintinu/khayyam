import { khayyamCD } from './khayyam'
import { createFakeLog } from './testlib'

describe('Khayyam domain', () => {
  describe('CD', () => {
    it('x package', async () => {
      const logger = createFakeLog()
      const workspace = logger.fakeSys.loadWorkspace('npm/x/deps')
      const p = khayyamCD(
        workspace,
        logger.jobManager
      )
      expect(logger.logged).toMatchSnapshot('log')
      expect(logger.tree()).toMatchSnapshot('tree')
      await p
    })
    it('xy package', async () => {
      const logger = createFakeLog()
      const workspace = logger.fakeSys.loadWorkspace('npm/xy/deps')
      const p = khayyamCD(
        workspace,
        logger.jobManager
      )
      expect(logger.logged).toMatchSnapshot('log')
      expect(logger.tree()).toMatchSnapshot('tree')
      await p
    })
    it('a,b,c,d,e packages', async () => {
      const logger = createFakeLog()
      const workspace = logger.fakeSys.loadWorkspace('npm/abcde/deps')
      const p = khayyamCD(
        workspace,
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
