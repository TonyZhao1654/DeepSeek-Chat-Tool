import * as sessionApi from "../../api/api-session.js";
import { sessionSelect, sessionTitle } from "../../ui/elements.js";
import { showError } from "../../ui/error-banner.js";
import { appState } from "../../state/app-state.js";

export async function refreshSessions() {
	const list = await sessionApi.listSessions();
	sessionSelect.innerHTML = "";
	(list || []).forEach((s) => {
		sessionSelect.add(new Option(s.title || s.id, s.id));
	});
	return list || [];
}

export function showSessionInBar(id, title) {
	sessionSelect.value = id;
	sessionTitle.value = title || "";
}

export function updateSessionOption(id, title) {
	const option = [...sessionSelect.options].find((o) => o.value === id);
	if (option) option.textContent = title;
}

export async function renameCurrent() {
	if (!appState.sessionId) return;
	const title = sessionTitle.value.trim();
	if (!title) return;
	try {
		await sessionApi.updateSession(appState.sessionId, { title });
		updateSessionOption(appState.sessionId, title);
	} catch (error) {
		showError("重命名失败：" + error.message);
	}
}

export async function autoTitle(text) {
	if (!appState.sessionId) return;
	const current = sessionTitle.value.trim();
	if (current && current !== "新会话") return;
	const sessionId = appState.sessionId;
	const title = text.replace(/\s+/g, " ").slice(0, 20);
	if (!title) return;
	try {
		await sessionApi.updateSession(sessionId, { title });
		if (appState.sessionId !== sessionId) return;
		sessionTitle.value = title;
		updateSessionOption(sessionId, title);
	} catch {}
}
