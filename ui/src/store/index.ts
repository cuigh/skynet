import { createStore, createLogger } from 'vuex'
import { Mutations } from "./mutations";

const debug = import.meta.env.DEV

export interface State {
    name: string | null;
    token: string | null;
    theme: string | null;
    ajaxLoading: boolean;
}

function initState(): State {
    return {
        name: localStorage.getItem("name"),
        token: localStorage.getItem("token"),
        theme: localStorage.getItem("theme"),
        ajaxLoading: false,
    }
}

export const store = createStore<State>({
    strict: debug,
    state: initState(),
    getters: {
        anonymous(state) {
            return !state.token
        }
    },
    mutations: {
        [Mutations.Login](state, user) {
            localStorage.setItem("name", user.name);
            localStorage.setItem("token", user.token);
            state.name = user.name;
            state.token = user.token;           
        },
        [Mutations.Logout](state) {
            localStorage.removeItem("name");
            localStorage.removeItem("token");
            state.name = null;
            state.token = null;
        },
        [Mutations.SetToken](state, token) {
            localStorage.setItem("token", token);
            state.token = token;
        },
        [Mutations.SetTheme](state, theme) {
            localStorage.setItem("theme", theme);
            state.theme = theme;
        },
        [Mutations.SetAjaxLoading](state, loading) {
            state.ajaxLoading = loading;
        },
    },
    plugins: debug ? [createLogger()] : [],
})
