import { h } from 'vue'
import { NIcon, MenuOption } from 'naive-ui'
import { RouteLocationNormalizedLoaded, RouterLink } from 'vue-router'
import {
    HomeOutline as HomeIcon,
    AlarmOutline as AlarmIcon,
    NotificationsOutline as NotificationsIcon,
    PersonOutline as PersonIcon,
    PeopleOutline as PeopleIcon,
    SettingsOutline as SettingsIcon,
    DocumentTextOutline as DocumentTextIcon,
    ConstructOutline as ConstructIcon,
    KeyOutline as KeyIcon,
} from "@vicons/ionicons5";

function renderIcon(icon: any) {
    return () => h(NIcon, null, { default: () => h(icon) });
}

export const renderMenuLabel = (option: any) => {
    if (!('path' in option)) {
        return option.label
    }
    return h(
        RouterLink,
        {
            to: option.path
        },
        {
            default: () => option.label
        }
    )
}

export function findMenuValue(route: RouteLocationNormalizedLoaded): string {
    var path = route.path;
    do {
        const option = findOption(menuOptions, path)
        if (option) {
            return option.key
        } else {
            const index = path.lastIndexOf("/")
            if (index <= 0) {
                return ""
            }
            path = path.substr(0, index)
        }
    } while (true)
}

function findOption(options: MenuOption[], path: string): any {
    for (const option of options) {
        if (option.path === path) {
            return option
        } else if (option.children) {
            const opt = findOption(option.children, path)
            if (opt) return opt
        }
    }
    return null
}

export function findActiveOptions(route: RouteLocationNormalizedLoaded): MenuOption[] {
    const result: MenuOption[] = []
    findOptions(result, menuOptions, route.path)
    return result
}

function findOptions(result: MenuOption[], options: MenuOption[], path: string): boolean {
    for (const option of options) {
        if (option.path) {
            if (option.path != "/" && path.startsWith(<string>option.path)) {
                result.push(option)
                return true
            }
        } else if (option.children) {
            result.push(option)
            if (findOptions(result, option.children, path)) {
                return true
            } else {
                result.pop()
            }
        }
    }
    return false
}

export const menuOptions: MenuOption[] = [
    {
        label: "首页",
        key: "home",
        path: "/",
        icon: renderIcon(HomeIcon),
    },
    {
        label: "任务管理",
        key: "tasks",
        path: "/tasks",
        icon: renderIcon(AlarmIcon),
    },
    {
        label: "作业管理",
        key: "jobs",
        path: "/jobs",
        icon: renderIcon(DocumentTextIcon),
    },
    {
        label: "账号管理",
        key: "account",
        // path: "/account",
        icon: renderIcon(PersonIcon),
        children: [
            {
                label: "用户管理",
                key: "users",
                path: "/account/users",
                icon: renderIcon(PersonIcon),
            },
            {
                label: "角色管理",
                key: "roles",
                path: "/account/roles",
                icon: renderIcon(PeopleIcon),
            },         
            // {
            //     label: "凭证管理",
            //     key: "tokens",
            //     path: "/account/tokens",
            //     icon: renderIcon(KeyIcon),
            // },            
        ],
    },
    {
        label: "系统设置",
        key: "config",
        icon: renderIcon(SettingsIcon),
        children: [
            {
                label: "通知设置",
                key: "notice",
                path: "/config/notice",
                icon: renderIcon(NotificationsIcon),
            },
            {
                label: "高级设置",
                key: "advance",
                path: "/config/advance",
                icon: renderIcon(ConstructIcon),
            },            
        ],
    },
]
