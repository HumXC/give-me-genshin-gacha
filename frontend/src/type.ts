// 配置文件
export class Option {
    showGacha: {
        roleUp: boolean;
        armsUp: boolean;
        permanent: boolean;
        start: boolean;
    } = {
        roleUp: false,
        armsUp: false,
        permanent: false,
        start: false,
    };
    otherOption: {
        autoSync: boolean;
        useProxy: boolean;
        darkTheme: boolean;
    } = {
        autoSync: false,
        useProxy: false,
        darkTheme: false,
    };
    controlBar: {
        selectedUid: string;
    } = {
        selectedUid: "",
    };
}

// 饼图使用的数据
export class PieData {
    usedCosts: Array<{
        gachaType: string;
        cost: number;
    }> = []; // 几发未出金
    tatols: {
        t301: Array<{
            total: number;
            itemType: string;
            rankType: string;
        }>;
        t302: Array<{
            total: number;
            itemType: string;
            rankType: string;
        }>;
        t200: Array<{
            total: number;
            itemType: string;
            rankType: string;
        }>;
        t100: Array<{
            total: number;
            itemType: string;
            rankType: string;
        }>;
    } = {
        t301: [],
        t302: [],
        t200: [],
        t100: [],
    };
}
// 一条完整的祈愿数据
export interface GachaLog {
    gacha_type: string;
    time: number;
    name: string;
    lang: string;
    item_type: string;
    rank_type: string;
    id: string;
}

// 祈愿类型
export type GachaType = "301" | "302" | "200" | "100" | "400";
export const GachaTypeWithName = new Map([
    // 前端一般不做 400 的处理
    ["400", "角色活动祈愿"],
    ["301", "角色活动祈愿"],
    ["302", "武器活动祈愿"],
    ["200", "常驻祈愿"],
    ["100", "新手祈愿"],
]);
export type Message = {
    type: string;
    msg: string;
};
