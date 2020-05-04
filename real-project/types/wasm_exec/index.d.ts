declare class Go {
  readonly importObject: Record<string, Record<string, WebAssembly.ImportValue>>
  run(_: WebAssembly.Instance): Promise<void>
}
