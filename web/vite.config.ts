import { reactRouter } from "@react-router/dev/vite";
import tailwindcss from "@tailwindcss/vite";
import path from "path";
import { defineConfig } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";

export default defineConfig({
    plugins: [reactRouter(), tailwindcss(), tsconfigPaths()],
    resolve: {
        alias: {
            "@": path.resolve(__dirname, "./app"),
        }
    },
    server: {
        proxy: {
            "/api": {
                target: "http://127.0.0.1:4200",
                changeOrigin: true,
            }
        }
    },
});
