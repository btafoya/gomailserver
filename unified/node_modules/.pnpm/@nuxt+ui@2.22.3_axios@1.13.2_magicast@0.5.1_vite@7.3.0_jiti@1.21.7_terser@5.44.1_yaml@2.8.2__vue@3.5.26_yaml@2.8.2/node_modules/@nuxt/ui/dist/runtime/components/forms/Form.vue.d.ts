import type { PropType } from 'vue';
import type { Schema as JoiSchema } from 'joi';
import type { ObjectSchema as YupObjectSchema } from 'yup';
import type { StandardSchemaV1 } from '@standard-schema/spec';
import type { Struct } from 'superstruct';
import type { FormError, FormEventType, ValidateReturnSchema } from '../../types/form.js';
type Schema = PropType<StandardSchemaV1> | PropType<YupObjectSchema<any>> | PropType<Struct<any, any>> | PropType<JoiSchema>;
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    schema: {
        type: Schema;
        default: any;
    };
    state: {
        type: ObjectConstructor;
        required: true;
    };
    validate: {
        type: PropType<(state: any) => Promise<FormError[]>> | PropType<(state: any) => FormError[]>;
        default: () => any[];
    };
    validateOn: {
        type: PropType<FormEventType[]>;
        default: () => string[];
    };
}>, {
    onSubmit: (payload: Event) => Promise<void>;
    errors: Readonly<import("vue").Ref<readonly {
        readonly path: string;
        readonly message: string;
    }[], readonly {
        readonly path: string;
        readonly message: string;
    }[]>>;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, ("error" | "submit")[], "error" | "submit", import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    schema: {
        type: Schema;
        default: any;
    };
    state: {
        type: ObjectConstructor;
        required: true;
    };
    validate: {
        type: PropType<(state: any) => Promise<FormError[]>> | PropType<(state: any) => FormError[]>;
        default: () => any[];
    };
    validateOn: {
        type: PropType<FormEventType[]>;
        default: () => string[];
    };
}>> & Readonly<{
    onSubmit?: (...args: any[]) => any;
    onError?: (...args: any[]) => any;
}>, {
    schema: StandardSchemaV1<unknown, unknown> | YupObjectSchema<any, import("yup").AnyObject, any, ""> | Struct<any, any> | import("joi").AnySchema<any>;
    validate: ((state: any) => Promise<FormError[]>) | ((state: any) => FormError[]);
    validateOn: FormEventType[];
}, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
export declare function isStandardSchema(schema: any): schema is StandardSchemaV1;
export declare function validateStandardSchema(state: any, schema: StandardSchemaV1): Promise<ValidateReturnSchema<typeof state>>;
