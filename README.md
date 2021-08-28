# khayyam

A clean repository automation for Typescript mono workspace projects. Includes yarn, jest and eslint predifined configurations. Customizations and new automations automations are planned.

# install

`npm i -g khayyam`

# Get stared

## declare you workspace `khayyam.yml`

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

### `khayyam deps`

It will create and maintain automatically package.json, tsconfig.json, eslint.json, jest.config.js...

### `khayyam run`

`khayyam run [PACKAGE]` use esbuild to fasterly run `src/index.ts` on desired `[PACKAGE]`

### `khayyam test`

`khayyam test [--watch] [--coverage]` use esbuild and jest to fasterly test all packages in the workspace

### `khayyam build`

`khayyam build` build everything in the repository

### another commands
- [ ] khayyam run
- [ ] khayyam watch --tray
-   [ ] https://github.com/getlantern/systray
- [ ] khayyam deploy
- [ ] component diagram
- [ ] class diagram
- [ ] commitizen
- [ ] clean architecture sample

```
public class Order {
    items
}
public class Accounting { 
    public {[tax:Double]} calculateTax() {...}
    // or 
    public produceInvocice()
}
public class SalesPolicy {
    public void applyPromotions() {...}
}
public class StockManagement {
    public Stock checkItemsAvailability() {...}
}
public class OrderReporting {
    public String describeOrder() {...}
}
```

## Notes
- This project is in alfa version, use responsibly.
- This project is based on https://github.com/deref/uni

