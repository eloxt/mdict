import { type RouteConfig, index, layout } from "@react-router/dev/routes";

export default [
    layout("./components/sidebar-layout/sidebar-layout.tsx", [
        index("./main.tsx")
    ]),
] satisfies RouteConfig;