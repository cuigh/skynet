import ajax, { Result } from './ajax'

export interface User {
    name: string;
    login_name: string;
    email: string;
    phone: string;
    password: string;
}

export interface State {
    fresh: boolean;
}

export interface Summary {
    version: string;
    goVersion: string;
    taskCount: number;
    jobCount: number;
    userCount: number;
}

export class SystemApi {
    checkState() {
        return ajax.get<State>('/system/check-state')
    }

    initDB() {
        return ajax.post<Result<Object>>('/system/init-db')
    }

    initUser(user: User) {
        return ajax.post<Result<Object>>('/system/init-user', user)
    }

    summarize() {
        return ajax.get<Summary>('/system/summarize')
    }
}

export default new SystemApi
