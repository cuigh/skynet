import ajax, { Result } from './ajax'

export interface Role {
    id: string;
    name: string;
    desc: string;
    perms: string[];
    create_time: number;
    modify_time: number;
}

export class RoleApi {
    find(id: string) {
        return ajax.get<Role>('/role/find', { id })
    }

    search(name?: string) {
        return ajax.get<Role[]>('/role/search', { name })
    }

    save(role: Role) {
        return ajax.post<Result<Object>>('/role/save', role)
    }

    delete(id: string) {
        return ajax.post<Result<Object>>('/role/delete',  { id })
    }
}

export default new RoleApi
