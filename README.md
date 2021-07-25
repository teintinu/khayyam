# monoclean

A clean repository automation for Typescript mono workspace projects. Includes yarn, jest and eslint predifined configurations. Customizations and new automations automations are planned.

# install

`npm i -g monoclean`

# Get stared

## declare you workspace `monoclean.yml`

```
workspace:
  name: "example"
  version: "1.0.0"

packages:
  "@example/a":
    folder: "a"
    dependencies:
      "@example/b": "*"

  "@example/b":
    folder: "b"
```

### `monoclean deps`

It will create and maintain automatically package.json, tsconfig.json, eslint.json, jest.config.js...

### `monoclean run`

`monoclean run [PACKAGE]` use esbuild to fasterly run `src/index.ts` on desired `[PACKAGE]`

### `monoclean test`

`monoclean test [PACKAGE]` use esbuild and jest to fasterly test all packages in the workspace

### another commands
- [ ] monoclean build
- [ ] monoclean deploy

## Notes
- This project is in alfa version, use responsibly.
- This project is based on https://github.com/deref/uni
