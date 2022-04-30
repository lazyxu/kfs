import './wasm_exec';

if (!WebAssembly.instantiateStreaming) { // polyfill
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return WebAssembly.instantiate(source, importObject);
  };
}

const go = new window.Go();
let mod, inst;

export async function run() {
  if (!inst) {
    const result = await WebAssembly.instantiateStreaming(fetch('extraResources/lib.wasm'), go.importObject);
    mod = result.module;
    inst = result.instance;
  }
  await go.run(inst);
  console.log('wasm', inst, mod);
  inst = await WebAssembly.instantiate(mod, go.importObject); // reset instance
}
