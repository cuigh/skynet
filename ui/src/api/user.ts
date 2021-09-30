import ajax, { Result } from './ajax'

export interface AuthUser {
    token: string;
    id: string;
    name: string;
}

export interface User {
    id: string;
    name: string;
    login_name: string;
    roles: string[];
    admin: boolean;
    status: number;
    email?: string;
    phone?: string;
    wecom?: string;
}

export interface LoginArgs {
    name: string;
    password: string;
}

export interface SearchArgs {
    name?: string;
    login_name?: string;
    filter?: string;
    page_index: number;
    page_size: number;
}

export interface SearchResult {
    items: User[];
    total: number;
}

export interface SetStatusArgs {
    id: string;
    status: number;
}

export interface ModifyPasswordArgs {
    old_pwd: string;
    new_pwd: string;
}

export class UserApi {
    login(args: LoginArgs) {
        return ajax.post<AuthUser>('/user/sign-in', args)
    }

    save(user: User) {
        return ajax.post<Result<Object>>('/user/save', user)
    }

    find(id: string) {
        return ajax.get<User>('/user/find', { id })
    }

    fetch(ids: string[]) {
        return ajax.get<User[]>('/user/fetch', { ids: ids.join(',') })
    }

    search(args: SearchArgs) {
        return ajax.get<SearchResult>('/user/search', args)
    }

    setStatus(args: SetStatusArgs) {
        return ajax.post<Result<Object>>('/user/set-status', args)
    }

    modifyPassword(args: ModifyPasswordArgs) {
        return ajax.post<Result<Object>>('/user/modify-password', args)
    }

    modifyProfile(user: User) {
        return ajax.post<Result<Object>>('/user/modify-profile', user)
    }
}

export default new UserApi