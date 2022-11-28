// 配置文件
export interface Option {
    showGacha: {
        roleUp: boolean;
        armsUp: boolean;
        permanent: boolean;
        start: boolean;
    };
    otherOption: {
        autoSync: boolean;
        useProxy: boolean;
        darkTheme: boolean;
    };
}

// 饼图使用的数据
export interface PieData {
    usedCost: number; // 几发未出金
    arms3Total: number; // 三星武器数量
    role4Total: number; // 四星角色数量
    arms4Total: number; // 四星武器数量
    role5Total: number; // 五星角色数量
    arms5Total: number; // 五星武器数量
    gachaType: string; // 祈愿类型
}

// 一条完整的祈愿数据
export interface GachaLog {
    uid: string;
    gacha_type: string;
    item_id: string;
    count: string;
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
