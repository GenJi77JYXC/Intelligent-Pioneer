import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const AGENT: AppRouteRecordRaw = {
    // --- 父级路由 (一级菜单) ---
    path: '/agent',       // URL 路径
    name: 'agent',        // 路由名称，保持唯一
    component: DEFAULT_LAYOUT, // 使用默认的布局组件 (带侧边栏和顶部导航)
    meta: {
        locale: 'menu.agent', // 国际化文本的 key (下一步会定义)
        icon: 'icon-desktop', // 菜单图标 (可以在 Arco 图标库中查找)
        requiresAuth: true,     // 表示需要登录才能访问
        order: 1,               // 菜单排序，数字越小越靠前
    },

    // --- 子级路由 (二级菜单) ---
    children: [
        {
            path: 'list', // 子路径，完整的 URL 是 /agent/list
            name: 'AgentList', // 子路由名称
            component: () => import('@/views/agent/list/index.vue'), // 懒加载页面组件 (下一步会创建)
            meta: {
                locale: 'menu.agent.list', // 子菜单的国际化 key
                requiresAuth: true,
                roles: ['*'], // 表示所有角色都可以访问
            },
        },
        // 未来如果还有其他页面，比如 “Agent详情”，可以继续在这里添加
    ],
};

export default AGENT;