type InputProps = {
    id?: string;
    size?: string | number | symbol;
    color?: string;
    name?: string;
    eagerValidation?: boolean;
    legend?: string | null;
};
export declare const useFormGroup: (inputProps?: InputProps, config?: any, bind?: boolean) => {
    inputId: import("vue").ComputedRef<string>;
    name: import("vue").ComputedRef<string>;
    size: import("vue").ComputedRef<any>;
    color: import("vue").ComputedRef<string>;
    emitFormBlur: () => void;
    emitFormInput: import("@vueuse/core").UseDebounceFnReturn<() => void>;
    emitFormChange: () => void;
};
export {};
