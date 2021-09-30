import ajax, { Result } from './ajax'

export interface Job {
    id: string;
    task: string;
    handler: string;
    scheduler: string;
    mode: number;
    fire_time: number;
    args?: {
        name: string;
        value: string;
    }[];
    dispatch: {
        status: number;
        error?: string;
        time: number;
    },
    execute: {
        status: number;
        error?: string;
        end_time: number;
        start_time: number;
    },
}

export interface SearchArgs {
    task?: string;
    mode?: number;
    dispatch_status?: number;
    execute_status?: number;
    page_index: number;
    page_size: number;
}

export interface SearchResult {
    items: Job[];
    total: number;
}

export class JobApi {
    find(id: string) {
        return ajax.get<Job>('/job/find', { id })
    }

    search(args: SearchArgs) {
        if (args.mode == null) {
            args.mode = -1
        }
        if (args.dispatch_status == null) {
            args.dispatch_status = -1
        }
        if (args.execute_status == null) {
            args.execute_status = -1
        }
        return ajax.get<SearchResult>('/job/search', args)
    }

    retry(id: string) {
        return ajax.post<Result<Object>>('/job/retry', { id })
    }
}

export default new JobApi
