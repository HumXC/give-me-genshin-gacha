export interface Option {
    isShow: boolean;
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

export interface PieData {
    几发未出金: number;
    三星武器: number;
    四星角色: number;
    四星武器: number;
    五星角色: number;
    五星武器: number;
    name: string;
}
