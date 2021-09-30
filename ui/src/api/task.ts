import ajax, { Result } from './ajax'

export interface Task {
    name: string;
    runner: string;
    handler: string;
    triggers: string[];
    desc?: string;
    args?: {
        name: string;
        value: string;
    }[];
    enabled: boolean;
    alerts: string[];
    maintainers?: string[];
}

export interface SearchArgs {
    name?: string;
    runner?: string;
    page_index: number;
    page_size: number;
}

export interface SearchResult {
    items: Task[];
    total: number;
}

export interface ExecuteArgs {
    name: string;
    args?: {
        name: string;
        value: string;
    }[];
}

export class TaskApi {
    find(name: string) {
        return ajax.get<Task>('/task/find', { name })
    }

    search(args: SearchArgs) {
        return ajax.get<SearchResult>('/task/search', args)
    }

    save(task: Task) {
        return ajax.post<Result<Object>>('/task/save', task)
    }

    delete(name: string) {
        return ajax.post<Result<Object>>('/task/delete', { name })
    }

    execute(args: ExecuteArgs) {
        return ajax.post<Result<Object>>('/task/execute', args)
    }
}

export default new TaskApi
