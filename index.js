export const wasmBrowserInstantiate = async (wasmModuleUrl, importObject) => {
  let response = undefined;

  if (WebAssembly.instantiateStreaming) {
    response = await WebAssembly.instantiateStreaming(
      fetch(wasmModuleUrl),
      importObject
    );
  } else {
    const fetchAndInstantiateTask = async () => {
      const wasmArrayBuffer = await fetch(wasmModuleUrl).then((response) =>
        response.arrayBuffer()
      );
      return WebAssembly.instantiate(wasmArrayBuffer, importObject);
    };

    response = await fetchAndInstantiateTask();
  }

  return response;
};

const encodeString = (name, instance) => {
  const bytes = new TextEncoder("utf8").encode(name);
  const ptr = instance.exports.alloc(bytes.length);
  const mem = new Uint8Array(instance.exports.memory.buffer, ptr, bytes.length);
  mem.set(new Uint8Array(bytes));

  return {
    ptr: ptr,
    length: bytes.length,
  };
};

const go = new Go();

const load = async () => {
  const wasmModule = await wasmBrowserInstantiate(
    "./main.wasm",
    go.importObject
  );

  go.run(wasmModule.instance);

  return {
    Match: (pattern, s, flags) => {
      const _pattern = encodeString(pattern, wasmModule.instance);
      const _s = encodeString(s, wasmModule.instance);
      return Boolean(
        wasmModule.instance.exports.Match(
          _pattern.ptr,
          _pattern.length,
          _s.ptr,
          _s.length,
          flags
        )
      );
    },
    FNM_NOESCAPE: wasmModule.instance.exports._FNM_NOESCAPE(),
    FNM_PATHNAME: wasmModule.instance.exports._FNM_PATHNAME(),
    FNM_PERIOD: wasmModule.instance.exports._FNM_PERIOD(),
    FNM_LEADING_DIR: wasmModule.instance.exports._FNM_LEADING_DIR(),
    FNM_CASEFOLD: wasmModule.instance.exports._FNM_CASEFOLD(),
    FNM_IGNORECASE: wasmModule.instance.exports._FNM_IGNORECASE(),
    FNM_FILE_NAME: wasmModule.instance.exports._FNM_FILE_NAME(),
  };
};

export const fnmatch = await load();

window.fnmatch = fnmatch
