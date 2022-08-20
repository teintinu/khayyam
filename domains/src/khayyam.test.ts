import { khayyamCD, khayyamDev } from './khayyam'
import { createFakeLog } from './testlib'

describe('Khayyam domain', () => {
  describe('CD', () => {
    it('x package', async () => {
      const logger = createFakeLog()
      const p = khayyamCD(
        logger.fakeSys,
        'npm/x/deps',
        logger.jobManager
      )
      expect(logger.logged).toMatchSnapshot('log')
      expect(logger.tree()).toMatchSnapshot('tree')
      await p
    })
    it('xy package', async () => {
      const logger = createFakeLog()
      const p = khayyamCD(
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
      const p = khayyamCD(
        logger.fakeSys,
        'npm/abcde/deps',
        logger.jobManager
      )
      expect(logger.logged).toMatchSnapshot('log')
      expect(logger.tree()).toMatchSnapshot('tree')
      await p
    }, 10000)
  })
  describe('dev', () => {
    it('x package', async () => {
      const logger = createFakeLog()
      const p = await khayyamDev(
        logger.fakeSys,
        'npm/x/deps',
        logger.jobManager
      )
      expect(logger.logged).toMatchSnapshot('01-log-before build')
      expect(logger.tree()).toMatchSnapshot('01-tree-before build')
      await p.waitBuild()
      expect(logger.logged).toMatchSnapshot('02-log-after build')
      expect(logger.tree()).toMatchSnapshot('02-tree-after build')

      logger.fakeSys.simulateChange('fakeBundlerNPM/x/deps/khayyam.yaml')
      await p.waitBuild()
      expect(logger.logged).toMatchSnapshot('03-log-after rebuild')
      expect(logger.tree()).toMatchSnapshot('03-tree-after rebuild')

      await p.stop()

      expect(logger.logged).toMatchSnapshot('04-log-after stop')
      expect(logger.tree()).toMatchSnapshot('04-tree-after stop')
      expect(logger.sysHandlersState()).toBe('clean')
    }, 15000)
    it('a,b,c,d,e packages', async () => {
      const logger = createFakeLog()
      const p = await khayyamDev(
        logger.fakeSys,
        'npm/abcde/deps',
        logger.jobManager
      )
      expect(logger.logged).toMatchSnapshot('01-log-before build')
      expect(logger.tree()).toMatchSnapshot('01-tree-before build')
      await p.waitBuild()
      expect(logger.logged).toMatchSnapshot('02-log-after build')
      expect(logger.tree()).toMatchSnapshot('02-tree-after build')

      logger.fakeSys.simulateChange('fakeBundlerNPM/abcde/deps/khayyam.yaml')
      await p.waitBuild()
      expect(logger.logged).toMatchSnapshot('03-log-after rebuild')
      expect(logger.tree()).toMatchSnapshot('03-tree-after rebuild')

      await p.stop()

      expect(logger.logged).toMatchSnapshot('04-log-after stop')
      expect(logger.tree()).toMatchSnapshot('04-tree-after stop')
      expect(logger.sysHandlersState()).toBe('clean')
    }, 35000)
  })
})
