import { createFakeJob, createFakeLog, walkedDebug } from './testlib'

describe('workspace', () => {
  describe('just one package', () => {
    it('create workspace', async () => {
      const logger = createFakeLog()
      const ws = logger.fakeSys.loadWorkspace('npm/x')
      expect(ws.packages.map((p) => p.name).sort().join()).toBe('x')
      expect(ws.bundlers.map((b) => b.name).sort().join()).toBe('fakeBundlerNPM')
      expect(logger.logged?.length).toBe(0)
      expect({
        layers: ws.layers,
        pkgs: ws.packages
      }).toMatchSnapshot()
    })
    it('findPackage workspace', async () => {
      const logger = createFakeLog()
      const ws = logger.fakeSys.loadWorkspace('npm/x')
      expect(ws.findPackage('x')?.name).toBe('x')
      expect(ws.findPackage('a')).toBeUndefined()
      expect(logger.logged?.length).toBe(0)
    })
    it('findBundler workspace', async () => {
      const logger = createFakeLog()
      const ws = logger.fakeSys.loadWorkspace('npm/x')
      expect(ws.findBundler('fakeBundlerNPM')?.name).toBe('fakeBundlerNPM')
      expect(ws.findBundler('invalid')).toBeUndefined()
      expect(logger.logged?.length).toBe(0)
    })
    describe('walk workspace', () => {
      it('all', () => {
        const logger = createFakeLog()
        const ws = logger.fakeSys.loadWorkspace('npm/x')
        const pkgs: string[] = []
        const walked = ws.walk('all', (pkg) => {
          pkgs.push(pkg.name)
          return createFakeJob(
            logger,
            10,
            pkg.name + ' walk',
            ['walk']
          )
        })
        expect(pkgs.join()).toBe('x')
        expect(logger.logged?.length).toBe(0)
        expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
      })
      it('filter x', () => {
        const logger = createFakeLog()
        const ws = logger.fakeSys.loadWorkspace('npm/x')
        const pkgs: string[] = []
        const walked = ws.walk({ x: true }, (pkg) => {
          pkgs.push(pkg.name)
          return createFakeJob(
            logger,
            10,
            pkg.name + ' walk',
            ['walk']
          )
        })
        expect(pkgs.join()).toBe('x')
        expect(logger.logged?.length).toBe(0)
        expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
      })
      it('filter none', () => {
        const logger = createFakeLog()
        const ws = logger.fakeSys.loadWorkspace('npm/x')
        const pkgs: string[] = []
        const walked = ws.walk({}, (pkg) => {
          pkgs.push(pkg.name)
          return createFakeJob(
            logger,
            10,
            pkg.name + ' walk',
            ['walk']
          )
        })
        expect(pkgs.join()).toBe('')
        expect(logger.logged?.length).toBe(0)
        expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
      })
    })
  })
  describe('just x,y packages', () => {
    it('create workspace', async () => {
      const logger = createFakeLog()
      const ws = logger.fakeSys.loadWorkspace('npm/xy/deps')
      expect(ws.packages.map((p) => p.name).sort().join()).toBe('x,y')
      expect(ws.bundlers.map((b) => b.name).sort().join()).toBe('fakeBundlerNPM')
      expect(logger.logged?.length).toBe(0)
      expect({
        layers: ws.layers,
        pkgs: ws.packages
      }).toMatchSnapshot()
    })
    it('findPackage workspace', async () => {
      const logger = createFakeLog()
      const ws = logger.fakeSys.loadWorkspace('npm/xy/deps')
      const x = ws.findPackage('x')
      expect(x?.name).toBe('x')
      expect(x?.dependencies.join()).toBe('y')
      expect(ws.findPackage('a')).toBeUndefined()
      expect(logger.logged?.length).toBe(0)
    })
    it('findBundler workspace', async () => {
      const logger = createFakeLog()
      const ws = logger.fakeSys.loadWorkspace('npm/xy/deps')
      expect(ws.findBundler('fakeBundlerNPM')?.name).toBe('fakeBundlerNPM')
      expect(ws.findBundler('invalid')).toBeUndefined()
      expect(logger.logged?.length).toBe(0)
    })
    describe('walk workspace', () => {
      it('all', () => {
        const logger = createFakeLog()
        const ws = logger.fakeSys.loadWorkspace('npm/xy/deps')
        const pkgs: string[] = []
        const walked = ws.walk('all', (pkg) => {
          pkgs.push(pkg.name)
          return createFakeJob(
            logger,
            10,
            pkg.name + ' walkXY',
            ['walk']
          )
        })
        expect(pkgs.join()).toBe('y,x')
        expect(logger.logged?.length).toBe(0)
        expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
      })
      it('filter x', () => {
        const logger = createFakeLog()
        const ws = logger.fakeSys.loadWorkspace('npm/xy/deps')
        const pkgs: string[] = []
        const walked = ws.walk({ x: true }, (pkg) => {
          pkgs.push(pkg.name)
          return createFakeJob(
            logger,
            10,
            pkg.name + ' walkXY',
            ['walk']
          )
        })
        expect(pkgs.join()).toBe('x')
        expect(logger.logged?.length).toBe(0)
        expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
      })
      it('filter y', () => {
        const logger = createFakeLog()
        const ws = logger.fakeSys.loadWorkspace('npm/xy/deps')
        const pkgs: string[] = []
        const walked = ws.walk({ y: true }, (pkg) => {
          pkgs.push(pkg.name)
          return createFakeJob(
            logger,
            10,
            pkg.name + ' walkXY',
            ['walk']
          )
        })
        expect(pkgs.join()).toBe('y,x')
        expect(logger.logged?.length).toBe(0)
        expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
      })
      it('filter none', () => {
        const logger = createFakeLog()
        const ws = logger.fakeSys.loadWorkspace('npm/xy/deps')
        const pkgs: string[] = []
        const walked = ws.walk({}, (pkg) => {
          pkgs.push(pkg.name)
          return createFakeJob(
            logger,
            10,
            pkg.name + ' walkXY',
            ['walk']
          )
        })
        expect(pkgs.join()).toBe('')
        expect(logger.logged?.length).toBe(0)
        expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
      })
    })
  })
  describe('just a,b,c,d,e packages', () => {
    it('create workspace', async () => {
      const logger = createFakeLog()
      const ws = logger.fakeSys.loadWorkspace('npm/abcde/deps')
      expect(ws.packages.map((p) => p.name).sort().join()).toBe('a,b,c,d,e')
      expect(ws.bundlers.map((b) => b.name).sort().join()).toBe('fakeBundlerNPM')
      expect(logger.logged?.length).toBe(0)
      expect({
        layers: ws.layers,
        pkgs: ws.packages
      }).toMatchSnapshot()
    })
    it('findPackage workspace', async () => {
      const logger = createFakeLog()
      const ws = logger.fakeSys.loadWorkspace('npm/abcde/deps')
      const a = ws.findPackage('a')
      expect(a?.name).toBe('a')
      expect(a?.dependencies.join()).toBe('b,c')
      expect(ws.findPackage('x')).toBeUndefined()
      expect(logger.logged?.length).toBe(0)
    })
    it('findBundler workspace', async () => {
      const logger = createFakeLog()
      const ws = logger.fakeSys.loadWorkspace('npm/abcde/deps')
      expect(ws.findBundler('fakeBundlerNPM')?.name).toBe('fakeBundlerNPM')
      expect(ws.findBundler('invalid')).toBeUndefined()
      expect(logger.logged?.length).toBe(0)
    })
    describe('walk workspace', () => {
      it('all', () => {
        const logger = createFakeLog()
        const ws = logger.fakeSys.loadWorkspace('npm/abcde/deps')
        const pkgs: string[] = []
        const walked = ws.walk('all', (pkg) => {
          pkgs.push(pkg.name)
          return createFakeJob(
            logger,
            10,
            pkg.name + ' walkABCDE',
            ['walk']
          )
        })
        expect(pkgs.join()).toBe('e,d,c,b,a')
        expect(logger.logged?.length).toBe(0)
        expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
      })
      it('filter a', () => {
        const logger = createFakeLog()
        const ws = logger.fakeSys.loadWorkspace('npm/abcde/deps')
        const pkgs: string[] = []
        const walked = ws.walk({ a: true }, (pkg) => {
          pkgs.push(pkg.name)
          return createFakeJob(
            logger,
            10,
            pkg.name + ' walkABCDE',
            ['walk']
          )
        })
        expect(pkgs.join()).toBe('a')
        expect(logger.logged?.length).toBe(0)
        expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
      })
      it('filter b', () => {
        const logger = createFakeLog()
        const ws = logger.fakeSys.loadWorkspace('npm/abcde/deps')
        const pkgs: string[] = []
        const walked = ws.walk({ b: true }, (pkg) => {
          pkgs.push(pkg.name)
          return createFakeJob(
            logger,
            10,
            pkg.name + ' walkABCDE',
            ['walk']
          )
        })
        expect(pkgs.join()).toBe('b,a')
        expect(logger.logged?.length).toBe(0)
        expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
      })
      it('filter c', () => {
        const logger = createFakeLog()
        const ws = logger.fakeSys.loadWorkspace('npm/abcde/deps')
        const pkgs: string[] = []
        const walked = ws.walk({ c: true }, (pkg) => {
          pkgs.push(pkg.name)
          return createFakeJob(
            logger,
            10,
            pkg.name + ' walkABCDE',
            ['walk']
          )
        })
        expect(pkgs.join()).toBe('c,b,a')
        expect(logger.logged?.length).toBe(0)
        expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
      })
      it('filter d', () => {
        const logger = createFakeLog()
        const ws = logger.fakeSys.loadWorkspace('npm/abcde/deps')
        const pkgs: string[] = []
        const walked = ws.walk({ d: true }, (pkg) => {
          pkgs.push(pkg.name)
          return createFakeJob(
            logger,
            10,
            pkg.name + ' walkABCDE',
            ['walk']
          )
        })
        expect(pkgs.join()).toBe('d,c,b,a')
        expect(logger.logged?.length).toBe(0)
        expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
      })
      it('filter c,d', () => {
        const logger = createFakeLog()
        const ws = logger.fakeSys.loadWorkspace('npm/abcde/deps')
        const pkgs: string[] = []
        const walked = ws.walk({ c: true, d: true }, (pkg) => {
          pkgs.push(pkg.name)
          return createFakeJob(
            logger,
            10,
            pkg.name + ' walkABCDE',
            ['walk']
          )
        })
        expect(pkgs.join()).toBe('d,c,b,a')
        expect(logger.logged?.length).toBe(0)
        expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
      })
      it('filter e', () => {
        const logger = createFakeLog()
        const ws = logger.fakeSys.loadWorkspace('npm/abcde/deps')
        const pkgs: string[] = []
        const walked = ws.walk({ e: true }, (pkg) => {
          pkgs.push(pkg.name)
          return createFakeJob(
            logger,
            10,
            pkg.name + ' walkABCDE',
            ['walk']
          )
        })
        expect(pkgs.join()).toBe('e,d,c,b,a')
        expect(logger.logged?.length).toBe(0)
        expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
      })
      it('filter none', () => {
        const logger = createFakeLog()
        const ws = logger.fakeSys.loadWorkspace('npm/abcde/deps')
        const pkgs: string[] = []
        const walked = ws.walk({}, (pkg) => {
          pkgs.push(pkg.name)
          return createFakeJob(
            logger,
            10,
            pkg.name + ' walkABCDE',
            ['walk']
          )
        })
        expect(pkgs.join()).toBe('')
        expect(logger.logged?.length).toBe(0)
        expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
      })
    })
  })
})
