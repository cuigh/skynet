import { isRef, onMounted, reactive } from "vue"

export function useDataTable(loader: Function, filter: Object | Function) {
    const state = reactive({
        loading: false,
        data: [],
    })
    const pagination = reactive({
        page: 1,
        pageCount: 1,
        pageSize: 10,
        itemCount: 0,
        prefix({ itemCount }: any) {
            return `共 ${itemCount} 项`
        }
    })
    const fetchData = async function (page: number = 1) {
        state.data = [];
        state.loading = true;
        try {
            let args = typeof filter === 'function' ? filter() : filter
            args = isRef(args) ? args.value : args
            let r = await loader({
                ...args,
                page_index: page,
                page_size: pagination.pageSize,
            });
            state.data = r.data?.items || [];
            pagination.itemCount = r.data?.total || 0
            pagination.page = page
            pagination.pageCount = Math.ceil(pagination.itemCount / pagination.pageSize)
        } finally {
            state.loading = false;
        }
    }

    onMounted(fetchData)

    return { state, pagination, fetchData }
}
