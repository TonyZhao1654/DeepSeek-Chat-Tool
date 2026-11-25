import {
	addContentButton,
	deleteSessionButton,
	inputElement,
	modelSelect,
	newSessionButton,
	sendButton,
	sessionSelect,
	sessionTitle,
} from "./ui/elements.js";
import { showError } from "./ui/error-banner.js";
import { loadModels } from "./features/models/model-selector.js";
import { refreshBalance } from "./features/account/balance.js";
import { refreshSessions, renameCurrent } from "./features/sessions/session-bar.js";
import { addContent, renderContents } from "./features/contents/content-panel.js";
import { send } from "./features/chat/chat-controller.js";
import * as sessionController from "./features/sessions/session-controller.js";

function bindEvents() {
	inputElement.addEventListener("keydown", (e) => {
		if (e.key === "Enter" && !e.shiftKey) {
			e.preventDefault();
			send();
		}
	});
	sendButton.addEventListener("click", send);
	sessionSelect.addEventListener("change", async () => {
		try { await sessionController.openSession(sessionSelect.value); }
		catch (err) { showError("打开会话失败：" + err.message); }
	});
	sessionTitle.addEventListener("change", renameCurrent);
	modelSelect.addEventListener("change", sessionController.changeModel);
	newSessionButton.addEventListener("click", async () => {
		try { await sessionController.createSession(); }
		catch (err) { showError("新建会话失败：" + err.message); }
	});
	deleteSessionButton.addEventListener("click", async () => {
		try { await sessionController.deleteCurrentSession(); }
		catch (err) { showError("删除会话失败：" + err.message); }
	});
	addContentButton.addEventListener("click", addContent);
}

async function init() {
	bindEvents();
	await loadModels();
	await refreshBalance();
	try {
		const list = await refreshSessions();
		if (list.length) await sessionController.openSession(list[0].id);
		else await sessionController.createSession();
	} catch (err) {
		showError("会话初始化失败：" + err.message);
		renderContents([]);
	}
}

init();
