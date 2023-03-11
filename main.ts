import { write } from "./src/archive/zip.ts";
import { connect } from "./src/azuredevops/index.ts";
// Learn more at https://deno.land/manual/examples/module_metadata#concepts
if (import.meta.main) {
  await connect()
}
