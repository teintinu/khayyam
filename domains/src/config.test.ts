describe.skip('config', () => {
  it.skip('config', () => {
    const yml = `workspace:
name: "khayyam"
version: "0.1.1"

domain:
"@khayyam/domain":
  folder: "domain"
  bundle: node-library

adapters:
"@khayyam/node-library-adapter":
  folder: "adapters/fake-repository"
  bundle: node-library
"@khayyam/node-application":
  folder: "adapters/fake-repository"
  bundle: node-library
"@khayyam/react-application":
  folder: "adapters/fake-repository"
  bundle: node-library

executables:
"@khayyam/server":
  folder: "apps/server"
  bundle: node-application

"@khayyam/client":
  folder: "apps/client"
  bundle: react-application
`
    expect(yml).toBe('x')
  })
  it.skip('custom bundlers', () => {
    const yml = `workspace:
name: "khayyam"
version: "0.1.1"

domain:
"@khayyam/domain":
  folder: "domain"
  type: node-library

adapters:
"@khayyam/node-library-adapter":
  folder: "adapters/fake-repository"
  type: node-library
"@khayyam/node-application":
  folder: "adapters/fake-repository"
  type: node-library
"@khayyam/react-application":
  folder: "adapters/fake-repository"
  type: node-library

executables:
"@khayyam/server":
  folder: "apps/server"
  type: node-application

"@khayyam/client":
  folder: "apps/client"
  type: react-application

bundlers:
  npm-bundler-sample:
    watch:
    - src/**/*.ts
    - src/**/*.tsx
    build:
    - npm run clean-dist
    - npm run build
    test:
    - npm test
`
    expect(yml).toBe('x')
  })
})
