import { createFakeJob, createFakeLog, walkedDebug } from './testlib'

describe('workspace', () => {
  describe('just one package', () => {
    it('create', async () => {
      const logger = createFakeLog()
      const ws = await logger.fakeSys.loadWorkspace('npm/x')
      expect(ws.packages.map((p) => p.name).sort().join()).toBe('x')
      expect(ws.bundlers.map((b) => b.name).sort().join()).toBe('fakeBundlerNPM')
      expect({
        layers: ws.layers,
        pkgs: ws.packages
      }).toMatchSnapshot('layers/pkg')
      expect(logger.logged).toMatchSnapshot('logged')
    })
    it('findPackage', async () => {
      const logger = createFakeLog()
      const ws = await logger.fakeSys.loadWorkspace('npm/x')
      expect(ws.findPackage('x')?.name).toBe('x')
      expect(ws.findPackage('a')).toBeUndefined()
      expect(logger.logged).toMatchSnapshot('logged')
    })
    it('findBundler', async () => {
      const logger = createFakeLog()
      const ws = await logger.fakeSys.loadWorkspace('npm/x')
      expect(ws.findBundler('fakeBundlerNPM')?.name).toBe('fakeBundlerNPM')
      expect(ws.findBundler('invalid')).toBeUndefined()
      expect(logger.logged).toMatchSnapshot('logged')
    })
    describe('walk ', () => {
      [false, true].forEach((treeDep) => {
        describe('treeDep=' + (treeDep ? 'true' : 'false'), () => {
          describe('one walk', () => {
            it('all', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/x')
              const pkgs: string[] = []
              const walked = ws.walk('all', true, treeDep, (pkg) => {
                pkgs.push(pkg.name)
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk',
                  ['walk']
                )
              })
              expect(pkgs.join()).toBe('x')
              expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter x', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/x')
              const pkgs: string[] = []
              const walked = ws.walk({ x: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name)
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk',
                  ['walk']
                )
              })
              expect(pkgs.join()).toBe('x')
              expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter none', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/x')
              const pkgs: string[] = []
              const walked = ws.walk({}, true, treeDep, (pkg) => {
                pkgs.push(pkg.name)
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk',
                  ['walk']
                )
              })
              expect(pkgs.join()).toBe('')
              expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
          })
          describe('two walks', () => {
            it('all', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/x')
              const pkgs: string[] = []
              const walked1 = ws.walk('all', true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1',
                  ['walk1']
                )
              })
              const walked2 = ws.walk('all', true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2',
                  ['walk2']
                )
              })
              expect(pkgs.join()).toBe('x1,x2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter x', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/x')
              const pkgs: string[] = []
              const walked1 = ws.walk({ x: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ x: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2',
                  ['walk2']
                )
              })
              expect(pkgs.join()).toBe('x1,x2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter none', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/x')
              const pkgs: string[] = []
              const walked1 = ws.walk({}, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({}, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2',
                  ['walk2']
                )
              })
              expect(pkgs.join()).toBe('')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
          })
          describe('walk all-deps', () => {
            it('all', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/x')
              const pkgs: string[] = []
              const walked1 = ws.walk('all', true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1',
                  ['walk1']
                )
              })
              const walked2 = ws.walk('all', true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2',
                  ['walk2']
                )
              })
              walked1.depends('all', [walked2])
              expect(pkgs.join()).toBe('x1,x2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter x', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/x')
              const pkgs: string[] = []
              const walked1 = ws.walk({ x: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ x: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2',
                  ['walk2']
                )
              })
              walked1.depends('all', [walked2])
              expect(pkgs.join()).toBe('x1,x2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter none', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/x')
              const pkgs: string[] = []
              const walked1 = ws.walk({}, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({}, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2',
                  ['walk2']
                )
              })
              walked1.depends('all', [walked2])
              expect(pkgs.join()).toBe('')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
          })
          describe('walk each-deps', () => {
            it('all', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/x')
              const pkgs: string[] = []
              const walked1 = ws.walk('all', true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1',
                  ['walk1']
                )
              })
              const walked2 = ws.walk('all', true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2',
                  ['walk2']
                )
              })
              walked1.depends('each', [walked2])
              expect(pkgs.join()).toBe('x1,x2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter x', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/x')
              const pkgs: string[] = []
              const walked1 = ws.walk({ x: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ x: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2',
                  ['walk2']
                )
              })
              walked1.depends('each', [walked2])
              expect(pkgs.join()).toBe('x1,x2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter none', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/x')
              const pkgs: string[] = []
              const walked1 = ws.walk({}, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({}, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2',
                  ['walk2']
                )
              })
              walked1.depends('each', [walked2])
              expect(pkgs.join()).toBe('')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
          })
        })
      })
    })
  })
  describe('just x,y packages', () => {
    it('create', async () => {
      const logger = createFakeLog()
      const ws = await logger.fakeSys.loadWorkspace('npm/xy/deps')
      expect(ws.packages.map((p) => p.name).sort().join()).toBe('x,y')
      expect(ws.bundlers.map((b) => b.name).sort().join()).toBe('fakeBundlerNPM')
      expect({
        layers: ws.layers,
        pkgs: ws.packages
      }).toMatchSnapshot('layers/pkg')
      expect(logger.logged).toMatchSnapshot('logged')
    })
    it('findPackage', async () => {
      const logger = createFakeLog()
      const ws = await logger.fakeSys.loadWorkspace('npm/xy/deps')
      const x = ws.findPackage('x')
      expect(x?.name).toBe('x')
      expect(x?.dependencies.join()).toBe('y')
      expect(ws.findPackage('a')).toBeUndefined()
      expect(logger.logged).toMatchSnapshot('logged')
    })
    it('findBundler', async () => {
      const logger = createFakeLog()
      const ws = await logger.fakeSys.loadWorkspace('npm/xy/deps')
      expect(ws.findBundler('fakeBundlerNPM')?.name).toBe('fakeBundlerNPM')
      expect(ws.findBundler('invalid')).toBeUndefined()
      expect(logger.logged).toMatchSnapshot('logged')
    })
    describe('walk', () => {
      [false, true].forEach((treeDep) => {
        describe('treeDep=' + (treeDep ? 'true' : 'false'), () => {
          describe('one walk', () => {
            it('all', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/xy/deps')
              const pkgs: string[] = []
              const walked = ws.walk('all', true, treeDep, (pkg) => {
                pkgs.push(pkg.name)
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walkXY',
                  ['walk']
                )
              })
              expect(pkgs.join()).toBe('y,x')
              expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter x', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/xy/deps')
              const pkgs: string[] = []
              const walked = ws.walk({ x: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name)
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walkXY',
                  ['walk']
                )
              })
              expect(pkgs.join()).toBe('x')
              expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter y', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/xy/deps')
              const pkgs: string[] = []
              const walked = ws.walk({ y: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name)
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walkXY',
                  ['walk']
                )
              })
              expect(pkgs.join()).toBe('y,x')
              expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter none', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/xy/deps')
              const pkgs: string[] = []
              const walked = ws.walk({}, true, treeDep, (pkg) => {
                pkgs.push(pkg.name)
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walkXY',
                  ['walk']
                )
              })
              expect(pkgs.join()).toBe('')
              expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
          })
          describe('two walks', () => {
            it('all', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/xy/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk('all', true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1XY',
                  ['walk1']
                )
              })
              const walked2 = ws.walk('all', true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2XY',
                  ['walk2']
                )
              })
              expect(pkgs.join()).toBe('y1,x1,y2,x2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter x', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/xy/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ x: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1XY',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ x: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2XY',
                  ['walk2']
                )
              })
              expect(pkgs.join()).toBe('x1,x2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter y', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/xy/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ y: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1XY',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ y: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2XY',
                  ['walk2']
                )
              })
              expect(pkgs.join()).toBe('y1,x1,y2,x2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter none', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/xy/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({}, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1XY',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({}, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2XY',
                  ['walk2']
                )
              })
              expect(pkgs.join()).toBe('')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
          })
          describe('walk all-deps', () => {
            it('all', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/xy/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk('all', true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1XY',
                  ['walk1']
                )
              })
              const walked2 = ws.walk('all', true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2XY',
                  ['walk2']
                )
              })
              walked1.depends('all', [walked2])
              expect(pkgs.join()).toBe('y1,x1,y2,x2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter x', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/xy/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ x: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1XY',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ x: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2XY',
                  ['walk2']
                )
              })
              walked1.depends('all', [walked2])
              expect(pkgs.join()).toBe('x1,x2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter y', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/xy/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ y: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1XY',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ y: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2XY',
                  ['walk2']
                )
              })
              walked1.depends('all', [walked2])
              expect(pkgs.join()).toBe('y1,x1,y2,x2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter none', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/xy/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({}, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1XY',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({}, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2XY',
                  ['walk2']
                )
              })
              walked1.depends('all', [walked2])
              expect(pkgs.join()).toBe('')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
          })
          describe('walk each-deps', () => {
            it('all', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/xy/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk('all', true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1XY',
                  ['walk1']
                )
              })
              const walked2 = ws.walk('all', true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2XY',
                  ['walk2']
                )
              })
              walked1.depends('each', [walked2])
              expect(pkgs.join()).toBe('y1,x1,y2,x2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter x', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/xy/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ x: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1XY',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ x: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2XY',
                  ['walk2']
                )
              })
              walked1.depends('each', [walked2])
              expect(pkgs.join()).toBe('x1,x2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter y', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/xy/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ y: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1XY',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ y: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2XY',
                  ['walk2']
                )
              })
              walked1.depends('each', [walked2])
              expect(pkgs.join()).toBe('y1,x1,y2,x2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter none', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/xy/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({}, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1XY',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({}, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2XY',
                  ['walk2']
                )
              })
              walked1.depends('each', [walked2])
              expect(pkgs.join()).toBe('')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
          })
        })
      })
    })
  })
  describe('just a,b,c,d,e packages', () => {
    it('create', async () => {
      const logger = createFakeLog()
      const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
      expect(ws.packages.map((p) => p.name).sort().join()).toBe('a,b,c,d,e')
      expect(ws.bundlers.map((b) => b.name).sort().join()).toBe('fakeBundlerNPM')
      expect({
        layers: ws.layers,
        pkgs: ws.packages
      }).toMatchSnapshot('layers/pkg')
      expect(logger.logged).toMatchSnapshot('logged')
    })
    it('findPackage', async () => {
      const logger = createFakeLog()
      const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
      const a = ws.findPackage('a')
      expect(a?.name).toBe('a')
      expect(a?.dependencies.join()).toBe('b,c')
      expect(ws.findPackage('x')).toBeUndefined()
      expect(logger.logged).toMatchSnapshot('logged')
    })
    it('findBundler', async () => {
      const logger = createFakeLog()
      const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
      expect(ws.findBundler('fakeBundlerNPM')?.name).toBe('fakeBundlerNPM')
      expect(ws.findBundler('invalid')).toBeUndefined()
      expect(logger.logged).toMatchSnapshot('logged')
    })
    describe('walk', () => {
      [false, true].forEach((treeDep) => {
        describe('treeDep=' + (treeDep ? 'true' : 'false'), () => {
          describe('one walk', () => {
            it('all', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked = ws.walk('all', true, treeDep, (pkg) => {
                pkgs.push(pkg.name)
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walkABCDE',
                  ['walk']
                )
              })
              expect(pkgs.join()).toBe('e,d,c,b,a')
              expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter a', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked = ws.walk({ a: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name)
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walkABCDE',
                  ['walk']
                )
              })
              expect(pkgs.join()).toBe('a')
              expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter b', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked = ws.walk({ b: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name)
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walkABCDE',
                  ['walk']
                )
              })
              expect(pkgs.join()).toBe('b,a')
              expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter c', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked = ws.walk({ c: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name)
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walkABCDE',
                  ['walk']
                )
              })
              expect(pkgs.join()).toBe('c,b,a')
              expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter d', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked = ws.walk({ d: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name)
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walkABCDE',
                  ['walk']
                )
              })
              expect(pkgs.join()).toBe('d,c,b,a')
              expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter c,d', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked = ws.walk({ c: true, d: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name)
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walkABCDE',
                  ['walk']
                )
              })
              expect(pkgs.join()).toBe('d,c,b,a')
              expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter e', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked = ws.walk({ e: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name)
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walkABCDE',
                  ['walk']
                )
              })
              expect(pkgs.join()).toBe('e,d,c,b,a')
              expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter none', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked = ws.walk({}, true, treeDep, (pkg) => {
                pkgs.push(pkg.name)
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walkABCDE',
                  ['walk']
                )
              })
              expect(pkgs.join()).toBe('')
              expect(walkedDebug(walked)).toMatchSnapshot('walkedDebug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
          })
          describe('two walks', () => {
            it('all', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk('all', true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk('all', true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              expect(pkgs.join()).toBe('e1,d1,c1,b1,a1,e2,d2,c2,b2,a2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter a', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ a: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ a: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              expect(pkgs.join()).toBe('a1,a2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter b', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ b: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ b: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              expect(pkgs.join()).toBe('b1,a1,b2,a2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter c', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ c: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ c: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              expect(pkgs.join()).toBe('c1,b1,a1,c2,b2,a2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter d', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ d: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ d: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              expect(pkgs.join()).toBe('d1,c1,b1,a1,d2,c2,b2,a2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter c,d', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ c: true, d: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ c: true, d: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              expect(pkgs.join()).toBe('d1,c1,b1,a1,d2,c2,b2,a2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter e', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ e: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ e: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              expect(pkgs.join()).toBe('e1,d1,c1,b1,a1,e2,d2,c2,b2,a2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter none', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({}, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({}, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              expect(pkgs.join()).toBe('')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
          })
          describe('walk all-deps', () => {
            it('all', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk('all', true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk('all', true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              walked1.depends('all', [walked2])
              expect(pkgs.join()).toBe('e1,d1,c1,b1,a1,e2,d2,c2,b2,a2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter a', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ a: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ a: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              walked1.depends('all', [walked2])
              expect(pkgs.join()).toBe('a1,a2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter b', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ b: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ b: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              walked1.depends('all', [walked2])
              expect(pkgs.join()).toBe('b1,a1,b2,a2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter c', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ c: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ c: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              walked1.depends('all', [walked2])
              expect(pkgs.join()).toBe('c1,b1,a1,c2,b2,a2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter d', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ d: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ d: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              walked1.depends('all', [walked2])
              expect(pkgs.join()).toBe('d1,c1,b1,a1,d2,c2,b2,a2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter c,d', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ c: true, d: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ c: true, d: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              walked1.depends('all', [walked2])
              expect(pkgs.join()).toBe('d1,c1,b1,a1,d2,c2,b2,a2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter e', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ e: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ e: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              walked1.depends('all', [walked2])
              expect(pkgs.join()).toBe('e1,d1,c1,b1,a1,e2,d2,c2,b2,a2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter none', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({}, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({}, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              walked1.depends('all', [walked2])
              expect(pkgs.join()).toBe('')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
          })
          describe('walk each-deps', () => {
            it('all', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk('all', true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk('all', true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              walked1.depends('each', [walked2])
              expect(pkgs.join()).toBe('e1,d1,c1,b1,a1,e2,d2,c2,b2,a2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter a', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ a: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ a: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              walked1.depends('each', [walked2])
              expect(pkgs.join()).toBe('a1,a2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter b', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ b: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ b: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              walked1.depends('each', [walked2])
              expect(pkgs.join()).toBe('b1,a1,b2,a2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter c', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ c: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ c: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              walked1.depends('each', [walked2])
              expect(pkgs.join()).toBe('c1,b1,a1,c2,b2,a2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter d', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ d: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ d: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              walked1.depends('each', [walked2])
              expect(pkgs.join()).toBe('d1,c1,b1,a1,d2,c2,b2,a2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter c,d', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ c: true, d: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ c: true, d: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              walked1.depends('each', [walked2])
              expect(pkgs.join()).toBe('d1,c1,b1,a1,d2,c2,b2,a2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter e', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({ e: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({ e: true }, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              walked1.depends('each', [walked2])
              expect(pkgs.join()).toBe('e1,d1,c1,b1,a1,e2,d2,c2,b2,a2')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
            it('filter none', async () => {
              const logger = createFakeLog()
              const ws = await logger.fakeSys.loadWorkspace('npm/abcde/deps')
              const pkgs: string[] = []
              const walked1 = ws.walk({}, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '1')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk1ABCDE',
                  ['walk1']
                )
              })
              const walked2 = ws.walk({}, true, treeDep, (pkg) => {
                pkgs.push(pkg.name + '2')
                return createFakeJob(
                  logger,
                  10,
                  pkg.name + ' walk2ABCDE',
                  ['walk2']
                )
              })
              walked1.depends('each', [walked2])
              expect(pkgs.join()).toBe('')
              expect(walkedDebug(walked1)).toMatchSnapshot('walked1Debug')
              expect(walkedDebug(walked2)).toMatchSnapshot('walked2Debug')
              expect(logger.logged).toMatchSnapshot('logged')
            })
          })
        })
      })
    })
  })
})
