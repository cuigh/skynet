import { Result } from "@/api/ajax";
import { Ref, ref } from "vue"
import { FormItemRule } from "naive-ui";

export function useForm<T>(form: Ref, action: () => Promise<Result<T>>, success?: (data: T) => void) {
    const submiting = ref(false)
    async function submit(e: Event) {
        e.preventDefault();
        form.value.validate(async (errors: any) => {
            if (errors) {
                return
            }

            submiting.value = true;
            try {
                let r = await action()
                success ? success(<T>r.data) : window.message.info("操作成功");
            } finally {
                submiting.value = false;
            }
        });
    }

    return { submit, submiting }
}

export function requiredRule(field?: string, message?: string): FormItemRule {
    return {
        required: true,
        message: formatMessage(field, message ?? "不能为空"),
        trigger: ["input", "blur"],
    }
}

export function customRule(validator: (rule: any, value: any) => boolean, message?: string, field?: string, required?: boolean): FormItemRule {
    return createRule(validator, message, field, required)
}


export function emailRule(field?: string): FormItemRule {
    const reg = /^([a-zA-Z0-9]+[-_\.]?)+@[a-zA-Z0-9]+\.[a-z]+$/;
    return regexRule(reg, "格式不正确", field)
}

export function phoneRule(field?: string): FormItemRule {
    const reg = /^[1][3,4,5,7,8][0-9]{9}$/;
    return regexRule(reg, "格式不正确", field)
}


export function lengthRule(min: number, max: number, field?: string): FormItemRule {
    return createRule((rule: any, value: string): boolean => {
        return value.length >= min && value.length <= max
    }, `长度必须为 ${min}-${max} 位`, field)
}

export function passwordRule(field?: string): FormItemRule {
    const reg = /^[a-zA-Z0-9_-]+$/;
    return regexRule(reg, "只能包含英文字母、数字、下划线和中划线", field)
}

export function regexRule(reg: RegExp, message?: string, field?: string): FormItemRule {
    return {
        message: formatMessage(field, message),
        trigger: ["input", "blur"],
        validator(rule: any, value: string): boolean {
            return !value || reg.test(value)
        },
    };
}

function createRule(validator: (rule: any, value: string) => boolean, message?: string, field?: string, required?: boolean): FormItemRule {
    return {
        required: required,
        message: formatMessage(field, message),
        trigger: ["input", "blur"],
        validator,
    };
}

function formatMessage(field?: string, message?: string) {
    return field ? `${field}: ${message}` : message
}