import { compress } from "../deps.ts";

export type ZipFileName = `${string}.zip`
export async function write(dir: string, destFileName: ZipFileName) {
  await compress([dir], destFileName);
}
