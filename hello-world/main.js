async function run () {
  const go = new Go()
  const fetching = fetch('main.wasm')
  const result = await WebAssembly.instantiateStreaming(fetching, go.importObject)
  await go.run(result.instance)
}

run().catch(console.error)
