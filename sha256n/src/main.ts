import { sha256 } from 'sha.js'

async function load (): Promise<[Promise<void>]> {
  const go = new Go()
  const fetching = fetch('main.wasm')
  const result = await WebAssembly.instantiateStreaming(fetching, go.importObject)
  const running = go.run(result.instance)

  while (!('go_wasm' in globalThis)) {
    await new Promise(resolve => setTimeout(resolve))
  }

  return [running]
}

function timeit(it: () => void): number {
  const start = Date.now()
  it()
  const end = Date.now()

  return end - start
}

async function run(): Promise<void> {
  const [running] = await load()

  const toHash = Uint8Array.of(1, 2, 3)
  const iterations = 100

  const inWasm = timeit(() => {
    const ret = go_wasm.sha256n(toHash, iterations)
    if (ret instanceof Error) {
      throw ret
    }
  })
  console.log(`time running in Wasm: ${inWasm}`)

  const inJS = timeit(() => {
    let disgested = Buffer.from(toHash);
    for (let i = 0; i < iterations; i++) {
      const hasher = new sha256();
      hasher.update(disgested)
      disgested = hasher.digest()
    }
  })
  console.log(`time running in JS: ${inJS}`)

  await running
}

run().catch(console.error)
