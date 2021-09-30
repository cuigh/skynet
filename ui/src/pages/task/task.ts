export const alerts = [
    { value: "email", text: "邮件" },
    { value: "wecom", text: "企业微信" },
]

export function alertText(type: string) {
    return alerts.find(a => a.value === type)?.text
}