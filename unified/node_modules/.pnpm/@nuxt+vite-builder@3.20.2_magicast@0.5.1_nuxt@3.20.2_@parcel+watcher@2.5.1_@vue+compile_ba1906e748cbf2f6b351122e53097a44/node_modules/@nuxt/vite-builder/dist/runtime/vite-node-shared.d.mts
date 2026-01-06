/** @typedef {import('node:net').Socket} Socket */
/** @typedef {import('../plugins/vite-node').ViteNodeFetch} ViteNodeFetch */
/** @type {import('../plugins/vite-node').ViteNodeServerOptions} */
export const viteNodeOptions: import("../plugins/vite-node").ViteNodeServerOptions;
/**
 * @type {ViteNodeFetch}
 */
export const viteNodeFetch: ViteNodeFetch;
export type Socket = import("node:net").Socket;
export type ViteNodeFetch = import("../plugins/vite-node").ViteNodeFetch;
