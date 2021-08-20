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

`monoclean test [--watch] [--coverage]` use esbuild and jest to fasterly test all packages in the workspace

### `monoclean build`

`monoclean build` build everything in the repository

### another commands
- [ ] monoclean run
- [ ] monoclean watch --tray
-   [ ] https://github.com/getlantern/systray
- [ ] monoclean deploy
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

