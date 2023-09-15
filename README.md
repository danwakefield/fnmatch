# fnmatch
Updated clone of kballards golang fnmatch gist (https://gist.github.com/kballard/272720)

## JavaScript interopability

This module has javascript<>go interop via wasm

### Build for wasm target

```
tinygo build -o main.wasm -target wasm ./main.go
```

### Run in browser
Serve `index.html`
```
npx serve
```

Open the dev tools and `fnmatch` is available on the global scope.


### Using in a JavaScript project
This module has a named export `fnmatch`, the library as is was built to work in the browser and might not be compatible with every build system.

```ts
import { fnmatch } from '..'

fnmatch.Match("*example.com", "slashid@example.com", fnmatch.FNM_CASEFOLD) // true
```
