async function run (): Promise<void> {
  const go = new Go()
  const fetching = fetch('main.wasm')
  const result = await WebAssembly.instantiateStreaming(fetching, go.importObject)
  const running = go.run(result.instance)

  while (!('go_wasm' in globalThis)) {
    await new Promise(resolve => setTimeout(resolve))
  }

  const ret = go_wasm.increment(0)
  if (ret instanceof Error) {
    throw ret
  }
  console.log(ret)

  await running
}

run().catch(console.error)
