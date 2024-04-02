<script lang="ts" setup>
import {
	NButton,
	NIcon,
	useDialog,
} from "naive-ui";
import {
	Dismiss16Regular,
} from "@vicons/fluent";
import axios from "axios";

const dialog = useDialog();

async function quitApp() {
	dialog.warning({
		title: "Confirm",
		content: "Are you sure you want to quit the app? You won't receive notifications until you start it again.",
		positiveText: "Quit",
		negativeText: "Cancel",
		onPositiveClick: async () => {
			await axios.post("/api/quit");
			window.location.href = "about:blank";
		},
	});
}
</script>

<template>
	<n-button type="error" @click="quitApp()">
		<n-icon>
			<Dismiss16Regular />
		</n-icon>
		Quit App
	</n-button>
</template>