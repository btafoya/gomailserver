import type { Strategy } from '../types/index.js';
export declare const twMerge: (...classLists: import("tailwind-merge").ClassNameValue[]) => string;
export declare function mergeConfig<T>(strategy: Strategy, ...configs: any[]): T;
export declare function parseConfigValue(value: string): string;
export declare function hexToRgb(hex: string): string;
export declare function getSlotsChildren(slots: any): any;
/**
 * "123-foo" will be parsed to 123
 * This is used for the .number modifier in v-model
 */
export declare function looseToNumber(val: any): any;
export * from './lodash.js';
export * from './link.js';
