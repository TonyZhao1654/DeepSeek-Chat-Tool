import { modelSelect } from "../../ui/elements.js";
import { showError } from "../../ui/error-banner.js";
import { appState } from "../../state/app-state.js";
import { getSelectedModel, setModelValue } from "../models/model-selector.js";
import {
	createSession as createSessionApi,
	deleteSession as deleteSessionApi,
	getSession as fetchSession,
	updateSession,
} from "../../api/api-session.js";
import { refreshSessions, showSessionInBar } from "./session-bar.js";
import { renderContents } from "../contents/content-panel.js";
import { renderMessages } from "../chat/chat-view.js";
import { stopStream } from "../chat/chat-controller.js";

export async function openSession(id) {
	if (!id) return;
	if (appState.streaming) stopStream();
	const session = await fetchSession(id);
	appState.sessionId = session.id;
	showSessionInBar(session.id, session.title);
	setModelValue(session.model);
	renderMessages(session.messages || []);
	renderContents(session.contents || []);
}

export async function createSession() {
	const session = await createSessionApi({ model: getSelectedModel() });
	await refreshSessions();
	await openSession(session.id);
}

export async function deleteCurrentSession() {
	if (!appState.sessionId) return;
	if (!confirm("删除当前会话及其内容块、聊天记录？")) return;
	if (appState.streaming) stopStream();
	await deleteSessionApi(appState.sessionId);
	appState.sessionId = null;
	const list = await refreshSessions();
	if (list.length) {
		await openSession(list[0].id);
	} else {
		await createSession();
	}
}

export async function changeModel() {
	if (!appState.sessionId) return;
	try {
		await updateSession(appState.sessionId, { model: modelSelect.value });
	} catch (error) {
		showError("保存模型设置失败：" + error.message);
	}
}
