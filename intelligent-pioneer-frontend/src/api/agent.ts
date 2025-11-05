import axios from 'axios';

// 定义 Agent 信息的 TypeScript 接口，与后端 AgentInfo 对应
export interface Agent {
    id: number;
    uuid: string;
    hostname: string;
    ip_address: string;
    os: string;
    status: 'online' | 'offline';
    created_at: string;
    updated_at: string;
}

// 定义 API 的返回类型
export type AgentListRes = Agent[];

// 封装获取 Agent 列表的请求函数
export function getAgentList() {
    return axios.get<AgentListRes>('/api/v1/agent');
}