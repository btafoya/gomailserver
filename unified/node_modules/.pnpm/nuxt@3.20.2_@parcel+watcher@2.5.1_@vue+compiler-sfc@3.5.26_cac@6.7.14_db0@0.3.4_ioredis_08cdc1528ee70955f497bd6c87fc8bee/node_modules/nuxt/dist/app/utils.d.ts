/** @since 3.9.0 */
export declare function toArray<T>(value: T | T[]): T[];
export declare function getUserTrace(): {
    source: string;
    column?: number;
    function?: string;
    line?: number;
}[];
export declare function getUserCaller(): {
    source: string;
    line: number | undefined;
    column: number | undefined;
} | null;
