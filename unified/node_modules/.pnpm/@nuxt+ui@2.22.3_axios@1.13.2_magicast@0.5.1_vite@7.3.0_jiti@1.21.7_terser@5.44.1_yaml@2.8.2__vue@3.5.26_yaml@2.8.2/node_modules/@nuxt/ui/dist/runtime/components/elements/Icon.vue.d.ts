import type { PropType } from 'vue';
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    name: {
        type: StringConstructor;
        required: true;
    };
    mode: {
        type: PropType<"svg" | "css">;
        required: false;
        default: any;
    };
    size: {
        type: (StringConstructor | NumberConstructor)[];
        required: false;
        default: any;
    };
    customize: {
        type: FunctionConstructor;
        required: false;
        default: any;
    };
}>, {}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    name: {
        type: StringConstructor;
        required: true;
    };
    mode: {
        type: PropType<"svg" | "css">;
        required: false;
        default: any;
    };
    size: {
        type: (StringConstructor | NumberConstructor)[];
        required: false;
        default: any;
    };
    customize: {
        type: FunctionConstructor;
        required: false;
        default: any;
    };
}>> & Readonly<{}>, {
    mode: "svg" | "css";
    size: string | number;
    customize: Function;
}, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
