import ajax, { Result } from './ajax'

export interface EmailOptions {
    enabled: string,
    sender: string,
    receiver?: string,
    smtp_address: string,
    smtp_username: string,
    smtp_password: string,
    title?: string,
    body?: string,
}

export interface WecomOptions {
    enabled: string,
    mode: string,
    robot_key: string,
    corp_id: string,
    app_id: string,
    app_secret: string,
    users?: string,
    parties?: string,
    tags?: string,
    chats?: string,
    title?: string,
    body?: string,
    msg_type: string,
}

export class ConfigApi {
    find(id: string) {
        return ajax.get<Object>('/config/find', { id })
    }

    save(id: string, options: Object) {
        return ajax.post<Result<Object>>('/config/save', { id, options })
    }
}

export default new ConfigApi
