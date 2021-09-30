// export interface BasicPageParams {
//     pageIndex: number;
//     pageSize: number;
// }

// export interface BasicFetchResult<T extends any> {
//     items: T[];
//     totalCount: number;
// }

// type Name = string; // 基本类型
// type NameResolver = () => string; // 函数
// type NameOrResolver = Name | NameResolver; // 联合类型

// interface StringArray {
//     [index: number]: string;
// }

// export type RolePageParams = BasicPageParams & RoleParams;

// export type AccountParams = BasicPageParams & {
//     account?: string;
//     nickname?: string;
//   };