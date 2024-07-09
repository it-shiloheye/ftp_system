import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react-swc";
import { TanStackRouterVite } from "@tanstack/router-plugin/vite";
import tailwindcss from "tailwindcss";

// https://vitejs.dev/config/
/** @type {import('vite').UserConfig} */
export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, process.cwd(), "");
	return {
		define: {
			"process.env.BROWSER_PORT": JSON.stringify(env.BROWSER_PORT),
			// 'process.env.YOUR_BOOLEAN_VARIABLE': env.YOUR_BOOLEAN_VARIABLE,
			// If you want to exposes all env variables, which is not recommended
			// 'process.env': env
		},
		envDir: ".",
		server: {
      strictPort:true,
			proxy: {
				"^/api": {
					target: `http://127.0.0.1:${env.BROWSER_PORT}`,
					changeOrigin: true,
				},
			},
		},
		plugins: [TanStackRouterVite(), react()],
		css: {
			postcss: {
				plugins: [tailwindcss()],
			},
		},
	};
});
