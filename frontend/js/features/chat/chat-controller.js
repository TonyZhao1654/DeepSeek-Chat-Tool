import { openChatStream } from "../../api/api-deepseek.js";
import { readErrorResponse } from "../../api/http.js";
import { inputElement, modelSelect, sendButton } from "../../ui/elements.js";
import { clearError, showError } from "../../ui/error-banner.js";
import { appState } from "../../state/app-state.js";
import { flushContentSaves } from "../contents/content-panel.js";
import { formatBalance, updateBalance } from "../account/balance.js";
import { autoTitle } from "../sessions/session-bar.js";
import {
	appendMessage,
	appendMeta,
	createAssistantBubble,
	renderAssistantContent,
	scrollToBottom,
} from "./chat-view.js";
import { readChatStream } from "./chat-stream.js";

export function stopStream() {
	if (appState.controller) appState.controller.abort();
}

export async function send() {
	if (appState.streaming) {
		stopStream();
		return;
	}
	const text = inputElement.value.trim();
	if (!text) return;
	if (!appState.sessionId) {
		showError("会话尚未就绪，请稍后再试");
		return;
	}
	if (appState.pendingSend) return;
	appState.pendingSend = true;
	try { await flushContentSaves(); } finally { appState.pendingSend = false; }

	clearError();
	appendMessage("user", text);
	inputElement.value = "";
	autoTitle(text);

	const sendSessionId = appState.sessionId;
	const bubble = createAssistantBubble();
	let content = "";
	let reasoning = "";
	let finished = false;

	const handleEvent = (parsed) => {
		if (parsed.done) {
			finished = true;
			applyDone(bubble, content, parsed.data);
			return;
		}
		const chunk = parsed.chunk;
		if (!chunk || !chunk.choices || !chunk.choices[0]) return;
		const delta = chunk.choices[0].delta || {};
		if (delta.reasoning_content) {
			reasoning += delta.reasoning_content;
			bubble.reasoning.textContent = reasoning;
			bubble.reasoning.hidden = false;
		}
		if (delta.content) {
			content += delta.content;
			renderAssistantContent(bubble, content);
		}
		scrollToBottom();
	};

	setStreaming(true);
	appState.controller = new AbortController();

	try {
		const response = await openChatStream(
			{
				session_id: sendSessionId,
				model: modelSelect.value,
				messages: [{ role: "user", content: text }],
			},
			appState.controller.signal
		);
		const contentType = response.headers.get("Content-Type") || "";
		if (!response.ok || !contentType.includes("text/event-stream")) {
			const msg = await readErrorResponse(response);
			bubble.wrapper.remove();
			throw new Error(msg);
		}
		await readChatStream(response, handleEvent);
	} catch (error) {
		if (error.name === "AbortError") {
			if (!content) bubble.wrapper.remove();
			if (appState.sessionId === sendSessionId) appendMeta("已停止生成");
		} else {
			showError(error.message);
		}
	} finally {
		if (!finished && bubble.thinking.isConnected) bubble.thinking.remove();
		setStreaming(false);
		appState.controller = null;
		inputElement.focus();
	}
}

function setStreaming(streaming) {
	appState.streaming = streaming;
	sendButton.textContent = streaming ? "停止" : "发送";
	sendButton.classList.toggle("stopping", streaming);
	inputElement.disabled = streaming;
}

function applyDone(bubble, content, data) {
	if (bubble.thinking.isConnected) bubble.thinking.remove();
	const finalContent = data && data.choices && data.choices[0] && data.choices[0].message
		? data.choices[0].message.content : "";
	if (!content && finalContent) renderAssistantContent(bubble, finalContent);

	const usage = (data && data.usage) || {};
	const parts = [];
	if (usage.total_tokens > 0) {
		parts.push(`输入 ${usage.prompt_tokens} tokens`);
		parts.push(`输出 ${usage.completion_tokens} tokens`);
		if (data && typeof data.cost === "number" && data.cost > 0) {
			parts.push(`花费 ¥${data.cost.toFixed(6)}`);
		}
	}
	if (data && data.balance && data.balance.total_balance) {
		updateBalance(data.balance.total_balance);
		parts.push(`余额 ${formatBalance(data.balance.total_balance)} 元`);
	}
	if (data && data.incomplete) {
		parts.push(data.stream_error ? `回答异常结束：${data.stream_error}` : "回答异常结束");
	}
	if (parts.length) appendMeta(parts.join(" · "));
}
