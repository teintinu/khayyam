import { Bundler } from './bundler'
import { Job } from './job'
import { System } from './sys'

export interface Workspace {
  readonly layers: ReadonlyArray<Layer>
  readonly packages: ReadonlyArray<Package>
  readonly bundlers: ReadonlyArray<Bundler>
  findLayer (layerName: string): Layer|undefined
  findPackage (packageName: string): Package|undefined
  findBundler (bundlerName: string): Bundler|undefined
  walk (
    filter: 'all'|ByPackage<boolean>,
    fn: (pkg: Package, bundler: Bundler)=>Job): WalkedJobs
}

export interface WalkedJobs {
  jobs: ByPackage<Job[]>
  depends(...deps:WalkedJobs[]): void
}
export interface Layer {
    readonly name: string;
    readonly dependencies: string[]
}

export interface Package {
    readonly name: string;
    readonly folder: string;
    readonly layer: string;
    readonly repository?: string
    readonly bundlers: ReadonlyArray<string>
    readonly dependencies: ReadonlyArray<string>
}

export interface ByPackage<T> {
  [pkgName:string]:T
}

export function createWorkspace ({
  sys,
  layers,
  packages,
  bundlers
}:{
  sys: System,
  layers: Layer[],
  packages: Package[],
  bundlers: Bundler[]
}) {
  const ws: Workspace = {
    get layers () {
      return layers
    },
    get packages () {
      return packages
    },
    get bundlers () {
      return bundlers
    },
    findLayer,
    findPackage,
    findBundler,
    walk
  }
  return ws
  function findLayer (layerName: string): Layer|undefined {
    return layers.find((layer) => layer.name === layerName)
  }
  function findBundler (bundlerName: string) {
    return bundlers.find((bundler) => bundler.name === bundlerName)
  }
  function findPackage (packageName: string) {
    return packages.find((pkg) => pkg.name === packageName)
  }
  function walk (
    filter: 'all'|ByPackage<boolean>,
    fn: (pkg: Package, bundler: Bundler)=>Job): WalkedJobs {
    const ret: WalkedJobs = {
      jobs: {},
      depends (...deps:WalkedJobs[]): void {
        Object.keys(ret.jobs).forEach(nDependant => {
          const tDependants = ret.jobs[nDependant]
          tDependants.forEach(tDependant => {
            deps.forEach(dep => {
              Object.keys(dep.jobs).forEach(nDependency => {
                const tDependencies = ret.jobs[nDependency]
                tDependencies.forEach(tDependency =>
                  tDependant.depends(tDependency)
                )
              })
            })
          })
        })
      }
    }
    const flat: ByPackage<Job[]> = {}
    const tree: ByPackage<true> = {}
    packages.forEach((pkg) => walkOn(pkg))
    return ret
    function walkOn (pkg: Package): boolean {
      let filtered = filter === 'all' || filter[pkg.name]
      if (tree[pkg.name]) {
        sys.notify(
          'Circular dependency on packages' + Object.keys(tree),
          pkg.name
        )
        return false
      }
      if (!flat[pkg.name]) {
        flat[pkg.name] = []
        tree[pkg.name] = true
        pkg.dependencies.forEach((depName) => {
          const dep = checkAndGetDependency(pkg, depName)
          if (dep && walkOn(dep)) {
            filtered = true
          }
        })
        if (filtered) {
          pkg.bundlers.forEach((bundlerName) => {
            const b = findBundler(bundlerName)
            if (b) {
              const nJob = fn(pkg, b)
              console.log('push flag ' + nJob.title + ' on ' + pkg.name, pkg)
              flat[pkg.name].push(nJob)
              pkg.dependencies.forEach((depName) => {
                if (flat[depName]) {
                  console.log('push deps on ' + depName + ' on ' + nJob.title)
                  nJob.depends(...flat[depName])
                }
              })
              if (ret.jobs[pkg.name]) {
                ret.jobs[pkg.name].push(nJob)
              } else {
                ret.jobs[pkg.name] = [nJob]
              }
            } else {
              sys.notify(
                'Bundler not found: ' + bundlerName + ' on package ' + pkg.name,
                pkg.name
              )
            }
          })
        }
        delete tree[pkg.name]
      }
      return filtered
    }
    function checkAndGetDependency (pkg:Package, depName: string): Package|null {
      const dep = findPackage(depName)
      if (!dep) {
        sys.notify(
          'Dependency not found: ' + depName + ' on package ' + pkg.name,
          pkg.name
        )
        return 1 as any as null
      }
      const pkgLayer = findLayer(pkg.layer)
      if (pkg.layer && !pkgLayer) {
        sys.notify(
          'Layer not found: ' + pkg.layer + ' on package ' + pkg.name,
          pkg.name
        )
        return 2 as any as null
      }
      const depLayer = findLayer(dep.layer)
      if (dep.layer && !depLayer) {
        sys.notify(
          'Layer not found: ' + dep.layer + ' on package ' + dep.name,
          dep.name
        )
        return 3 as any as null
      }
      if (pkgLayer && depLayer && checkLayers(pkgLayer, depLayer)) {
        sys.notify(
          'The package ' + pkg.name + ' on layer ' + pkg.layer +
          " can't depends on package " + dep.name + ' on layer ' + dep.layer,
          pkg.name
        )
        return 4 as any as null
      }
      return dep
      function checkLayers (usingLayer: Layer, usedLayer: Layer): boolean {
        const sameLayer = usingLayer === usedLayer
        if (sameLayer) return true
        return usingLayer.dependencies.some((n1) =>
          usedLayer.dependencies.some((n2) => {
            const l1 = findLayer(n1)
            if (!l1) {
              sys.notify('Layer ' + n1 + ' not exist but is referenced by layer' + usingLayer, pkg.name)
              return false
            }
            const l2 = findLayer(n2)
            if (!l2) {
              sys.notify('Layer ' + n2 + ' not exist but is referenced by layer' + usedLayer, pkg.name)
              return false
            }
            return checkLayers(l1, l2)
          })
        )
      }
    }
  }
}
