import { createRouter, createWebHistory } from "vue-router";

import LoginView from "./views/LoginView.vue";
import InfoView from "./views/InfoView.vue";
import RegisterDeviceView from "./views/RegisterDeviceView.vue";

const router = createRouter({
	history: createWebHistory(),
	routes: [
		{
			path: "/login",
			component: LoginView,
			name: "login",
		},
		{
			path: "/info",
			component: InfoView,
			name: "info",
		},
		{
			path: "/register-device",
			component: RegisterDeviceView,
			name: "register-device",
		}
	]
});

export default router;