import { h, PropType } from "vue";
import Anchor from "../components/Anchor.vue";
import { NButton, NPopconfirm, NSpace, NTag, NTime } from "naive-ui";

/**
 * Format duration
 * @param d milliseconds
 */
export function formatDuration(d: number): string {
    // 3h4m21.392s / 3h0.050s / 3s / 0.003s
    var h: any = Math.floor(d / 3600000)
    var m: any = Math.floor((d - h * 3600000) / 60000)
    var s: any = Math.floor((d - h * 3600000 - m * 60000) / 1000)
    var ms: any = d % 1000
    var r = ''
    if (h > 0) { r += h + 'h' }
    if (m > 0) { r += m + 'm' }
    if (s > 0 || ms > 0) {
        r += s
        if (ms > 0) {
            if (ms < 10) {
                r += ".00" + ms;
            } else if (ms < 100) {
                r += ".0" + ms;
            } else {
                r += "." + ms;
            }
        }
        r += 's'
    }
    return r

    //-- 03:04:21.392 / 00:00:00.003 / 00:00:00.300 --
    // var h: any = Math.floor(d / 3600000)
    // var m: any = Math.floor((d - h * 3600000) / 60000)
    // var s: any = Math.floor((d - h * 3600000 - m * 60000) / 1000)
    // var ms: any = d % 1000
    // if (h < 10) { h = "0" + h; }
    // if (m < 10) { m = "0" + m; }
    // if (s < 10) { s = "0" + s; }
    // if (ms < 10) {
    //     ms = "00" + ms;
    // } else if (ms < 100) {
    //     ms = "0" + ms;
    // }
    // return h + ':' + m + ':' + s + '.' + ms
}

export function renderLink(href: string, text: string) {
    return h(Anchor, { href }, { default: () => text })
}

export function renderTag(
    text: string,
    type: "default" | "error" | "info" | "success" | "warning" = "default",
    size: "small" | "medium" | "large" = "small"
) {
    return h(NTag, { type, size, round: true }, { default: () => text, })
}

export function renderTime(time: number) {
    return h(NTime, { time, format: "y-MM-dd HH:mm:ss" })
}

export interface Button {
    type: "default" | "primary" | "error" | "info" | "success" | "warning",
    text: string,
    action: (e: MouseEvent) => void,
    prompt?: string,
}

export function renderButtons(btns: Array<Button>) {
    return h(
        NSpace,
        { size: 4 },
        {
            default: () => btns.map(btn => renderButton(btn.type, btn.text, btn.action, btn.prompt)),
        }
    );
}

export function renderButton(
    type: "default" | "primary" | "error" | "info" | "success" | "warning",
    text: string,
    action: (e: MouseEvent) => void,
    prompt?: string,
) {
    if (prompt) {
        return h(
            NPopconfirm,
            {
                onPositiveClick: action,
            },
            {
                default: () => prompt,
                trigger: () => renderBtn(type, text),
            }
        )
    }
    return renderBtn(type, text, action)
}

function renderBtn(
    type: "default" | "primary" | "error" | "info" | "success" | "warning",
    text: string,
    action?: (e: MouseEvent) => void,
) {
    return h(
        NButton,
        {
            size: "tiny",
            ghost: true,
            type: type,
            onClick: action,
        },
        { default: () => text }
    )
}
