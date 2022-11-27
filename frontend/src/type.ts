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
    几发未出金: number;
    三星武器: number;
    四星角色: number;
    四星武器: number;
    五星角色: number;
    五星武器: number;
    name: string;
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
