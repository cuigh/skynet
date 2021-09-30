export function statusType(status: number) {
    switch (status) {
        case 0:
            return "default"
        case 1:
            return "success"
        case 2:
            return "error"
        default:
            return "warning"
    }
}

export function statusText(status: number) {
    switch (status) {
        case 0:
            return "未知"
        case 1:
            return "成功"
        case 2:
            return "失败"
        default:
            return `异常[${status}]`
    }
}