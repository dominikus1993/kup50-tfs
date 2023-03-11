import { write } from "./src/archive/zip.ts";
// Learn more at https://deno.land/manual/examples/module_metadata#concepts
if (import.meta.main) {
  await write("kup", "test.zip");
  console.log("Add 2 + 3 =");
}
