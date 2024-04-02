<script lang="ts" setup>
import {
	NCard,
	NIcon,
	NSpace,
	NButton,
	NAlert,
	NDescriptions,
	NDescriptionsItem,
	NA,
	NPopconfirm,
} from "naive-ui";
import {
	SignOut20Regular,
	Open16Regular,
} from "@vicons/fluent";
import { computed, ref } from "vue";
import axios from "axios";
import { useRouter } from "vue-router";

import QuitAppButton from "@/components/QuitAppButton.vue";
import type { UserInfoResponse } from "@/lib/responses";

const router = useRouter();

const userInfo = ref<null | UserInfoResponse>(null);

async function refresh() {
	const resp = await axios.get("/api/userinfo");
	userInfo.value = resp.data as UserInfoResponse;
}

async function logout() {
	await axios.post("/api/logout");
	router.push({ name: "login" });
}

refresh();
</script>

<template>
	<n-card title="Information" :segmented="true">
		<template #header-extra>
			<quit-app-button />
		</template>

		<n-descriptions v-if="userInfo !== null && userInfo.loggedin" title="Currently logged in as">
			<n-descriptions-item label="Email">
				{{ userInfo.username }}
			</n-descriptions-item>
			<n-descriptions-item label="Device Name">
				{{ userInfo.devicename }}
			</n-descriptions-item>
		</n-descriptions>

		<n-alert v-if="userInfo !== null && !userInfo.loggedin" type="warning">
			You are not logged in.
		</n-alert>

		<template #footer>
			<n-space justify="end" v-if="userInfo !== null">
				<n-a href="https://pushover.net/settings">
					<n-button v-if="userInfo.loggedin">
						<template #icon>
							<n-icon><Open16Regular /></n-icon>
						</template>
						Manage account
					</n-button>
				</n-a>
				<n-popconfirm v-if="userInfo.loggedin" :on-positive-click="logout">
					<template #trigger>
						<n-button type="error" >
							<template #icon>
								<n-icon><SignOut20Regular /></n-icon>
							</template>
							Logout
						</n-button>
					</template>
					Are you sure you want to log out?
				</n-popconfirm>
				
				<router-link v-if="!userInfo.loggedin" to="/login">
					<n-button type="primary">
						<template #icon>
							<n-icon><Open16Regular /></n-icon>
						</template>
						Login
					</n-button>
				</router-link>
			</n-space>
		</template>
	</n-card>
</template>